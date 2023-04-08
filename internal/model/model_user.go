package model

import (
	"github.com/Zhoangp/User-Service/pkg/common"
	"github.com/Zhoangp/User-Service/pkg/utils"
)

type Users struct {
	common.SQLModel
	Email        string          `json:"email" gorm:"column:email"`
	Password     string          `json:"password" gorm:"column:password"`
	FirstName    string          `json:"firstName" gorm:"column:firstName"`
	LastName     string          `json:"lastName" gorm:"column:lastName"`
	Phone        string          `json:"phone" gorm:"column:phoneNumber"`
	Role         string          `json:"role" gorm:"column:role"`
	Address      string          `json:"address" gorm:"column:address"`
	IsInstructor bool `json:"is_instructor"`
	Avatar       *common.Image   `json:"avatar" `
}
type UserRegister struct {
	Credential   string `json:"credential"`
	FirstName    string `json:"firstName" gorm:"column:firstName"`
	LastName     string `json:"lastName" gorm:"column:lastName"`
	Phone        string `json:"phoneNumber" gorm:"column:phoneNumber"`
	Role         string `json:"role" gorm:"column:role"`
	Address      string `json:"address"`
	IsInstructor int    `json:"is_instructor"`
}
type UserLogin struct {
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
}
type UserChangePassword struct {
	Email       string `json:"email" gorm:"column:email"`
	OldPassword string `json:"password" gorm:"column:password"`
	NewPass     string `json:"newpassword"`
}

func (Users) TableName() string {
	return "Users"
}

func (u *Users) GetUserId() int {
	return u.Id
}

func (u *Users) GetUserEmail() string {
	return u.Email
}

func (u *Users) GetUserRole() string {
	return u.Role
}
func (u *Users) PrepareCreate() error {
	passHashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = passHashed
	u.Role = "guest"
	return nil
}
