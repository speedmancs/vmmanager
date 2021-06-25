package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/speedmancs/vmmanager/util"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//A sample use
var user = User{
	ID:       1,
	Username: "harry",
	Password: "$2a$04$8Ox3kwx6Af6FRQrHVVI4UOH7/gqmFfhl.ibIyjxirPVtMMc1l.8Z2",
}

func createToken(userId uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userId
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	return token, nil
}

func checkPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hashedPassword))
	if err != nil {
		return false
	}

	return true
}

func hashAndSalt(password string) string {
	pwd := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	log.Println("bearToken:", bearToken)
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func tokenValid(r *http.Request) error {
	token, err := verifyToken(r)
	log.Println("token:", token)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func Login(r io.Reader) (string, error) {
	var loginUser User
	err := json.NewDecoder(r).Decode(&loginUser)
	if err != nil {
		return "", fmt.Errorf("Invalid json provided")
	}

	if loginUser.Username != user.Username || !checkPassword(loginUser.Password, user.Password) {
		return "", fmt.Errorf("Please provide valid login details")
	}

	return createToken(user.ID)
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth...")
		err := tokenValid(r)
		if err != nil {
			util.RespondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		h.ServeHTTP(w, r)
	})
}
