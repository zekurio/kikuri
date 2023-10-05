package discordoauth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type stateClaims struct {
	jwt.StandardClaims
	Payload map[string]string `json:"pld,omitempty"`
}

func (d *DiscordOAuth) getHandler() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return d.stateSigningKey, nil
	}
}

func (d *DiscordOAuth) encodeAndSignWithPayload(payload map[string]string) (string, error) {
	now := time.Now()
	claims := stateClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    d.clientID,
			IssuedAt:  now.Unix(),
			NotBefore: now.Unix(),
			ExpiresAt: now.Add(5 * time.Minute).Unix(),
		},
		Payload: payload,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(d.stateSigningKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (d *DiscordOAuth) decodeAndVerifyStateToken(token string) (stateClaims, error) {
	claims := stateClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, d.getHandler())
	if err != nil {
		return stateClaims{}, err
	}
	return claims, nil
}
