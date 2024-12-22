# Auth Server

## 요구사항
1. 회원가입: 사용자 이름과 비밀번호를 받아서 저장
2. 로그인: 사용자 이름과 비밀번호를 검증하여 로그인 성공/실패 여부 판단 
 - 성공 시 JWT 토큰 반환
 - 실패 시 에러 메시지 반환
3. 데이터 저장: 메모리 내 맵(Map)을 사용 (*영구 저장소로 변경 예정*)
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
```sh

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
--- PASS: TestRegisterHandler (0.15s)
=== RUN   TestRegisterHandlerDuplicate
--- PASS: TestRegisterHandlerDuplicate (0.07s)
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
```