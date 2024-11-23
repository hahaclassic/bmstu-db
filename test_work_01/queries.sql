-- 1. Таблица "Виды блюд"
create table dish_types (
    id SERIAL PRIMARY KEY,
    dish_name VARCHAR(100) not null,
    description TEXT,
    rating DECIMAL(2, 1) check (rating >= 0 and rating <= 5)
);

drop table dish_types;

-- 2. Таблица "Продукты"
create table products (
    id SERIAL PRIMARY KEY,
    product_name VARCHAR(100) not null,
    production_date DATE,
    shelf_life INTERVAL,
    supplier_name VARCHAR(100)
);

-- 3. Таблица "Меню"
create table menu (
    id SERIAL PRIMARY KEY,
    menu_name VARCHAR(100) not null,
    meal_type VARCHAR(20) check (meal_type in ('завтрак', 'обед', 'ужин')),
    description TEXT
);

-- 4. Таблица для связи блюд и продуктов
create table dish_products (
    dish_id INT references dish_types(id) on delete cascade,
    product_id INT references products(id) on delete cascade,
    PRIMARY KEY (dish_id, product_id)
);

drop table dish_products;

-- 5. Таблица для связи меню и блюд
create table menu_dishes (
    menu_id INT references menu(id) on delete cascade,
    dish_id INT references dish_types(id) on delete cascade,
    PRIMARY KEY (menu_id, dish_id)
);

----------------------
-- Тестовые данные 
----------------------
insert into dish_types (dish_name, description, rating)
VALUES 
    ('Пицца', 'Традиционная итальянская пицца', 4.5),
    ('Пицца с ананасами', 'Нетрадиционная итальная пицца', 2),
    ('Салат Цезарь', 'Классический салат с курицей', 5),
    ('Борщ', 'Вкусный', 4.8),
    ('Суши', 'С рыбой и рисом', 4.0),
    ('Паста Карбонара', 'Итальянская паста с соусом карбонара', 3.9),
    ('Стейк', 'говяжий medium rare so tasty so delishios', 4.7),
    ('Омлет', 'Завтрак из яиц', 3.5),
    ('Блины', 'Лучше, чем в теремке', 4.3),
    ('Чизбургер', 'Бургер с сыром', 2.7),
    ('Суп Том Ям', 'Острый тайский суп с креветками', 4.6);

   
insert into products (product_name, production_date, shelf_life, supplier_name)
VALUES 
    ('Мука', '2024-11-01', '90 days', 'Поставщик А'),
    ('Курица', '2024-11-10', '10 days', 'Поставщик Б'),
    ('Томаты', '2024-11-09', '5 days', 'Поставщик В'),
    ('Сыр', '2024-10-20', '30 days', 'Поставщик Г'),
    ('Рис', '2024-10-25', '120 days', 'Поставщик Д'),
    ('Картофель', '2024-11-08', '20 days', 'Поставщик Е'),
    ('Лосось', '2024-11-11', '7 days', 'Поставщик Ж'),
    ('Свекла', '2024-11-05', '14 days', 'Поставщик З'),
    ('Хлеб', '2024-11-13', '3 days', 'Поставщик И'),
    ('Бекон', '2024-11-10', '15 days', 'Поставщик К');

insert into menu (menu_name, meal_type, description)
VALUES 
    ('Меню Завтрак', 'завтрак', 'Меню для завтрака'),
    ('Меню Обед', 'обед', 'Меню для обеда'),
    ('Меню Ужин', 'ужин', 'Меню для ужина'),
    ('Меню Бизнес-ланч', 'обед', 'Меню для бизнес-ланча'),
    ('Меню Праздничное', 'ужин', 'Праздничное меню'),
    ('Меню Вегетарианское', 'обед', 'Вегетарианское меню'),
    ('Меню Детское', 'завтрак', 'Детское меню на утро'),
    ('Меню Фитнес', 'обед', 'Меню для спортиков'),
    ('Меню Экзотическое', 'ужин', 'Меню с экзотическими блюдами'),
    ('Меню Летнее', 'обед', 'Летнее меню с кока-колой'),
    ('Меню Зимнее', 'ужин', 'Зимнее меню с чаем теплым');


select * from products;
select * from dish_types;

INSERT INTO dish_products (dish_id, product_id) VALUES
    (1, 1),
    (1, 2),
    (2, 3),
    (2, 4),
    (3, 5),
    (3, 6),
    (4, 7),
    (5, 8),
    (6, 9),
    (7, 10);

INSERT INTO menu_dishes (menu_id, dish_id) VALUES
    (1, 1),
    (1, 2),
    (2, 3),
    (2, 4),
    (3, 5),
    (3, 6),
    (4, 7),
    (4, 1),
    (5, 2),
    (5, 3);

   
-----------------------
-- ЗАДАНИЕ №2. ЗАПРОСЫ
-----------------------

-- 1. Инструкция select, использующая поисковое выражение case.
-- Выставление категории качества в зависимости от рейтинга блюда.
SELECT 
    dish_name,
    rating,
    CASE 
	    when rating > 4.5 then 'Отлично'
        WHEN rating > 3.5 THEN 'Хорошо'
        WHEN rating > 3 THEN 'Приемлемо'
        ELSE 'Плохо'
    END AS rating_quality
FROM 
    dish_types;
   
 
-- 2. Инструкция select, использующая скалярные подзапросы в выражениях столбцов.
-- Посчет количества имеющихся продуктов, необходимых для каждого блюда.
SELECT 
    dt.dish_name,
    dt.rating,
    (SELECT COUNT(*) FROM dish_products as dp WHERE dt.id = dp.dish_id) AS number_of_products
FROM 
    dish_types as dt;

   
-- 3. Инструция с накопительный оконной функцией.
-- 
SELECT 
    dish_name,
    rating,
    AVG(rating) OVER (ORDER BY id) AS cumulative_avg_rating
FROM 
    dish_types
ORDER BY 
    id;
   
----------------------------------
-- ЗАДАНИЕ №3. ХРАНИМАЯ ФУНКЦИЯ
----------------------------------
SELECT nspname || '.' || relname AS "table_name",
	    pg_total_relation_size(C.oid) / 1024 AS "total_size in kb"
	FROM pg_class C
	LEFT JOIN pg_namespace N ON (N.oid = C.relnamespace)
	WHERE nspname NOT IN ('pg_catalog', 'information_schema')
	    AND C.relkind <> 'i'
	    AND nspname !~ '^pg_toast'
	ORDER BY pg_total_relation_size(C.oid) ASC
	LIMIT 5;


CREATE OR REPLACE FUNCTION get_smallest_tables_in_kb()
RETURNS TABLE(table_name TEXT, total_size_in_kb BIGINT) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        nspname || '.' || relname AS "table_name",
        pg_total_relation_size(C.oid) / 1024 AS "total_size_in_kb"
    FROM 
        pg_class C
    LEFT JOIN 
        pg_namespace N ON (N.oid = C.relnamespace)
    WHERE 
        nspname NOT IN ('pg_catalog', 'information_schema')
        AND C.relkind <> 'i'
        AND nspname !~ '^pg_toast'
    ORDER BY 
        pg_total_relation_size(C.oid) ASC
    LIMIT 5;
END;
$$ LANGUAGE plpgsql;

select get_smallest_tables_in_kb();
