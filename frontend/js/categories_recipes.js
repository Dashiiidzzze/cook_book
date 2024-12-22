let currentCategoryId = null; // Глобальная переменная для хранения ID категории

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
            const response = await fetch(`/categories-recipes/filters?category_id=${currentCategoryId}`); // Запрос к серверу
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
        query.append('category', currentCategoryId); //так как не выбрана категория
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

document.addEventListener('DOMContentLoaded', async () => {
    // Получаем параметры из URL
    const urlParams = new URLSearchParams(window.location.search);
    const categoryId = urlParams.get('category_id'); // Получаем ID категории
    const categoryName = urlParams.get('category_name'); // Получаем ID категории

    if (!categoryId) {
        console.error('Параметр category_id не найден');
        alert('Не указана категория для поиска рецептов.');
        return;
    }
    currentCategoryId = categoryId; // Сохраняем category_id глобально

    try {
        // Делаем запрос к серверу для получения рецептов по категории
        const response = await fetch(`/categories-recipes/recipes?category_id=${categoryId}`);
        if (!response.ok) throw new Error('Ошибка при загрузке рецептов категории');

        const recipes = await response.json(); // Получаем список рецептов
        displayRecipes(recipes, categoryName); // Отображаем рецепты с помощью уже готовой функции
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить рецепты для этой категории.');
    }
});

// Функция для отображения рецептов
function displayRecipes(recipes, categoryName) {
    const recipesSection = document.getElementById('recipes-section');
    if (!recipesSection) {
        console.error('Секция #recipes-section не найдена.');
        return;
    }

    // Если заголовок еще не создан, создаём его
    let title = recipesSection.querySelector('h2');
    if (!title) {
        title = document.createElement('h2');
        recipesSection.appendChild(title);
        // Обновляем текст заголовка
        title.textContent = 'Рецепты в категории ' + categoryName;
    }

    // Удаляем предыдущие рецепты, но не заголовок
    const oldRecipesContainer = recipesSection.querySelector('.recipe-container');
    if (oldRecipesContainer) oldRecipesContainer.remove();

    // Создаём контейнер для рецепта
    const recipesContainer = document.createElement('div');
    recipesContainer.id = 'recipe';
    recipesContainer.className = 'recipe-container';

    // Создаём карточки для каждого рецепта
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

    recipesSection.appendChild(recipesContainer);
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