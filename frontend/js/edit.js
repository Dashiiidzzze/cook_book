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
        const formData = collectFormData();
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

    // Устанавливаем главное фото
    if (recipe.photo) {
        setMainPhoto(recipe.photo);
    }

    document.getElementById('public').checked = recipe.public;
}

// // Добавление категорий: загрузка всех и выделение нужных
// async function loadCategories(selectedCategories) {
//     try {
//         const response = await fetch('/categories/all'); // Запрос всех категорий
//         const categories = await response.json();

//         const select = document.getElementById('categories');
//         categories.forEach(category => {
//             const option = document.createElement('option');
//             option.value = category.id;
//             option.textContent = category.name;
//             if (selectedCategories.includes(category.id)) {
//                 option.selected = true;
//             }
//             select.appendChild(option);
//         });
//     } catch (error) {
//         console.error('Ошибка загрузки категорий:', error);
//     }
// }

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

// async function loadCategories(selectedCategories) {
//     try {
//         const response = await fetch('/categories/all'); // Запрос всех категорий
//         const categories = await response.json();

//         const container = document.getElementById('category-container');
//         container.innerHTML = ''; // Очищаем контейнер

//         categories.forEach(category => {
//             const div = document.createElement('div');
//             div.classList.add('category-item');
//             div.textContent = category.name;
//             div.dataset.id = category.id;

//             // Приведение типов для корректного сравнения
//             if (selectedCategories.map(Number).includes(Number(category.id))) {
//                 div.classList.add('selected');
//             }

//             // Добавляем обработчик клика для переключения состояния
//             div.addEventListener('click', () => {
//                 div.classList.toggle('selected');
//             });

//             container.appendChild(div);
//         });
//     } catch (error) {
//         console.error('Ошибка загрузки категорий:', error);
//     }
// }


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
    const formData = {
        name: document.getElementById('name').value,
        cook_time: {
            hours: document.getElementById('cook_time_hours').value,
            minutes: document.getElementById('cook_time_minutes').value
        },
        instructions: document.getElementById('instructions').value,
        ingredients: [],
        steps: [],
        categories: Array.from(document.getElementById('categories').selectedOptions).map(opt => opt.value),
        public: document.getElementById('public').checked,
        photo: null
    };

    // Главное фото
    const photoInput = document.getElementById('photo');
    if (photoInput.files.length > 0) {
        formData.photo = await fileToBase64(photoInput.files[0]);
    } else if (photoInput.dataset.existingPhoto) {
        formData.photo = photoInput.dataset.existingPhoto;
    }

    // Ингредиенты
    document.querySelectorAll('.ingredients-container').forEach(container => {
        const name = container.querySelector('.ingredient').value;
        const quantity = container.querySelector('.quantity').value;
        formData.ingredients.push({ name, quantity });
    });

    // Этапы
    for (const container of document.querySelectorAll('.steps-container')) {
        const stepText = container.querySelector('.step').value;
        const photoInput = container.querySelector('.step-photo');
        const photoPreview = container.querySelector('.step-photo-preview');

        const stepData = { step: stepText };
        if (photoInput.files.length > 0) {
            stepData.photo = await fileToBase64(photoInput.files[0]);
        } else if (photoPreview) {
            stepData.photo = photoPreview.src.split(',')[1]; // Удаляем prefix base64
        }

        formData.steps.push(stepData);
    }

    return formData;
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
    const formData = new FormData();

    formData.append('name', data.name);
    formData.append('cook_time_hours', data.cook_time.hours);
    formData.append('cook_time_minutes', data.cook_time.minutes);
    formData.append('instructions', data.instructions);

    if (data.photo) {
        formData.append('photo', data.photo);
    }

    data.ingredients.forEach((ingredient, index) => {
        formData.append(`ingredients[${index}][name]`, ingredient.name);
        formData.append(`ingredients[${index}][quantity]`, ingredient.quantity);
    });

    data.steps.forEach((step, index) => {
        formData.append(`steps[${index}][step]`, step.step);
        formData.append(`steps[${index}][photo]`, step.photo);
    });

    data.categories.forEach((category, index) => {
        formData.append(`categories[${index}]`, category);
    });

    formData.append('public', data.public);

    try {
        const response = await fetch(`/edit/save?recipe_id=${recipeId}`, {
            method: 'POST',
            body: formData
        });
        if (!response.ok) throw new Error('Ошибка сохранения данных');
        alert('Рецепт успешно сохранен!');
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось сохранить данные рецепта.');
    }
}
