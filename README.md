# HealthCheck — веб‑сервис проверки доступности URL

## старт

```bash
go run .
```

Сервер слушает `:8080`.
Состояние хранилища сохраняется в файл `storage.json` в рабочей директории.

## API

### Отправка ссылок на проверку

`POST /urls`
Тело запроса (любой из вариантов):

```json
{"url":"https://example.com"}
```

```json
{"urls":["https://a.com","https://b.com"]}
```

Ответ:

```json
{
  "id": 1,
  "statuses": {
    "https://a.com": true,
    "https://b.com": false
  }
}
```

Где true — ресурс доступен, false — недоступен

### Получение статуса по id

`GET /status?id=1`
Ответ:

```json
{
  "id": 1,
  "urls": ["https://a.com","https://b.com"],
  "status": {
    "https://a.com": true,
    "https://b.com": false
  }
}
```

### Генерация PDF‑отчета по списку id

`POST /report/pdf`
Тело запроса:

```json
{"ids":[1,2,3]}
```

Ответ: файл `application/pdf`

## Проверки доступности

Запрос `GET` с таймаутом `5s`.
Успешным считается ответ с кодом `200–399`.

## Остановка и перезапуск

Все заявки и их статусы сохраняются в `storage.json` сразу после обработки.
При старте сервис загружает состояние из `storage.json`, чтобы не терять уже созданные `id` и проверенные ссылки.
При получении сигнала остановки выполняется **graceful shutdown**: сервер перестаёт принимать новые соединения, дожидается завершения активных запросов, затем сохраняет состояние на диск.

## Примеры запросов (curl)

```bash
# Отправить один URL
curl -X POST http://localhost:8080/urls \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'

# Отправить несколько URL
curl -X POST http://localhost:8080/urls \
  -H "Content-Type: application/json" \
  -d '{"urls":["https://a.com","https://b.com"]}'

# Получить статусы по id
curl "http://localhost:8080/status?id=1"

# Получить PDF отчёт по нескольким id
curl -X POST http://localhost:8080/report/pdf \
  -H "Content-Type: application/json" \
  -d '{"ids":[1,2,3]}' --output report.pdf
```