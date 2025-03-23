# Gravitum user service

Сервис ответственный за обработку запросов к user REST API. 

Предоставляет доступ к сущности пользователя посредством http c интерфейсом и сообщениями описанными в [api](/api/openapi/user/v1/api.yaml).

### Запуск
- `make service`
- `make dev-compose-up`
- [url rest api](http://localhost:3000/api/web/v1)

### Остановка
- `make dev-compose-down`