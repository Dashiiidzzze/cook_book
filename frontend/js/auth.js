document.addEventListener("DOMContentLoaded", () => {
    const form = document.getElementById("wrapper");

    form.addEventListener("submit", async (event) => {
        event.preventDefault(); // Отключаем стандартное поведение формы

        // Получаем логин, пароль
        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;
        //const rememberMe = document.getElementById("rememberMe").checked;
        const action = event.submitter.value; // Определяем, какую кнопку нажали

        let url = "";
        if (action === "login") {
            url = "/auth/login"; // Маршрут для входа
        } else if (action === "register") {
            url = "/auth/register"; // Маршрут для регистрации
        }

        try {
            // Отправляем POST-запрос
            const response = await fetch(url, {
                method: "POST",
                headers: {
                    "Content-Type": "application/x-www-form-urlencoded",
                },
                body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
            });

            if (response.ok) {
                const message = action === "login" ? "Вход выполнен!" : "Регистрация успешна!\n\nВойдите в систему";
                Swal.fire({
                    icon: 'success',
                    title: 'Успех!',
                    text: message,
                    timer: 2000, // Окно исчезнет через 2 секунды
                    showConfirmButton: false,
                });
                //alert(message);

                // Для входа сохраняем токен и переходим на главную страницу
                setTimeout(() => {
                if (action === "login") {
                    window.location.href = "/main";
                } else {
                    // При регистрации перенаправляем на страницу входа
                    window.location.href = "/auth";
                }
                }, 2000);
            } else {
                const errorText = await response.text();
                Swal.fire({
                    icon: 'error',
                    title: 'Ошибка!',
                    text: `Ошибка: ${errorText}`,
                    timer: 2000, // Окно исчезнет через 2 секунды
                    showConfirmButton: false,
                });
                //alert(`Ошибка: ${errorText}`);
            }
        } catch (error) {
            console.error("Ошибка запроса:", error);
            alert("Произошла ошибка при выполнении запроса.");
        }
    });
});
