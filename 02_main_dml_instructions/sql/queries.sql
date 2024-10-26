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
case EXTRACT(YEAR FROM release_date)
    WHEN EXTRACT(YEAR FROM CURRENT_DATE) THEN 'This Year'
    WHEN EXTRACT(YEAR FROM CURRENT_DATE) - 1 THEN 'Last Year'
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
-- Самый прослушиваемый трек среди всех в жанре "Pop" и "Rock"
SELECT genre_table.genre, track_table.name, track_table.stream_count
FROM (
    SELECT genre, MAX(stream_count) AS max_streams
    FROM tracks
    WHERE genre IN ('Pop', 'Rock') 
    GROUP BY genre
) AS genre_table
JOIN (
    SELECT t.genre, t.name, t.stream_count
    FROM tracks t
    WHERE t.stream_count = (
        SELECT MAX(t2.stream_count)
        FROM tracks t2
        WHERE t2.genre = t.genre
    )
) AS track_table
ON genre_table.genre = track_table.genre AND genre_table.max_streams = track_table.stream_count;

-- 13. Инструкция SELECT, использующая вложенные подзапросы с уровнем вложенности 3
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
-- Добавление в плейлист "Best 100 rock tracks" 100 самых прослушиваемых треков в жанре "Рок"

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

delete from playlist_tracks pt where pt.playlist_id = 'dbd88664-8fa1-45d7-89b1-363da25054d4'; -- id of 'Best 100 rock tracks'

select t.id, t.name, t.stream_count, pt.track_order from tracks t join playlist_tracks pt on t.id = pt.track_id
join playlists p on p.id = pt.playlist_id where p.title = 'Best 100 rock tracks' order by pt.track_order;

-- 18. Простая инструкция UPDATE
-- Обновление никнейма пользователя с id=2ee6dbe9-6777-4443-9789-016c36cc41cd
UPDATE users SET name = 'Tompson777' WHERE id = '2ee6dbe9-6777-4443-9789-016c36cc41cd';

-- 19. Инструкция UPDATE со скалярным подзапросом в предложении SET.
-- Установка рейтинга плейлиста в зависимости от количества его слушателей (сколько людей добавили его себе в коллекцию)

-- все пользователи, родившиеся с '1960-01-01' по '1970-01-01', добавляют к себе плейлист "Best 100 rock tracks";
INSERT INTO user_playlists (playlist_id, user_id, is_favorite, access_level)
SELECT 'dbd88664-8fa1-45d7-89b1-363da25054d4', u.id, false, 2 
FROM users u 
WHERE u.birth_date BETWEEN '1960-01-01' AND '1970-01-01';

select * from user_playlists up where playlist_id = 'dbd88664-8fa1-45d7-89b1-363da25054d4';

UPDATE playlists p 
SET rating = (
    SELECT count(*) 
    FROM user_playlists up 
    WHERE p.id = up.playlist_id
)
WHERE p.rating = 0;

UPDATE playlists p set rating = 0 where p.id = 'dbd88664-8fa1-45d7-89b1-363da25054d4';

select * from playlists p where rating = 0;

-- 20. Простая инструкция DELETE
-- Удаление треков с нулевым количеством прослушиваний
DELETE FROM tracks WHERE stream_count = 0;

-- 21. Инструкция DELETE с вложенным коррелированным подзапросом в предложении WHERE
-- Удаление плейлистов, у которых нет треков
select p.id, p.title from playlists p left join playlist_tracks pt ON p.id = pt.playlist_id where pt.playlist_id is null;

select pt.playlist_id, count(*) from playlist_tracks pt group by pt.playlist_id;

insert into playlists (id, title, description, private, last_updated, rating)
VALUES (uuid_generate_v4(), '!!!! My Empty !!!!!', 'Empty playlist', FALSE, CURRENT_TIMESTAMP, 0);

DELETE FROM playlists 
WHERE id not IN (
	select pt.playlist_id from playlist_tracks pt
);

-- 22. Инструкция SELECT, использующая простое обобщенное табличное выражение (CTE)
-- Подсчет среднего количества треков в альбомах
WITH album_track_count (album_id, track_count) AS (
    SELECT album_id, COUNT(*) AS track_count
    FROM tracks
    GROUP BY album_id
)
SELECT AVG(track_count) AS average_tracks_per_album
FROM album_track_count;

-- 23. Инструкция SELECT, использующая рекурсивное обобщенное табличное выражение.
-- Добавим пару фитов
INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('3531f9f1-8901-4e6f-bf52-e69bbb07a9f4', 'e5df1125-a73a-4eae-b4bf-37b15db47cee');

INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('7ca63761-d6bf-436c-8fef-409218f664f8', '2346bafb-7e0f-45af-8005-0f2fbdf6a545');

select id, name from artists a where id in ('e5df1125-a73a-4eae-b4bf-37b15db47cee', 
	'2346bafb-7e0f-45af-8005-0f2fbdf6a545', 'a2dfddfb-5290-4847-b5c5-553beb767c4e');
	

INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('e7295720-bc46-4633-8b91-48362d0730e7', '35a3ebf8-1d61-4d26-b8c8-1281e917299f');

INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('efa63961-78cd-4d21-80c5-70c2dcfcf916', '6b2a1220-f717-4f52-92d8-6087db0a0a3f');

INSERT INTO tracks_by_artists (track_id, artist_id)
VALUES ('1b27012b-6848-48f5-932f-052084c92baa', '6b2a1220-f717-4f52-92d8-6087db0a0a3f');

-- Поиск кратчайшего "feat-path" - количества фитов между двумя исполнителями

-- Version 1
--WITH RECURSIVE artist_connections (artist1, artist2, num_feats) AS (
--    -- Базовый случай: прямой фит между артистами
--    SELECT 
--	    ta1.artist_id artist1, 
--	    ta2.artist_id artist2, 
--	    1 num_feats
--	FROM tracks_by_artists ta1 
--	JOIN tracks_by_artists ta2 
--	    ON ta1.track_id = ta2.track_id 
--	    AND ta1.artist_id <> ta2.artist_id
--
--    UNION ALL
--
--    -- Рекурсивная часть: поиск фитов через других артистов
--    select artist1, artist2, num_feats from (
--    	WITH artist_connections2 AS (
--    		SELECT * FROM artist_connections
--    	)
--    	select * from (
--    		select * from artist_connections2
--	    	union
--	    	select ac1.artist2 as artist1, ac2.artist2 as artist2, ac1.num_feats + ac2.num_feats as num_feats 
--	    		from artist_connections2 ac1 join artist_connections2 ac2 
--	    		on ac1.artist1 = ac2.artist1 and ac1.artist2 <> ac2.artist2 
--	    		
--	    	where (ac1.num_feats + ac2.num_feats < (select ac3.num_feats from artist_connections2 ac3
--				where ac3.artist1 = ac1.artist2 and ac3.artist2 = ac2.artist2) 
--	    		or (select count(*) from artist_connections2 ac4 
--	    			where ac4.artist1 = ac1.artist2 and ac4.artist2 = ac2.artist2) = 0)
--    	) WHERE EXISTS (
--		    SELECT 1 
--		    FROM artist_connections2 ac1 
--		    JOIN artist_connections2 ac2 
--		        ON ac1.artist1 = ac2.artist1 AND ac1.artist2 <> ac2.artist2 
--		    WHERE (ac1.num_feats + ac2.num_feats < 
--		           (SELECT ac3.num_feats 
--		            FROM artist_connections2 ac3 
--		            WHERE ac3.artist1 = ac1.artist2 AND ac3.artist2 = ac2.artist2) 
--		           OR (SELECT COUNT(*) 
--		               FROM artist_connections2 ac4 
--		               WHERE ac4.artist1 = ac1.artist2 AND ac4.artist2 = ac2.artist2) = 0)
--		)
--    )
--)
--SELECT artist1, artist2, MIN(num_feats) AS shortest_feat_path FROM artist_connections 
--GROUP BY artist1, artist2 
--ORDER BY shortest_feat_path;

-- Version 2
WITH RECURSIVE artist_connections (artist1, artist2, num_feats) AS (
    -- Базовый случай: прямой фит между артистами
    SELECT 
	    ta1.artist_id artist1, 
	    ta2.artist_id artist2, 
	    1 num_feats
	FROM tracks_by_artists ta1 
	JOIN tracks_by_artists ta2 
	    ON ta1.track_id = ta2.track_id 
	    AND ta1.artist_id <> ta2.artist_id

    UNION ALL

    -- Рекурсивная часть: поиск фитов через других артистов
    select artist1, artist2, num_feats from (
    	WITH artist_connections2 AS (
		    SELECT * FROM artist_connections
		),
		new_connections AS (
		    SELECT ac1.artist2 AS artist1, 
		           ac2.artist2 AS artist2, 
		           ac1.num_feats + ac2.num_feats AS num_feats 
		    FROM artist_connections2 ac1 
		    JOIN artist_connections2 ac2 
		        ON ac1.artist1 = ac2.artist1 AND ac1.artist2 <> ac2.artist2 
		    WHERE (ac1.num_feats + ac2.num_feats < 
		           (SELECT MIN(ac3.num_feats) 
		            FROM artist_connections2 ac3 
		            WHERE ac3.artist1 = ac1.artist2 AND ac3.artist2 = ac2.artist2) 
		           OR (SELECT COUNT(*) 
		               FROM artist_connections2 ac4 
		               WHERE ac4.artist1 = ac1.artist2 AND ac4.artist2 = ac2.artist2) = 0)
		)
		SELECT artist1, artist2, num_feats 
		FROM (
		    SELECT * FROM artist_connections2
		    UNION ALL
		    SELECT artist1, artist2, num_feats FROM new_connections
		) t
		WHERE EXISTS (SELECT 1 FROM new_connections)
    )
)
SELECT artist1, artist2, MIN(num_feats) AS shortest_feat_path FROM artist_connections 
GROUP BY artist1, artist2 
ORDER BY shortest_feat_path;


-- 24. Оконные функции. Использование конструкци MIN/MAX/AVG OVER()
-- Получение минимального, максимального и среднего числа прослушиваний для каждого жанра (для сравнения)
SELECT genre, name, stream_count,
    AVG(stream_count) OVER(PARTITION BY genre) AS avg_stream_count,
    MIN(stream_count) OVER(PARTITION BY genre) AS min_stream_count,
    MAX(stream_count) OVER(PARTITION BY genre) AS max_stream_count
FROM tracks;

-- 25. Оконные функции для устранения дублей
CREATE TEMPORARY TABLE temp_tracks (
    id UUID,
    name VARCHAR(100),
    order_in_album INT,
    album_id UUID,
    explicit BOOLEAN,
    duration INT,
    genre VARCHAR(50),
    stream_count BIGINT
);

-- Создание дублей
INSERT INTO temp_tracks (id, name, order_in_album, album_id, explicit, duration, genre, stream_count)
SELECT 
    id, 
    name, 
    order_in_album, 
    album_id, 
    explicit, 
    duration, 
    genre, 
    stream_count
FROM tracks;

-- Проверка, что дубли действительно есть
SELECT id, count(*) FROM temp_tracks group by id;

SELECT id, ROW_NUMBER() OVER (
            PARTITION by id, name, album_id, order_in_album, explicit, duration, genre, stream_count 
            ORDER BY id
        ) AS row_num
    FROM temp_tracks;

-- Удаление дублей
WITH enumerated_tracks AS (
    SELECT ctid, ROW_NUMBER() OVER (
               PARTITION by id, name, album_id, order_in_album, explicit, duration, genre, stream_count 
           ) AS row_num
    FROM temp_tracks
)
DELETE FROM temp_tracks tt
WHERE ctid IN (
    SELECT ctid FROM enumerated_tracks et WHERE et.row_num > 1 
);


-- DEFENCE 
-- Для каждого несовершеннолетноего пользователя, зарегистрированного за последний год, вывести
-- 1. кол-во плейлистов
-- 2. самый часто встречаемый жанр в плейлисте
-- 3. страна самого часто встречаемого артиста в данном жанре
   
WITH user_genre AS (
    SELECT 
        u.id AS user_id,
        u.name AS username,
        AGE(u.birth_date) AS age,
        COUNT(DISTINCT up.playlist_id) AS num_of_playlists,
        mode() WITHIN GROUP (ORDER BY t.genre) AS most_frequent_genre
    FROM 
        users u
    JOIN 
        user_playlists up ON u.id = up.user_id
    JOIN 
        playlist_tracks pt ON up.playlist_id = pt.playlist_id
    JOIN 
        tracks t ON pt.track_id = t.id
    WHERE 
        u.birth_date > CURRENT_DATE - INTERVAL '18 years' 
        AND u.registration_date > CURRENT_DATE - INTERVAL '1 year'
    GROUP BY 
        u.id, u.name
)
SELECT 
    ug.user_id,
    ug.username,
    ug.age,
    ug.num_of_playlists,
    ug.most_frequent_genre,
    mode() WITHIN GROUP (ORDER BY a.name) AS most_frequent_artist,
    mode() WITHIN GROUP (ORDER BY a.country) AS top_artist_country
FROM 
    user_genre ug
JOIN 
    user_playlists up ON ug.user_id = up.user_id
JOIN 
    playlist_tracks pt ON up.playlist_id = pt.playlist_id
JOIN 
    tracks t ON pt.track_id = t.id AND t.genre = ug.most_frequent_genre
JOIN 
    tracks_by_artists ta ON t.id = ta.track_id
JOIN 
    artists a ON ta.artist_id = a.id
GROUP BY 
    ug.user_id, ug.username, ug.age, ug.num_of_playlists, ug.most_frequent_genre;

   





