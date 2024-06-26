# Golang набор на все случаи жизни

## Описание

Проект использует **Go (Gin)** на **Air**, **Postgres**, **Redis**, **gRPC**, **Websocket**, **RabbitMQ**, **JWT** и **Oauth2** для создания мощного и безопасного приложения.

Для установки программного обеспечения вам понадобится:

- Docker

### Установка

Чтобы запустить среду разработки, выполните следующие шаги:

```bash
docker-compose up --build
```

### Использование

Для авторизации через Google перейдите по адресу:

```bash
GET localhost:8080/api/login-gl
```
После авторизации вы получите cookie `access_token`, который позволит вам обращаться к `/api/admin/`, при условии, что ваша `role` в таблице `users` = `admin`.

Для регистрации через номер телефона отправьте запрос на:

```bash
POST localhost:8080/api/register
```
с таким телом запроса:

```bash
{
    "phone": "7XXXXXXXXXX",
    "password": "your_password"
}
```
После регистрации вы получите смс с кодом подтверждения. Отправьте этот код на:

```bash
POST localhost:8080/api/confirm-sms
```
с таким телом запроса:

```bash
{
    "phone": "7XXXXXXXXXX",
    "code": your_code
}
```

Для входа через номер телефона отправьте запрос на:

```bash
POST localhost:8080/api/login
```
с таким телом запроса:

```bash
{
    "phone": "7XXXXXXXXXX",
    "password": "your_password"
}
```

В результате вы получите `token` для авторизации. Используйте его при POST и GET запросах.
Вы можете передать в `Headers` ключ `Authorization` с параметром `Bearer <token>`

### Чат в реальном времени на Websocket & RabbitMQ

- Авторизуйтесь через смс или Google
- Перейдите в браузере по ссылке `http://localhost:8080/chat/<id пользователя, которому пишем>`

### Go Sql модуль

При сборке приложения скрипт выполняет триггеры из директории `sql/triggers` и отслеживает, изменились ли они с момента последней сборки. Если изменились, то выполняет их

### Unit-тесты

Выполняются при сборке образа. Или можно вручную, выполнив 
```bash
go test ./...
```

## GUI

### pgAdmin

- Перейдите в браузере по ссылке `http://localhost/`
- Авторизуйтесь. Емейл и пароль указаны в `.env`
- Выполните `docker ps`, найдите IMAGE с NAME `eda-db-1`, скопируйте CONTAINER ID
- Выполните `docker inspect <CONTAINER ID>`, найдите `Gatevay` и скопируйте значение
- Используйте это значение как хост для подключения к базе данных. Остальные значения можно найти в `.env`

### RabbitMQ

- Перейдите в браузере по ссылке `http://localhost:15672/`
- Имя пользователя и пароль в `.env`

### Redis

- Скачайте любой Redis GUI, например RESP.app 
- Ссылка для подключения `redis://localhost:6379`
- Пароль в `.env`

## Qodana

Чтобы использовать Qodana, вам нужно запустить `docker-compose up qodana`. Это запустит анализ кода Qodana для вашего проекта, и результаты будут сохранены в директории `qodana-results` в вашем локальном проекте.
Вы можете просмотреть HTML-отчет, открыв файл `qodana-results/report/index.html` в браузере. Этот отчет представляет собой интерактивный интерфейс, который позволяет просматривать и анализировать результаты анализа кода.

Автор:
Батов Григорий - [Daemondik](https://github.com/Daemondik)
