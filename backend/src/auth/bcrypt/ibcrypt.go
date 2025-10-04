package bcrypt

type Bcrypt interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hash, password string) error
}
