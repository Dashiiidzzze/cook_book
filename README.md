# cook_book

Онлайн ресурс, который предоставляет возможность создания собственной книги рецептов (как в публичном, так и в приватном формате) и обмена рецептами между пользователями.

## Установка

Прискачивании репозитория и запуске сайтана на OCUbuntuнеобходимо выполнить команды:

-git clone https://github.com/Dashiiidzzze/cook_book.git

-sudo apt install postgresql -y (дляустановкиpostgresql);

-sudo apt install postgresql-contrib;

-sudo apt install pgadmin4 (для установки  pgAdmin4);

-sudoaptinstallgolang(для установки языка го);

-curl-fsSLhttps://get.docker.com-oget-docker.sh(скачивание скрипта установки Docker)

-sudoshget-docker.sh(запуск скритпа установки)

-sudousermod-aGdocker$USER(добавление пользователя для запускабез root)

-sudo apt install docker-compose (для установки dockercompose);

-в директори проекта создать файл .env, который должен содержать примернотакие данные:

```apache
DB_HOST=db

DB_PORT=5431

DB_USER=имяпользователя бд

DB_PASSWORD=парольдля бд

DB_NAME=имябд

API_PORT=8080
```

-dockercomposeup- -build(в директории проекта для сборки сервера);

-dockercomposeup(в директории проекта для запускасервера).
