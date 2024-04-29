package hasher

type Hasher interface {
	GenerateFromPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, password string) error
}

var hasherImpl Hasher

func SetHasher(h Hasher) {
	hasherImpl = h
}

func GetHasher() Hasher {
	return hasherImpl
}

func GenerateFromPassword(password string) (string, error) {
	return hasherImpl.GenerateFromPassword(password)
}

func CompareHashAndPassword(hashedPassword, password string) error {
	return hasherImpl.CompareHashAndPassword(hashedPassword, password)
}
