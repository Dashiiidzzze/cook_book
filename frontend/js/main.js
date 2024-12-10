// Показываем/скрываем фильтры
document.getElementById('filterButton').addEventListener('click', () => {
    const filterPopup = document.getElementById('filters');
    filterPopup.style.display = filterPopup.style.display === 'block' ? 'none' : 'block';
});

// Функция для получения последних 10 рецептов
async function fetchLatestRecipes() {
    try {
        const response = await fetch('/main/recipes'); // Запрос к API
        if (!response.ok) throw new Error('Ошибка при загрузке рецептов');

        const recipes = await response.json();
        displayRecipes(recipes); // Отображаем рецепты
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить рецепты.');
    }
}

// Отображение рецептов
function displayRecipes(recipes) {
    const recipesContainer = document.getElementById('recipes');
    recipesContainer.innerHTML = ''; // Очищаем контейнер

    recipes.forEach(recipe => {
        const recipeCard = document.createElement('div');
        recipeCard.className = 'recipe-card';
        recipeCard.innerHTML = `
            <img src="${recipe.photo}" alt="Фото рецепта" class="recipe-photo">
            <h3>${recipe.name}</h3>
            <p>Время приготовления: ${recipe.cook_time}</p>
            <p>Ингредиенты: ${recipe.ingredients
            .map(ing => `${ing.name}`) // Форматируем каждый ингредиент
            .join(', ')}</p>
            
        `;
        // Добавляем обработчик клика на всю карточку
        recipeCard.addEventListener('click', () => {
            window.location.href = `/recipe?recipe_id=${recipe.id}`; // Переход на страницу рецепта
        });
        // <button class="recipe-button" data-id="${recipe.id}">Перейти к рецепту</button>
        // <p>Ингредиенты: ${recipe.ingredients.join(', ')}</p>             <p>Ингредиенты: ${recipe.ingredients}</p>
        recipesContainer.appendChild(recipeCard);
    });
}

// Загружаем рецепты при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchLatestRecipes);
