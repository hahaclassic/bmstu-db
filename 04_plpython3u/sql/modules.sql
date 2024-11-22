CREATE EXTENSION IF NOT EXISTS plpython3u;

-- 1. Определяемая пользователем скалярная функция
-- Общая продолжительность треков в плейлисте
CREATE OR REPLACE FUNCTION plpy_playlist_total_duration(playlist_id UUID)
RETURNS INT
LANGUAGE plpython3u
AS $$
    query = f"""
        SELECT duration FROM tracks t JOIN playlist_tracks pt ON t.id = pt.track_id WHERE pt.playlist_id = '{playlist_id}'
    """

    durations = plpy.execute(query)
    total_duration = sum([row['duration'] for row in durations])
    return total_duration
$$;


select plpy_playlist_total_duration('dea4927c-963d-4363-896a-ef87d669963f');


-- 2. Пользовательская агрегатная функция
-- Средняя продолжительность треков в плейлисте
CREATE OR REPLACE FUNCTION plpy_avg_duration_in_playlist(playlist_id UUID)
RETURNS float
LANGUAGE plpython3u
AS $$
    query = f"""
        SELECT duration FROM tracks t JOIN playlist_tracks pt ON t.id = pt.track_id WHERE pt.playlist_id = '{playlist_id}'
    """
    durations = plpy.execute(query)
    total_duration = sum([row['duration'] for row in durations])

    return total_duration / len(durations) if durations else 0
$$;


select plpy_avg_duration_in_playlist('dea4927c-963d-4363-896a-ef87d669963f');


-- 3. Определяемая пользователем табличная функция
-- Получение треков из альбома
CREATE OR REPLACE FUNCTION plpy_get_album_tracks(album_id UUID)
RETURNS TABLE (track_id UUID, track_name VARCHAR)
LANGUAGE plpython3u
AS $$
    results = []
    tracks = plpy.execute(f"SELECT id, name FROM tracks WHERE album_id = '{album_id}'")
    for track in tracks:
        results.append((track['id'], track['name']))
    return results
$$;


SELECT * FROM plpy_get_album_tracks('a7087655-c700-4130-a67b-1f4df2aefa28');


-- 4. Хранимая процедура CLR
-- Процедура обновляет статус премиум для всех пользователей старше 65 лет, добавляя им один год премиум-подписки
CREATE OR REPLACE PROCEDURE plpy_update_premium_for_seniors()
LANGUAGE plpython3u
AS $$
    query = """
        UPDATE users
        SET 
            premium = TRUE,
            premium_expiration = CASE
                WHEN premium = TRUE THEN premium_expiration + INTERVAL '1 year'
                ELSE CURRENT_DATE + INTERVAL '1 year'
            END
        WHERE EXTRACT(YEAR FROM AGE(birth_date)) >= 65
    """
    
    plpy.execute(query)
$$;

call plpy_update_premium_for_seniors();

select * from users where EXTRACT(YEAR FROM AGE(birth_date)) >= 65;


-- 5. Триггер
-- Триггер предотвращает удаление исполнителя, если у него есть альбомы.
CREATE OR REPLACE FUNCTION plpy_prevent_artist_deletion()
RETURNS TRIGGER
LANGUAGE plpython3u
AS $$
    query = f"""
        SELECT 1 
        FROM albums_by_artists 
        WHERE artist_id = '{TD["old"]["id"]}'
    """

    result = plpy.execute(query)
    if len(result) > 0:
        plpy.error('Cannot delete artist with albums.')
$$;

CREATE TRIGGER plpy_instead_of_delete_artist
INSTEAD OF DELETE ON artists_view
FOR EACH ROW
EXECUTE FUNCTION plpy_prevent_artist_deletion();

CREATE VIEW artists_view AS
SELECT * FROM artists;

delete from artists_view where id = '6f29c8af-9453-4ac1-b090-8f95cf687276';

delete from albums_by_artists aba where aba.artist_id = '6f29c8af-9453-4ac1-b090-8f95cf687276';

select * from artists_view av where (select count(*) from albums_by_artists aba where aba.artist_id = av.id) = 0;


-- 6. Определяемый пользователем тип данных CLR
-- тип album_info с информацией об альбоме
CREATE TYPE album_info AS (
    album_id UUID,
    title VARCHAR,
    release_date DATE
);

-- Табличная функция, использующая тип album_info
CREATE OR REPLACE FUNCTION plpy_get_albums_info(artist_id UUID)
RETURNS SETOF album_info
LANGUAGE plpython3u
AS $$
    query = f"""
        SELECT a.id AS album_id, a.title, a.release_date
        FROM albums a
        JOIN albums_by_artists aa ON a.id = aa.album_id
        WHERE aa.artist_id = '{artist_id}'
        ORDER BY a.release_date DESC
    """
    
    result = plpy.execute(query)
    return [row for row in result]
$$;

select plpy_get_albums_info('c0c4215f-8da9-4dc0-af9c-38c0c8ac69ea');
