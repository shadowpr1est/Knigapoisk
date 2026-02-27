package auth

type RegisterInput struct {
	Username string
	Email    string
	Password string
}
type RegisterOutput struct {
	AccessToken  string
	RefreshToken string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type RefreshInput struct {
	RefreshToken string
}

type RefreshOutput struct {
	AccessToken  string
	RefreshToken string
}
