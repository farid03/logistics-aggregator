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

// OrderState : Состояние заказа
type OrderState string

// List of OrderState
const (
	REQUESTED   OrderState = "REQUESTED"
	IN_PROGRESS OrderState = "IN_PROGRESS"
	COMPLETED   OrderState = "COMPLETED"
)
