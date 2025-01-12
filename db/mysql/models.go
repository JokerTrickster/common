package mysql

import "gorm.io/gorm"

// Tokens 테이블
type Tokens struct {
	gorm.Model
	UserID           uint   `json:"userID" gorm:"column:user_id"`
	AccessToken      string `json:"accessToken" gorm:"column:access_token"`
	RefreshToken     string `json:"refreshToken" gorm:"column:refresh_token"`
	RefreshExpiredAt int64  `json:"refreshExpiredAt" gorm:"column:refresh_expired_at"`
}

// Users 테이블
type Users struct {
	gorm.Model
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
	Birth    string `json:"birth" gorm:"column:birth"`
	Name     string `json:"name" gorm:"column:name"`
	Sex      string `json:"sex" gorm:"column:sex"`
	Provider string `json:"provider" gorm:"column:provider"`
	Push     *bool  `json:"push" gorm:"column:push"`
	Image    string `json:"image" gorm:"column:image"`
}

// Foods 테이블
type Foods struct {
	gorm.Model
	Name    string `json:"name" gorm:"column:name"`
	ImageID int    `json:"imageID" gorm:"column:image_id"`
}

// Categories 테이블
type Categories struct {
	gorm.Model
	Name   string `json:"name" gorm:"column:name"`
	TypeID int    `json:"typeID" gorm:"column:type_id"`
}

// FoodCategories 테이블
type FoodCategories struct {
	gorm.Model
	FoodID     int `json:"foodID" gorm:"column:food_id"`
	CategoryID int `json:"categoryID" gorm:"column:category_id"`
}

// Nutrients 테이블
type Nutrients struct {
	gorm.Model
	FoodName     string  `json:"foodName" gorm:"column:food_name"`
	Amount       string  `json:"amount" gorm:"column:amount"`
	Kcal         float64 `json:"kcal" gorm:"column:kcal"`
	Carbohydrate float64 `json:"carbohydrate" gorm:"column:carbohydrate"`
	Protein      float64 `json:"protein" gorm:"column:protein"`
	Fat          float64 `json:"fat" gorm:"column:fat"`
}
