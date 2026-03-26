package rtc

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type LiveKit struct {
	key    string
	secret string
	url    string
}

func NewLiveKit(key, secret, url string) *LiveKit {
	return &LiveKit{key: key, secret: secret, url: url}
}

func (l *LiveKit) GetURL() string {
	return l.url
}

// videoGrant mirrors LiveKit's VideoGrant claim structure.
type videoGrant struct {
	RoomJoin bool   `json:"roomJoin"`
	Room     string `json:"room"`
}

type livekitClaims struct {
	Video videoGrant `json:"video"`
	jwt.RegisteredClaims
}

// GetToken generates a LiveKit-compatible JWT access token.
func (l *LiveKit) GetToken(room, identity string) (string, error) {
	now := time.Now()
	claims := livekitClaims{
		Video: videoGrant{
			RoomJoin: true,
			Room:     room,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    l.key,
			Subject:   identity,
			ID:        uuid.New().String(),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.secret))
}
