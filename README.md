# Học Go qua TODO API

Bộ tài liệu này chia bài thực hành Gin + PostgreSQL thành năm phần ngắn. Làm xong một tệp rồi mới mở tệp kế tiếp.

1. [01 — Khởi tạo dự án và PostgreSQL](01-khoi-tao.md)
2. [02 — Model và repository](02-model-va-repository.md)
3. [03 — Service và test](03-service-va-test.md)
4. [04 — Handler và router Gin](04-handler-va-router.md)
5. [05 — Chạy API và kiểm thử](05-chay-va-kiem-thu.md)

## Lộ trình

```text
Docker Compose (PostgreSQL)
          ↑
repository ← service ← Gin handler ← curl / trình duyệt
```

Mục tiêu là hiểu từng lớp, không chỉ sao chép mã. Mỗi phần đều có lệnh kiểm tra riêng.

- Go chạy trực tiếp trong Codespaces để bạn chỉnh sửa và chạy lại nhanh.
- PostgreSQL chạy trong Docker Compose.
- API dùng Gin, GORM và cổng 8080.
- JWT để sau; hiện tại mọi TODO dùng chung một danh sách.

Tài liệu đầy đủ ban đầu vẫn ở [đây](../huong-dan-todo-gin-postgresql.md).

