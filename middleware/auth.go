package middleware

import (
	"log"

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

func HashAndSalt(password string) string {
	pwd := []byte(password)
	// pwd []byte = ([]byte)password
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}
