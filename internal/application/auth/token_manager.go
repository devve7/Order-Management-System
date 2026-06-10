package auth

type TokenManager interface {
	CreateAccessToken(actor Actor) (AccessToken, error)
	ParseAccessToken(token AccessToken) (Actor, error)
}
