--- using system catalog

COMMENT ON FUNCTION playlist_total_duration(playlist_id UUID) IS 'Общая продолжительность треков в плейлисте';

CREATE OR REPLACE FUNCTION search_functions_and_procedures(search_text TEXT)
RETURNS TABLE (object_name TEXT, object_type TEXT, description TEXT) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.proname::TEXT AS object_name,
        CASE
            WHEN p.prokind = 'f' THEN 'Function'
            WHEN p.prokind = 'p' THEN 'Procedure'
            ELSE NULL
        END::TEXT AS object_type,
        obj_description(p.oid, 'pg_proc')::TEXT AS description
    FROM 
        pg_proc AS p
    JOIN 
        pg_namespace AS n ON p.pronamespace = n.oid
    WHERE 
        (p.prosrc ILIKE '%' || search_text || '%')
        AND n.nspname NOT IN ('pg_catalog', 'information_schema') -- исключаем системные схемы
        AND p.prokind IN ('f', 'p'); -- выбираем только функции и процедуры
END;
$$ LANGUAGE plpgsql;


select object_name, object_type, description from search_functions_and_procedures('playlist');


--CREATE OR REPLACE PROCEDURE remove_duplicates(table_name TEXT)
--LANGUAGE plpgsql
--AS $$
--DECLARE
--    column_names TEXT;
--    query TEXT;
--BEGIN
--    -- Получаем список столбцов для указанной таблицы
--    SELECT string_agg(col.column_name, ', ') INTO column_names
--    FROM information_schema.columns as col
--    WHERE col.table_name = $1
--      AND col.table_schema = 'public';  -- предполагаем, что таблица находится в схеме public
--
--    -- Формируем запрос для удаления дубликатов
--    query := 'WITH unique_records AS ('
--           || 'SELECT DISTINCT ON (' || column_names || ') * FROM ' || quote_ident($1) || ' ORDER BY ' || column_names || ') '
--           || 'DELETE FROM ' || quote_ident($1) || ' '
--           || 'WHERE ctid NOT IN (SELECT ctid FROM unique_records);';
--
--    -- Выполняем запрос
--    EXECUTE query;
--END;
--$$;
--
--
--select t1.name, t1.id, t2.name, t2.id from tracks t1 join tracks t2 on t1.name = t2.name;
--
--call remove_duplicates('tracks');
--
-- SELECT string_agg('tracks', ', ')
--    FROM information_schema.columns as col
--    WHERE col.table_name = $1
--      AND col.table_schema = 'public';
      
     
     
     
     
     
     
     
     
     
     
     
     