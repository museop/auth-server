package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL 드라이버
)

// PostgreSQLUserStore는 PostgreSQL 기반 UserStore 구현체입니다.
type PostgreSQLUserStore struct {
	db *sql.DB
}

// NewPostgreSQLUserStore는 새로운 PostgreSQLUserStore 인스턴스를 생성합니다.
func NewPostgreSQLUserStore(dsn string) (*PostgreSQLUserStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgreSQLUserStore{db: db}, nil
}

// SaveUser는 사용자 이름과 해시된 비밀번호를 PostgreSQL에 저장합니다.
func (s *PostgreSQLUserStore) SaveUser(username, hashedPassword string) error {
	_, err := s.db.Exec(`INSERT INTO users (username, password) VALUES ($1, $2)`, username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to save user: %v", err)
	}
	return nil
}

// GetUser는 사용자 이름으로 해시된 비밀번호를 PostgreSQL에서 가져옵니다.
func (s *PostgreSQLUserStore) GetUser(username string) (string, error) {
	var hashedPassword string
	err := s.db.QueryRow(`SELECT password FROM users WHERE username = $1`, username).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("failed to get user: %v", err)
	}
	return hashedPassword, nil
}
