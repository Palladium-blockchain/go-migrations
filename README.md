# go-migrations

Библиотека и простой CLI для управления SQL-миграциями в PostgreSQL.

## Возможности

- Создание шаблонов миграций (`.up.sql` / `.down.sql`).
- Применение миграций вверх с блокировкой через `pg_advisory_lock`.
- Хранение применённых версий в таблице `schema_migrations`.
- Программный интерфейс для выполнения `Up`/`Down`.

## Установка

```bash
go get github.com/Palladium-blockchain/go-migrations
```

## Формат миграций

Файлы миграций хранятся в одной директории и имеют формат:

- `<id>.up.sql`
- `<id>.down.sql`

`<id>` — произвольный идентификатор миграции. Встроенный генератор использует формат
`<timestamp>_<name>` (UTC, `YYYYMMDDHHMMSS`). При загрузке миграций порядок
определяется лексикографически по `<id>`, поэтому timestamp помогает сохранить
правильную последовательность.

> `down.sql` файл опционален для применения вверх, но обязателен для отката.

## Использование как библиотека

### Подключение драйвера и источника

```go
package main

import (
    "context"
    "database/sql"
    "os"

    migratepostgres "github.com/Palladium-blockchain/go-migrations/pkg/driver/postgres"
    "github.com/Palladium-blockchain/go-migrations/pkg/migrator"
    migratefs "github.com/Palladium-blockchain/go-migrations/pkg/source/fs"

    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
    if err != nil {
        panic(err)
    }
    defer db.Close()

    driver := migratepostgres.NewDriver(db)
    source := migratefs.NewSource(os.DirFS("migrations"))

    if err := migrator.NewMigrator(driver, source).Up(context.Background()); err != nil {
        panic(err)
    }
}
```

Если в базе уже есть записи о миграциях, которых больше нет локально, и их нужно
игнорировать при `Up`, можно явно включить это поведение:

```go
if err := migrator.NewMigrator(
    driver,
    source,
    migrator.WithAllowOrphanedMigrations(),
).Up(context.Background()); err != nil {
    panic(err)
}
```

### Откат последней миграции

```go
if err := migrator.NewMigrator(driver, source).Down(ctx); err != nil {
    if errors.Is(err, migrate.ErrNoChange) {
        // нечего откатывать
    }
    return err
}
```

## Использование CLI

### Сборка

```bash
go build -o migration ./cmd/migrations
```

### Создание миграции

```bash
MIGRATIONS_PATH="migrations" \
./migration create add_users_table
```

Команда создаст пару файлов:

- `migrations/<timestamp>_add_users_table.up.sql`
- `migrations/<timestamp>_add_users_table.down.sql`

### Применение миграций

```bash
DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
MIGRATIONS_PATH="migrations" \
./migration migrate
```

Если нужно игнорировать миграции, которые уже есть в таблице `schema_migrations`,
но отсутствуют локально, добавьте флаг:

```bash
DATABASE_URL="postgres://user:pass@localhost:5432/dbname?sslmode=disable" \
MIGRATIONS_PATH="migrations" \
./migration migrate --allow-orphaned-migrations
```

Параметры окружения:

- `DATABASE_URL` — строка подключения к PostgreSQL (обязательна).
- `MIGRATIONS_PATH` — путь к директории миграций (по умолчанию `migrations`).

Параметры CLI:

- `--allow-orphaned-migrations` — игнорировать уже применённые миграции,
  которых нет в локальном каталоге миграций. На `Down` это не влияет: для
  отката локальный `down.sql` всё равно обязателен.

## Как работает хранение миграций

При первом запуске создаётся таблица:

```sql
CREATE TABLE IF NOT EXISTS schema_migrations (
    version TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
)
```

Для сериализации запуска используется `pg_advisory_lock(424242)`, чтобы не допустить
параллельного применения миграций.
