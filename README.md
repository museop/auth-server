# Auth Server

## 요구사항
1. 회원가입: 사용자 이름과 비밀번호를 받아서 저장
2. 로그인: 사용자 이름과 비밀번호를 검증하여 로그인 성공/실패 여부 반환
3. 데이터 저장: 메모리 내 **맵(Map)**을 사용
4. HTTP API: JSON 형식의 요청과 응답

## 테스트

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

응답:
```
Login successful for user testuser
```

### 로그인 실패 (잘못된 비밀번호)
```
curl -X POST -H "Content-Type: application/json" -d '{"username":"testuser", "password":"wrongpassword"}' http://localhost:8080/login
```

응답:
```
Invalid password
```
