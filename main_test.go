package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
}

// 회원가입 중복 체크 테스트
func TestRegisterHandlerDuplicate(t *testing.T) {
	// 준비
	userStore = make(map[string]string)
	userStore["testuser"] = "password123" // 이미 존재하는 사용자
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
	userStore["testuser"] = "password123" // 테스트 유저 추가
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
}

// 로그인 실패 테스트 (잘못된 비밀번호)
func TestLoginHandlerInvalidPassword(t *testing.T) {
	// 준비
	userStore = make(map[string]string)
	userStore["testuser"] = "password123"
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
