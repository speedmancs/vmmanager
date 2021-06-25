package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

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
