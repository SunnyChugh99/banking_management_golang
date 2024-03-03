package gapi

import (
	"context"
	"fmt"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/lib/pq"
	"github.com/lib/val"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (server *Server) CreateUser(ctx context.Context,req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	hashedPassword, err := util.HashPassword(req.GetPassword())
	fmt.Println("here-3")

	if err!=nil{
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}


	args := db.CreateUserParams{
		Username: req.Username,
		HashedPassword: hashedPassword,
		FullName: req.FullName,
		Email: req.Email,
	}

	fmt.Println("here-4")

	user, err := server.store.CreateUser(ctx, args)
	if err!=nil{
		fmt.Println("in error")

		if pqErr,ok := err.(*pq.Error); ok{
			switch pqErr.Code.Name(){
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "User already exists")
			}
		}
		return nil, status.Errorf(codes.Internal, "Unable to create user : %s", err)

	}

	fmt.Println("here-5")
	// userResponse := newUserResponse(user)
	// ctx.JSON(http.StatusOK, userResponse)


	rsp := &pb.CreateUserResponse{User: convertUser(user)}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation){
	if err:= val.ValidateUsername(req.GetUsername()); err!=nil{


	}




}