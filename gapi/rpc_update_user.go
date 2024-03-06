package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"github.com/SunnyChugh99/banking_management_golang/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)



func (server *Server) UpdateUser(ctx context.Context,req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	//1- authorization.go   ==> input - context , o/tp ==> authpayload, err
	authPayload, err := server.authorizeUser(ctx)

	if err!=nil{
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}
	

	if authPayload.Username != req.GetUsername(){
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info : %s", err)

	}

	args := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid: req.FullName != nil,
		}, 
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid: req.Email != nil,
		}, 
	}

	if req.Password != nil{
		hashedPassword, err := util.HashPassword(req.GetPassword())
		fmt.Println("here-3")
	
		if err!=nil{
			return nil, status.Errorf(codes.Internal, "failed to hash password")
		}
		args.HashedPassword =  sql.NullString{
			String: hashedPassword,
			Valid: true,
		}
		args.PasswordChangedAt =  sql.NullTime{
			Time: time.Now(),
			Valid: true,
		}

	}

	fmt.Println("here-4")

	user, err := server.store.UpdateUser(ctx, args)
	if err!=nil{
		fmt.Println("in error")
		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound, "User not found : %s", err)
		}
		return nil, status.Errorf(codes.Internal, "Unable to create user : %s", err)
	}

	fmt.Println("here-5")


	rsp := &pb.UpdateUserResponse{User: convertUser(user)}

	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if req.Password != nil{
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	
	}
	if req.FullName != nil{
		if err := val.ValidateFullName(req.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	if req.Email != nil{
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}
	return violations
}