
## Mục tiêu

Sử dụng Redis để tối ưu hiệu suất xác thực người dùng và kiểm tra quyền truy cập:

* ✅ Lưu thông tin người dùng sau khi đăng nhập vào Redis.
* ✅ Mỗi request từ client đều kiểm tra JWT hợp lệ.
* ✅ Nếu hợp lệ → lấy thông tin user (ID, role, permissions) từ Redis thay vì DB.

---

## Quy trình xử lý

### 1. Đăng nhập

**Endpoint:** `POST /api/auth/login`

* Sau khi đăng nhập thành công:

  * Backend tạo JWT token.
  * Lưu session vào Redis với TTL 72h.

```json
Key: session:<token>
Value:
{
  "id": 1,
  "email": "test@example.com",
  "role": "user",
  "permissions": ["read:products"]
}
```
`Post http://localhost:8080/admin/products`
*sau khi đăng nhập:
*Check API lấy thông tin để gọi phần Admin và kiểm tra phân quyền nếu là Admin sẽ log thành công và user thì sẽ từ chối không được phép truy cập
json
{
    "error": "Chỉ admin được phép truy cập"
}


---

### 2. Middleware xử lý mỗi request

File: `middlewares/auth.go`

```go
func JWTAuthMiddleware() gin.HandlerFunc {
  Lấy token từ header Authorization
  → Validate token
  → Nếu hợp lệ:
      → Kiểm tra Redis với key session:<token>
      → Nếu có: dùng thông tin trong Redis
      → Nếu không: truy DB và lưu lại Redis
  → Gán vào context: user_id, role, permissions
}
```

---

## 🔧 File liên quan

| File                   | Vai trò                                  |
| ---------------------- | ---------------------------------------- |
| `middlewares/auth.go`  | Kiểm tra JWT + Redis session             |
| `controllers/auth.go`  | Đăng ký, đăng nhập và tạo token          |
| `config/redis.go`      | Cấu hình Redis client                    |
| `utils/jwt.go`         | Hàm tạo và validate JWT                  |
| `models/permission.go` | (Tùy chọn) Cấu hình permission theo role |

---

##  Kiểm tra Redis thủ công

```bash
docker exec -it perfume-api-redis-1 redis-cli

# Xem tất cả session:
keys session:*

# Xem thông tin user lưu theo token
get session:<your_jwt_token>
```

---

## 📌 Ghi chú bổ sung

* TTL Redis sẽ được cập nhật mỗi lần request hợp lệ.
* Nếu Redis mất kết nối, hệ thống sẽ fallback dùng DB.
