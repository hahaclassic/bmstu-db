CREATE EXTENSION IF NOT EXISTS plpython3u;

-- 1. Определяемая пользователем скалярная функция
-- Общая продолжительность треков в плейлисте
CREATE OR REPLACE FUNCTION playlist_total_duration(playlist_id UUID)
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

select playlist_total_duration('dea4927c-963d-4363-896a-ef87d669963f');


-- 2. Пользовательская агрегатная функция
-- Средняя продолжительность треков в плейлисте
CREATE OR REPLACE FUNCTION avg_duration_in_playlist(playlist_id UUID)
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

select avg_duration_in_playlist('dea4927c-963d-4363-896a-ef87d669963f');


-- 3. Определяемая пользователем табличная функция
CREATE OR REPLACE FUNCTION get_album_tracks(album_id UUID)
RETURNS TABLE (track_id UUID, track_name VARCHAR)
LANGUAGE plpython3u
AS $$
    results = []
    tracks = plpy.execute(f"SELECT id, name FROM tracks WHERE album_id = '{album_id}'")
    for track in tracks:
        results.append((track['id'], track['name']))
    return results
$$;

SELECT * FROM get_album_tracks('a7087655-c700-4130-a67b-1f4df2aefa28');
