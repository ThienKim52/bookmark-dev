7
hj# Hướng Dẫn Toàn Diện Về Unit Test Các Tầng (Layers) Trong Golang

Tài liệu này giải thích chi tiết **cú pháp (syntax)**, **bản chất tư duy (logic)** khi viết Unit Test cho từng tầng kiến trúc trong Golang (Handler, Service, Repository), sự khác biệt giữa các tầng test, và giải đáp thắc mắc cốt lõi: *"Nếu ở tầng Service ta tự biên tự diễn kết quả trả về bằng `.Return()`, thì làm sao biết code Service thực sự hoạt động?"*.

---

## 1. Giải Đáp Thắc Mắc Cốt Lỗi: "Nếu tự mock `.Return()`, làm sao biết test đang hoạt động?"

> *"Nếu ở Service layer mình dùng `setupMock` rồi `.Return(...)` tự định nghĩa kết quả trả về theo ý mình, thì chẳng phải mình đang tự lừa mình sao? Làm sao biết hàm Service có thật sự chạy đúng ngoài đời?"*

Đây là câu hỏi kinh điển nhất của mọi kỹ sư khi tiếp cận Unit Test và Mocking. Để hiểu được điều này, hãy nắm 3 luận điểm sau:

### 1.1. Ẩn dụ: "Phòng thử nghiệm va chạm ô tô" (Crash Test)
Hãy tưởng tượng bạn là kỹ sư chế tạo **Hệ thống Túi Khí Ô tô** (tương đương tầng `Service`):
- Để kiểm tra túi khí có bung khi va chạm không, bạn không cần phải lái chiếc xe ra đường cao tốc thật và đâm vào một chiếc xe tải thật (`Database / Third-party API thật`).
- Bạn đặt hệ thống túi khí vào **phòng lab**, và gửi vào đó **tín hiệu xung điện giả lập** một vụ va chạm (`Mock Repo .Return(va_cham)`).
- **Mục tiêu của bạn không phải là kiểm tra xem "vụ va chạm có thật không"**, mà là kiểm tra: **"Khi nhận tín hiệu va chạm đó, vi mạch túi khí của bạn có kích nổ đúng 0.05 giây hay không?"**.

### 1.2. Đối tượng được kiểm tra (Subject Under Test) là ai?
- **Mock KHÔNG PHẢI là đối tượng được test.** Mock chỉ là **Môi trường giả lập (Inputs/Conditions)**.   
- **Service MỚI LÀ đối tượng được test.**
- Bạn dùng `.Return(...)` để tạo ra **các kịch bản thế giới khác nhau**:
  1. *Kịch bản 1:* Giả sử DB tìm thấy User, Pass đúng -> Service có nhả Token không?
  2. *Kịch bản 2:* Giả sử DB báo lỗi sập mạng (`errors.New("db dead")`) -> Service có crash không hay trả về error an toàn?
  3. *Kịch bản 3:* Giả sử Hacker nhập sai pass (Hasher trả về `false`) -> **Service có vô tình cấp Token không?**

### 1.3. Bằng chứng: Khi Service bị bug, Test SẼ BỊ FAIL ngay lập tức!
Xem hàm [Login](file:///d:/CODE/Golang/golang-dev/internal/app/service/user/login.go#L26):
```go
if !s.hasher.Compare(user.Password, password) {
    return "", ErrInvalidCredentials
}
// Tiếp tục tạo JWT...
```
Giả sử một ngày nào đó, một lập trình viên vô tình sửa nhầm code:
```go
// BUG NGUY HIỂM: Bỏ quên dấu phủ định '!'
if s.hasher.Compare(user.Password, password) {
    return "", ErrInvalidCredentials
}
```
Lúc này, dù bạn có setup mock hasher `.Return(false)`, thì khi chạy test case `"Wrong password"`, hàm `Login` sẽ không dừng lại mà tiếp tục đi xuống sinh JWT token. Test sẽ thấy:
- Kết quả mong đợi: `ErrInvalidCredentials`
- Kết quả thực tế: `Token: "jwt-token-abc", err: nil`
👉 **TEST FAILED NGAY LẬP TỨC!**
Test đã cứu hệ thống của bạn khỏi một lỗ hổng bảo mật nghiêm trọng. Đó chính là bằng chứng Unit Test đang bảo vệ bạn!

### 1.4. Behavior Verification (Xác minh hành vi)
Với Mockery, bạn không chỉ kiểm tra output trả về, bạn còn kiểm tra **hành vi**:
- *Nếu sai password, hàm `GenerateJWT` có bị gọi trộm không?* (Nếu gọi, Mockery lập tức ném lỗi `Unexpected Call`).
- *Hàm `GetUserByUsername` có bị gọi 2 lần gây tốn tài nguyên không?* (`.Once()`).


---

## 3. Cú Pháp & Chi Tiết Từng Tầng Test

---

### TẦNG 1: Repository Layer Unit Test (Tầng Truy Cập Dữ Liệu)

#### Nhiệm vụ của tầng này:
- Kiểm tra các câu truy vấn SQL / GORM có đúng cú pháp không.
- Kiểm tra việc mapping struct Go với các cột trong DB.
- Kiểm tra các ràng buộc: Unique constraint, Foreign key, Default values.

#### Phương pháp test:
- **KHÔNG DÙNG MOCK INTERFACE Ở ĐÂY.** 
- Vì Repository chính là tầng giao tiếp trực tiếp với database. Nếu bạn mock database thì bạn chẳng test được gì cả (bạn đang mock chính cái bạn cần test).
- **Giải pháp chuẩn:** Dùng **SQLite In-Memory Database** (`file::memory:?cache=shared`). Tốc độ thực thi chỉ vài mili-giây, độc lập, tự dọn dẹp sau khi chạy.

#### Mẫu code thực tế trong dự án ([create_test.go](file:///d:/CODE/Golang/golang-dev/internal/app/repository/user/create_test.go)):
```go
func TestSQLRepository_CreateUser(t *testing.T) {
    // 1. Khởi tạo DB SQLite ảo trong RAM và AutoMigrate bảng
    db := sqldb.InitMockDB(t)
    err := db.AutoMigrate(&model.User{})
    assert.NoError(t, err)

    // 2. Khởi tạo Repository thật gắn với DB ảo
    repo := NewSqlRepository(db)

    // 3. Thực thi hành động cần test
    user, err := repo.CreateUser(context.Background(), &model.User{
        Username: "test1",
        Email:    "test1@gmail.com",
    })

    // 4. Assert kết quả trả về từ hàm
    assert.NoError(t, err)
    assert.NotNil(t, user)

    // 5. Query trực tiếp DB ảo để kiểm tra dữ liệu ĐÃ THỰC SỰ LƯU ĐÚNG CỘT CHƯA
    var savedUser model.User
    err = db.First(&savedUser, "username = ?", "test1").Error
    assert.NoError(t, err)
    assert.Equal(t, "test1@gmail.com", savedUser.Email)
}
```

---

### TẦNG 2: Service Layer Unit Test (Tầng Nghiệp Vụ)

#### Nhiệm vụ của tầng này:
- Kiểm tra **luồng nghiệp vụ thuần túy (Business Logic)**.
- Rẽ nhánh điều kiện: User có tồn tại không? Mật khẩu có khớp không? Tài khoản có bị khóa không?
- Xử lý và chuyển đổi mã lỗi (Error Mapping: từ DB error sang Business error).

#### Phương pháp test:
- **Mock tất cả các dependencies bên ngoài** (Repository, Hasher, JWT, Redis, Mailer...).
- Sử dụng thư viện **Mockery + testify/mock**.

#### Cú pháp quan trọng của Mockery:
1. `mockObj.On("MethodName", args...).Return(returnValues...).Once()`
2. `mock.Anything`: Dùng cho tham số không cần kiểm tra chính xác giá trị (như `context.Context`).
3. `mock.MatchedBy(predicateFunc)`: Kiểm tra điều kiện logic của tham số (ví dụ claims chứa đúng `sub`).

#### Mẫu code thực tế trong dự án ([login_test.go](file:///d:/CODE/Golang/golang-dev/internal/app/service/user/login_test.go)):
```go
func TestLogin(t *testing.T) {
    t.Parallel()

    testCases := []struct {
        name          string
        username      string
        password      string
        setupMocks    func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator)
        expectedToken string
        expectedErr   error
    }{
        {
            name:     "Login thành công",
            username: "johndoe",
            password: "password123",
            setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
                // Kịch bản: Repo tìm thấy user
                repo.On("GetUserByUsername", mock.Anything, "johndoe").
                    Return(&model.User{ID: "123", Username: "johndoe", Password: "hashed-pw"}, nil).
                    Once()

                // Kịch bản: Hasher so sánh khớp password
                hasher.On("Compare", "hashed-pw", "password123").
                    Return(true).
                    Once()

                // Kịch bản: JWT Generator sinh token
                jwtGen.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
                    return claims["sub"] == "123"
                })).Return("fake-token-xyz", nil).Once()
            },
            expectedToken: "fake-token-xyz",
            expectedErr:   nil,
        },
        {
            name:     "Sai password -> Chặn đứng, không sinh JWT",
            username: "johndoe",
            password: "wrong-password",
            setupMocks: func(repo *mocks_userrepo.Repository, hasher *mocks_utils.Hasher, jwtGen *mocks_jwtutils.JWTGenerator) {
                repo.On("GetUserByUsername", mock.Anything, "johndoe").
                    Return(&model.User{ID: "123", Password: "hashed-pw"}, nil).
                    Once()

                // Hasher trả về false
                hasher.On("Compare", "hashed-pw", "wrong-password").
                    Return(false).
                    Once()

                // CHÚ Ý: jwtGen KHÔNG ĐƯỢC SETUP Ở ĐÂY!
                // Nếu code Service bị lỗi mà vẫn gọi jwtGen, Mockery sẽ báo FAIL ngay lập tức!
            },
            expectedToken: "",
            expectedErr:   ErrInvalidCredentials,
        },
    }

    for _, tc := range testCases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            // 1. Tạo các Mock sạch cho từng test case
            repoMock := mocks_userrepo.NewRepository(t)
            hasherMock := mocks_utils.NewHasher(t)
            jwtGenMock := mocks_jwtutils.NewJWTGenerator(t)

            // 2. Nạp kịch bản vào Mock
            tc.setupMocks(repoMock, hasherMock, jwtGenMock)

            // 3. Khởi tạo Service với các Mock
            svc := NewService(repoMock, hasherMock, jwtGenMock)

            // 4. Chạy hàm cần test
            token, err := svc.Login(context.Background(), tc.username, tc.password)

            // 5. Kiểm tra kết quả
            assert.Equal(t, tc.expectedToken, token)
            if tc.expectedErr != nil {
                assert.EqualError(t, err, tc.expectedErr.Error())
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

### TẦNG 3: Handler Layer Unit Test (Tầng Giao Thức HTTP)

#### Nhiệm vụ của tầng này:
- Kiểm tra việc bóc tách dữ liệu từ HTTP Request: Đọc JSON Body (`ShouldBindJSON`), đọc URL Param (`c.Param`), đọc Header.
- Kiểm tra Validation đầu vào tầng Transport (ví dụ: password tối thiểu 8 ký tự).
- Kiểm tra việc trả về đúng **HTTP Status Code** (200, 400, 401, 403, 500) và format JSON response.

#### Phương pháp test:
- Dùng `net/http/httptest`:
  - `httptest.NewRecorder()`: Đóng vai trò như "trình duyệt ảo" để hứng HTTP Response từ Gin.
  - `gin.CreateTestContext(w)`: Tạo context Gin giả lập cho bài test.
- **Mock tầng Service**: Vì Handler chỉ điều hướng request tới Service, ta không cần chạy Service thật.

#### Mẫu code thực tế trong dự án ([handler/user/login_test.go](file:///d:/CODE/Golang/golang-dev/internal/app/handler/user/login_test.go)):
```go
func TestHandler_Login(t *testing.T) {
    gin.SetMode(gin.TestMode)

    testCases := []struct {
        name               string
        reqBody            string
        setupMock          func(c *gin.Context) *mocks_userservice.Service
        expectedStatusCode int
        expectedBody       string
    }{
        {
            name:    "Login thành công -> Trả về 200 OK",
            reqBody: `{"username":"test123","password":"password123"}`,
            setupMock: func(c *gin.Context) *mocks_userservice.Service {
                mockSvc := mocks_userservice.NewService(t)
                mockSvc.On("Login", c, "test123", "password123").
                    Return("jwt-token-string", nil).
                    Once()
                return mockSvc
            },
            expectedStatusCode: http.StatusOK,
            expectedBody:       `"Logged in successfully!"`,
        },
        {
            name:    "Validation fail (password < 8 ký tự) -> Trả về 400 Bad Request",
            reqBody: `{"username":"test123","password":"123"}`, // Quá ngắn
            setupMock: func(c *gin.Context) *mocks_userservice.Service {
                // Service KHÔNG ĐƯỢC PHÉP ĐƯỢC GỌI vì Handler đã chặn từ lúc validate JSON
                return mocks_userservice.NewService(t)
            },
            expectedStatusCode: http.StatusBadRequest,
            expectedBody:       "",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // 1. Tạo ResponseRecorder giả lập response client
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)

            // 2. Tạo Request giả lập
            c.Request = httptest.NewRequest("POST", "/v1/users/login", bytes.NewBufferString(tc.reqBody))
            c.Request.Header.Set("Content-Type", "application/json")

            // 3. Khởi tạo Handler với mock service
            mockSvc := tc.setupMock(c)
            hdl := NewHandler(mockSvc)

            // 4. Gọi hàm handler
            hdl.Login(c)

            // 5. Assert mã HTTP và nội dung body
            assert.Equal(t, tc.expectedStatusCode, w.Code)
            if tc.expectedBody != "" {
                assert.Contains(t, w.Body.String(), tc.expectedBody)
            }
        })
    }
}
```

---

### TẦNG 4: Integration Test (Kiểm Thử Tích Hợp)

#### Sự khác biệt giữa Unit Test và Integration Test:
| Đặc điểm | Unit Test (Từng tầng riêng lẻ) | Integration Test (Toàn hệ thống) |
| :--- | :--- | :--- |
| **Phạm vi** | Chỉ test 1 hàm / 1 struct cô lập | Test toàn bộ luồng đi từ HTTP -> Handler -> Service -> Repo -> DB |
| **Mục đích** | Kiểm tra logic ngóc ngách, edge cases, error cases | Kiểm tra xem các tầng có "khớp bánh răng" với nhau không |
| **Tốc độ** | Siêu nhanh (vài mili-giây) | Chậm hơn (vài trăm mili-giây đến vài giây) |
| **Database** | Không dùng hoặc dùng SQLite RAM | SQLite / Docker DB / Redis thật |
| **Mocking** | Mock triệt để các dependency | Chỉ mock các third-party bên ngoài (như SendGrid, Stripe, SMS...) |

---

## 4. Bảng So Sánh Tổng Hợp Các Tầng Unit Test

| Tầng | Đối tượng test (SUT) | Thứ được MOCK | Công cụ chính | Trọng tâm kiểm tra |
| :--- | :--- | :--- | :--- | :--- |
| **Repository** | `SqlRepository` struct | Không mock (dùng SQLite RAM) | `sqldb.InitMockDB`, `GORM` | SQL syntax, Transaction, Mapping model, DB constraints |
| **Service** | `service` struct | `Repository`, `Hasher`, `JWT`, v.v. | `mockery`, `testify/mock`, `assert` | **Business logic**, Phân nhánh if/else, Xử lý lỗi, Gọi đúng hàm |
| **Handler** | `userHandler` struct | `Service` | `httptest.NewRecorder`, `gin.CreateTestContext` | HTTP Status code, Parse JSON, Validate DTO, Header |
| **Integration** | Toàn bộ ứng dụng | Chỉ mock external API (Google, Stripe) | `httptest.NewServer`, Router thật, DB | Toàn bộ vòng đời request từ lúc gửi đến lúc ghi DB |

---

## 5. Quy Tắc Vàng Khi Viết Test

1. **Đừng test lại framework:** Không cần viết test để xem Gin có parse được route không, hoặc GORM có chạy được `First()` không. Hãy test **code bạn viết**.
2. **Kiểm tra cả 2 nhánh (Happy Path & Sad Path):** Với mỗi hàm, luôn có ít nhất 1 test case chạy thành công và các test case bao phủ tất cả các nhánh `if err != nil`.
3. **Mỗi tầng một trách nhiệm:**
   - Đừng cố test validation độ dài password ở tầng Service nếu nó đã được validate ở tầng Handler.
   - Đừng cố test cú pháp SQL ở tầng Service, hãy để tầng Repository đảm nhiệm.
4. **Mock phải phản ánh thế giới thực:** Giá trị trả về trong `.Return(...)` phải giống với những gì mà dependency thật sẽ trả về (ví dụ: nếu DB không thấy record thì trả về `dbutils.ErrRecordNotFound`, chứ không tự bịa ra một lỗi lạ).
