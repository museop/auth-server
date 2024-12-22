package store

// UserStore 인터페이스 정의
type UserStore interface {
	SaveUser(username, hashedPassword string) error
	GetUser(username string) (hashedPassword string, err error)
}
