# Plane CLI: CLI-утилита для AI-агентов на Go

**Дата:** 2026-05-16
**Статус:** На ревью
**Автор:** OpenCode Agent

## 1. Обзор

### 1.1 Цель

Создать CLI-утилиту `plane` на Go для взаимодействия AI-агентов (PM, Tech Lead, Developer, Code Reviewer, QA) с Plane API. Утилита заменяет прямые HTTP-запросы к Plane из n8n-воркфлоу, предоставляя единый интерфейс для работы с задачами, циклами, модулями, лейблами, стейтами, участниками и страницами.

### 1.2 Ключевые решения

| Параметр | Выбор |
|---|---|
| Язык | Go (1.22+) |
| CLI-фреймворк | Cobra |
| Аутентификация | API-токен (`X-API-Key`) |
| Формат вывода | JSON |
| Retry-логика | Экспоненциальный backoff + Circuit Breaker |
| HTTP-клиент | go-retryablehttp |
| Конфигурация | `~/.config/plane/config.json` |
| Алиасы команд | `wi` ⇔ `work-item` ⇔ `work-items` |

### 1.3 Границы системы

**Входит в scope:**
- CLI-интерфейс с иерархическими командами (Cobra)
- Аутентификация через API-токен
- Контекст (workspace + project + token) с сохранением в конфиг
- CRUD-операции над Work Items, комментариями
- Операции чтения/создания над Cycles, Modules, Labels, States
- Операции чтения над Members, Pages
- Retry с экспоненциальным backoff и circuit breaker
- Сборка статических бинарников под linux/darwin (amd64 + arm64)

**Не входит в scope:**
- Интерактивный режим (TUI)
- OAuth-аутентификация (только API-токен)
- WebSocket/real-time подписки
- Админские операции (instance admin API)
- Webhook-управление
- Автодополнение shell (отложено)
- Интеграция с Git API (используется отдельный инструмент)

---

## 2. Архитектура

### 2.1 Структура проекта

```
plane/
├── cmd/plane/
│   └── main.go                # точка входа
├── internal/
│   ├── cli/                   # Cobra-команды
│   │   ├── root.go            # root command + глобальные флаги (--token, --api-url, --max-retries, --timeout)
│   │   ├── context.go         # plane context set/show/unset
│   │   ├── workitem.go        # plane wi {list,get,create,update,delete,comment}
│   │   ├── cycle.go           # plane cycle {list,get,create}
│   │   ├── module.go          # plane module {list,get}
│   │   ├── state.go           # plane state {list,get}
│   │   ├── label.go           # plane label {list,create}
│   │   ├── member.go          # plane member list
│   │   └── page.go            # plane page {list,get,create}
│   ├── api/
│   │   ├── client.go          # HTTP-клиент: buildRequest, execute, retry/cb-логика
│   │   ├── transport.go       # Обёртка над hashicorp/go-retryablehttp
│   │   └── endpoints.go       # Построение URL (конкатенация api_url + path)
│   ├── config/
│   │   └── context.go         # read/write ~/.config/plane/config.json
│   ├── models/
│   │   └── models.go          # Domain-типы: WorkItem, Cycle, Module, State, Label, Member, Page, Comment
│   └── output/
│       └── json.go            # JSON-сериализация (indent + color)
├── go.mod
├── go.sum
├── Makefile                    # build для linux/darwin amd64/arm64
└── README.md
```

### 2.2 Зависимости

| Библиотека | Назначение |
|---|---|
| `github.com/spf13/cobra` | CLI-фреймворк |
| `github.com/hashicorp/go-retryablehttp` | HTTP-клиент с retry |
| Стандартная библиотека Go | `encoding/json`, `net/http`, `os`, `flag` (через cobra/pkg) |

Никаких других внешних зависимостей — бинарник остаётся минимальным.

### 2.3 Схема взаимодействия

```
┌──────────┐     ┌─────────┐     ┌──────────────┐
│ n8n      │────▶│ plane   │────▶│ Plane API     │
│ workflow │     │ CLI     │     │ (Django REST) │
└──────────┘     └─────────┘     └──────────────┘
                      │
                      ▼
              ┌──────────────┐
              │ ~/.config/   │
              │ plane/       │
              │ config.json  │
              └──────────────┘
```

Агент в n8n вызывает `plane` через Execute Command node. Результат — JSON в stdout. Ошибки — в stderr + код возврата.

---

## 3. Аутентификация

### 3.1 Приоритет источников токена

```
1. Флаг --token               # plane --token "plane_key_xxx" wi list
2. Переменная PLANE_TOKEN     # PLANE_TOKEN=plane_key_xxx plane wi list
3. Контекст (config.json)     # plane context set --token "plane_key_xxx"; plane wi list
```

Значение из источника с более высоким приоритетом переопределяет нижестоящие.

### 3.2 Конфигурационный файл

Расположение: `~/.config/plane/config.json`

```json
{
  "workspace": "myorg",
  "project": "abc123-uuid",
  "token": "plane_key_abc123",
  "api_url": "https://plane.example.com"
}
```

- По умолчанию `api_url` = `https://api.plane.so`
- Файл создаётся командой `plane context set --workspace <slug> --project <uuid> --token <key>`
- `plane context set --api-url <url>` меняет URL
- `plane context show` выводит текущий контекст (без токена)
- `plane context unset` удаляет конфиг

### 3.3 Переменные окружения

| Переменная | Назначение |
|---|---|
| `PLANE_TOKEN` | API-токен |
| `PLANE_API_URL` | Базовый URL API |
| `PLANE_WORKSPACE` | Slug workspace по умолчанию |
| `PLANE_PROJECT` | UUID проекта по умолчанию |

---

## 4. Команды

### 4.1 Глобальные флаги (root)

```
--token         string   API-токен (переопределяет PLANE_TOKEN и контекст)
--api-url       string   Базовый URL API (переопределяет PLANE_API_URL и контекст)
--workspace     string   Slug workspace (переопределяет контекст)
--project       string   UUID проекта (переопределяет контекст)
--max-retries   int      Максимальное число повторных попыток (default: 4)
--timeout       duration Общий таймаут на все retry-попытки (default: 30s)
--no-color               Отключить ANSI-цвета в выводе
```

### 4.2 Context

```
plane context set    --workspace <slug> --project <uuid> --token <key> [--api-url <url>]
plane context show
plane context unset
```

### 4.3 Work Item

```
# CRUD
plane wi list        [--state <name>] [--assignee <email>] [--label <name>] [--cycle <id>] [--search <query>]
plane wi get         <id>
plane wi create      --title <text> [--description <text>] [--state <name>] [--priority <urgent|high|medium|low|none>] [--assignee <email>] [--labels <a,b,c>]
plane wi update      <id> [--title <text>] [--state <name>] [--description <text>] [--assignee <email>] [--priority <urgent|high|medium|low|none>] [--labels <a,b,c>]
plane wi delete      <id>

# Комментарии
plane wi comment add <id> --content <text>
plane wi comment list <id>
```

Алиасы: `wi` / `work-item` / `work-items` — все три взаимозаменяемы.

### 4.4 Cycle

```
plane cycle list      [--sort-by <name|start_date|end_date>]
plane cycle get       <id>
plane cycle create    --name <text> --project <uuid> [--start-date <date>] [--end-date <date>] [--description <text>]
```

### 4.5 Module

```
plane module list
plane module get      <id>
```

### 4.6 State

```
plane state list
plane state get       <id>
```

### 4.7 Label

```
plane label list
plane label create    --name <text> --color <hex> [--parent <id>]
```

### 4.8 Member

```
plane member list
```

### 4.9 Page

```
plane page list
plane page get        <id>
plane page create     --name <text> [--description <text>] [--content <text>]
```

---

## 5. Вывод (Output)

### 5.1 JSON-формат

Все команды выводят JSON в stdout:

```json
{
  "items": [
    {
      "id": "abc-123",
      "name": "Implement login",
      "description": "...",
      "state": "In Progress",
      "priority": "high",
      "assignees": ["user@example.com"],
      "labels": ["frontend", "auth"],
      "cycle": null,
      "module": null,
      "created_at": "2026-05-15T10:00:00Z",
      "updated_at": "2026-05-16T12:00:00Z",
      "created_by": "admin@example.com",
      "url": "https://plane.example.com/myorg/projects/abc/issues/abc-123"
    }
  ],
  "total": 1
}
```

Ответ-обёртка всегда содержит `items` + `total`. Одиночные сущности (get) также упаковываются в `items` массив из одного элемента.

### 5.2 Ошибки

```json
{
  "error": "not_found",
  "message": "Work item abc-456 not found"
}
```

Ошибки всегда выводятся как JSON в stdout. Человекочитаемое сообщение дублируется в stderr.

### 5.3 Коды возврата

| Код | Значение |
|---|---|
| 0 | Успех |
| 1 | Ошибка API (4xx, 5xx) |
| 2 | Circuit breaker open |
| 3 | Ошибка конфигурации (нет токена, нет workspace) |
| 4 | Ошибка сети (не удалось соединиться после всех retry) |

---

## 6. Retry и Circuit Breaker

### 6.1 Параметры по умолчанию

| Параметр | Значение |
|---|---|
| MaxRetries | 4 (5 попыток всего) |
| Backoff | Экспоненциальный: 1s, 2s, 4s, 8s |
| MaxTotalTimeout | 30s (на все попытки + запрос) |
| CircuitBreakerThreshold | 5 последовательных ошибок |
| CircuitBreakerReset | Сброс счётчика при успешном запросе |

### 6.2 Retry-стратегия по HTTP-кодам

| HTTP статус | Retry? | Поведение |
|---|---|---|
| 200-399 | Нет | Успех |
| 429 (Rate Limit) | Да | Ожидание Retry-After заголовка |
| 500, 502, 503, 504 | Да | Retry с backoff |
| Connection timeout/refused | Да | Retry с backoff |
| 400, 401, 403, 404, 422 | Нет | Ошибка — возвращается сразу |

### 6.3 Circuit Breaker

Счётчик последовательных ошибок инкрементируется при любом неуспешном запросе (включая исчерпанные retry). При достижении порога (5) CLI выводит:

```json
{
  "error": "circuit_breaker_open",
  "message": "Circuit breaker open after 5 consecutive failures. Plane API may be unavailable."
}
```

И завершается с кодом 2. Успешный запрос сбрасывает счётчик в 0.

---

## 7. Модели данных (Domain Types)

```go
type WorkItem struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    State       string   `json:"state,omitempty"`
    Priority    string   `json:"priority,omitempty"`
    Assignees   []string `json:"assignees,omitempty"`
    Labels      []string `json:"labels,omitempty"`
    Cycle       string   `json:"cycle,omitempty"`
    Module      string   `json:"module,omitempty"`
    CreatedAt   string   `json:"created_at"`
    UpdatedAt   string   `json:"updated_at"`
    CreatedBy   string   `json:"created_by"`
    URL         string   `json:"url"`
}

type Cycle struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    StartDate   string `json:"start_date,omitempty"`
    EndDate     string `json:"end_date,omitempty"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type Module struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}

type State struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
    Group string `json:"group"`
}

type Label struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Color string `json:"color"`
}

type Member struct {
    ID          string `json:"id"`
    Email       string `json:"email"`
    DisplayName string `json:"display_name"`
    Role        string `json:"role"`
}

type Page struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    Content     string `json:"content,omitempty"`
    CreatedAt   string `json:"created_at"`
    UpdatedAt   string `json:"updated_at"`
}

type Comment struct {
    ID        string `json:"id"`
    Content   string `json:"content"`
    CreatedBy string `json:"created_by"`
    CreatedAt string `json:"created_at"`
}
```

---

## 8. API-клиент

### 8.1 Базовый URL

Все запросы идут на `{api_url}/api/v1/`. Базовый URL определяется из:
1. Флаг `--api-url`
2. Переменная `PLANE_API_URL`
3. Конфиг `config.json`
4. Значение по умолчанию `https://api.plane.so`

### 8.2 Заголовки

Каждый запрос включает:
```
X-API-Key: <token>
Content-Type: application/json
Accept: application/json
User-Agent: plane-cli/1.0.0
```

### 8.3 Эндпоинты

| Команда | Метод | Путь |
|---|---|---|
| `wi list` | GET | `/workspaces/{ws}/projects/{pr}/work-items/` |
| `wi get` | GET | `/workspaces/{ws}/projects/{pr}/work-items/{id}/` |
| `wi create` | POST | `/workspaces/{ws}/projects/{pr}/work-items/` |
| `wi update` | PATCH | `/workspaces/{ws}/projects/{pr}/work-items/{id}/` |
| `wi delete` | DELETE | `/workspaces/{ws}/projects/{pr}/work-items/{id}/` |
| `wi comment list` | GET | `/workspaces/{ws}/projects/{pr}/work-items/{id}/comments/` |
| `wi comment add` | POST | `/workspaces/{ws}/projects/{pr}/work-items/{id}/comments/` |
| `cycle list` | GET | `/workspaces/{ws}/projects/{pr}/cycles/` |
| `cycle get` | GET | `/workspaces/{ws}/projects/{pr}/cycles/{id}/` |
| `cycle create` | POST | `/workspaces/{ws}/projects/{pr}/cycles/` |
| `module list` | GET | `/workspaces/{ws}/projects/{pr}/modules/` |
| `module get` | GET | `/workspaces/{ws}/projects/{pr}/modules/{id}/` |
| `state list` | GET | `/workspaces/{ws}/projects/{pr}/states/` |
| `state get` | GET | `/workspaces/{ws}/projects/{pr}/states/{id}/` |
| `label list` | GET | `/workspaces/{ws}/projects/{pr}/issue-labels/` |
| `label create` | POST | `/workspaces/{ws}/projects/{pr}/issue-labels/` |
| `member list` | GET | `/workspaces/{ws}/projects/{pr}/members/` |
| `page list` | GET | `/workspaces/{ws}/projects/{pr}/pages/` |
| `page get` | GET | `/workspaces/{ws}/projects/{pr}/pages/{id}/` |
| `page create` | POST | `/workspaces/{ws}/projects/{pr}/pages/` |

### 8.4 Параметры запроса (query params)

| Команда | Параметр | Query-ключ в API |
|---|---|---|
| `wi list --state` | Фильтр по стейту | `?state=<name>` |
| `wi list --assignee` | Фильтр по назначенному | `?assignees=<email>` |
| `wi list --label` | Фильтр по лейблу | `?labels=<name>` |
| `wi list --cycle` | Фильтр по циклу | `?cycle_id=<id>` |
| `wi list --search` | Текстовый поиск | `?search=<query>` |

Для `cycle list --sort-by`: параметр `order_by` со значениями `name`, `start_date`, `-start_date`, `end_date`, `-end_date`.

---

## 9. Обработка ошибок

| Ситуация | Код возврата | Сообщение в JSON |
|---|---|---|
| Нет токена ни во флаге, ни в env, ни в конфиге | 3 | `{ "error": "missing_token", "message": "..." }` |
| Нет workspace / project | 3 | `{ "error": "missing_context", "message": "..." }` |
| Plane вернул 404 | 1 | `{ "error": "not_found", "message": "..." }` |
| Plane вернул 400/422 | 1 | `{ "error": "validation", "message": "..." }` |
| Plane вернул 401/403 | 1 | `{ "error": "unauthorized", "message": "..." }` |
| Plane вернул 500/502/503 | 1 | `{ "error": "server_error", "message": "..." }` (после retry) |
| Таймаут соединения | 4 | `{ "error": "connection_timeout", "message": "..." }` |
| DNS/сеть недоступна | 4 | `{ "error": "network", "message": "..." }` |
| Circuit breaker open | 2 | `{ "error": "circuit_breaker_open", "message": "..." }` |

---

## 10. Makefile и сборка

```makefile
APP_NAME = plane
GO = go
LDFLAGS = -ldflags="-s -w"

.PHONY: build build-linux build-darwin build-all clean test

build:
	$(GO) build $(LDFLAGS) -o bin/$(APP_NAME) ./cmd/plane

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-amd64 ./cmd/plane

build-linux-arm64:
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-linux-arm64 ./cmd/plane

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-amd64 ./cmd/plane

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(APP_NAME)-darwin-arm64 ./cmd/plane

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64

clean:
	rm -rf bin/

test:
	$(GO) test ./...
```

---

## 11. Примеры использования (для агентов)

### 11.1 Инициализация контекста (делается один раз в начале пайплайна)

```bash
plane context set \
  --workspace "myorg" \
  --project "550e8400-e29b-41d4-a716-446655440000" \
  --token "$PLANE_API_TOKEN"
```

### 11.2 Создание задачи (PM создаёт task)

```bash
plane wi create \
  --title "Add dark mode toggle" \
  --description "Implement dark mode toggle in settings" \
  --priority "high" \
  --labels "frontend,ui"
```

### 11.3 Поиск задач (Tech Lead ищет sub-issues)

```bash
plane wi list --label "frontend" --state "Backlog"
```

### 11.4 Обновление статуса (Developer начинает работу)

```bash
plane wi update <id> --state "In Progress" --assignee "dev@example.com"
```

### 11.5 Добавление комментария (QA пишет баг-репорт)

```bash
plane wi comment add <id> --content "Found: button color not updating on toggle"
```

### 11.6 Завершение задачи (Developer закрывает)

```bash
plane wi update <id> --state "Done"
```

---

## 12. План тестирования

| Уровень | Что тестируется | Инструмент |
|---|---|---|
| Unit | Парсинг аргументов Cobra, сериализация JSON | `go test` |
| Unit | Retry-логика и circuit breaker | `go test` с `httptest.Server` |
| Unit | Контекст (read/write config.json) | `go test` с временной директорией |
| Integration | Полный цикл команд против тестового инстанса Plane | Docker Compose + `go test` |
| Smoke | Бинарник запускается, `--help` работает | Makefile target |

---

## 13. Примечания

- Название бинарника: `plane`
- Модуль Go: `github.com/makeplane/plane-cli`
- Релизы: GitHub Releases, бинарники прикрепляются к тегу
- Версионирование: SemVer, начальная версия `v0.1.0`
- AGENTS.md обновляется после реализации: добавляются примеры вызова `plane` из агентов
