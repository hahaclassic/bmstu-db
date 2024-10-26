-------------------------------------------------------
-- FUNCTIONS
-------------------------------------------------------

-- 1. Скалярная функция
-- Подсчет количества треков исполнителя 
create or replace function count_artist_tracks(artist_id UUID)
returns int
LANGUAGE sql
AS $$
	select count(*) from tracks_by_artists as ta where ta.artist_id = count_artist_tracks.artist_id;
$$;

select a.id, a.name, count_artist_tracks(a.id) as number_of_tracks from artists a;


-- 2. Подставляемая табличная функция
-- Получение альбомов исполнителя
create or replace function get_albums_by_artist(artist_id UUID)
RETURNS TABLE (album_id UUID, title VARCHAR, release_date DATE)
LANGUAGE sql
AS $$
    SELECT a.id, a.title, a.release_date
    FROM albums a
    JOIN albums_by_artists aa ON a.id = aa.album_id
    WHERE aa.artist_id = artist_id
    ORDER BY a.release_date DESC;
$$;

select album_id, title, release_date from get_albums_by_artist('e42cff8f-9c2f-4b2c-8308-e110c83e7c6d');


-- 3. Многооператорная табличная функция
-- Получение информации о плейлистах: id, название, количество треков, общая продолжительность
CREATE OR REPLACE FUNCTION get_playlist_summary()
RETURNS TABLE (
    playlist_id UUID,
    playlist_title VARCHAR(100),
    total_tracks INT,
    total_duration INT
) AS $$
DECLARE
    playlist_record RECORD;
BEGIN
    FOR playlist_record IN
        SELECT id, title FROM playlists
    LOOP
        SELECT 
            COUNT(pt.track_id), 
            COALESCE(SUM(t.duration), 0)
        INTO 
            total_tracks, 
            total_duration
        FROM 
            playlist_tracks pt
        LEFT JOIN 
            tracks t ON pt.track_id = t.id
        WHERE 
            pt.playlist_id = playlist_record.id;

        playlist_id := playlist_record.id;
        playlist_title := playlist_record.title;

        RETURN NEXT;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

select playlist_id, playlist_title, total_tracks, total_duration from get_playlist_summary();


-- 4. Функция с рекурсивным ОТВ
-- Функция, использующая рекурсивное ОТВ, для нахождения всех треков, созданных артистом и теми, кто с ним фитовал
CREATE FUNCTION find_tracks_by_collaborators(artist_id UUID)
RETURNS TABLE (track_id UUID, track_name VARCHAR)
LANGUAGE SQL
AS $$
WITH RECURSIVE artist_collaborators AS (
    SELECT DISTINCT ta1.artist_id AS artist_id
    FROM tracks_by_artists ta1
    WHERE ta1.artist_id = find_tracks_by_collaborators.artist_id

    UNION

    SELECT DISTINCT ta2.artist_id
    FROM tracks_by_artists ta1
    JOIN tracks_by_artists ta2 ON ta1.track_id = ta2.track_id
    WHERE ta1.artist_id = find_tracks_by_collaborators.artist_id AND ta2.artist_id <> ta1.artist_id
)
SELECT t.id, t.name
FROM artist_collaborators ac
JOIN tracks_by_artists ta ON ac.artist_id = ta.artist_id
JOIN tracks t ON ta.track_id = t.id;
$$;

INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('aed38846-d6b4-45fe-90d6-9a9b1e98907c', 'c12a849b-5a06-4df8-92f5-168d201bb8fd');

select track_id, track_name from find_tracks_by_collaborators('e42cff8f-9c2f-4b2c-8308-e110c83e7c6d') where track_id not in 
(select tba.track_id from tracks_by_artists tba where tba.artist_id = 'e42cff8f-9c2f-4b2c-8308-e110c83e7c6d');

-------------------------------------------------------
-- PROCEDURES
-------------------------------------------------------

-- 5. Хранимая процедура без параметров
-- Процедура для выдачи премиума всем пользователям от 65 лет.
CREATE OR REPLACE PROCEDURE free_premium_to_pensioners()
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE users 
    SET premium = true, 
        premium_expiration = CURRENT_DATE + INTERVAL '1 year'
    WHERE premium = false AND EXTRACT(YEAR FROM AGE(birth_date)) >= 65;
END;
$$;

select id, name, birth_date, premium, premium_expiration from users 
WHERE premium = false and EXTRACT(YEAR FROM AGE(birth_date)) >= 65;

call free_premium_to_pensioners();


-- 6. Рекурсивная хранимая процедура 
-- Удаление всех артистов, сотрудничающих с заданным артистом,
CREATE OR REPLACE PROCEDURE delete_artist_and_collaborators(artist_uuid UUID, recursion_level INT DEFAULT 1)
LANGUAGE plpgsql
AS $$
DECLARE
    collaborator UUID;
BEGIN
    IF recursion_level > 5 THEN
        RAISE NOTICE 'Max recursion level reached for artist: %', artist_uuid;
        RETURN;
    END IF;

    FOR collaborator IN
        SELECT ta2.artist_id
        FROM tracks_by_artists ta1
        JOIN tracks_by_artists ta2 ON ta1.track_id = ta2.track_id
        WHERE ta1.artist_id = artist_uuid AND ta2.artist_id <> artist_uuid
    LOOP
        call delete_artist_and_collaborators(collaborator, recursion_level + 1);
    END LOOP;

    DELETE FROM artists WHERE id = artist_uuid;
END;
$$;

INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('aed38846-d6b4-45fe-90d6-9a9b1e98907c', 'c12a849b-5a06-4df8-92f5-168d201bb8fd');

CALL delete_artist_and_collaborators('39e9ce7d-1058-41e0-a27d-386bbd84f49a', 1);


-- 7. Хранимая процедура с курсором
-- Деактивация истекших премиумов
CREATE OR REPLACE PROCEDURE deactivate_expired_premium_users()
LANGUAGE plpgsql
AS $$
DECLARE
    user_cursor CURSOR FOR 
        SELECT id, premium 
        FROM users 
        WHERE premium = TRUE AND premium_expiration < NOW()
        FOR UPDATE;
BEGIN
    FOR user_record IN user_cursor LOOP
        UPDATE users 
        SET premium = FALSE 
        WHERE CURRENT OF user_cursor; 
    END LOOP;
END;
$$;

SELECT id, premium, premium_expiration FROM users WHERE premium = TRUE AND premium_expiration < NOW();

call deactivate_expired_premium_users();


-- 8. Хранимая процедура доступа к метаданным
-- Вывод таблиц, атрибутов и их типов данных
CREATE OR REPLACE PROCEDURE get_table_metadata()
LANGUAGE plpgsql
AS $$
DECLARE
    table_record RECORD;
    column_record RECORD;
BEGIN
    FOR table_record IN 
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
    LOOP
        RAISE NOTICE 'Table: %', table_record.table_name;

        FOR column_record IN 
            SELECT column_name, data_type
            FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = table_record.table_name
        LOOP
            RAISE NOTICE '	| % | % |', column_record.column_name, column_record.data_type;
        END LOOP;
    END LOOP;
END;
$$;

call get_table_metadata();


-------------------------------------------------------
-- DML TRIGGERS
-------------------------------------------------------

-- 9. Trigger AFTER
-- Обновление даты последнего редактирования плейлиста
CREATE OR REPLACE FUNCTION update_playlist_timestamp() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE playlists
    SET last_updated = CURRENT_TIMESTAMP
    WHERE id = NEW.playlist_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql; 

-- Создание триггера для таблицы playlist_tracks
CREATE TRIGGER trigger_update_playlist_timestamp
AFTER INSERT OR UPDATE ON playlist_tracks
FOR EACH ROW
EXECUTE FUNCTION update_playlist_timestamp();


-- 10. Trigger INSTEAD OF
-- Триггер, который предотвращает удаление артистов, у которых есть альбомы
CREATE OR REPLACE FUNCTION prevent_artist_deletion()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF EXISTS (SELECT 1 FROM albums_by_artists WHERE artist_id = OLD.id) THEN
        RAISE EXCEPTION 'Cannot delete artist with albums.';
    END IF;
    
    DELETE FROM artists WHERE id = OLD.id;
    RETURN NULL; 
END;
$$;

CREATE TRIGGER instead_of_delete_artist
INSTEAD OF DELETE ON artists_view
FOR EACH ROW
EXECUTE FUNCTION prevent_artist_deletion();

CREATE VIEW artists_view AS
SELECT * FROM artists;

delete from artists_view where id = '1718e0ef-29da-400e-b242-1954539a7c0c';

delete from albums_by_artists where artist_id = '1718e0ef-29da-400e-b242-1954539a7c0c';
