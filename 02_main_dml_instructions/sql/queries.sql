CREATE extension IF NOT EXISTS "uuid-ossp";

-- 1. Инструкция SELECT, использующая предикат сравнения. 
-- Получение альбомов исполнителей из Франции.
SELECT DISTINCT 
	albums.id AS album_id, 
    albums.title AS album_title, 
    artists.name AS artist_name
FROM albums
INNER JOIN albums_by_artists ON albums.id = albums_by_artists.album_id
INNER JOIN artists ON albums_by_artists.artist_id = artists.id
WHERE artists.country = 'France';


-- 2. Инструкция SELECT, использующая предикат BETWEEN. 
-- Получение треков, у которых от 100к до 200к прослушиваний и они explicit.
select distinct t.name, t.genre, t.stream_count from tracks as t 
where t.stream_count between 100000 and 200000 and t.explicit;


-- 3. Инструкция SELECT, использующая предикат LIKE. 
-- Получение всех треков, у которых в названии есть слово 'book'
select distinct t.id as track_id, t.name as track_name from tracks as t where t.name like '%book%';

-- 4. Инструкция SELECT, использующая предикат IN с вложенным подзапросом
-- Получение альбомов исполнителей из Франции
SELECT id, title
FROM albums
WHERE id IN (
    SELECT album_id
    FROM albums_by_artists
    WHERE artist_id IN (
        SELECT id
        FROM artists
        WHERE country = 'France'
    )
);

-- 5. Инструкция SELECT, использующая предикат EXISTS с вложенным подзапросом.
-- Получение всех прейлистов, в которых есть хотя бы 1 трек, добавленный позже 10 октября
SELECT playlists.id, playlists.title
FROM playlists
WHERE EXISTS (
    SELECT 1
    FROM playlist_tracks
    WHERE playlist_tracks.playlist_id = playlists.id
    AND playlist_tracks.date_added > '2024-10-10'
);


-- 6. Инструкция SELECT, использующая предикат сравнения с квантором
-- Получение треков, у которых количество прослушиваний больше, чем у всех треков жанра "Pop"
SELECT id, name, genre
FROM tracks
WHERE stream_count > ALL (
    SELECT stream_count
    FROM tracks
    WHERE genre = 'Pop'
);

-- 7. Инструкция SELECT, использующая агрегатные функции в выражениях столбцов
-- Запрос для подсчета среднего, макс. и мин. количества прослушиваний для треков жанра "Pop"
SELECT AVG(stream_count) AS "Actual AVG", SUM(stream_count) / COUNT(id) AS "Calc AVG", 
	MIN(stream_count) as "min", MAX(stream_count) as "max" FROM tracks where genre = 'Pop';

-- 8. нструкция SELECT, использующая скалярные подзапросы в выражениях столбцов
-- Получение среднего, максимального и минимального количества прослушиваний на треке для каждого исполнителя
SELECT 
    artists.name AS artist_name,
    (SELECT AVG(stream_count)
     FROM tracks t2
     JOIN tracks_by_artists tba ON t2.id = tba.track_id
     WHERE tba.artist_id = artists.id) AS avg_stream_count,
    (SELECT MIN(stream_count)
     FROM tracks t2
     JOIN tracks_by_artists tba ON t2.id = tba.track_id
     WHERE tba.artist_id = artists.id) AS min_stream_count,
    (SELECT MAX(stream_count)
     FROM tracks t2
     JOIN tracks_by_artists tba ON t2.id = tba.track_id
     WHERE tba.artist_id = artists.id) AS max_stream_count
FROM artists;


-- 9. Инструкция SELECT, использующая простое выражение CASE
-- Получение альбомов с делением по дате выпуска на "этот год", "прошлый", "за последнее десятилетие" и "ранее"
SELECT title, release_date,
CASE
    WHEN EXTRACT(YEAR FROM release_date) = EXTRACT(YEAR FROM CURRENT_DATE) THEN 'This Year'
    WHEN EXTRACT(YEAR FROM release_date) = EXTRACT(YEAR FROM CURRENT_DATE) - 1 THEN 'Last Year'
    when EXTRACT(year from release_date) >= EXTRACT(year from CURRENT_DATE) - 10 then 'Last Decade'
    ELSE 'Earlier'
END AS release_status
FROM albums;

-- 10. Инструкция SELECT, использующая поисковое выражение CASE
-- Запрос для классификации треков по продолжительности
SELECT name,
CASE
    WHEN duration < 200 THEN 'Short'
    WHEN duration < 250 THEN 'Medium'
    ELSE 'Long'
END AS duration_category
FROM tracks;

-- 11. Создание временной таблицы из результирующего набора данных инструкции SELECT
-- Создание временной таблицы с треками с наибольшим количеством прослушиваний
SELECT id, name, stream_count
INTO TEMPORARY TABLE top_streamed_tracks
FROM tracks
ORDER BY stream_count DESC
LIMIT 10;

select * from top_streamed_tracks;

drop table if exists top_streamed_tracks;

-- 12. Инструкция SELECT, использующая вложенные коррелированные подзапросы в качестве производных таблиц в предложении FROM
-- хуета какая-то

-- 13. Инструкция SELECT, использующая вложенные подзапросы с уровнем вложенности 3
select p.id, p.title from playlists p where p.id (select )

-- 14. Инструкция SELECT, консолидирующая данные с помощью GROUP BY, но без HAVING
-- Запрос для получения средней продолжительности треков в альбоме
SELECT a.id, a.title, AVG(t.duration) AS avg_duration FROM tracks as t inner join albums a on t.album_id = a.id GROUP BY a.id, a.title;


-- 15. Инструкция SELECT, консолидирующая данные с помощью GROUP BY и HAVING
-- Запрос для получения альбомов, у которых средняя продолжительность треков выше общей средней
SELECT a.id, a.title, AVG(t.duration) AS avg_duration FROM tracks as t inner join albums a on t.album_id = a.id GROUP BY a.id, a.title
HAVING AVG(duration) > (SELECT AVG(duration) FROM tracks);

-- 16. Однострочная инструкция INSERT, выполняющая вставку в таблицу одной строки значений.
-- Запрос для добавления нового пользователя
INSERT INTO users (id, name, registration_date, birth_date, premium)
VALUES (uuid_generate_v4(), 'John Wick', CURRENT_TIMESTAMP, '1970-01-01', FALSE);


-- 17. Многострочная инструкция INSERT, выполняющая вставку в таблицу
-- результирующего набора данных вложенного подзапроса.
-- Добавление в плейлист "Best 100 rock tracks" 100 самых прослушиваех треков в жанре "Рок"

CREATE OR REPLACE FUNCTION update_last_updated() 
RETURNS TRIGGER AS $$
BEGIN
    UPDATE playlists
    SET last_updated = CURRENT_TIMESTAMP
    WHERE id = NEW.playlist_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql; 

CREATE TRIGGER update_playlist_timestamp
AFTER INSERT OR UPDATE ON playlist_tracks
FOR EACH ROW
EXECUTE FUNCTION update_last_updated();


INSERT INTO playlists (id, title, description, private, last_updated, rating)
VALUES (uuid_generate_v4(), 'Best 100 rock tracks', 'Top 100 rock tracks based on stream count', FALSE, CURRENT_TIMESTAMP, 0);

select * from playlists p where title = 'Best 100 rock tracks';

WITH playlist_info AS (
    SELECT id FROM playlists WHERE title = 'Best 100 rock tracks'
),
top_rock_tracks AS (
    SELECT id
    FROM tracks
    WHERE genre = 'Rock'
    ORDER BY stream_count DESC
    LIMIT 100
)
INSERT INTO playlist_tracks (track_id, playlist_id, date_added, track_order)
SELECT t.id, pi.id, CURRENT_TIMESTAMP, 
       COALESCE(
           (SELECT MAX(track_order) FROM playlist_tracks WHERE playlist_id = pi.id), 0
       ) + ROW_NUMBER() OVER ()
FROM top_rock_tracks t, playlist_info pi;

delete from playlist_tracks pt where pt.playDELETE Orders
WHERE CustomerID IS NULLlist_id = '1308b55c-a84a-4d7d-ba6a-8d69be6cd25e';

select t.id, t.name, t.stream_count, pt.track_order from tracks t join playlist_tracks pt on t.id = pt.track_id
join playlists p on p.id = pt.playlist_id where p.title = 'Best 100 rock tracks';


-- 18. Простая инструкция UPDATE
-- Обновление никнейма пользователя с id=2ee6dbe9-6777-4443-9789-016c36cc41cd
UPDATE users SET name = 'Tompson777' WHERE id = '2ee6dbe9-6777-4443-9789-016c36cc41cd';

-- 19. Инструкция UPDATE со скалярным подзапросом в предложении SET.
-- Продление премиум-подписки на 1 год
UPDATE users
SET premium_expiration = (SELECT CURRENT_TIMESTAMP + INTERVAL '1 year')
WHERE id = 'user-uuid';


-- 20. Простая инструкция DELETE.
-- 
delete 



