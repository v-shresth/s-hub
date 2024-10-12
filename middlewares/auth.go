package middlewares

import (
	handler2 "cms/handler"
	"cms/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"time"
)

var exemptedMethods = map[string]struct{}{
	"/pb.UserService/Register": {},
	"/pb.UserService/Login":    {},
}

func AuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	newCtx := ctx
	if _, ok := exemptedMethods[info.FullMethod]; !ok {
		server, ok := info.Server.(*handler2.Server)
		if !ok {
			return nil, status.Error(codes.Internal, "Server Error")
		}
		userId, sessionId, err := authorize(ctx, *server)
		if err != nil {
			return nil, err
		}

		md, ok := metadata.FromIncomingContext(newCtx)
		if ok {
			md.Append(utils.AuthedUserId, fmt.Sprintf("%d", userId))
			md.Append(utils.AuthedSessionId, fmt.Sprintf("%d", sessionId))
		}
		newCtx = metadata.NewIncomingContext(newCtx, md)
	}
	h, err := handler(newCtx, req)
	return h, err
}

func authorize(ctx context.Context, server handler2.Server) (uint, uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, 0, status.Errorf(codes.Internal, "retrieving metadata failed")
	}

	tokens := md.Get("authorization")
	if len(tokens) == 0 {
		return 0, 0, status.Error(codes.Unauthenticated, "token not sent in header")
	}

	claims, err := server.TokenMaker.ValidateToken(tokens[0])
	if err != nil {
		return 0, 0, err
	}

	sessionId := claims.SessionId
	session, err := server.UserSvc.FindUserSession(ctx, sessionId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, 0, status.Error(codes.NotFound, "session not found")
		}
		return 0, 0, status.Error(codes.Internal, "something went wrong. Please try again in some time")
	}

	currentTime := time.Now()
	if session.EndedAt.Before(currentTime) {
		return 0, 0, status.Error(codes.Unauthenticated, "You are not logged in")
	}

	return session.UserId, session.ID, nil
}
