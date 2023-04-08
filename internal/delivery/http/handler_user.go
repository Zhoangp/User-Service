package userhttp

import (
	"context"
	"github.com/Zhoangp/User-Service/internal/model"
	"github.com/Zhoangp/User-Service/pb"
	"github.com/Zhoangp/User-Service/pkg/common"
)

type userHandler struct {
	UC UserUseCase
	pb.UnimplementedUserServiceServer
}

type UserUseCase interface {
	ChangeUser(data *model.Users) error
	ChangePassword(data *model.UserChangePassword) error
	SendToken(email string) error
}

func NewUserHandler(userUC UserUseCase) *userHandler {
	return &userHandler{UC: userUC}
}

func (userHandler *userHandler) ChangePassword(ctx context.Context, request *pb.ChangePasswordRequest) (*pb.ChangePasswordResponse, error) {
	data := model.UserChangePassword{
		Email:       request.Email,
		OldPassword: request.Password,
		NewPass:     request.NewPassword,
	}
	if err := userHandler.UC.ChangePassword(&data); err != nil {
		return &pb.ChangePasswordResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.ChangePasswordResponse{}, nil

}

func (userHandler *userHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	data := model.Users{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.PhoneNumber,
		Address:   req.Address,
	}

	if err := userHandler.UC.ChangeUser(&data); err != nil {
		return &pb.UpdateUserResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.UpdateUserResponse{}, nil

}
func (userHandler *userHandler) GetTokenResetPass(ctx context.Context, req *pb.GetTokenResetPassRequest) (*pb.GetTokenResetPassResponse, error) {

	err := userHandler.UC.SendToken(req.Email)
	if err != nil {
		return &pb.GetTokenResetPassResponse{
			Error: HandleError(err),
		}, nil
	}
	return &pb.GetTokenResetPassResponse{}, nil
}
func HandleError(err error) *pb.ErrorResponse {
	if errors, ok := err.(*common.AppError); ok {
		return &pb.ErrorResponse{
			Code:    int64(errors.StatusCode),
			Message: errors.Message,
		}
	}
	appErr := common.ErrInternal(err.(error))
	return &pb.ErrorResponse{
		Code:    int64(appErr.StatusCode),
		Message: appErr.Message,
	}
}
