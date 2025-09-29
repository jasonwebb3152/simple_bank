package gapi

import (
	"context"
	"database/sql"
	"time"

	db "github.com/jasonwebb3152/simplebank/db/sqlc"
	"github.com/jasonwebb3152/simplebank/pb"
	"github.com/jasonwebb3152/simplebank/util"
	"github.com/jasonwebb3152/simplebank/val"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: add authorization layer
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, InvalidArgumentError(violations)
	}

	if authPayload.Username != req.GetUsername() {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			Valid:  req.FullName != nil,
			String: req.GetFullName(),
		},
		Email: sql.NullString{
			Valid:  req.Email != nil,
			String: req.GetEmail(),
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}

		arg.HashedPassword = sql.NullString{
			Valid:  true,
			String: hashedPassword,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Valid: true,
			Time:  time.Now(),
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "case_not_found":
				return nil, status.Errorf(codes.NotFound, "username not found: %s", err)
			case "unique_violation": // unique username and email stuff
				return nil, status.Errorf(codes.AlreadyExists, "email already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {

	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if email := req.GetEmail(); email != "" {
		if err := val.ValidateEmail(email); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	if password := req.GetPassword(); password != "" {
		if err := val.ValidatePassword(password); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if fullName := req.GetFullName(); fullName != "" {
		if err := val.ValidateFullName(fullName); err != nil {
			violations = append(violations, fieldViolation("fullName", err))
		}
	}
	return
}
