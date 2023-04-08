package usecase

import (
	"github.com/Zhoangp/User-Service/config"
	"github.com/Zhoangp/User-Service/internal/model"
	"github.com/Zhoangp/User-Service/internal/repo"
	"github.com/Zhoangp/User-Service/pkg/common"
	"github.com/Zhoangp/User-Service/pkg/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	NewUsers(data *model.Users) error
	FindDataWithCondition(conditions map[string]any) (*model.Users, error)
	UpdateUser(user model.Users, newInformation map[string]any) error
	StoreToken(user *model.Users, token string) (error, *gorm.DB)
}

type userUseCase struct {
	cf       *config.Config
	userRepo *repo.UserRepository
}

func NewUserUseCase(userRepo *repo.UserRepository, cf *config.Config) *userUseCase {
	return &userUseCase{cf, userRepo}
}
func (uc *userUseCase) Register(data *model.Users) error {
	if user, _ := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email}); user != nil {
		return model.ErrEmailExisted
	}
	if err := data.PrepareCreate(); err != nil {
		return err
	}
	if err := uc.userRepo.NewUsers(data); err != nil {
		return err
	}
	return nil
}

func (uc *userUseCase) ChangePassword(data *model.UserChangePassword) error {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email, "password": data.OldPassword})
	if err != nil {
		return model.ErrEmailOrPasswordInvalid
	}
	err = utils.ComparePassword(data.OldPassword, data.NewPass)
	if err == nil {
		return common.NewCustomError(err, "The new password cannot be the same as the old password")
	}
	passHashed, _ := utils.HashPassword(data.NewPass)
	if err := uc.userRepo.UpdateUser(user, map[string]any{"password": passHashed}); err != nil {
		return err
	}
	return nil
}

func (uc *userUseCase) ChangeUser(data *model.Users) error {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email})
	if err != nil {
		return model.ErrEmailOrPasswordInvalid
	}

	if err := uc.userRepo.UpdateUser(user, map[string]any{"email": data.Email,
		"firstName":   data.FirstName,
		"lastName":    data.LastName,
		"phoneNumber": data.Phone,
		"address":     data.Address,
		"picture": data.Avatar,
	}); err != nil {
		return model.ErrCannotUpdateUser
	}
	return nil
}

func (uc *userUseCase) SendToken(email string) error {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": email})
	if err != nil {
		return model.ErrEmailOrPasswordInvalid
	}

	token, err := utils.GenerateToken(utils.TokenPayload{Email: user.Email, Role: user.Role, Password: user.Password}, uc.cf.Service.PasswordTokenExpired, uc.cf.Service.Secret)
	if err != nil {
		return err
	}

	err = utils.SendToken(uc.cf, user.Email, token.AccessToken)
	if err != nil {
		return common.ErrBadRequest(err)
	}
	return nil
}
