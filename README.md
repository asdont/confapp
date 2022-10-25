## ConfApp - конфигурации сервисов

---

### Документация

Документация(методы) - http://127.0.0.1:22952/doc/index.html

---

### Компиляция

Приложение - `go build -o confapp cmd/app/main.go`

Swagger-документация - `swag init -g cmd/app/main.go`

---

### Postgres(таблицы, доступ)

`scripts/create_tables.sql`

авторизация - `config/conf.toml`

---

### Docker(база снаружи)

`docker build . -t confapp:v0.1`

`docker run --network=host confapp:v0.1`

### Docker-compose(база внутри)

`docker-compose build`

`docker-compose up`


