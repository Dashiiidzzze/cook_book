-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

--Таблица категорий
CREATE TABLE dish_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Таблица рецептов
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    dish_type_id INT REFERENCES dish_types(id),
    cook_time INTERVAL,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_private BOOLEAN DEFAULT TRUE,
    holiday VARCHAR(100),
    instructions TEXT,
    photo BYTEA
);

-- Таблица ингредиентов
CREATE TABLE ingredients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Таблица связи рецептов и ингредиентов
CREATE TABLE recipe_ingredients (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id INT REFERENCES ingredients(id) ON DELETE CASCADE,
    quantity VARCHAR(100)
);

-- Заполнение таблицы dish_types
INSERT INTO dish_types (name) VALUES
('Салаты'),
('Супы'),
('Основные блюда'),
('Закуски'),
('Десерты'),
('Напитки'),
('Выпечка'),
('Соусы'),
('Завтраки'),
('Гриль'),
('Детские блюда'),
('Консервация');