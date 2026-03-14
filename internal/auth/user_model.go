package auth

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"unique;not null"`
	Name     string `json:"name"`
	Password string `json:"password" gorm:"not null"` // hashed

}
