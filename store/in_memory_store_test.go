package store

import (
	"testing"

	_ "github.com/lib/pq"
)

func TestInMemoryUserStore(t *testing.T) {
	store := NewInMemoryUserStore()

	// 사용자 저장 테스트
	err := store.SaveUser("testuser", "hashedpassword123")
	if err != nil {
		t.Fatalf("failed to save user: %v", err)
	}

	// 중복 사용자 저장 테스트
	err = store.SaveUser("testuser", "hashedpassword123")
	if err == nil {
		t.Fatalf("expected error when saving duplicate user")
	}

	// 사용자 가져오기 테스트
	hashedPassword, err := store.GetUser("testuser")
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	if hashedPassword != "hashedpassword123" {
		t.Errorf("expected hashed password to be %q, got %q", "hashedpassword123", hashedPassword)
	}

	// 존재하지 않는 사용자 가져오기 테스트
	_, err = store.GetUser("nonexistent")
	if err == nil {
		t.Fatalf("expected error when getting nonexistent user")
	}
}
