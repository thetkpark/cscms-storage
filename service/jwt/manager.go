package jwt

type Manager interface {
	Generate(userID string) (string, error)
	Validate(tokenString string) (string, error)
}
