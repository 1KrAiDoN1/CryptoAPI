package repository

type TokenRepository interface {
	GenerateJWToken(email string, password string) (string, error)
	GenerateRefreshToken() (string, error)
	GetTokens(email string, password string) (string, string, error)
	ParseToken(access_token string) (int, error)
	HashToken(Password string) string
}
