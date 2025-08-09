<p align="center">
  <picture>
    <source height="125" media="(prefers-color-scheme: dark)" srcset="assets/teacher.png">
    <img height="125" alt="Fiber" src="assets/light-teacher.png">
  </picture>
</p>

<p align="center">
  <strong>MS-teacher</strong> — это <strong>микросервис для учителей</strong>, который предоставляет API для взаимодействия с оценками, а также отправкой сообщений к привязанным соцсетям.
</p>

# 💡Основные возможности

- **Получение оценок**: Получение всех оценок, получение оценок для каждого предмета по отдельности.  
- **Аналитика**: Аналитика оценок студента.

# 🤖 Используемые технологии

- **Golang** — основной язык программирования
- **MS-database** — основая библиотека для микросервиса
- **Fiber** — фреймворк для написания REST API
- **MongoDB** — основная база данных
- **Redis** — для управления доступом
- **Docker** — развертывание проекта

# ⚠️Важно
Перед стартом необходимо перейти в [MS-database](https://github.com/Muraddddddddd9/ms-database) и поднять MongoDB, Redis, S3

# ⚡️ Быстрый старт
Перейти в env и поменять конфигурацию
```env
MONGO_NAME=diary
# MONGO_HOST=localhost # <- для локального запуска
MONGO_HOST=host.docker.internal # <- для запуска в Docker 
MONGO_PORT=27018 # <- ваш порт (27018 для Docker)
MONGO_USERNAME=college # <- username для MongoDB
MONGO_PASSWORD=BIM_LOCAL1 # <- пароль для MongoDB
MONGO_AUTH_SOURCE=admin

REDIS_DB=0
# REDIS_HOST=localhost # <- для локального запуска
REDIS_HOST=host.docker.internal # <- для запуска в Docker 
REDIS_PASSWORD=BIM_LOCAL1 # <- пароль для MongoDB
REDIS_PORT=6380 # <- ваш порт (6380 для Docker)

ORIGIN_URL=http://localhost:5173 # <- адрес сайта
PROJECT_PORT=:8082 # <- порт приложения
```

## CMD
Клонирование репозитория
```bash
git clone https://github.com/Muraddddddddd9/ms-teacher.git
```
Установка всех пакетов
```bash
go get .
```
Запустить программу
```bash
go run .
```
## Docker
Клонирование репозитория
```bash
git clone https://github.com/Muraddddddddd9/ms-teacher.git
```
Билд Docker container 
```bash
docker-compose build
```
Поднятие Docker container 
```bash
docker-compose up
```

# 🧬 API
- <strong>get_evaluation/:group/:object<strong> - Get, получение оценок определённого предмета для определённой группы
- <strong>send_evaluation<strong> - Post, выставление оценки студенту
- <strong>delete_evaluation/:id<strong> - Delete, удаление оценки по id
- <strong>get_my_classroom_group<strong> - Get, получение списка группы учителя
- <strong>get_my_classroom_object/:group<strong> - Get, получение списка предметов по id группы
- <strong>message_contest<strong> - Post, отправка сообщения студенту на одну из соц. сетей

# 🧩 Остальные
- <strong>[MS-admin](https://github.com/Muraddddddddd9/ms-admin)</strong> - микросервис (необходимый)
- <strong>[MS-database](https://github.com/Muraddddddddd9/ms-database)</strong> - микросервис (необходимый)
- <strong>[MS-student](https://github.com/Muraddddddddd9/ms-student)</strong> - микросервис (необходимый)
- <strong>[MS-common](https://github.com/Muraddddddddd9/ms-common)</strong> - микросервис (необходимый)
- <strong>[MS-telegram](https://github.com/Muraddddddddd9/ms-telegram)</strong> - микросервис (необходимый)
- <strong>[MDiary](https://github.com/Muraddddddddd9/MDiary)</strong> - Вебсайт (необходимый)