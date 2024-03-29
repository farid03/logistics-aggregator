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

// LoadingPlaces : Места погрузки прицепа автомобиля
type LoadingPlaces string

// List of LoadingPlaces
const (
	BACK  LoadingPlaces = "BACK"
	LEFT  LoadingPlaces = "LEFT"
	RIGHT LoadingPlaces = "RIGHT"
	TOP   LoadingPlaces = "TOP"
)
