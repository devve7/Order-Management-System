package auth

type AccessToken string

func NewAccessToken(token string) (AccessToken, error) {
	if token == "" {
		return "", ErrInvalidToken
	}
	return AccessToken(token), nil
}

func (t AccessToken) String() string {
	return string(t)
}
