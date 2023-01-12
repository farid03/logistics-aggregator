package handler

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	storage "logistics-aggregator/src/go"
	"logistics-aggregator/src/go/model"
	"net/http"
	"path/filepath"
	"strings"
)

var header = filepath.Join("resources", "components", "header.html")
var main = filepath.Join("resources", "main.html")
var index = filepath.Join("resources", "index.html")
var profile = filepath.Join("resources", "profile.html")
var advert = filepath.Join("resources", "advert.html")
var registration = filepath.Join("resources", "registration.html")

func loadHtmlTemplate(w http.ResponseWriter, path ...string) {
	tmpl, err := template.ParseFiles(path...)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	loadHtmlTemplate(w, index, header)
}

func Main(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var users []model.User
	storage.PG.Find(&users)
	fmt.Println(users)

	//создаем html-шаблон
	tmpl, err := template.ParseFiles(main, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func AuthGet(w http.ResponseWriter, r *http.Request) {
	loadHtmlTemplate(w, profile, header)
}

func AuthPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	username := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))
	cookie, err := SignIn(model.AuthBody{Username: username, Password: password})
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusBadRequest)

	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

func LogoutGet(w http.ResponseWriter, r *http.Request) {
	authToken, err := readCookie("authToken", r)
	if err != nil {
		delete(storage.Cache, authToken)
	}
	Index(w, r)
}

func AddAdvertGet(w http.ResponseWriter, r *http.Request) {
	loadHtmlTemplate(w, advert, header)
}

func AddAdvertPost(w http.ResponseWriter, r *http.Request) {

}

func AddCarPost(w http.ResponseWriter, r *http.Request) {

}

func OrderIdDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func OrderIdGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("aboba"))
	w.WriteHeader(http.StatusOK)
}

func RegistrationGet(w http.ResponseWriter, r *http.Request) {
	loadHtmlTemplate(w, registration, header)
}

func RegistrationPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	pass, _ := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("password")), bcrypt.DefaultCost)
	var userType model.UserType
	switch r.PostFormValue("userType") {
	case "Customer":
		userType = model.CUSTOMER
	case "Cargo carrier":
		userType = model.EXECUTOR
	default: // админом зарегистрироваться через форму нельзя
		log.Println("RegistrationPost: Invalid UserType in request!")
		w.WriteHeader(http.StatusBadRequest)
	}

	user := model.User{
		Name:      r.PostFormValue("name"),
		Surname:   r.PostFormValue("surname"),
		Username:  r.PostFormValue("username"),
		Password:  string(pass),
		UserState: model.ACTIVE,
		UserType:  userType,
	}

	cookie, err := SignUp(user)
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusBadRequest)
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
