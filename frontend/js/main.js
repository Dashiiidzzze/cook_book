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
            const response = await fetch('/main/filters'); // Запрос к серверу
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
        query.append('myrecipe', 0);

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
        <div class="recipe-text">
            <h3>${recipe.name}</h3>
            <p>Время приготовления: ${recipe.cook_time}</p>
            <p>Ингредиенты: ${recipe.ingredients
            .map(ing => `${ing.name}`) // Форматируем каждый ингредиент
            .join(', ')}</p>
            <button class="add-to-myrecipes" data-id="${recipe.id}">Добавить в мои рецепты</button>
        </div>
        <img src="data:image/jpeg;base64,${recipe.photo}" alt="Фото рецепта" class="recipe-photo">
        `;

        // Добавляем обработчик клика на кнопку
        recipeCard.querySelector('.add-to-myrecipes').addEventListener('click', (e) => {
            e.stopPropagation();  // Останавливаем всплытие события
            const recipeId = e.target.dataset.id;
            addToFavorites(recipeId); // Функция добавления в избранное
        });

        // Добавляем обработчик клика на всю карточку
        recipeCard.addEventListener('click', () => {
            window.location.href = `/recipe?recipe_id=${recipe.id}`; // Переход на страницу рецепта
        });
        recipesContainer.appendChild(recipeCard);
    });
}

async function addToFavorites(recipeId) {
    try {
        // Делаем запрос на сервер для добавления рецепта в личные рецепты
        const response = await fetch(`/api/add-to-myrecipes`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ recipeId: recipeId }), // Отправляем ID рецепта
        });

        if (response.redirected) {
            window.location.href = response.url; // Выполнить редирект
            return;
        }

        // Проверяем статус ответа
        if (response.status === 409) { // Обработка статуса "Conflict"
            Swal.fire({
                icon: 'warning',
                title: 'Внимание!',
                text: 'Это ваш рецепт. Вы не можете добавить его повторно.',
                confirmButtonColor: '#ff7c00',
            });
            return;
        }

        if (!response.ok) {
            throw new Error('Ошибка при добавлении рецепта в личные рецепты');
        }
        Swal.fire({
            icon: 'success',
            title: 'Успех!',
            text: 'Рецепт добавлен в ваши рецепты.',
            timer: 2000, // Окно исчезнет через 2 секунды
            showConfirmButton: false,
        });
        //alert('Рецепт добавлен в мои рецепты!'); // Показать сообщение об успехе
    } catch (error) {
        console.error('Ошибка:', error);
        // Показать красивое окно об ошибке
        Swal.fire({
            icon: 'error',
            title: 'Ошибка!',
            text: 'Не удалось добавить рецепт в ваши рецепты. Попробуйте снова.',
            confirmButtonColor: '#ff7c00',
        });
        //alert('Не удалось добавить рецепт в мои рецепты.');
    }
}


document.addEventListener('DOMContentLoaded', fetchLatestRecipes);
