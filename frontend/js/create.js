// Показываем/скрываем фильтры
document.getElementById('filterButton').addEventListener('click', () => {
    const filterPopup = document.getElementById('filters');
    filterPopup.style.display = filterPopup.style.display === 'block' ? 'none' : 'block';
});

// Функция для получения последних 10 рецептов
async function fetchLatestRecipes() {
    // try {
    //     const response = await fetch('/api/recipes/public'); // Запрос к API
    //     if (!response.ok) throw new Error('Ошибка при загрузке рецептов');

    //     const recipes = await response.json();
    //     displayRecipes(recipes); // Отображаем рецепты
    // } catch (error) {
    //     console.error('Ошибка:', error);
    //     alert('Не удалось загрузить рецепты.');
    // }
    const recipes = [
        {
            "name": "Десерты"
        },
        {
            "name": "Первое"
        },
        {
            "name": "Супы"
        },
        {
            "name": "Паста Болоньезе"
        },
        {
            "name": "Шоколадный торт"
        },
        {
            "name": "Паста Болоньезе"
        }
    ];
    displayRecipes(recipes); // Отображаем рецепты
}

// Отображение рецептов
function displayRecipes(recipes) {
    const recipesContainer = document.getElementById('recipes');
    recipesContainer.innerHTML = ''; // Очищаем контейнер

    recipes.forEach(recipe => {
        const recipeCard = document.createElement('div');
        recipeCard.className = 'recipe-card';
        recipeCard.innerHTML = `
            <h3>${recipe.name}</h3>
        `;
        recipesContainer.appendChild(recipeCard);
    });
}

// Загружаем рецепты при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchLatestRecipes);
