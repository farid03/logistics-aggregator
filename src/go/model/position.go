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

// Position Координаты точки на карте
type Position struct {
	// Уникальный идентификатор записи координат
	Id string `json:"id"`
	// Широта
	Latitude float64 `json:"latitude"`
	// Долгота
	Longitude float64 `json:"longitude"`
}
