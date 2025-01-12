package mysql

import "gorm.io/gorm"

type Tokens struct {
	gorm.Model
	UserID           uint   `json:"userID" gorm:"column:user_id"`
	AccessToken      string `json:"accessToken" gorm:"column:access_token"`
	RefreshToken     string `json:"refreshToken" gorm:"column:refresh_token"`
	RefreshExpiredAt int64  `json:"refreshExpiredAt" gorm:"column:refresh_expired_at"`
}

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

type MetaTables struct {
	gorm.Model
	TableName        string `json:"tableName" gorm:"column:table_name"`
	TableDescription string `json:"tableDescription" gorm:"column:table_description"`
}

type Scenarios struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type Times struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type Types struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type Themes struct {
	gorm.Model
	Name        string `json:"name" gorm:"column:name"`
	Image       string `json:"image" gorm:"column:image"`
	Description string `json:"description" gorm:"column:description"`
}

type UserAuths struct {
	gorm.Model
	Email    string `json:"email" gorm:"column:email"`
	AuthCode string `json:"authCode" gorm:"column:auth_code"`
	Type     string `json:"type" gorm:"column:type"`
}

type FoodImages struct {
	gorm.Model
	Name  string `json:"name" gorm:"column:name"`
	Image string `json:"image" gorm:"column:image"`
}

type Reports struct {
	gorm.Model
	UserID int    `json:"userID" gorm:"column:user_id"`
	Reason string `json:"reason" gorm:"column:reason"`
}

type Nutrients struct {
	gorm.Model
	FoodName     string  `json:"foodName" gorm:"column:food_name"`
	Amount       string  `json:"amount" gorm:"column:amount"`
	Kcal         float64 `json:"kcal" gorm:"column:kcal"`
	Carbohydrate float64 `json:"carbohydrate" gorm:"column:carbohydrate"`
	Protein      float64 `json:"protein" gorm:"column:protein"`
	Fat          float64 `json:"fat" gorm:"column:fat"`
}

type UserTokens struct {
	gorm.Model
	UserID uint   `json:"userID" gorm:"column:user_id"`
	Token  string `json:"token" gorm:"column:token"`
}

type Foods struct {
	gorm.Model
	Name    string `json:"name" gorm:"column:name"`
	ImageID int    `json:"imageID" gorm:"column:image_id"`
}
type CategoryTypes struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name"`
}

type Categories struct {
	gorm.Model
	Name   string `json:"name" gorm:"column:name"`
	TypeID int    `json:"typeID" gorm:"column:type_id"`
}

type FoodCategories struct {
	gorm.Model
	FoodID     int `json:"foodID" gorm:"column:food_id"`
	CategoryID int `json:"categoryID" gorm:"column:category_id"`
}

type FoodHistories struct {
	gorm.Model
	UserID int `json:"userID" gorm:"column:user_id"`
	FoodID int `json:"foodID" gorm:"column:food_id"`
}
