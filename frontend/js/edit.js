document.addEventListener('DOMContentLoaded', () => {
    const recipeForm = document.getElementById('recipeForm');
    const recipeId = getRecipeIdFromURL(); // Функция получения ID из URL

    // Загружаем данные рецепта при загрузке страницы
    loadRecipeData(recipeId);

    // Добавляем обработчик кнопки "Добавить ингредиент"
    document.getElementById('addIngredient').addEventListener('click', addIngredientField);

    // Добавляем обработчик кнопки "Добавить этап"
    document.getElementById('addStep').addEventListener('click', addStepField);

    // Обработчик формы на отправку
    recipeForm.addEventListener('submit', async (event) => {
        event.preventDefault();
        const formData = await collectFormData();
        await saveRecipeData(recipeId, formData);
    });
});

// Функция получения ID рецепта из URL
function getRecipeIdFromURL() {
    const params = new URLSearchParams(window.location.search);
    return params.get('recipe_id'); // например, /edit?id=123
}

// Загрузка данных рецепта с сервера
async function loadRecipeData(recipeId) {
    try {
        const response = await fetch(`/recipe/view?recipe_id=${recipeId}`);
        if (!response.ok) throw new Error('Ошибка загрузки данных рецепта');

        const data = await response.json();
        fillFormWithData(data); // Заполняем форму
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить данные рецепта.');
    }
}

// Заполнение формы данными рецепта
function fillFormWithData(data) {
    const { recipe, comments } = data;
    document.getElementById('name').value = recipe.name;
    // Обработка cook_time (формат "00:00:00")
    const [hours, minutes] = recipe.cook_time.split(':');
    document.getElementById('cook_time_hours').value = parseInt(hours, 10);
    document.getElementById('cook_time_minutes').value = parseInt(minutes, 10);
    // document.getElementById('cook_time_hours').value = recipe.cook_time;
    // document.getElementById('cook_time_minutes').value = recipe.cook_time;
    document.getElementById('instructions').value = recipe.instructions;

    // Добавляем ингредиенты
    const ingredientsContainer = document.getElementById('ingredients');
    recipe.ingredients.forEach(ingredient => {
        addIngredientField(ingredient.name, ingredient.quantity);
    });

    // Добавляем этапы и фото
    const stepsContainer = document.getElementById('steps');
    recipe.steps.forEach(step => {
        addStepField(step.step, step.photo);
    });

    // Устанавливаем категории
    loadCategories(recipe.categories);

    // // Устанавливаем главное фото
    // if (recipe.photo) {
    //     setMainPhoto(recipe.photo);
    // }
    // Устанавливаем главное фото
    if (recipe.photo) {
        setMainPhoto(recipe.photo);
        const photoInput = document.getElementById('photo');
        photoInput.dataset.existingPhoto = recipe.photo; // Устанавливаем существующее фото
    }

    document.getElementById('public').checked = recipe.public;
}

async function loadCategories(selectedCategories) {
    try {
        const response = await fetch('/categories/all'); // Запрос всех категорий
        const categories = await response.json();

        const select = document.getElementById('categories');
        select.innerHTML = ''; // Очищаем старые опции, если они есть

        categories.forEach(category => {
            const option = document.createElement('option');
            option.value = String(category.id);
            option.textContent = category.name;

            // Устанавливаем опцию как выбранную, если она есть в selectedCategories
            if (selectedCategories.map(String).includes(String(category.id))) {
                option.selected = true;
            }
            select.appendChild(option);
        });

        // Принудительно обновляем select для корректного отображения
        select.dispatchEvent(new Event('change')); 
    } catch (error) {
        console.error('Ошибка загрузки категорий:', error);
    }
}

// Добавление полей для ингредиентов
function addIngredientField(name = '', quantity = '') {
    const container = document.getElementById('ingredients');
    const div = document.createElement('div');
    div.classList.add('ingredients-container');
    div.innerHTML = `
        <input type="text" class="ingredient" placeholder="Ингредиент" value="${name}" required>
        <input type="text" class="quantity" placeholder="Количество" value="${quantity}" required>
    `;
    container.appendChild(div);
}

// Добавление полей для этапов с фото
function addStepField(description = '', photo = '') {
    const container = document.getElementById('steps');
    const div = document.createElement('div');
    div.classList.add('steps-container');
    div.innerHTML = `
        <textarea type="text" class="step" placeholder="Этап приготовления" required>${description}</textarea>
        <input type="file" class="step-photo" accept="image/*">
        ${photo ? `<img class="step-photo-preview" src="data:image/jpeg;base64,${photo}" alt="Предыдущее фото" style="max-width: 200px; display: block;">` : ''}
    `;
    container.appendChild(div);
}

// Установка главного фото
function setMainPhoto(photo) {
    const photoInput = document.getElementById('photo');
    const img = document.createElement('img');
    img.src = `data:image/jpeg;base64,${photo}`;
    img.style.maxWidth = '200px';
    img.style.display = 'block';
    photoInput.insertAdjacentElement('afterend', img);
}

// Сбор данных формы
async function collectFormData() {
    const formElement = document.getElementById('recipeForm');
    const formData = new FormData(formElement);
    const recipeData = {};

    // Сбор основных данных
    recipeData.name = formData.get('name');

    // Время приготовления
    const hours = parseInt(formData.get('cook_time_hours'), 10) || 0;
    const minutes = parseInt(formData.get('cook_time_minutes'), 10) || 0;
    recipeData.cook_time = `${hours}:${minutes < 10 ? '0' : ''}${minutes}`;

    // Ингредиенты
    recipeData.ingredients = [];
    document.querySelectorAll('.ingredients-container').forEach(container => {
        const ingredient = container.querySelector('.ingredient').value.trim();
        const quantity = container.querySelector('.quantity').value.trim();
        if (ingredient && quantity) {
            recipeData.ingredients.push({ name: ingredient, quantity });
        }
    });

    // Инструкции
    recipeData.instructions = formData.get('instructions');

    // Этапы приготовления
    recipeData.steps = [];
    const stepContainers = document.querySelectorAll('.steps-container');
    for (const container of stepContainers) {
        const stepText = container.querySelector('.step').value.trim();
        const photoInput = container.querySelector('.step-photo');
        const photoPreview = container.querySelector('.step-photo-preview');
        const stepData = { step: stepText }; // Исправлено на ключ "step"
        //????????????????????/
        if (photoInput && photoInput.files.length > 0) {
            stepData.photo = await fileToBase64(photoInput.files[0]);
        } else if (photoPreview && photoPreview.src) {
            console.log('Передаем старое фото этапа:', photoPreview.src);
            stepData.photo = photoPreview.src.split(',')[1]; // Удаляем prefix base64
        }

        // if (photoPreview && photoPreview.src) {
        //     console.log('Старое фото этапа:', photoPreview.src); // Проверяем, что содержится в photoPreview.src
        //     stepData.photo = photoPreview.src.split(',')[1]; // Убираем префикс base64
        // } else {
        //     console.log('Нет фото этапа');
        // }

        // // Обработка фото этапа
        // if (stepPhoto) {
        //     stepData.photo = await fileToBase64(stepPhoto);
        // }

        recipeData.steps.push(stepData);
    }

    // Главное фото
    const photoInput = document.getElementById('photo');
    if (photoInput && photoInput.files.length > 0) {
        recipeData.photo = await fileToBase64(photoInput.files[0]);
    } else if (photoInput.dataset.existingPhoto) {
        recipeData.photo = photoInput.dataset.existingPhoto;
    }

    // Категории
    recipeData.categories = Array.from(document.getElementById('categories').selectedOptions)
        .map(option => parseInt(option.value, 10)); // Преобразование к числу

    // Обработка галочки "выложить в общий доступ"
    recipeData.public = document.getElementById('public').checked;
    console.log('Собранные данные:', recipeData);
    return recipeData;
}

// Конвертация файла в Base64
function fileToBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onload = () => resolve(reader.result.split(',')[1]);
        reader.onerror = error => reject(error);
        reader.readAsDataURL(file);
    });
}

// Сохранение рецепта
async function saveRecipeData(recipeId, data) {
    console.log("Collected Data:", data);
    try {
        const response = await fetch(`/edit/save?recipe_id=${recipeId}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            throw new Error('Ошибка при сохранении рецепта');
        }

        //alert('Рецепт успешно сохранен!');
        Swal.fire({
            icon: 'success',
            title: 'Успех!',
            text: 'Рецепт успешно сохранен!',
            timer: 2000, // Окно исчезнет через 2 секунды
            showConfirmButton: false,
        });

        // Перенаправление после задержки
        setTimeout(() => {
            window.location.href = '/myrecipes';
        }, 2000);
    } catch (error) {
        console.error('Ошибка:', error);
        Swal.fire({
            icon: 'error',
            title: 'Ошибка!',
            text: 'Не удалось сохранить рецепт.',
            timer: 2000, // Окно исчезнет через 2 секунды
            showConfirmButton: false,
        });
        //alert('Не удалось сохранить рецепт.');
    }
}
