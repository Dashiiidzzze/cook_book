// Загружаем категории при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchCategories);

// Функция для получения rкатегорий
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

// Отображение категорий
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
            window.location.href = `/categories-recipes?category_id=${category.id}&category_name=${category.name}`;
        });

        categoriesContainer.appendChild(categoryCard);
    });

    categoriesSection.appendChild(categoriesContainer);
}