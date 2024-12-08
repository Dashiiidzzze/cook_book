// Показываем/скрываем фильтры
// document.getElementById('filterButton').addEventListener('click', () => {
//     const filterPopup = document.getElementById('filters');
//     filterPopup.style.display = filterPopup.style.display === 'block' ? 'none' : 'block';
// });

// Функция для получения последних 10 рецептов
async function fetchRecipe() {
    try {
        // Получаем recipe_id из текущего URL
        const params = new URLSearchParams(window.location.search);
        const recipeId = params.get('recipe_id');

        if (!recipeId) {
            console.error('Параметр recipe_id отсутствует в URL');
            return;
        }
        const response = await fetch(`/recipe/view?recipe_id=${recipeId}`); // Запрос к API
        if (!response.ok) throw new Error('Ошибка при загрузке рецептов');

        const recipe = await response.json();
        displayRecipe(recipe); // Отображаем рецепты
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить рецепты.');
    }
}

// Отображение рецептов
function displayRecipe(recipe) {
    const recipesContainer = document.getElementById('recipe');
    if (!recipesContainer) {
        console.error('Ошибка: элемент с ID "recipe" не найден');
        return;
    }
    //recipesContainer.innerHTML = ''; // Очищаем контейнер

    const recipeCard = document.createElement('div');
    recipeCard.className = 'recipe-card';
    recipeCard.innerHTML = `
        <h3>${recipe.name}</h3>
        <p>Время приготовления: ${recipe.cook_time}</p>
        <p>Категория: ${recipe.dish_type}</p>
        <p>Ингредиенты: ${recipe.ingredients}</p>
        <p>Рецепт: ${recipe.instructions}</p>
        <img src="${recipe.photo}" alt="Фото рецепта" class="recipe-photo">
    `;
    // recipeCard.innerHTML = `
    //     <h3>${recipe.name || 'Название отсутствует'}</h3>
    //     <p>Время приготовления: ${recipe.cook_time || 'Не указано'}</p>
    //     <p>Категория: ${recipe.dish_type || 'Не указана'}</p>
    //     <p>Ингредиенты: ${recipe.ingredients || 'Нет данных'}</p>
    //     <p>Рецепт: ${recipe.instructions || 'Инструкции отсутствуют'}</p>
    //     <img src="${recipe.photo || 'placeholder.jpg'}" 
    //             onerror="this.src='placeholder.jpg'" 
    //             alt="Фото рецепта" class="recipe-photo">
    // `;
    recipesContainer.appendChild(recipeCard);
}

// Загружаем рецепты при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchRecipe);