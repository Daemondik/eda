# Golang набор на все случаи жизни

## Описание

Проект использует **Go (Gin)** на **Air**, **Postgres**, **Redis**, **Websocket**, **RabbitMQ**, **JWT** и **Oauth2** для создания мощного и безопасного приложения.

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
После авторизации вы получите cookie `access_token`, который позволит вам обращаться к `/api/admin/`.

Для регистрации через номер телефона отправьте запрос на:

```bash
POST localhost:8080/api/register
```
с таким телом запроса:

```bash
{
    "email": "your_email@example.com",
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
    "phone": "your_phone_number",
    "code": your_code
}
```

Автор:
Батов Григорий - [Daemondik](https://github.com/Daemondik)