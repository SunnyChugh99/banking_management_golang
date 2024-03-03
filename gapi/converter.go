package gapi

import (
	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func convertUser(user db.User)(*pb.User){
	return &pb.User{
		Username: user.Username,
		FullName: user.FullName,          
		Email: user.Email, 
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}
}