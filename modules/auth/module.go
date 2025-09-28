package auth

import (
	// "fmt"
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/pedrokoblitz/opus-go/actions"
)

func Register(registry *actions.ActionRegistry) {
	registry.RegisterPayload("default:id", func() actions.Payload { return &actions.DefaultID{} })

	registry.RegisterPayload("password:request", func() actions.Payload { return &ResetRequest{} })
	registry.RegisterPayload("password:reset", func() actions.Payload { return &Reset{} })
	registry.RegisterPayload("auth:register", func() actions.Payload { return &SignUp{} })
	registry.RegisterPayload("auth:login", func() actions.Payload { return &Login{} })
	registry.RegisterPayload("user", func() actions.Payload { return &User{} })
	registry.RegisterPayload("auth:token", func() actions.Payload { return &Token{} })

	registry.RegisterAction("auth:before", BeforeShowHook)

	registry.RegisterAction("auth:register", SignUpAction)
	registry.RegisterAction("auth:confirm", ConfirmAction)
	registry.RegisterAction("auth:login", LoginAction)
	registry.RegisterAction("auth:logout", LogoutAction)
	registry.RegisterAction("auth:show", ShowUserAction)
	registry.RegisterAction("auth:update", UpdateUserAction)
	registry.RegisterAction("auth:remove", RemoveAccountAction)

	registry.RegisterAction("password:request", RequestPasswordResetAction)
	registry.RegisterAction("password:reset", ResetPasswordAction)
}

type AuthBasePayload struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"userId"`
	Email         string `json:"email"`
	Token         string `json:"token"`
	Password      string `json:"password" yaml:"password"`
	InputPassword string
	ExpiresAt     time.Time `json:"expires_at"`
}

// type AuthBaseRequest struct {}
// type AuthBaseResponse struct {}

var jwtSecret = []byte("your-secret-key") // keep this private

func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type Claims struct {
	UserID uint `json:"userId"`
	jwt.StandardClaims
}

func GetClaims(userID uint) *Claims {
	expirationTime := time.Now().Add(24 * time.Hour)
	return &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
}

func GenerateRandomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes)
}

func GenerateBearerToken(userID uint) (string, error) {
	claims := GetClaims(userID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseBearerToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GenerateCSRFToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(1 * time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

func ParseCSRFToken(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})
	return err == nil && token.Valid
}
