package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

// User 구조체
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 메모리 데이터 저장소
var (
	userStore = make(map[string]string) // username: hashedPassword
	mutex     = sync.Mutex{}            // 동시 접근 방지용 뮤텍스
)

// 비밀번호 해싱 함수
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 비밀번호 검증 함수
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 회원가입 핸들러
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// 이미 존재하는 유저 체크
	if _, exists := userStore[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// 비밀번호 해싱
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// 해시된 비밀번호 저장
	userStore[user.Username] = hashedPassword
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s registered successfully", user.Username)
}

// 로그인 핸들러
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	// 유저 확인
	if hashedPassword, exists := userStore[user.Username]; exists {
		// 비밀번호 검증
		if checkPasswordHash(user.Password, hashedPassword) {
			fmt.Fprintf(w, "Login successful for user %s", user.Username)
			return
		}
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	http.Error(w, "User does not exist", http.StatusNotFound)
}

func main() {
	// 라우트 설정
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	// 서버 실행
	port := ":8080"
	log.Printf("Server is running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
