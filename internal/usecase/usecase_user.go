package usecase

import (
	"github.com/Zhoangp/User-Service/config"
	"github.com/Zhoangp/User-Service/internal/model"
	"github.com/Zhoangp/User-Service/internal/repo"
	"github.com/Zhoangp/User-Service/pb"
	"github.com/Zhoangp/User-Service/pkg/common"
	"github.com/Zhoangp/User-Service/pkg/utils"
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type UserRepository interface {
	NewUsers(data *model.Users) error
	FindDataWithCondition(conditions map[string]any) (*model.Users, error)
	UpdateUser(user model.Users, newInformation map[string]any) error
	StoreToken(user *model.Users, token string) (error, *gorm.DB)
	NewInstructor(user *model.Users, intructor *model.Instructor) error
	GetInstructor(condition map[string]any) (*model.Instructor, error)
}

type userUseCase struct {
	cf       *config.Config
	h        *utils.Hasher
	userRepo *repo.UserRepository
}

func NewUserUseCase(userRepo *repo.UserRepository, cf *config.Config, h *utils.Hasher) *userUseCase {
	return &userUseCase{cf: cf, userRepo: userRepo, h: h}
}
func (uc *userUseCase) GetInstructor(id, key string) (*pb.GetInstructorInformationResponse, error) {
	idDecode, err := uc.h.Decode(id)
	if err != nil {
		return nil, err
	}
	var instructor *model.Instructor
	if key == "user" {
		instructor, err = uc.userRepo.GetInstructor(map[string]any{"user_id": idDecode})
		if err != nil {
			return nil, err
		}
	} else {
		instructor, err = uc.userRepo.GetInstructor(map[string]any{"id": idDecode})
		if err != nil {
			return nil, err
		}
	}
	return &pb.GetInstructorInformationResponse{
		Information: &pb.Instructor{
			Id:           uc.h.Encode(instructor.Id),
			FirstName:    instructor.User.FirstName,
			LastName:     instructor.User.LastName,
			Email:        instructor.User.Email,
			Website:      instructor.Website,
			Linkedin:     instructor.LinkedIn,
			Youtube:      instructor.Youtube,
			Bio:          instructor.Bio,
			UserId:       uc.h.Encode(instructor.UserId),
			NumStudents:  instructor.NumStudents,
			NumReviews:   instructor.NumReviews,
			TotalCourses: instructor.TotalCourses,
			Avt: &pb.Picture{
				Url:    instructor.User.Avatar.Url,
				Width:  instructor.User.Avatar.Width,
				Height: instructor.User.Avatar.Height,
			},
		},
	}, nil
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
	}); err != nil {
		return model.ErrCannotUpdateUser
	}
	return nil
}

func (uc *userUseCase) ChangeAvatar(data *model.Users) error {
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"email": data.Email})
	if err != nil {
		return model.ErrEmailOrPasswordInvalid
	}

	if err := uc.userRepo.UpdateUser(user, map[string]any{"email": data.Email,
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

func (uc *userUseCase) NewInstructor(data *model.Instructor, userId string) error {
	userIdDecoded, err := uc.h.Decode(userId)
	user, err := uc.userRepo.FindDataWithCondition(map[string]any{"id": userIdDecoded})
	if err != nil {
		return model.ErrEmailOrPasswordInvalid
	}

	data.UserId = userIdDecoded
	if _, err = govalidator.ValidateStruct(data); err != nil {
		return common.NewCustomError(err, err.Error())
	}
	if err := uc.userRepo.NewInstructor(user, data); err != nil {
		return model.ErrCannotCreateInstructor
	}
	return nil
}
