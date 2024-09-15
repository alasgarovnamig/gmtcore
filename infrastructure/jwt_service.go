package infrastructure

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type JwtService struct {
	SecretKey                           []byte
	Issuer                              string
	TokenValidityInSeconds              int
	TokenValidityInSecondsForRememberMe int
}
type UserCredential struct {
	Claims       jwt.MapClaims
	UserParentID uint
	UserRoleID   uint
}

func (j *JwtService) GenerateToken(userID uint, subject string, rememberMe bool) (string, error) {
	claims := jwt.MapClaims{
		"iss": j.Issuer,
		"sub": subject,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(
			func() int {
				if rememberMe {
					return j.TokenValidityInSecondsForRememberMe
				}
				return j.TokenValidityInSeconds
			}(),
		) * time.Second).Unix(),
		"user_id": userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(j.SecretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *JwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.SecretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	return parsedToken, nil
}

func (j *JwtService) TokenClaimsSetToContext(context *fiber.Ctx, token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("could not extract claims from token")
	}

	userCredential := UserCredential{
		Claims: claims,
	}

	context.Locals("UserCredential", userCredential)
	return nil
}

func (j *JwtService) GetValueByKeyFromJWTClaims(jwtClaims jwt.MapClaims, key string) (interface{}, error) {
	//Claims'in jwt.MapClaims olduğunu varsayıyoruz, çünkü jwt.Claims indeksleme desteklemiyor
	//jwtClaims, ok := jwtClaims.(jwt.MapClaims)
	//if !ok {
	//	return nil, fmt.Errorf("invalid claims type")
	//}

	// İstenen key'in jwtClaims içinde olup olmadığını kontrol ediyoruz
	if value, exists := jwtClaims[key]; exists {
		return value, nil
	}
	return nil, fmt.Errorf("key %s not found in context", key)
}
