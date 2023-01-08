package handler

import (
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"logistics-aggregator/src/go/model"
	"net/http"
	"path/filepath"
	"strings"
)

func loadHtmlTemplate(w http.ResponseWriter, path string) {
	tmpl, err := template.ParseFiles(path)
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
	// здесь юзеру отдаем страницу регистрации
	loadHtmlTemplate(w, filepath.Join("resources", "static", "index.html"))

	//fs := http.FileServer(http.Dir("logistics-aggregator/resources/static"))
	//http.Handle("/static/", http.StripPrefix("/static/", fs))
	// нужно будет чтобы сделать доступными файлы с css стилями
	//fmt.Fprintf(w, "hehhehe")
}

func AuthPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	username := strings.TrimSpace(r.PostFormValue("username"))
	password := strings.TrimSpace(r.PostFormValue("password"))
	cookie, res := SignIn(model.AuthBody{Username: username, Password: password})
	if !res {
		w.WriteHeader(http.StatusBadRequest)
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
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
	loadHtmlTemplate(w, filepath.Join("resources", "static", "registration.html"))
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
	case "Заказчик":
		userType = model.CUSTOMER
	case "Грузоперевозчик":
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

	cookie, res := SignUp(user)
	if !res {
		w.WriteHeader(http.StatusBadRequest)
	}

	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
}
