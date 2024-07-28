package auth

import (
	"time"
	"log"
	"net/http"
	"context"
	//"strings"
	
	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Username string `json: "username"`
	jwt.RegisteredClaims
}


func NewClaims(id int, timeExpired time.Time) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": id,
		"exp": jwt.NewNumericDate(timeExpired),
	}
}

func GetTokenString(userId int) (string, error) {
	timeExpired := time.Now().Add(5 * time.Minute) //take from env file or ...
	secretString := []byte("se number 1 from cr go and et pass") //take from env file or ...

	claims := NewClaims(userId, timeExpired)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString(secretString)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}

func CheckUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Println("empty token")
			return
		}
		
		var claims jwt.MapClaims
	
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
			return []byte("se number 1 from cr go and et pass"), nil
		})
		
		if err != nil {
			log.Println(err)
			return
		}
	
		log.Println(token)
		
		ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
