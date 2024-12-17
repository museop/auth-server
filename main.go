package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// User 구조체
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// 메모리 데이터 저장소
var (
	userStore = make(map[string]string) // username: password
	mutex     = sync.Mutex{}            // 동시 접근 방지용 뮤텍스
)

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

	// 메모리에 저장
	userStore[user.Username] = user.Password
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

	// 유저 확인 및 비밀번호 검증
	if password, exists := userStore[user.Username]; exists {
		if password == user.Password {
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
