# Auth Server

## 요구사항
1. 회원가입: 사용자 이름과 비밀번호를 받아서 저장
2. 로그인: 사용자 이름과 비밀번호를 검증하여 로그인 성공/실패 여부 반환
3. 데이터 저장: 메모리 내 **맵(Map)**을 사용
4. HTTP API: JSON 형식의 요청과 응답

## 동작 확인

### 회원가입
```
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"password123"}' http://localhost:8080/register
```
응답:
```
User testuser registered successfully
```

### 로그인 성공
```
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"password123"}' http://localhost:8080/login
```

응답 시:
```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNjg1MzgxNzI4fQ.TyIgU2YmIQrX"
}
```


### 로그인 실패 (잘못된 비밀번호)
```
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"wrongpassword"}' http://localhost:8080/login
```

응답:
```
Invalid password
```

### 보호된 리소스 접근

```
curl -X GET -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNjg1MzgxNzI4fQ.TyIgU2YmIQrX" http://localhost:8080/protected
```

응답:
```
You have accessed a protected resource!
```

## 테스트

테스트 실행:

```
go test
```

출력 예시:

```
PASS
ok  	_/path/to/auth-server	0.005s
```