package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/museop/auth-server/store"
)

// 사용자 정보를 JSON으로 직렬화하기 위한 헬퍼 함수
func createUserPayload(username, password string) []byte {
	user := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: username,
		Password: password,
	}
	data, _ := json.Marshal(user)
	return data
}

func setupTestUserStore() {
	userStore = store.NewInMemoryUserStore()
}

// 회원가입 테스트
func TestRegisterHandler(t *testing.T) {
	setupTestUserStore() // 테스트 환경 초화

	// 준비
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

	// 저장소에 사용자 확인
	hashedPassword, err := userStore.GetUser("testuser")
	if err != nil {
		t.Errorf("failed to find user in store: %v", err)
	}

	// 비밀번호 해싱 검증
	if !checkPasswordHash("password123", hashedPassword) {
		t.Errorf("stored password hash does not match")
	}
}

// 회원가입 중복 사용자 체크 테스트
func TestRegisterHandlerDuplicate(t *testing.T) {
	setupTestUserStore() // 테스트 환경 초기화

	// 첫 번째 사용자 등록
	reqBody := createUserPayload("testuser", "password123")
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	registerHandler(rec, req)

	// 동일 사용자로 다시 등록
	req = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	rec = httptest.NewRecorder()
	registerHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusConflict {
		t.Errorf("expected status %d; got %d", http.StatusConflict, res.StatusCode)
	}
}

// 로그인 성공 테스트
func TestLoginHandlerSuccess(t *testing.T) {
	setupTestUserStore() // 테스트 환경 초기화

	// 사용자 등록
	reqBody := createUserPayload("testuser", "password123")
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	registerHandler(rec, req)

	// 로그인
	reqBody = createUserPayload("testuser", "password123")
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rec = httptest.NewRecorder()
	loginHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; got %d", http.StatusOK, res.StatusCode)
	}

	// 응답 본문에서 JWT 토큰 추출 및 검증
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
	setupTestUserStore() // 테스트 환경 초기화

	// 사용자 등록
	reqBody := createUserPayload("testuser", "password123")
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
	registerHandler(rec, req)

	// 잘못된 비밀번호로 로그인
	reqBody = createUserPayload("testuser", "wrongpassword")
	req = httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rec = httptest.NewRecorder()
	loginHandler(rec, req)

	// 결과 검증
	res := rec.Result()
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %d; got %d", http.StatusUnauthorized, res.StatusCode)
	}
}

// 로그인 실패 테스트 (존재하지 않는 사용자)
func TestLoginHandlerUserNotFound(t *testing.T) {
	setupTestUserStore() // 테스트 환경 초기화

	// 존재하지 않는 사용자로 로그인
	reqBody := createUserPayload("nonexistent", "password123")
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
	rec := httptest.NewRecorder()
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
