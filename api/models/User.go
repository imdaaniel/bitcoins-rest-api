package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/imdaaniel/bitcoins-rest-api/api/utils/date"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Name        string    `gorm:"size:60;not null;" json:"name"`
	Email       string    `gorm:"size:100;not null;" json:"email"`
	Password    string    `gorm:"size:100;not null;" json:"password"`
	DateOfBirth string    `gorm:"type:date;not null" json:"dateofbirth"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdat"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedat"`
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.DateOfBirth = html.EscapeString(strings.TrimSpace(u.DateOfBirth))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	if strings.ToLower(action) == "login" {
		if u.Password == "" {
			return errors.New("Required Password")
		}
		if u.Email == "" {
			return errors.New("Required Email")
		} else if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid Email")
		}

		return nil
	}

	if u.Name == "" {
		return errors.New("Required Name")
	}
	if u.Password == "" {
		return errors.New("Required Password")
	}
	if u.DateOfBirth == "" {
		return errors.New("Required Date of Birth")
	} else if date.ValiDate(u.DateOfBirth) == false {
		return errors.New("Invalid Date of Birth")
	}
	if u.Email == "" {
		return errors.New("Required Email")
	} else if err := checkmail.ValidateFormat(u.Email); err != nil {
		return errors.New("Invalid Email")
	}

	return nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error

	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindUsers(db *gorm.DB) (*[]User, error) {
	var err error

	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error

	if err != nil {
		return &[]User{}, err
	}

	return &users, err
}
