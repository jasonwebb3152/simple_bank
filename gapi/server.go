package gapi

import (
	"fmt"

	db "github.com/jasonwebb3152/simplebank/db/sqlc"
	"github.com/jasonwebb3152/simplebank/pb"
	"github.com/jasonwebb3152/simplebank/token"
	"github.com/jasonwebb3152/simplebank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// NewServer creates a new gRPC server (gRPC doesn't use routing)
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
