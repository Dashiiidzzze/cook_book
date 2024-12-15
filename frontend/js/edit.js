document.addEventListener('DOMContentLoaded', async function() {
    // Получаем ID рецепта из URL или другого источника
    const recipeId = getRecipeIdFromURL(); // Реализуйте эту функцию

    // Загружаем данные рецепта с сервера
    const recipeData = await fetchRecipeData(recipeId);
    populateForm(recipeData);

    // Остальной код по добавлению ингредиентов и шагов остается прежним...
});

// Функция для получения ID рецепта из URL (пример)
function getRecipeIdFromURL() {
    const params = new URLSearchParams(window.location.search);
    return params.get('recipe_id'); // Предполагаем, что ID передаётся как ?id=123
}

// Запрос данных рецепта с сервера
async function fetchRecipeData(recipeId) {
    try {
        const response = await fetch(`/edit/${recipeId}`);
        if (!response.ok) {
            throw new Error('Ошибка загрузки данных рецепта');
        }
        return await response.json();
    } catch (error) {
        console.error(error);
        alert('Не удалось загрузить рецепт');
        return null;
    }
}

// Заполнение формы данными рецепта
function populateForm(data) {
    if (!data) return;

    document.getElementById('name').value = data.name;
    const [hours, minutes] = data.cook_time.split(':');
    document.getElementById('cook_time_hours').value = parseInt(hours, 10);
    document.getElementById('cook_time_minutes').value = parseInt(minutes, 10);
    document.getElementById('instructions').value = data.instructions;

    // Заполняем ингредиенты
    const ingredientsContainer = document.getElementById('ingredients');
    data.ingredients.forEach(ingredient => {
        const container = document.createElement('div');
        container.className = 'ingredients-container';
        container.innerHTML = `
            <input type="text" class="ingredient" value="${ingredient.name}" placeholder="Ингредиент" required>
            <input type="text" class="quantity" value="${ingredient.quantity}" placeholder="Количество" required>
        `;
        ingredientsContainer.appendChild(container);
    });

    // Заполняем шаги
    const stepsContainer = document.getElementById('steps');
    data.steps.forEach(step => {
        const container = document.createElement('div');
        container.className = 'steps-container';
        container.innerHTML = `
            <textarea class="step" placeholder="Этап приготовления" required>${step.step}</textarea>
            <input type="file" class="step-photo" accept="image/*">
        `;
        stepsContainer.appendChild(container);
    });

    // Заполняем категории
    const categorySelect = document.getElementById('categories');
    Array.from(categorySelect.options).forEach(option => {
        if (data.categories.includes(parseInt(option.value))) {
            option.selected = true;
        }
    });

    // Обновляем чекбокс "Выложить в общий доступ"
    document.getElementById('public').checked = data.public;
}