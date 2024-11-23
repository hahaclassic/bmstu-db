--- 1. Выгрузить данные в json 
SELECT json_agg(row_to_json(a)) AS artists_json
FROM artists a;


COPY (
	SELECT json_agg(row_to_json(a)) AS artists_json
	FROM artists a) 
TO '/data/artists_json.json';


--- 2. Загрузить данные из json
CREATE temp TABLE temp_artists (
    id UUID PRIMARY KEY,
    name VARCHAR(100),
    genre VARCHAR(50),
    country VARCHAR(100),
    debut_year INT
);

select * from temp_artists;
select * from artists a;

drop table temp_artists;

CREATE TEMP TABLE temp_json_data (
    data JSON
);

COPY temp_json_data(data)
FROM '/data/artists_json.json';

select * from temp_json_data;

INSERT INTO temp_artists(id, name, genre, country, debut_year)
SELECT 
    (artist->>'id')::UUID,
    artist->>'name',
    artist->>'genre',
    artist->>'country',
    (artist->>'debut_year')::INT
FROM (
    SELECT json_array_elements(data) AS artist
	FROM temp_json_data
) subquery;


--- 3. Добавление атрибута с типом JSON
ALTER TABLE tracks ADD COLUMN metadata JSONB;

UPDATE tracks SET metadata = '{"producer": "John Wick", "studio": "Big Sound Studio"}';


--- 4.1. Извеление JSON фрагмента
SELECT id, name, metadata->'producer' AS producer FROM tracks;


--- 4.2 Извлечение значений узлов или атрибутов
SELECT id, name, metadata->>'studio' AS studio FROM tracks;


--- 4.3 Проверка существования узла или атрибута
SELECT id, name, metadata->'producer' AS producer FROM tracks where metadata ? 'producer';


--- 4.4 Добавление нового ключа json
UPDATE tracks
SET metadata = jsonb_set(metadata, '{location}', '"New York"');


--- 4.5 Разделить JSON документ на несколько строк по узлам
SELECT metadata->>'producer' as producer, 
		metadata->>'studio' as studio, 
		metadata->>'location' as location
FROM tracks;

WITH data AS (
    SELECT '{"phones": ["123-456-7890", "987-654-3210", "555-666-7777"]}'::jsonb AS json_data
)
SELECT phone
FROM data,
jsonb_array_elements(json_data->'phones') AS phone;



