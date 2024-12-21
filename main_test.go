package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

// 사용자 정보를 JSON으로 직렬화하기 위한 헬퍼 함수
func createUserPayload(username, password string) []byte {
	user := User{
		Username: username,
		Password: password,
	}
	data, _ := json.Marshal(user)
	return data
}

// 회원가입 테스트
func TestRegisterHandler(t *testing.T) {
	// 준비
	userStore = make(map[string]string) // 테스트를 위해 메모리 저장소 초기화
	reqBody := createUserPayload("testuser", "password123")

	// HTTP 요청 생성
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	// 핸들러 호출
	registerHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d; got %d", http.StatusCreated, res.StatusCode)
	}

	// 비밀번호가 해싱되어 저장되었는지 검증
	if hashedPassword, exists := userStore["testuser"]; exists {
		if checkPasswordHash("password123", hashedPassword) == false {
			t.Error("password hash does not match")
		}
	} else {
		t.Error("user not found in userStore")
	}
}

// 회원가입 중복 체크 테스트
func TestRegisterHandlerDuplicate(t *testing.T) {
	// 준비
	userStore = make(map[string]string)
	hashedPassword, _ := hashPassword("password123")
	userStore["testuser"] = hashedPassword // 이미 존재하는 사용자
	reqBody := createUserPayload("testuser", "password123")

	// HTTP 요청 생성
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	// 핸들러 호출
	registerHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d; got %d", http.StatusConflict, res.StatusCode)
	}
}

// 로그인 성공 테스트
func TestLoginHandlerSuccess(t *testing.T) {
	// 준비
	userStore = make(map[string]string)
	hashedPassword, _ := hashPassword("password123")
	userStore["testuser"] = hashedPassword // 비밀번호 해싱 후 저장
	reqBody := createUserPayload("testuser", "password123")

	// HTTP 요청 생성
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	// 핸들러 호출
	loginHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	// 응답 본문에서 JWT 토큰 추출
	var responseBody map[string]string
	err := json.NewDecoder(res.Body).Decode(&responseBody)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	token, exists := responseBody["token"]
	if !exists || token == "" {
		t.Fatalf("token not found in response body")
	}

	// JWT 토큰 검증
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecretKey, nil
	})
	if err != nil || !parsedToken.Valid {
		t.Fatalf("invalid token: %v", err)
	}

	// 클레임 확인 (예: username)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("failed to parse claims")
	}

	if claims["username"] != "testuser" {
		t.Errorf("expected username %s; got %v", "testuser", claims["username"])
	}
}

// 로그인 실패 테스트 (잘못된 비밀번호)
func TestLoginHandlerInvalidPassword(t *testing.T) {
	// 준비
	userStore = make(map[string]string)
	hashedPassword, _ := hashPassword("password123")
	userStore["testuser"] = hashedPassword
	reqBody := createUserPayload("testuser", "wrongpassword")

	// HTTP 요청 생성
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	// 핸들러 호출
	loginHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d; got %d", http.StatusUnauthorized, res.StatusCode)
	}
}

// 로그인 실패 테스트 (존재하지 않는 사용자)
func TestLoginHandlerUserNotFound(t *testing.T) {
	// 준비
	userStore = make(map[string]string)
	reqBody := createUserPayload("nonexistent", "password123")

	// HTTP 요청 생성
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()

	// 핸들러 호출
	loginHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected status %d; got %d", http.StatusNotFound, res.StatusCode)
	}
}

func TestProtectedHandler(t *testing.T) {
	// JWT 생성
	token, err := generateJWT("testuser")
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	// 테스트 케이스 정의
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid token",
			authHeader:     "Bearer " + token,
			expectedStatus: http.StatusOK,
			expectedBody:   "You have accessed a protected resource!",
		},
		{
			name:           "Missing token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing or invalid Authorization header\n",
		},
		{
			name:           "Invalid token",
			authHeader:     "Bearer invalidtoken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// HTTP 요청 생성
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			rec := httptest.NewRecorder()

			// 미들웨어와 핸들러 호출
			handler := jwtMiddleware(protectedHandler)
			handler.ServeHTTP(rec, req)

			// 결과 검증
			res := rec.Result()
			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, res.StatusCode)
			}

			body := rec.Body.String()
			if body != tt.expectedBody {
				t.Errorf("expected body %q; got %q", tt.expectedBody, body)
			}
		})
	}
}
