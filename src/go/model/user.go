package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

/*
 * logistics-aggregator
 *
 * Aгрегатор логистических заказов
 *
 * API version: 1.0
 * Contact: f.kurbanov120303@yandex.ru
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

// User Пользователь системы
type User struct {
	// Уникальный идентификатор пользователя
	ID uint32 `json:"id" gorm:"primary_key;auto_increment"`
	// Уникальный никнейм пользователя
	Username string `json:"username" gorm:"size:100;not null;unique"`
	// Пароль пользователя
	Password string `json:"password" gorm:"size:255;not null"`
	// Имя пользователя
	Name string `json:"name" gorm:"size:255;not null"`
	// Фамилия пользователя
	Surname string `json:"surname" gorm:"size:255;not null"`

	UserState UserState `json:"userState,omitempty" gorm:"not null"`

	UserType UserType `json:"userType" gorm:"not null"`
}

func (u *User) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		return json.Unmarshal([]byte(v), u)
	case []byte:
		return json.Unmarshal(v, u)
	}
	return fmt.Errorf("cannot convert %T to My struct", src)
}

//nolint:hugeParam
func (u User) Value() (driver.Value, error) {
	return json.Marshal(u)
}
