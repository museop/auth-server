package store

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupPostgreSQLContainer() (func(), *sql.DB, error) {
	ctx := context.Background()

	// PostgreSQL 컨테이너 설정
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpassword",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	// 컨테이너 생성
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %v", err)
	}

	// 컨테이너의 호스트 및 포트 가져오기
	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container port: %v", err)
	}

	// PostgreSQL DSN 구성
	dsn := fmt.Sprintf("postgres://testuser:testpassword@%s:%s/testdb?sslmode=disable", host, port.Port())

	// PostgreSQL 연결
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// 정리 함수 정의
	tearDown := func() {
		db.Close()
		container.Terminate(ctx)
	}

	return tearDown, db, nil
}

func TestPostgreSQLUserStore(t *testing.T) {
	tearDown, db, err := setupPostgreSQLContainer()
	if err != nil {
		t.Fatalf("failed to setup PostgreSQL container: %v", err)
	}
	defer tearDown()

	// 테스트 데이터베이스 초기화
	err = resetTestDatabase(db)
	if err != nil {
		t.Fatalf("failed to reset test database: %v", err)
	}

	// PostgreSQL UserStore 생성
	store := &PostgreSQLUserStore{db: db}

	// 사용자 저장 테스트
	err = store.SaveUser("testuser", "hashedpassword123")
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

// 테스트 데이터베이스 초기화 함수
func resetTestDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
	`)
	return err
}
