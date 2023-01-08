package handler

import (
	"crypto/md5"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"log"
	storage "logistics-aggregator/src/go"
	"logistics-aggregator/src/go/model"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func findByUserName(username string) (user model.User, exist bool) {
	result := storage.PG.Where("Username = ?", username).First(&user)
	if user.Username == "" || result.Error != nil {
		log.Printf("findByUserName (user not found): username{%s}, storage_err_message{%s}\n", username, result.Error)
		return
	}
	exist = true
	return
}

// User.ID не требуется
func authenticateUser(user model.User) http.Cookie {
	time64 := time.Now().Unix()
	timeInt := strconv.FormatInt(time64, 10)
	token := user.Username + user.Password + timeInt
	hashToken := md5.Sum([]byte(token))
	hashedToken := hex.EncodeToString(hashToken[:])
	storage.Cache[hashedToken] = user
	livingTime := 60 * time.Minute // кука будет жить 1 час
	expiration := time.Now().Add(livingTime)

	return http.Cookie{Name: "token", Value: url.QueryEscape(hashedToken), Expires: expiration}
}

func SignIn(authUser model.AuthBody) (http.Cookie, bool) {
	user, exist := findByUserName(authUser.Username)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password)); exist && err == nil {
		return authenticateUser(user), true
	}

	return http.Cookie{}, false
}

func SignUp(signUpUser model.User) (http.Cookie, bool) {
	_, exist := findByUserName(signUpUser.Username)
	if exist {
		log.Println("SingUp: (user with same username already exists) ", signUpUser)
		return http.Cookie{}, false
	}

	result := storage.PG.Create(&signUpUser)
	if result.Error != nil {
		log.Println("SingUp: Cannot create user: ", signUpUser, "error: ", result.Error)
		return http.Cookie{}, false
	}

	return authenticateUser(signUpUser), true
}
