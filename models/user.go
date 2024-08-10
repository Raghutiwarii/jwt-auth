package models

type User struct{
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `json:"name"`
	Email string `gorm:"unique" json:"email"`
	Password []byte `json:"-"`
	Address string `json:"address"`
}