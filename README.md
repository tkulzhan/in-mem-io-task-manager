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

## Пример использования:
1. Создать задачу
Метод: POST /tasks

Описание: Создаёт новую задачу и сразу запускает её выполнение.

Пример запроса:

http
POST /tasks
Content-Type: application/json

```json
{
  "type": "default",
  "data": {
    "title": "Сделать отчёт",
    "description": "Подготовить отчёт по продажам"
  }
}
```
Параметры:

type — тип задачи (не обязательный, по умолчнию "default") ("default" — поддерживается сейчас).

data:

title — название задачи.

description (опционально) — описание.

Пример успешного ответа:

```json
{
    "id": "0197c592-2c7d-7fc8-98da-d8e0184ca914",
    "title": "1st task",
    "description": "task description",
    "status": "pending",
    "created_at": "2025-07-01T15:39:40+05:00"
}
```

2. Получить задачу по ID
Метод: GET /tasks/{id}

Описание: Получить данные о задаче по ID.

Пример запроса:

http
GET /tasks/your-task-id
Пример успешного ответа:

```json
{
  "id": "your-task-id",
  "title": "Сделать отчёт",
  "description": "Подготовить отчёт по продажам",
  "status": "running", // pending, running, finished
  "created_at": "2025-07-01T12:34:56Z",
  "started_at": "2025-07-01T12:35:00Z",
  "finished_at": "2025-07-01T12:39:56Z", // if finished execution
  "processing_time": "4.5m" // if finished execution
}
```

3. Удалить задачу
Если у тебя есть endpoint удаления (например DELETE /tasks/{id} — зависит от твоего кода), то можно вызывать:

http
DELETE /tasks/your-task-id
Это остановит выполнение, если задача ещё выполняется.

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

### Для создания новых типов задач:
1. Создай структуру удовлетварающую интерфейсу:
```go
type Task interface {
	Execute(ctx context.Context, l *logger.Logger) error
	GetID() string
	MarshalJSON() ([]byte, error)
}
```
2. При наличии новых вводных создай функцию по типу:
```go
func parseDefaultTaskInput(data any) (entity.DefaultTaskInput, error) {
	var input entity.DefaultTaskInput

	raw, err := json.Marshal(data)
	if err != nil {
		return input, errors.NewBadRequestError("invalid data format")
	}

	if err := json.Unmarshal(raw, &input); err != nil {
		return input, errors.NewBadRequestError("invalid fields in task data")
	}

	if input.Title == "" {
		return input, errors.NewBadRequestError("title is required")
	}

	return input, nil
}
```

3. Добавь новый тип здесь вместе с париснгом req.Data:
```go
func (s *TaskManagerService) ToTaskObject(ctx context.Context, req generated.CreateTaskRequest) (Task, error) {
	taskType := "default"
	if req.Type != nil && *req.Type != "" {
		taskType = *req.Type
	}

	switch taskType {
	case "default":
		input, err := parseDefaultTaskInput(req.Data)
		if err != nil {
			return nil, err
		}
		return entity.NewDefaultTask(input.Title, deref(input.Description)), nil
    /// Пример:
    	case "custom":
		input, err := parseCustomTaskInput(req.Data)
		if err != nil {
			return nil, err
		}
		return entity.NewCustomTask(input.Title, deref(input.Description)), nil

	default:
		return nil, errors.NewBadRequestError(fmt.Sprintf("unsupported task type: %s", taskType))
	}
}
```