async function fetchUserName() {
    try {
        const response = await fetch('/profile/username'); // Запрос к API
        if (!response.ok) throw new Error('Ошибка при загрузке имени пользователя');

        const username = await response.json();
        displayUsername(username); // Отображаем рецепты
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось загрузить имя пользователя.');
    }
}

function displayUsername(data) {
    const usernameElement = document.querySelector('.profile-name');
    usernameElement.textContent = data.username || 'Неизвестный пользователь';
}

// Вызываем функцию при загрузке страницы
document.addEventListener('DOMContentLoaded', fetchUserName);
// Добавляем обработчик кнопки
//document.querySelector('.edit-button').addEventListener('click', changePassword);

document.addEventListener("DOMContentLoaded", () => {
    const changeButton = document.querySelector(".edit-form");

    changeButton.addEventListener("submit", async (event) => {
        event.preventDefault(); // Отключаем стандартное поведение кнопки

        // Получаем значения полей
        const oldPassword = document.getElementById("oldpass").value;
        const newPassword = document.getElementById("newpass").value;

        if (!oldPassword || !newPassword) {
            return;
        }

        try {
            // Отправляем POST-запрос с данными
            const response = await fetch("/profile/changepass", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    old_password: oldPassword,
                    new_password: newPassword,
                }),
            });

            if (response.ok) {
                Swal.fire({
                    icon: 'success',
                    title: 'Успех!',
                    text: 'Пароль успешно изменён!',
                    timer: 2000, // Окно исчезнет через 2 секунды
                    showConfirmButton: false,
                });
                document.getElementById("oldpass").value = "";
                document.getElementById("newpass").value = "";
            } else {
                const errorText = await response.text();
                Swal.fire({
                    icon: 'error',
                    title: 'Ошибка!',
                    text: `Ошибка: ${errorText}`,
                    timer: 2000, // Окно исчезнет через 2 секунды
                    showConfirmButton: false,
                });
            }
        } catch (error) {
            console.error("Ошибка запроса:", error);
            alert("Произошла ошибка при смене пароля.");
        }
    });
});


document.addEventListener("DOMContentLoaded", () => {
    const changeButton = document.querySelector(".logout-button");

    changeButton.addEventListener("click", async (event) => {
        event.preventDefault(); // Отключаем стандартное поведение кнопки

        try {
            // Отправляем POST-запрос с данными
            const response = await fetch("/profile/logout", {
                method: "POST",
                credentials: 'same-origin'
            });

            if (response.ok) {
                Swal.fire({
                    icon: 'success',
                    title: 'Успех!',
                    text: 'Вы успешно вышли из профиля!',
                    timer: 2000, // Окно исчезнет через 2 секунды
                    showConfirmButton: false,
                });

                // Перенаправление после задержки
                setTimeout(() => {
                    window.location.href = '/auth';
                }, 2000);
            } else {
                alert('Ошибка при выходе из профиля');
            }
        } catch (error) {
            console.error("Ошибка запроса:", error);
            alert("Произошла ошибка при выходе из профиля.");
        }
    });
});