document.addEventListener("DOMContentLoaded", () => {
    const registerBtn = document.getElementById("register-btn");

    // Обработчик для кнопки "Зарегистрироваться"
    registerBtn.addEventListener("click", () => {
        window.location.href = "/registration";
    });

    const form = document.getElementById("wrapper");

    form.addEventListener("submit", async (event) => {
        event.preventDefault(); // Отключаем стандартное поведение формы

        const username = document.getElementById("username").value;
        const password = document.getElementById("password").value;

        // Логика входа
        try {
            const response = await fetch("/auth/login", {
                method: "POST",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`
            });

            if (response.ok) {
                Swal.fire({
                    icon: 'success',
                    title: 'Успех!',
                    text: "Вход выполнен!",
                    timer: 2000,
                    showConfirmButton: false,
                });
                setTimeout(() => {
                    window.location.href = "/main";
                }, 2000);
            } else {
                const errorText = await response.text();
                Swal.fire({
                    icon: 'error',
                    title: 'Ошибка!',
                    text: `Ошибка: ${errorText}`,
                    timer: 2000,
                    showConfirmButton: false,
                });
            }
        } catch (error) {
            console.error("Ошибка запроса:", error);
            alert("Произошла ошибка при выполнении запроса.");
        }
    });
});
