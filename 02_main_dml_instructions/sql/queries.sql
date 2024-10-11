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
-- 


-- 6. Инструкция SELECT, использующая предикат сравнения с квантором
-- Найти все треки, у количество прослушиваний больше, чем у всех треков жанра "Pop"
SELECT id, name, stream_count 
FROM tracks
WHERE stream_count > ALL (
    SELECT stream_count 
    FROM tracks
    WHERE genre = 'Pop'
);
