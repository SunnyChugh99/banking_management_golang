package gapi

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/SunnyChugh99/banking_management_golang/db/sqlc"
	"github.com/SunnyChugh99/banking_management_golang/pb"
	"github.com/SunnyChugh99/banking_management_golang/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)


func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err !=nil{
		if err == sql.ErrNoRows{
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Unable to Login user : %s", err)

	}	

	err = util.CheckPassword(req.Password, user.HashedPassword)
		
	if err!=nil{
		return nil, status.Errorf(codes.NotFound, "Incorrect password : %s", err)

	}


	accessToken,accessPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.AccessTokenDuration,
	)
	if err!=nil{
		return nil, status.Errorf(codes.Unauthenticated, "Failed to create access token : %s", err)

	}
	

	refreshToken,refreshPayload, err := server.tokenMaker.CreateToken(
		user.Username,
		server.config.RefreshTokenDuration,
	)

	if err!=nil{
		return nil, status.Errorf(codes.Unauthenticated, "Failed to create refresh token : %s", err)

	}
	
	mtd := server.extractMetadata(ctx)

	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID: refreshPayload.ID,
		Username: user.Username,     
		RefreshToken: refreshToken,
		UserAgent: mtd.UserAgent,   
		ClientIp: mtd.ClientIp,     
		IsBlocked: false,    
		ExpiresAt: refreshPayload.ExpiredAt,    
	})

	if err!=nil{
		return nil, status.Errorf(codes.Unauthenticated, "Unable to Login user : %s", err)

	}
	fmt.Println(session)

	rsp := &pb.LoginUserResponse{
	User: convertUser(user),
	SessionId: session.ID.String(),
	AccessToken: accessToken, 
	RefreshToken: refreshToken,
	AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}