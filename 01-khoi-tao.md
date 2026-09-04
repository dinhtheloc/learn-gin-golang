> [Mục lục](README.md) · [02 — Model và repository →](02-model-va-repository.md)

# Khởi tạo dự án và database

**Tạo các tệp:** .gitignore, .env.example, docker-compose.yml.

- [ ] **Kiểm tra môi trường Codespaces**

~~~bash
go version
docker --version
docker compose version
~~~

Cả ba lệnh đều cần in ra phiên bản.

- [ ] **Tạo module và thư mục**

~~~bash
go mod init example.com/todo-api
mkdir -p cmd/api internal/database internal/server internal/todo
~~~

example.com/todo-api là module cục bộ hợp lệ. Khi muốn, bạn có thể đổi thành địa chỉ repository GitHub của mình.

- [ ] **Tạo .gitignore**

~~~gitignore
.env
tmp/
bin/
coverage.out
~~~

- [ ] **Tạo .env.example, rồi sao chép thành .env bằng trình quản lý tệp của VS Code**

~~~dotenv
DATABASE_URL=postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable
PORT=8080
~~~

- [ ] **Tạo docker-compose.yml**

~~~yaml
services:
  db:
    image: postgres:16-alpine
    container_name: todo-postgres
    environment:
      POSTGRES_USER: todo_user
      POSTGRES_PASSWORD: todo_password
      POSTGRES_DB: todo_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U todo_user -d todo_db"]
      interval: 5s
      timeout: 5s
      retries: 10

volumes:
  postgres_data:
~~~

- [ ] **Khởi động PostgreSQL**

~~~bash
docker compose up -d
docker compose ps
~~~

Kết quả mong đợi: dịch vụ db trở thành healthy. Volume postgres_data giữ dữ liệu khi container được tạo lại.

- [ ] **Cài các thư viện Go**

~~~bash
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go get gorm.io/driver/postgres
go get gorm.io/gorm
go mod tidy
~~~

- [ ] **Lưu checkpoint**

~~~bash
git add .gitignore .env.example docker-compose.yml go.mod go.sum
git commit -m "chore: set up Go and PostgreSQL"
~~~

---

Tiếp theo: [02 — Model và repository →](02-model-va-repository.md)
