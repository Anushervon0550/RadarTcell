# ❓ FAQ: Часто задаваемые вопросы по RadarTcell

## 🎯 ОБЩИЕ ВОПРОСЫ

### Q: Что такое RadarTcell?
**A:** Это REST API для управления каталогом технологий. Позволяет просматривать, фильтровать и администрировать технологии, тренды, SDG, теги и организации.

### Q: На каком языке написан проект?
**A:** Go (Golang) — современный язык от Google для backend-разработки.

### Q: Какая база данных используется?
**A:** PostgreSQL — мощная реляционная БД с поддержкой расширений и full-text search.

### Q: Нужно ли знать Go, чтобы понять проект?
**A:** Базовое знание Go поможет, но эта документация объясняет все с нуля. Начните с чтения `ПОЛНОЕ_ОБЪЯСНЕНИЕ_ПРОЕКТА.md`.

---

## 📁 ВОПРОСЫ О СТРУКТУРЕ

### Q: Зачем так много папок в `internal/`?
**A:** Это Clean Architecture — разделение по слоям ответственности:
- `domain` — модели данных (что хранится)
- `ports` — интерфейсы (контракты)
- `service` — бизнес-логика (правила)
- `repository` — работа с БД (SQL)
- `httpapi` — HTTP обработчики (API)

### Q: Почему код в папке `internal`, а не `src`?
**A:** В Go принято класть внутренний код в `internal/`. Это делает пакеты приватными — их нельзя импортировать из других проектов.

### Q: Что такое `cmd/api/`?
**A:** Это точка входа (entry point) приложения. Здесь находится `main.go` — файл, который запускает программу.

### Q: Зачем нужна папка `ports`?
**A:** Там хранятся интерфейсы — контракты между слоями. Это позволяет менять реализацию (например, заменить PostgreSQL на MySQL) без изменения остального кода.

---

## 🗄️ ВОПРОСЫ О БАЗЕ ДАННЫХ

### Q: Что такое миграции?
**A:** Это версионированные скрипты изменения структуры БД. Вместо ручного изменения таблиц, мы создаем файлы типа `000001_create_users.up.sql`.

### Q: Как применить миграции?
**A:**
```powershell
migrate -path migrations -database $env:DATABASE_URL up
```

### Q: Как откатить миграцию?
**A:**
```powershell
migrate -path migrations -database $env:DATABASE_URL down 1
```

### Q: Что такое `.up.sql` и `.down.sql`?
**A:**
- `.up.sql` — применить изменение (например, `CREATE TABLE`)
- `.down.sql` — откатить изменение (например, `DROP TABLE`)

### Q: Зачем нужны индексы в БД?
**A:** Ускоряют поиск и фильтрацию. Например, индекс на `name` ускоряет запросы с `WHERE name = '...'`.

### Q: Что такое TRL?
**A:** Technology Readiness Level — уровень готовности технологии (от 1 до 9):
- 1-3: Исследования
- 4-6: Разработка
- 7-9: Готово к производству

### Q: Что такое SDG?
**A:** Sustainable Development Goals — 17 целей устойчивого развития ООН (борьба с бедностью, чистая энергия, и т.д.).

---

## 🔄 ВОПРОСЫ О РАБОТЕ КОДА

### Q: Как данные идут от клиента до БД?
**A:**
```
HTTP Request → Router → Handler → Service → Repository → PostgreSQL
```
И обратно:
```
PostgreSQL → Repository → Service → Handler → HTTP Response
```

### Q: В чем разница между Service и Repository?
**A:**
- **Service** — бизнес-логика (валидация, правила, оркестрация)
- **Repository** — работа с БД (SQL запросы, транзакции)

### Q: Зачем нужны интерфейсы (ports)?
**A:** Чтобы слои не зависели от конкретных реализаций. Например, Service зависит от интерфейса `TechnologyRepository`, а не от `PostgresTechnologyRepository`. Так мы можем легко заменить БД или написать mock для тестов.

### Q: Что делает `main.go`?
**A:**
1. Загружает переменные окружения (`.env`)
2. Подключается к БД
3. Создает все компоненты (роутер, сервисы, репозитории)
4. Запускает HTTP сервер
5. Ждет сигнал завершения и корректно останавливается

### Q: Что такое `app.BuildRouter()`?
**A:** Фабрика (Factory), которая создает все компоненты приложения:
- Репозитории
- Сервисы
- Обработчики
- Роутер

Это называется "композиционный корень" (Composition Root).

### Q: Зачем нужен `context.Context`?
**A:** Для управления временем жизни операций:
- Таймауты (операция не должна длиться больше 5 секунд)
- Отмена (можно отменить запрос)
- Передача метаданных (например, request ID)

---

## 🔐 ВОПРОСЫ ОБ АВТОРИЗАЦИИ

### Q: Как работает авторизация?
**A:**
1. Админ отправляет логин/пароль на `/admin/login`
2. Сервер проверяет и генерирует JWT токен
3. Клиент сохраняет токен
4. При следующих запросах клиент отправляет `Authorization: Bearer <token>`
5. Middleware проверяет токен

### Q: Что такое JWT?
**A:** JSON Web Token — токен авторизации. Это зашифрованная строка с информацией о пользователе и сроком действия.

### Q: Где хранится пароль админа?
**A:** В переменных окружения (`.env` файл):
```
ADMIN_USER=admin
ADMIN_PASSWORD=admin123
```

### Q: Безопасно ли хранить пароль в .env?
**A:** На продакшене нужно использовать:
- Vault (HashiCorp Vault, AWS Secrets Manager)
- Переменные окружения на сервере (не в файле)
- Хеширование паролей (bcrypt)

### Q: Как получить JWT токен через PowerShell?
**A:**
```powershell
$body = @{username="admin"; password="admin123"} | ConvertTo-Json
$login = Invoke-RestMethod -Method POST -Uri "http://localhost:8080/admin/login" -Body $body -ContentType "application/json"
$token = $login.token
```

### Q: Как использовать токен в запросах?
**A:**
```powershell
$headers = @{Authorization="Bearer $token"}
Invoke-RestMethod -Method POST -Uri "http://localhost:8080/admin/technologies" -Headers $headers -Body $techBody -ContentType "application/json"
```

---

## 📡 ВОПРОСЫ ОБ API

### Q: Какие есть публичные эндпоинты?
**A:**
- `GET /api/technologies` — список технологий
- `GET /api/technologies/{slug}` — одна технология
- `GET /api/trends` — список трендов
- `GET /api/sdgs` — список SDG
- `GET /api/tags` — список тегов
- `GET /api/organizations` — список организаций

### Q: Какие есть админские эндпоинты?
**A:**
- `POST /admin/technologies` — создать технологию
- `PUT /admin/technologies/{id}` — обновить
- `DELETE /admin/technologies/{id}` — удалить
- Аналогично для trends, sdgs, tags, organizations

### Q: Как работает пагинация?
**A:**
```
GET /api/technologies?page=2&limit=20
```
Возвращает:
```json
{
  "data": [...],
  "total": 150,
  "page": 2,
  "limit": 20
}
```

### Q: Как фильтровать технологии?
**A:**
```
GET /api/technologies?search=AI&trend_id=...&readiness_level_min=5&page=1
```

### Q: Как получить технологии конкретного тренда?
**A:**
```
GET /api/trends/ai/technologies
```
или
```
GET /api/technologies?trend_id=trend-1
```

### Q: Где посмотреть документацию API?
**A:** Swagger UI доступен по адресу:
```
http://localhost:8080/swagger/
```

### Q: Что возвращается при ошибке?
**A:**
```json
{
  "error": "technology not found"
}
```
С соответствующим HTTP статусом (400, 404, 500).

---

## 🚀 ВОПРОСЫ О ЗАПУСКЕ

### Q: Как запустить проект локально?
**A:**
1. Создать `.env` файл
2. Запустить PostgreSQL (`docker-compose up -d db`)
3. Применить миграции (`migrate ... up`)
4. Запустить приложение (`go run cmd/api/main.go`)

### Q: Что должно быть в `.env` файле?
**A:**
```
DATABASE_URL=postgres://user:pass@localhost:5432/radardb
APP_PORT=8080
ADMIN_USER=admin
ADMIN_PASSWORD=admin123
JWT_SECRET=my-secret-key-change-in-production
```

### Q: Как запустить только БД?
**A:**
```powershell
docker-compose up -d db
```

### Q: Как запустить всё (БД + API) в Docker?
**A:**
```powershell
docker-compose up -d
```

### Q: Порт 8080 занят, как изменить?
**A:** В `.env` файле:
```
APP_PORT=3000
```

### Q: Как проверить, что API работает?
**A:**
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/healthz"
# Ожидаем: {"status": "ok"}
```

---

## 🧪 ВОПРОСЫ О ТЕСТИРОВАНИИ

### Q: Как запустить тесты?
**A:**
```powershell
go test ./...
```

### Q: Как посмотреть покрытие тестами?
**A:**
```powershell
go test -cover ./...
```

### Q: Как тестировать Service без реальной БД?
**A:** Использовать mock repository:
```go
mockRepo := &MockTechnologyRepo{
    FindAllFunc: func(ctx, params) {
        return []Technology{{ID: "1"}}, 1, nil
    },
}
service := NewTechnologyService(mockRepo)
```

### Q: Как тестировать Repository с реальной БД?
**A:** Создать тестовую БД и наполнить данными:
```go
func TestRepo_FindAll(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    seedTestData(db, t)
    
    repo := NewTechnologyRepo(db)
    items, _, err := repo.FindAll(ctx, params)
    
    assert.NoError(t, err)
    assert.Greater(t, len(items), 0)
}
```

### Q: Что такое E2E тесты?
**A:** End-to-End тесты — проверка всего приложения целиком. Запускается реальный сервер и БД, делаются реальные HTTP запросы. Пример: `scripts/smoke.ps1`.

---

## 🛠️ ВОПРОСЫ О РАЗРАБОТКЕ

### Q: Как добавить новое поле к технологии?
**A:**
1. Создать миграцию (`ALTER TABLE technologies ADD COLUMN ...`)
2. Добавить поле в `domain.Technology`
3. Обновить SQL запросы в repository
4. Обновить Swagger документацию

### Q: Как добавить новый эндпоинт?
**A:**
1. Добавить метод в интерфейс (`ports`)
2. Реализовать в repository (SQL)
3. Реализовать в service (логика)
4. Создать handler (HTTP)
5. Зарегистрировать маршрут в `router.go`
6. Обновить Swagger

### Q: Как добавить новый фильтр?
**A:**
1. Добавить параметр в `TechnologyListParams`
2. Добавить `WHERE` условие в SQL
3. Добавить парсинг в handler
4. Обновить Swagger

### Q: Как добавить валидацию?
**A:** В service:
```go
func (s *Service) Create(ctx, req) error {
    if req.Name == "" {
        return errors.New("name is required")
    }
    if len(req.Name) > 255 {
        return errors.New("name too long")
    }
    // ...
}
```

### Q: Как использовать транзакции?
**A:**
```go
tx, _ := r.db.Begin(ctx)
defer tx.Rollback(ctx)

// выполняем операции
tx.Exec(ctx, query1, ...)
tx.Exec(ctx, query2, ...)

tx.Commit(ctx)
```

### Q: Как добавить логирование?
**A:** Использовать библиотеку типа `zap` или `logrus`:
```go
logger.Info("processing request",
    zap.String("method", "GET"),
    zap.String("path", "/api/technologies"),
)
```

---

## 🐛 ВОПРОСЫ ОБ ОТЛАДКЕ

### Q: Ошибка "DATABASE_URL is required"
**A:** Создайте `.env` файл с переменной `DATABASE_URL`.

### Q: Ошибка "db ping error"
**A:** PostgreSQL не запущен. Запустите: `docker-compose up -d db`

### Q: Ошибка "relation 'technologies' does not exist"
**A:** Миграции не применены. Выполните: `migrate ... up`

### Q: Ошибка "401 Unauthorized"
**A:** Нужен JWT токен. Получите через `/admin/login`.

### Q: Ошибка "404 Not Found"
**A:** Проверьте URL. Возможно, технология с таким slug не существует.

### Q: Как посмотреть логи приложения?
**A:**
- Если через `go run` — логи в консоли
- Если через Docker — `docker-compose logs -f api`

### Q: Как подключиться к БД напрямую?
**A:**
```powershell
# Через Docker
docker-compose exec db psql -U radar -d radardb

# Локально
psql -U radar -d radardb -h localhost
```

### Q: Как проверить, какие миграции применены?
**A:**
```powershell
migrate -path migrations -database $env:DATABASE_URL version
```

---

## 📚 ВОПРОСЫ О КОНЦЕПЦИЯХ

### Q: Что такое Clean Architecture?
**A:** Архитектурный паттерн, разделяющий код на слои с зависимостями, направленными внутрь. Внешние слои (HTTP, БД) зависят от внутренних (бизнес-логика), а не наоборот.

### Q: Что такое Dependency Injection?
**A:** Паттерн, при котором зависимости передаются в компонент извне (через конструктор), а не создаются внутри.

### Q: Что такое Repository Pattern?
**A:** Паттерн, абстрагирующий работу с данными. Repository предоставляет простые методы типа `FindAll()`, скрывая SQL детали.

### Q: Что такое Middleware?
**A:** Функция-обертка, выполняющаяся до/после основного обработчика. Используется для логирования, авторизации, обработки ошибок.

### Q: Что такое slug?
**A:** URL-friendly идентификатор. Например, для технологии "Нейронные сети" slug будет `neural-networks`.

### Q: Зачем нужен `context.Context`?
**A:** Для управления временем жизни операций: таймауты, отмена, передача метаданных.

---

## 🎓 ВОПРОСЫ ОБ ОБУЧЕНИИ

### Q: С чего начать изучение проекта?
**A:**
1. Прочитайте `ПОЛНОЕ_ОБЪЯСНЕНИЕ_ПРОЕКТА.md`
2. Посмотрите `АРХИТЕКТУРА_ВИЗУАЛЬНО.md`
3. Изучите `ШПАРГАЛКА.md`
4. Попробуйте запустить проект
5. Прочитайте `ПРАКТИЧЕСКОЕ_РУКОВОДСТВО.md`
6. Попробуйте добавить простую фичу

### Q: Какие файлы читать в первую очередь?
**A:**
1. `cmd/api/main.go` — старт приложения
2. `internal/app/app.go` — композиция компонентов
3. `internal/httpapi/router.go` — маршруты API
4. `internal/domain/technology.go` — главная модель
5. `internal/service/technology_service.go` — пример сервиса
6. `internal/repository/postgres/technology_repo.go` — пример репозитория

### Q: Нужно ли знать Go для понимания проекта?
**A:** Базовое понимание Go поможет, но документация объясняет все ключевые моменты. Рекомендую параллельно изучать Go на [gobyexample.com](https://gobyexample.com/).

### Q: Где найти примеры кода?
**A:** В файле `ПРАКТИЧЕСКОЕ_РУКОВОДСТВО.md` есть множество примеров для разных задач.

### Q: Как проверить, что я всё понял?
**A:** Пройдите чеклист в конце `ШПАРГАЛКА.md`. Попробуйте:
- Запустить проект
- Получить JWT токен
- Создать технологию через API
- Добавить новое поле в миграцию

---

## 🔧 ТЕХНИЧЕСКИЕ ВОПРОСЫ

### Q: Почему используется chi, а не gin или echo?
**A:** chi — легковесный роутер, хорошо интегрируется со стандартной библиотекой Go. Gin и Echo тоже хороши, выбор — вопрос предпочтений.

### Q: Почему pgx, а не database/sql?
**A:** pgx — специализированный драйвер для PostgreSQL с лучшей производительностью и дополнительными фичами (prepared statements, copy, и т.д.).

### Q: Можно ли заменить PostgreSQL на MySQL?
**A:** Да, благодаря Repository Pattern. Нужно создать новую реализацию интерфейсов в `ports/repositories.go`.

### Q: Как масштабировать приложение?
**A:**
- Горизонтально: запустить несколько инстансов за Load Balancer
- Вертикально: увеличить ресурсы сервера
- Кэширование: добавить Redis
- Connection pooling: уже есть (pgxpool)

### Q: Как добавить Redis для кэширования?
**A:**
```go
import "github.com/go-redis/redis/v8"

type CatalogService struct {
    repo  ports.CatalogRepository
    redis *redis.Client
}

func (s *CatalogService) ListTrends(ctx) {
    // Пытаемся из кэша
    val, err := s.redis.Get(ctx, "trends").Result()
    if err == nil {
        var trends []domain.Trend
        json.Unmarshal([]byte(val), &trends)
        return trends, nil
    }
    
    // Если нет, идем в БД
    trends, _ := s.repo.FindAllTrends(ctx)
    
    // Сохраняем в кэш
    data, _ := json.Marshal(trends)
    s.redis.Set(ctx, "trends", data, 5*time.Minute)
    
    return trends, nil
}
```

---

## 💡 ЛУЧШИЕ ПРАКТИКИ

### Q: Как правильно обрабатывать ошибки?
**A:**
1. В repository — возвращать ошибку как есть
2. В service — оборачивать с контекстом: `fmt.Errorf("find technology: %w", err)`
3. В handler — преобразовывать в HTTP статус

### Q: Как логировать?
**A:**
- Использовать структурированное логирование (zap, logrus)
- Логировать уровни: DEBUG, INFO, WARN, ERROR
- Не логировать чувствительные данные (пароли, токены)

### Q: Как именовать переменные?
**A:**
- Go-стиль: короткие имена для локальных переменных (`i`, `err`, `ctx`)
- Экспортируемые (public) — с большой буквы
- Неэкспортируемые (private) — с маленькой

### Q: Нужны ли комментарии к коду?
**A:**
- Экспортируемые функции — обязательны (для godoc)
- Внутренняя логика — только если не очевидна
- Хороший код должен быть самодокументируемым

---

## 🎯 ИТОГОВЫЕ ВОПРОСЫ

### Q: Можно ли использовать этот проект как шаблон?
**A:** Да! Это отличная база для REST API на Go с Clean Architecture.

### Q: Где задать вопрос, если что-то не понятно?
**A:** Изучите документацию в папке проекта:
- `ПОЛНОЕ_ОБЪЯСНЕНИЕ_ПРОЕКТА.md` — общий обзор
- `АРХИТЕКТУРА_ВИЗУАЛЬНО.md` — диаграммы
- `ПРАКТИЧЕСКОЕ_РУКОВОДСТВО.md` — примеры
- `ШПАРГАЛКА.md` — краткая справка
- `FAQ.md` — этот файл

### Q: Как внести свой вклад в проект?
**A:**
1. Форкнуть репозиторий
2. Создать feature branch
3. Внести изменения
4. Написать тесты
5. Создать Pull Request

### Q: Что почитать для углубления знаний?
**A:**
- **Go:** [Effective Go](https://go.dev/doc/effective_go)
- **Архитектура:** [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- **PostgreSQL:** [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- **REST API:** [REST API Best Practices](https://restfulapi.net/)

---

**🎉 Если вашего вопроса нет в этом FAQ, изучите остальную документацию проекта!**

