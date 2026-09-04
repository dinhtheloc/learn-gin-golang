> [← 04 — Handler và router Gin](04-handler-va-router.md) · [Hoàn tất → Mục lục](README.md)

# Ghép các thành phần và chạy API

**Tạo tệp:** cmd/api/main.go.

- [ ] Tạo cmd/api/main.go:

~~~go
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"example.com/todo-api/internal/database"
	"example.com/todo-api/internal/server"
	"example.com/todo-api/internal/todo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found; using process environment")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := database.Open(dsn)
	if err != nil {
		log.Fatalf("database startup failed: %v", err)
	}

	if err := db.AutoMigrate(&todo.Todo{}); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	repository := todo.NewGormRepository(db)
	service := todo.NewService(repository)
	handler := todo.NewHandler(service)
	router := server.NewRouter(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("server listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
~~~

main là composition root: nó tạo phụ thuộc cụ thể đúng một lần. Service vẫn không phụ thuộc vào Gin hay GORM.

- [ ] **Định dạng, test và chạy**

~~~bash
gofmt -w cmd/api/main.go
go test ./...
go run ./cmd/api
~~~

Kết quả mong đợi: server listening on :8080. Giữ terminal này chạy và mở terminal thứ hai để thử API.

- [ ] **Commit**

~~~bash
git add cmd/api/main.go
git commit -m "feat: run TODO API"
~~~

+## 6. Kiểm thử thủ công mọi endpoint

- [ ] **Health check**

~~~bash
curl -i http://localhost:8080/health
~~~

Mong đợi: 200 OK và body {"status":"ok"}.

- [ ] **Tạo TODO**

~~~bash
curl -i -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"Learn Gin","description":"Build my first CRUD API"}'
~~~

Mong đợi: 201 Created. Các lệnh sau dùng ID 1 mà response vừa trả về.

- [ ] **Liệt kê và lấy chi tiết**

~~~bash
curl -i http://localhost:8080/todos
curl -i http://localhost:8080/todos/1
~~~

Cả hai trả 200 OK.

- [ ] **Cập nhật**

~~~bash
curl -i -X PATCH http://localhost:8080/todos/1 \
  -H 'Content-Type: application/json' \
  -d '{"completed":true}'
~~~

Mong đợi: 200 OK và completed là true.

- [ ] **Kiểm tra lỗi được xử lý có chủ đích**

~~~bash
curl -i -X POST http://localhost:8080/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"   "}'

curl -i http://localhost:8080/todos/999
curl -i http://localhost:8080/todos/not-a-number
~~~

Các status mong đợi lần lượt là 400, 404 và 400.

- [ ] **Xóa**

~~~bash
curl -i -X DELETE http://localhost:8080/todos/1
curl -i http://localhost:8080/todos/1
~~~

Các status mong đợi lần lượt là 204 và 404.

- [ ] **Xác minh cuối cùng**

~~~bash
go test ./...
docker compose ps
~~~

Go test phải thành công và db phải healthy.

## Ghi chú Codespaces

- Trong tab Ports, chuyển tiếp cổng 8080 nếu muốn truy cập từ trình duyệt bên ngoài Codespaces; để private khi đang học.
- Dừng database nhưng giữ dữ liệu: docker compose down.
- Xóa toàn bộ database để bắt đầu lại: docker compose down -v. Lệnh này xóa vĩnh viễn volume PostgreSQL.
- Bài học tiếp theo hợp lý: thêm migrations tường minh, rồi users và JWT. Khi thêm JWT, tạo bảng users và trường user_id trên TODO trước khi viết login route.

## Checklist hoàn tất

- [ ] docker compose ps hiển thị db healthy.
- [ ] go test ./... thành công.
- [ ] POST, GET, PATCH, DELETE hoạt động với PostgreSQL.
- [ ] .env bị Git bỏ qua; .env.example được commit.
- [ ] Bạn giải thích được vì sao Completed là *bool và vì sao Service nhận Repository interface.

---

Bạn đã hoàn tất toàn bộ CRUD API.
