// Показываем/скрываем фильтры
document.getElementById('filterButton').addEventListener('click', async () => {
    const filterPopup = document.getElementById('filters');
        // Если фильтры уже загружены, просто показываем/скрываем
        if (filterPopup.dataset.loaded === "true") {
            filterPopup.style.display = filterPopup.style.display === 'block' ? 'none' : 'block';
            return;
        }
    
        // Загружаем фильтры с сервера
        try {
            const response = await fetch('/myrecipes/filter'); // Запрос к серверу
            if (!response.ok) throw new Error('Ошибка при загрузке фильтров');
    
            const filters = await response.json();
            displayFilters(filters);
            filterPopup.dataset.loaded = "true"; // Помечаем, что фильтры загружены
            filterPopup.style.display = 'block';
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Не удалось загрузить фильтры.');
        }
});

// Функция для отображения загруженных фильтров
function displayFilters(filters) {
    const filterPopup = document.getElementById('filters');
    filterPopup.innerHTML = ''; // Очищаем старые фильтры

    filters.forEach(filter => {
        const label = document.createElement('label');
        label.innerHTML = `
            <input type="checkbox" value="${filter.id}"> ${filter.name}
        `;
        filterPopup.appendChild(label);
    });
}

// Назначаем обработчик события на кнопку "Поиск"
document.getElementById('searchButton').addEventListener('click', fetchRecipesWithFilters);

// Функция для поиска рецептов с фильтрами
async function fetchRecipesWithFilters() {
    const searchText = document.getElementById('searchBar').value; // Получаем текст из поля поиска
    const selectedFilters = Array.from(document.querySelectorAll('#filters input:checked'))
        .map(input => input.value); // Собираем отмеченные фильтры

    try {
        // Формируем параметры запроса
        const query = new URLSearchParams();
        if (searchText) query.append('search', searchText); // Добавляем текст поиска
        if (selectedFilters.length > 0) query.append('filters', selectedFilters.join(',')); // Добавляем фильтры
        query.append('category', 0); //так как не выбрана категория
        query.append('myrecipe', 1);

        // Выполняем запрос на сервер
        const response = await fetch(`/api/search?${query.toString()}`);
        if (!response.ok) throw new Error('Ошибка при поиске рецептов');

        // Обрабатываем и отображаем полученные рецепты
        const recipes = await response.json();
        displayRecipes(recipes);
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось выполнить поиск рецептов.');
    }
}

// Функция для получения последних 10 рецептов
async function fetchLatestRecipes() {
    try {
        const response = await fetch('/myrecipes/recipes'); // Запрос к API
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
        <div class="recipe-text">
            <h3>${recipe.name}</h3>
            <p>Время приготовления: ${recipe.cook_time}</p>
            <p>Ингредиенты: ${recipe.ingredients
            .map(ing => `${ing.name}`) // Форматируем каждый ингредиент
            .join(', ')}</p>
            <button class="delete-myrecipes" data-id="${recipe.id}" data-name="${recipe.name}">Удалить рецепт</button>
            <button class="edit-myrecipes" data-id="${recipe.id}">Редактировать рецепт</button>
        </div>
        <img src="data:image/jpeg;base64,${recipe.photo}" alt="Фото рецепта" class="recipe-photo">
        `;

        // Обработчик клика на кнопку "Удалить рецепт"
        recipeCard.querySelector('.delete-myrecipes').addEventListener('click', (event) => {
            event.stopPropagation(); // Остановка всплытия события
            const recipeId = event.target.dataset.id;
            const recipeName = event.target.dataset.name;

            // Показываем подтверждающее окно
            showConfirmModal(recipeName, (isConfirmed) => {
                if (isConfirmed) {
                    deleteRecipe(recipeId); // Удаляем рецепт, если подтверждено
                }
            }); 
        });

        // Добавляем обработчик клика на кнопку "Редактировать рецепт"
        recipeCard.querySelector('.edit-myrecipes').addEventListener('click', (event) => {
            event.stopPropagation(); // Остановка всплытия события
            const recipeId = event.target.dataset.id;
            window.location.href = `/edit?recipe_id=${recipeId}`; // Переход на страницу редактирования
        });

        // Добавляем обработчик клика на всю карточку
        recipeCard.addEventListener('click', () => {
            window.location.href = `/recipe?recipe_id=${recipe.id}`; // Переход на страницу рецепта
        });
        recipesContainer.appendChild(recipeCard);
    });
}

// Функция для отображения подтверждающего окна
function showConfirmModal(recipeName, callback) {
    const modal = document.getElementById('confirmModal');
    const confirmText = document.getElementById('confirmText');
    confirmText.textContent = `Вы точно хотите удалить рецепт "${recipeName}"?`;

    modal.style.display = 'flex';

    // Обработчик кнопки "Да"
    document.getElementById('confirmYes').onclick = () => {
        callback(true); // Подтверждение удаления
        modal.style.display = 'none';
    };

    // Обработчик кнопки "Отмена"
    document.getElementById('confirmNo').onclick = () => {
        callback(false); // Отмена удаления
        modal.style.display = 'none';
    };
}

// Функция для отправки запроса на удаление рецепта
async function deleteRecipe(recipeId) {
    try {
        const response = await fetch(`/myrecipes/recipes/${recipeId}`, { 
            method: 'DELETE' 
        });
        if (!response.ok) throw new Error('Не удалось удалить рецепт');
        Swal.fire({
            icon: 'success',
            title: 'Успех!',
            text: 'Рецепт удален.',
            timer: 2000, // Окно исчезнет через 2 секунды
            showConfirmButton: false,
        });
        fetchLatestRecipes(); // Перезагрузка списка рецептов
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка при удалении рецепта.');
    }
}

// Загружаем рецепты при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchLatestRecipes);
