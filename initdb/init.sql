-- Таблица пользователей
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL
);

-- Таблица рецептов
CREATE TABLE recipes (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    cook_time INTERVAL,
    is_favorite BOOLEAN DEFAULT FALSE,
    is_private BOOLEAN DEFAULT TRUE,
    instructions TEXT,
    photo BYTEA
);

-- Таблица категорий
CREATE TABLE dish_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    photo BYTEA
);

-- Таблица связи рецептов и категорий
CREATE TABLE recipe_dish_types (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    dish_types_id INT REFERENCES dish_types(id) ON DELETE CASCADE
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

-- Таблица связи рецептов и этапов приготовления                добавлена таблица
CREATE TABLE recipe_step (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    instructions TEXT,
    photo BYTEA
);

-- Таблица комментариев к рецепту                               добавлена таблица
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    recipe_id INT REFERENCES recipes(id) ON DELETE CASCADE,
    username VARCHAR(100),
    comment TEXT
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