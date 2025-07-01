# In-Memory IO Task Manager
Простой менеджер задач с REST API на Go, выполняющий задачи в памяти с поддержкой ограничения количества одновременно выполяемых задач и отмены задач.

## Подготовка окружения
1. Убедись, что установлен Go
2. Склонируй репозиторий
3. Создай файл .env в папке config, пример .env:
```env
REST_ADDRESS=0.0.0.0:8080
LOG_LEVEL=trace
HTTP_EXPOSE_API_SPECIFICATION=true
HTTP_READ_TIME_OUT=30
HTTP_WRITE_TIME_OUT=30
MAX_CONCURRENT_TASKS=1000
```
4. Установи зависимости:
```sh
go mod tidy
```

## Запуск приложения:
```sh
go run cmd/main.go
```

## Для разработчиков:
### Для запуска unit-тестов:
```sh
go test ./...
```

### Для генерации серверного кода из спецификаций ./api/openapi3/api.yaml:
```sh
go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=internal/infrastructure/http/server/config.yaml ./api/openapi3/api.yaml
```

### Для генерации моков под тесты:
```sh
go generate ./...
```