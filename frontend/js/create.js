document.addEventListener('DOMContentLoaded', async function() {
    // Загружаем категории
    const categories = await fetchCategories();
    populateCategories(categories);

    // Функция добавления ингредиента
    document.getElementById('addIngredient').addEventListener('click', function() {
        // Создаем контейнер для одного набора полей
        const ingredientsContainer = document.createElement('div');
        ingredientsContainer.className = 'ingredients-container';

        // Создаем поле для ингредиента
        const ingredientInput = document.createElement('input');
        ingredientInput.type = 'text';
        ingredientInput.classList.add('ingredient'); 
        ingredientInput.placeholder = 'Ингредиент';

        // Создаем поле для количества
        const quantityInput = document.createElement('input');
        quantityInput.type = 'text';
        quantityInput.classList.add('quantity');
        quantityInput.placeholder = 'Количество';

        // Добавляем поля в контейнер
        ingredientsContainer.appendChild(ingredientInput);
        ingredientsContainer.appendChild(quantityInput);

        // Добавляем контейнер в основной контейнер ingredients
        document.getElementById('ingredients').appendChild(ingredientsContainer);
    });

    // Функция добавления этапа
    document.getElementById('addStep').addEventListener('click', function() {
        // Создаем контейнер для одного шага
        const stepContainer = document.createElement('div');
        stepContainer.className = 'steps-container'; // Используйте корректное название класса

        // Создаем текстовое поле для описания шага
        const stepTextArea = document.createElement('textarea');
        stepTextArea.classList.add('step');
        stepTextArea.placeholder = 'Этап приготовления';

        // Создаем поле для загрузки изображения
        const stepPhotoInput = document.createElement('input');
        stepPhotoInput.type = 'file';
        stepPhotoInput.classList.add('step-photo');
        stepPhotoInput.name = 'photo';
        stepPhotoInput.accept = 'image/*';

        // Добавляем элементы в контейнер шага
        stepContainer.appendChild(stepTextArea);
        stepContainer.appendChild(stepPhotoInput);

        // Добавляем контейнер шага в основной блок "steps"
        document.getElementById('steps').appendChild(stepContainer);
    });

    document.getElementById('recipeForm').addEventListener('submit', async function (e) {
        e.preventDefault();
    
        const formData = new FormData(e.target);
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
            const stepPhoto = container.querySelector('.step-photo').files[0];
            const stepData = { step: stepText }; // Исправлено на ключ "step"
    
            // Обработка фото этапа
            if (stepPhoto) {
                stepData.photo = await fileToBase64(stepPhoto);
            }
    
            recipeData.steps.push(stepData);
        }
    
        // Главное фото
        const mainPhoto = document.getElementById('photo').files[0];
        if (mainPhoto) {
            recipeData.photo = await fileToBase64(mainPhoto);
        }
    
        // Категории
        recipeData.categories = Array.from(document.getElementById('categories').selectedOptions)
            .map(option => parseInt(option.value, 10)); // Преобразование к числу
    
        // Обработка галочки "выложить в общий доступ"
        recipeData.public = document.getElementById('public').checked;
    
        // Отправка на сервер
        await saveRecipe(recipeData);
    });
    
    // Конвертация файла в base64
    function fileToBase64(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result.split(',')[1]); // Убираем prefix
            reader.onerror = error => reject(error);
            reader.readAsDataURL(file);
        });
    }
    
    // Отправка данных на сервер
    async function saveRecipe(data) {
        try {
            const response = await fetch('/create/save', {
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
});

// Функция загрузки категорий
async function fetchCategories() {
    try {
        const response = await fetch('/categories/all'); // Запрос к API для получения категорий
        if (!response.ok) throw new Error('Ошибка при загрузке категорий');
        return await response.json();
    } catch (error) {
        console.error('Ошибка:', error);
        return [];
    }
}

// Заполнение выпадающего списка категориями
function populateCategories(categories) {
    const selectElement = document.getElementById('categories');
    // Очищаем существующие категории (если есть)
    selectElement.innerHTML = '';

    categories.forEach(category => {
        const option = document.createElement('option');
        option.value = category.id;
        option.textContent = category.name;
        selectElement.appendChild(option);
    });
}
