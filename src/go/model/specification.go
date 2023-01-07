package model

/*
 * logistics-aggregator
 *
 * Aгрегатор логистических заказов
 *
 * API version: 1.0
 * Contact: f.kurbanov120303@yandex.ru
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

type Specification struct {
	// Уникальный идентификатор автомобиля
	Id string `json:"id"`
	// Длина прицепа автомобиля
	Length float64 `json:"length"`
	// Высота прицепа автомобиля
	Height float64 `json:"height"`
	// Ширина прицепа автомобиля
	Width float64 `json:"width"`
	// Цвет прицепа автомобиля
	Color string `json:"color"`

	BodyType *TrailerType `json:"bodyType"`

	LoadingPlaces *LoadingPlaces `json:"loadingPlaces"`
}
