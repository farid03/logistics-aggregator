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
	"strconv"
	"strings"
)

var header = filepath.Join("resources", "components", "header.html")
var apikey = filepath.Join("resources", "components", "yandex-maps-api.key")
var main = filepath.Join("resources", "main.html")
var index = filepath.Join("resources", "index.html")
var profile = filepath.Join("resources", "profile.html")
var advert = filepath.Join("resources", "advert.html")
var registration = filepath.Join("resources", "registration.html")

func loadHtmlTemplate(w http.ResponseWriter, data any, path ...string) {
	tmpl, err := template.ParseFiles(path...)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	loadHtmlTemplate(w, nil, index, header)
}

func Main(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var orders []model.Order
	storage.PG.Preload("From").Preload("To").Find(&orders)
	fmt.Println(orders)

	loadHtmlTemplate(w, orders, main, header, apikey)
}

func AuthGet(w http.ResponseWriter, r *http.Request) { // возможно стоит перенести в /profile
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	authToken, err := readCookie("authToken", r)
	userID := storage.Cache[authToken].ID
	type UserData struct {
		Cars   []model.Car
		Orders []model.Order
	}
	var data UserData
	storage.PG.Preload("Position").Where("owner_id = ?", strconv.Itoa(int(userID))).Find(&data.Cars)
	storage.PG.Preload("From").Preload("To").Where("executor_id_refer = ?", strconv.Itoa(int(userID))).Find(&data.Orders)
	fmt.Println(data)

	loadHtmlTemplate(w, data, profile, header, apikey)
}

func AuthPost(w http.ResponseWriter, r *http.Request) {
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
	if err == nil {
		delete(storage.Cache, authToken)
	}
	Index(w, r)
}

func AddAdvertGet(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}

	loadHtmlTemplate(w, nil, advert, header, apikey)
}

func AddAdvertPost(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	authToken, _ := readCookie("authToken", r)
	userID := storage.Cache[authToken].ID
	var errs = make([]error, 8)
	var price, startLatitude, startLongitude, endLatitude, endLongitude, length, height, width float64
	price, errs[0] = strconv.ParseFloat(r.PostFormValue("price"), 64)
	startLatitude, errs[1] = strconv.ParseFloat(r.PostFormValue("from-latitude"), 64)
	startLongitude, errs[2] = strconv.ParseFloat(r.PostFormValue("from-longitude"), 64)
	endLatitude, errs[3] = strconv.ParseFloat(r.PostFormValue("to-latitude"), 64)
	endLongitude, errs[4] = strconv.ParseFloat(r.PostFormValue("to-longitude"), 64)
	length, errs[5] = strconv.ParseFloat(r.PostFormValue("length"), 64)
	height, errs[6] = strconv.ParseFloat(r.PostFormValue("height"), 64)
	width, errs[7] = strconv.ParseFloat(r.PostFormValue("width"), 64)

	for i, e := range errs {
		if e != nil {
			log.Println("AddAdvert: Cannot create parse line: ", i, "error: ", e)
			w.Write([]byte(fmt.Sprintln(e)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var start = model.Position{
		Latitude:  startLatitude,
		Longitude: startLongitude,
	}
	result := storage.PG.Create(&start)

	var end = model.Position{
		Latitude:  endLatitude,
		Longitude: endLongitude,
	}
	result = storage.PG.Create(&end)

	var specification = model.Specification{
		Length:        length,
		Height:        height,
		Width:         width,
		Color:         "nop",                                                 // FIXME нет валидации
		BodyType:      model.TrailerType(r.PostFormValue("trailerType")),     //  |
		LoadingPlaces: model.LoadingPlaces(r.PostFormValue("loadingPlaces")), //  |_
	}
	result = storage.PG.Create(&specification)

	var order = model.Order{
		OwnerID:       userID,
		Title:         r.PostFormValue("title"),
		Description:   r.PostFormValue("description"),
		Price:         price,
		State:         model.REQUESTED,
		From:          start,
		To:            end,
		Specification: specification,
	}

	result = storage.PG.Create(&order)
	if result.Error != nil {
		log.Println("AddAdvert: Cannot create order: ", order, "error: ", result.Error)
		w.Write([]byte(fmt.Sprintln(result.Error)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	loadHtmlTemplate(w, nil, advert, header, apikey) // FIXME переделать на нормальны респонс с ajax
}

func AddCarPost(w http.ResponseWriter, r *http.Request) {
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	authToken, err := readCookie("authToken", r)
	userID := storage.Cache[authToken].ID
	var length, height, width, startLatitude, startLongitude float64
	errs := make([]error, 5)
	length, errs[0] = strconv.ParseFloat(r.PostFormValue("length"), 64)
	height, errs[1] = strconv.ParseFloat(r.PostFormValue("height"), 64)
	width, errs[2] = strconv.ParseFloat(r.PostFormValue("width"), 64)
	startLatitude = 0 // изменение значений должно произойти потом
	startLongitude = 0

	for i, e := range errs {
		if e != nil {
			log.Println("AddCarPost: Cannot create parse line: ", i, "error: ", e)
			w.Write([]byte(fmt.Sprintln(e)))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var start = model.Position{
		Latitude:  startLatitude,
		Longitude: startLongitude,
	}
	result := storage.PG.Create(&start)
	if result.Error != nil {
		log.Println("AddCar: Cannot create position: ", start, "error: ", result.Error)
		w.Write([]byte(fmt.Sprintln(result.Error)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var specification = model.Specification{
		Length:        length,
		Height:        height,
		Width:         width,
		Color:         r.PostFormValue("color"),
		BodyType:      model.TrailerType(r.PostFormValue("trailerType")),
		LoadingPlaces: model.LoadingPlaces(r.PostFormValue("loadingPlaces")),
	}
	result = storage.PG.Create(&specification)
	if result.Error != nil {
		log.Println("AddCar: Cannot create specification: ", specification, "error: ", result.Error)
		w.Write([]byte(fmt.Sprintln(result.Error)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var car = model.Car{
		OwnerID:       userID,
		LicensePlate:  r.PostFormValue("license"),
		Position:      start,
		Specification: specification,
	}
	result = storage.PG.Create(&car)
	if result.Error != nil {
		log.Println("AddCar: Cannot create car: ", car, "error: ", result.Error)
		w.Write([]byte(fmt.Sprintln(result.Error)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	AuthGet(w, r)
}

func OrderIdDelete(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func OrderIdGet(w http.ResponseWriter, r *http.Request) { // этой задачей должен заниматься не GET запрос
	err := checkAuth(r) // проверяем авторизован ли пользователя
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusForbidden)
		return
	}
	authToken, err := readCookie("authToken", r)
	userID := storage.Cache[authToken].ID
	// логика взятия заказа на исполнение
	var order model.Order
	id := strings.TrimPrefix(r.URL.Path, "/order/") // не самый лучший способ получить параметр id из запроса
	result := storage.PG.Where("ID = ?", id).First(&order)
	if result.Error != nil {
		log.Println("OrderIdGet: Cannot create get order with id: ", id, "error: ", result.Error)
		w.Write([]byte(fmt.Sprintln(result.Error)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if order.State != model.REQUESTED {
		w.Write([]byte("order already executed"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	order.ExecutorIDRefer = userID
	order.State = model.IN_PROGRESS // FIXME нет валидации по конкретному автомобилю исполнителя
	result = storage.PG.Save(&order)
	if result.Error != nil {
		log.Println("OrderIdGet: Cannot save order: ", order, "error: ", result.Error)
		w.Write([]byte(fmt.Sprintln(result.Error)))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	Main(w, r)
}

func RegistrationGet(w http.ResponseWriter, r *http.Request) {
	loadHtmlTemplate(w, nil, registration, header)
}

func RegistrationPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	pass, _ := bcrypt.GenerateFromPassword([]byte(r.PostFormValue("password")), bcrypt.DefaultCost)

	user := model.User{
		Name:      r.PostFormValue("name"),
		Surname:   r.PostFormValue("surname"),
		Username:  r.PostFormValue("username"), // FIXME нет валидации на уникальность полученного значения
		Password:  string(pass),
		UserState: model.ACTIVE,
		UserType:  model.STANDARD, // админом зарегистрироваться через форму нельзя
	}

	cookie, err := SignUp(user)
	if err != nil {
		w.Write([]byte(fmt.Sprintln(err)))
		w.WriteHeader(http.StatusBadRequest)
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}
