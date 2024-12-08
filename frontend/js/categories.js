// Показываем/скрываем фильтры
document.getElementById('filterButton').addEventListener('click', () => {
    const filterPopup = document.getElementById('filters');
    filterPopup.style.display = filterPopup.style.display === 'block' ? 'none' : 'block';
});

// Загружаем категории при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchCategories);

// Функция для получения последних 10 рецептов
async function fetchCategories() {
    try {
        const response = await fetch('/categories/all'); // Запрос к API
        if (!response.ok) throw new Error('Ошибка при загрузке рецептов');

        const categories = await response.json();
        displayCategories(categories); // Отображаем рецепты
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить рецепты.');
    }
}

// Отображение рецептов
function displayCategories(categories) {
    const categoriesSection = document.getElementById('categories-section');

    // Убедимся, что секция существует
    if (!categoriesSection) {
        console.error('Секция #categories-section не найдена');
        return;
    }

    // Создаём заголовок
    const title = document.createElement('h2');
    title.textContent = 'Категории рецептов';
    categoriesSection.appendChild(title);

    // Создаём контейнер для категорий
    const categoriesContainer = document.createElement('div');
    categoriesContainer.id = 'categories';
    categoriesContainer.className = 'categories-container';

    // Создаём карточки для каждой категории
    categories.forEach(category => {
        const categoryCard = document.createElement('div');
        categoryCard.className = 'categories-card';
        categoryCard.setAttribute('data-id', category.id);  // Добавляем ID категории для использования при запросе

        const categoryName = document.createElement('h3');
        categoryName.textContent = category.name;

        categoryCard.appendChild(categoryName);

        // Добавляем обработчик события для клика по категории
        categoryCard.addEventListener('click', () => {
            fetchRecipesByCategory(category.id, category.name);  // Отправляем запрос с ID категории
        });

        categoriesContainer.appendChild(categoryCard);
    });

    categoriesSection.appendChild(categoriesContainer);
}

// Функция для получения рецептов по категории
async function fetchRecipesByCategory(categoryId, categoryName) {
    try {
        const response = await fetch(`/categories/recipes?category_id=${categoryId}`);
        if (!response.ok) throw new Error('Ошибка загрузки рецептов');
        const recipes = await response.json();

        // Отображаем рецепты
        displayRecipes(recipes, categoryName);
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

// Функция для отображения рецептов
function displayRecipes(recipes, categoryName) {
    // Скрываем секцию категорий и показываем секцию рецептов
    const categoriesSection = document.getElementById('categories-section');
    if (categoriesSection) categoriesSection.style.display = 'none';

    const recipesSection = document.getElementById('recipes-section');
    if (!recipesSection) {
        console.error('Секция #recipes-section не найдена.');
        return;
    }
    //recipesSection.style.display = 'block';

    // Создаём заголовок
    const title = document.createElement('h2');
    title.textContent = 'Рецепты в категории ' + categoryName;
    recipesSection.appendChild(title);

    // Создаём контейнер для категорий
    const recipesContainer = document.createElement('div');
    recipesContainer.id = 'recipe';
    recipesContainer.className = 'recipe-container';

    // Создаём карточки для каждой категории
    recipes.forEach(recipe => {
        const recipeCard = document.createElement('div');
        recipeCard.className = 'recipe-card';

        recipeCard.innerHTML = `
            <img src="${recipe.photo}" alt="Фото рецепта" class="recipe-photo">
            <h3>${recipe.name}</h3>
            <p>Время приготовления: ${recipe.cook_time}</p>
            <p>Ингредиенты: ${recipe.ingredients}</p>
        `;
        // Добавляем обработчик клика на всю карточку
        recipeCard.addEventListener('click', () => {
            window.location.href = `/recipe?recipe_id=${recipe.id}`; // Переход на страницу рецепта
        });

        // // Добавляем обработчик события для клика по рецепту
        // categoryCard.addEventListener('click', () => {
        //     fetchRecipesByCategory(category.id, category.name);  // Отправляем запрос с ID
        // });

        recipesContainer.appendChild(recipeCard);
    });

    recipesSection.appendChild(recipesContainer);
}