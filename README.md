# plane-cli

**plane-cli** — CLI-утилита на Go для взаимодействия с [Plane](https://github.com/makeplane/plane) API. Предназначена для AI-агентов (PM, Dev, TL, CR, QA), работающих в автоматизированных пайплайнах разработки.

## Установка

```bash
go install github.com/kazakovdmitriy/plane-cli/cmd/plane@latest
```

Или скачай готовый бинарник из [GitHub Releases](https://github.com/kazakovdmitriy/plane-cli/releases).

## Быстрый старт

```bash
# Настройка контекста (один раз)
plane context set \
  --workspace "myorg" \
  --project "550e8400-e29b-41d4-a716-446655440000" \
  --token "$PLANE_API_TOKEN"

# Создать задачу
plane wi create \
  --title "Add dark mode toggle" \
  --description "Implement dark mode toggle in settings" \
  --priority "high" \
  --labels "frontend"

# Найти задачи
plane wi list --state "Backlog" --label "frontend"

# Взять в работу
plane wi update <id> --state "In Progress" --assignee "dev@example.com"

# Добавить комментарий
plane wi comment add <id> --content "Button color not updating on toggle"

# Закрыть задачу
plane wi update <id> --state "Done"
```

## Аутентификация

Токен ищется в порядке приоритета:

1. Флаг `--token`
2. Переменная окружения `PLANE_TOKEN`
3. Конфиг `~/.config/plane/config.json`

## Команды

### Глобальные флаги

```
--token         string   API-токен
--api-url       string   Базовый URL API (default: https://api.plane.so)
--workspace, -w string   Slug workspace
--project, -p   string   UUID проекта
--max-retries   int      Ретраи (default: 4)
--timeout       duration Таймаут на все попытки (default: 30s)
```

### Context

```
plane context set    --workspace <slug> --project <uuid> --token <key>
plane context show
plane context unset
```

### Work Items

```
plane wi list       [--state ...] [--assignee ...] [--label ...] [--cycle ...] [--search ...]
plane wi get        <id>
plane wi create     --title "..." [--description "..."] [--state "..."] [--priority "..."] [--assignee "..."] [--labels "a,b"]
plane wi update     <id> [--title "..."] [--state "..."] [--description "..."] [--assignee "..."] [--priority "..."] [--labels "a,b"]
plane wi delete     <id>
plane wi comment list <id>
plane wi comment add <id> --content "..."
```

Алиасы: `wi`, `work-item`, `work-items`.

### Cycles

```
plane cycle list    [--sort-by name|start_date|-start_date|end_date|-end_date]
plane cycle get     <id>
plane cycle create  --name "..." [--start-date YYYY-MM-DD] [--end-date YYYY-MM-DD] [--description "..."]
```

### Modules

```
plane module list
plane module get    <id>
```

### States

```
plane state list
plane state get     <id>
```

### Labels

```
plane label list
plane label create  --name "..." --color "#ff0000"
```

### Members

```
plane member list
```

### Pages

```
plane page list
plane page get      <id>
plane page create   --name "..." [--description "..."] [--content "..."]
```

## Вывод

Все команды выводят JSON:

```json
{
  "items": [
    {
      "id": "abc-123",
      "name": "Add dark mode toggle",
      "state": "In Progress",
      "priority": "high",
      "assignees": ["dev@example.com"],
      "labels": ["frontend"],
      "url": "https://plane.example.com/myorg/projects/abc/issues/abc-123"
    }
  ],
  "total": 1
}
```

Ошибки:

```json
{
  "error": "not_found",
  "message": "Work item abc-456 not found"
}
```

## Коды возврата

| Код | Значение |
|-----|----------|
| 0 | Успех |
| 1 | Ошибка API (4xx, 5xx) |
| 2 | Circuit breaker open (нет связи с API) |
| 3 | Ошибка конфигурации (нет токена) |
| 4 | Ошибка сети |

## Retry и Circuit Breaker

При сетевых сбоях (5xx, 429, таймауты) клиент повторяет запрос с экспоненциальной задержкой (1s → 2s → 4s → 8s). После 5 последовательных ошибок включается circuit breaker — CLI завершается с кодом 2.

Настроить: `--max-retries` и `--timeout`.

## Сборка из исходников

```bash
make build              # текущая платформа
make build-all          # linux/darwin × amd64/arm64
go test ./...           # запуск тестов
```

## Лицензия

AGPL-3.0
