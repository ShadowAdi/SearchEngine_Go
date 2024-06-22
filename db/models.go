package db

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        string     `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Email     string     `gorm:"unique" json:"email"`
	Password  string     `json:"-"`
	IsAdmin   bool       `gorm:"default:false" json:"isAdmin"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u *User) CreateAdmin() error {
	user := User{
		Email:    "your email",
		Password: "your password",
		IsAdmin:  true,
	}
	newPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return errors.New("Error Hashing Password")
	}
	user.Password = string(newPassword)
	if err := DBConn.Create(&user).Error; err != nil {
		return errors.New("Error Creating User")
	}
	return nil

}

func (u *User) LoginAdmin(email string, password string) (*User, error) {
	if err := DBConn.Where("email = ? AND is_admin = ?", email, true).First(&u).Error; err != nil {
		return nil, errors.New("User Not Found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("Invalid Credentials")
	}
	return u, nil
}
