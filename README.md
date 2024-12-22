# Auth Server

## 요구사항
1. 회원가입: 사용자 이름과 비밀번호를 받아서 저장
2. 로그인: 사용자 이름과 비밀번호를 검증하여 로그인 성공/실패 여부 판단 
 - 성공 시 JWT 토큰 반환
 - 실패 시 에러 메시지 반환
3. 데이터 저장: 사용자 정보를 Postgres DB에 저장
4. HTTP API: JSON 형식의 요청과 응답


## 동작 확인

### DB 및 서버 실행 

Postgres 실행
```sh
docker-compose up -d
```

서버 실행
```sh
go run main.go
```

### 회원가입
```sh
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"password123"}' http://localhost:8080/register
```

응답:
```sh
User testuser registered successfully
```

### 로그인 성공
```sh
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"password123"}' http://localhost:8080/login
```

응답 예시:
```sh
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNjg1MzgxNzI4fQ.TyIgU2YmIQrX"
}
```


### 로그인 실패 (잘못된 비밀번호)
```sh
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"wrongpassword"}' http://localhost:8080/login
```

응답:
```sh
Invalid password
```

### 보호된 리소스 접근

```sh
curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNjg1MzgxNzI4fQ.TyIgU2YmIQrX" http://localhost:8080/protected
```

응답:
```sh
You have accessed a protected resource!
```

## 테스트

테스트 실행:

```sh
go test -v ./...
```

출력 예시:

```sh
=== RUN   TestRegisterHandler
--- PASS: TestRegisterHandler (0.14s)
=== RUN   TestRegisterHandlerDuplicate
--- PASS: TestRegisterHandlerDuplicate (0.13s)
=== RUN   TestLoginHandlerSuccess
--- PASS: TestLoginHandlerSuccess (0.13s)
=== RUN   TestLoginHandlerInvalidPassword
--- PASS: TestLoginHandlerInvalidPassword (0.13s)
=== RUN   TestLoginHandlerUserNotFound
--- PASS: TestLoginHandlerUserNotFound (0.00s)
=== RUN   TestProtectedHandler
=== RUN   TestProtectedHandler/Valid_token
=== RUN   TestProtectedHandler/Missing_token
=== RUN   TestProtectedHandler/Invalid_token
--- PASS: TestProtectedHandler (0.00s)
    --- PASS: TestProtectedHandler/Valid_token (0.00s)
    --- PASS: TestProtectedHandler/Missing_token (0.00s)
    --- PASS: TestProtectedHandler/Invalid_token (0.00s)
PASS
ok      github.com/museop/auth-server   (cached)
=== RUN   TestInMemoryUserStore
--- PASS: TestInMemoryUserStore (0.00s)
=== RUN   TestPostgreSQLUserStore
2024/12/22 22:46:20 github.com/testcontainers/testcontainers-go - Connected to docker: 
  Server Version: 24.0.6
  API Version: 1.43
  Operating System: Docker Desktop
  Total Memory: 7851 MB
  Testcontainers for Go Version: v0.34.0
  Resolved Docker Host: unix:///var/run/docker.sock
  Resolved Docker Socket Path: /var/run/docker.sock
  Test SessionID: 9d002d6a58afb8671710b71f815fcfd21526ab4a74b87957c1fb2e4ae630d0b4
  Test ProcessID: dd8064b3-41a1-4a27-ba54-27f191968a15
2024/12/22 22:46:27 🐳 Creating container for image testcontainers/ryuk:0.10.2
2024/12/22 22:46:27 ✅ Container created: ee85a4ecc575
2024/12/22 22:46:27 🐳 Starting container: ee85a4ecc575
2024/12/22 22:46:27 ✅ Container started: ee85a4ecc575
2024/12/22 22:46:27 ⏳ Waiting for container id ee85a4ecc575 image: testcontainers/ryuk:0.10.2. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false}
2024/12/22 22:46:27 🔔 Container is ready: ee85a4ecc575
2024/12/22 22:46:27 🐳 Creating container for image postgres:15
2024/12/22 22:46:27 ✅ Container created: 32fb4754b0db
2024/12/22 22:46:27 🐳 Starting container: 32fb4754b0db
2024/12/22 22:46:27 ✅ Container started: 32fb4754b0db
2024/12/22 22:46:27 ⏳ Waiting for container id 32fb4754b0db image: postgres:15. Waiting for: &{Port:5432/tcp timeout:<nil> PollInterval:100ms skipInternalCheck:false}
2024/12/22 22:46:28 🔔 Container is ready: 32fb4754b0db
2024/12/22 22:46:28 🐳 Stopping container: 32fb4754b0db
2024/12/22 22:46:28 ✅ Container stopped: 32fb4754b0db
2024/12/22 22:46:28 🐳 Terminating container: 32fb4754b0db
2024/12/22 22:46:28 🚫 Container terminated: 32fb4754b0db
--- PASS: TestPostgreSQLUserStore (8.36s)
PASS
ok      github.com/museop/auth-server/store     (cached)
```
