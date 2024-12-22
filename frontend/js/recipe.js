// Функция для получения рецепта по id
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
        if (!response.ok) throw new Error('Ошибка при загрузке рецепта');

        const recipe = await response.json();
        displayRecipe(recipe); // Отображаем рецепты
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить рецепт.');
    }
}

function displayRecipe(data) {
    const recipesContainer = document.getElementById('recipe');
    const commentsContainer = document.getElementById('comments'); // Находим контейнер для комментариев

    if (!recipesContainer) {
        console.error('Ошибка: элемент с ID "recipe" не найден');
        return;
    }
    if (!commentsContainer) {
        console.error('Ошибка: элемент с ID "comments" не найден');
        return;
    }
    
    const { recipe, comments } = data; // Извлекаем рецепт и комментарии из данных

    recipesContainer.innerHTML = '';
    commentsContainer.innerHTML = '';

    // Создаем карточку рецепта
    const recipeCard = document.createElement('div');
    recipeCard.className = 'recipe-card';

    // Проверка на наличие главного фото, заглушка если отсутствует
    const recipePhoto = recipe.photo
        ? `data:image/jpeg;base64,${recipe.photo}`
        : '../images/default-placeholder.png'; // Путь к заглушке
    
    recipeCard.innerHTML = `
        <h3>${recipe.name}</h3>
        <img src="${recipePhoto}" alt="Фото рецепта" class="recipe-photo">
        <p><strong>Время приготовления:</strong> ${recipe.cook_time}</p>
        <p><strong>Категории:</strong> ${recipe.categories.join(', ')}</p> <!-- Категории -->
        <p><strong>Ингредиенты:</strong></p>
        <ul>${recipe.ingredients
            .map(ing => `<li>${ing.name} - ${ing.quantity}</li>`) // Форматируем ингредиенты
            .join('')}</ul>
        <p><strong>Рецепт:</strong> ${recipe.instructions}</p>
        
    `;
    recipesContainer.appendChild(recipeCard);

    // Вставляем этапы приготовления внутри карточки рецепта
    if (recipe.steps.length > 0) {
        const stepsSection = document.createElement('div');
        stepsSection.className = 'steps-section';
        stepsSection.innerHTML = '<p><strong>Этапы приготовления:</strong></p>';
        
        recipe.steps.forEach((step, index) => {
            const stepCard = document.createElement('div');
            stepCard.className = 'step-card';
            stepCard.innerHTML = `
                <p><strong>Шаг ${index + 1}:</strong> ${step.step}</p>
                ${step.photo ? `<img src="data:image/jpeg;base64,${step.photo}" alt="Фото шага" class="step-photo">` : ''}
            `;
            stepsSection.appendChild(stepCard);
        });

        recipeCard.appendChild(stepsSection); // Добавляем блок с этапами в карточку рецепта
    } else {
        const noSteps = document.createElement('p');
        noSteps.textContent = 'Нет этапов приготовления для этого рецепта.';
        recipeCard.appendChild(noSteps); // Сообщение о том, что этапов нет
    }

    recipesContainer.appendChild(recipeCard); // Добавляем карточку рецепта в контейнер

    // Отображаем комментарии в отдельном контейнере
    if (comments.length > 0) {
        commentsContainer.innerHTML = '<h4>Комментарии:</h4>'; // Очистим контейнер для комментариев и добавим заголовок

        comments.forEach(comment => {
            const commentCard = document.createElement('div');
            commentCard.className = 'comment-card';
            commentCard.innerHTML = `
                <p><strong>${comment.username}:</strong> ${comment.text}</p>
            `;
            commentsContainer.appendChild(commentCard); // Добавляем каждый комментарий в контейнер
        });
    } else {
        commentsContainer.innerHTML = '<p>Пока никто не написал комментарии к рецепту, будьте первым!</p>'; // Сообщение о том, что комментариев нет
    }
}

// Ждем полной загрузки DOM
document.addEventListener('DOMContentLoaded', () => {
    fetchRecipe();

    // Находим кнопку по ID
    const submitButton = document.getElementById('submit-comment');

    // Добавляем обработчик события на клик
    submitButton.addEventListener('click', async () => {
        const commentText = document.getElementById('comment-text').value;

        // Проверяем, что поле ввода не пустое
        if (!commentText) {
            Swal.fire({
                icon: 'warning',
                title: 'Внимание!',
                text: 'Нельзя добавить пустой комментарий.',
                confirmButtonColor: '#ff7c00',
            });
            //alert('Пожалуйста, введите комментарий.');
            return;
        }

        console.log('Отправка комментария:', commentText);

        // Логика отправки комментария
        await submitCommentToServer(commentText);
    });
});

// Функция для отправки комментария на сервер
async function submitCommentToServer(commentText) {
    try {
        const recipeId = new URLSearchParams(window.location.search).get('recipe_id');
        if (!recipeId) {
            console.error('ID рецепта отсутствует.');
            return;
        }

        if (!commentText.trim()) {
            Swal.fire({
                icon: 'warning',
                title: 'Внимание!',
                text: 'Нельзя добавить пустой комментарий.',
                confirmButtonColor: '#ff7c00',
            });
            //alert('Нельзя отправить пустой комментарий');
            return;
        }
        
        // Отправляем запрос на сервер
        const response = await fetch('/recipe/add-comment', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ recipe_id: recipeId, comment: commentText }),
        });

        if (response.redirected) {
            window.location.href = response.url; // Выполнить редирект
            return;
        }

        if (!response.ok) {
            throw new Error('Ошибка при отправке комментария.');
        }

        // Очистка поля ввода
        document.getElementById('comment-text').value = '';

        // Обновление комментариев на странице
        Swal.fire({
            icon: 'success',
            title: 'Успех!',
            text: 'Комментарий добавлен.',
            timer: 2000, // Окно исчезнет через 2 секунды
            showConfirmButton: false,
        });
        //alert('Комментарий успешно отправлен!');
        console.log('Комментарий добавлен.');

        fetchRecipe() // перезагрузка страницы
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось отправить комментарий.');
    }
}