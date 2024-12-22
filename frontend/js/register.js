document.addEventListener("DOMContentLoaded", () => {
    const registerForm = document.getElementById("register-form");

    registerForm.addEventListener("submit", async (event) => {
        event.preventDefault(); // Отключаем стандартное поведение формы

        const username = document.getElementById("username").value.trim();
        const password = document.getElementById("password").value.trim();
        const confirmPassword = document.getElementById("confirm-password").value.trim();

        // Проверка совпадения паролей
        if (password !== confirmPassword) {
            Swal.fire({
                icon: 'error',
                title: 'Ошибка!',
                text: 'Пароли не совпадают!',
                timer: 2000,
                showConfirmButton: false
            });
            return;
        }

        try {
            // Отправка данных на сервер
            const response = await fetch("/auth/register", {
                method: "POST",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
            });

            if (response.ok) {
                Swal.fire({
                    icon: 'success',
                    title: 'Успех!',
                    text: "Регистрация успешна! Теперь войдите в систему.",
                    timer: 2000,
                    showConfirmButton: false
                });

                // Редирект на страницу входа
                setTimeout(() => {
                    window.location.href = "/auth";
                }, 2000);
            } else {
                const errorText = await response.text();
                Swal.fire({
                    icon: 'error',
                    title: 'Ошибка!',
                    text: `Ошибка: ${errorText}`,
                    timer: 2000,
                    showConfirmButton: false
                });
            }
        } catch (error) {
            console.error("Ошибка запроса:", error);
            Swal.fire({
                icon: 'error',
                title: 'Ошибка!',
                text: 'Произошла ошибка при выполнении запроса.',
                timer: 2000,
                showConfirmButton: false
            });
        }
    });
});
