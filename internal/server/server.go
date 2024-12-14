package server

import (
	"context"
	"kitchen/internal/manager"
	kitchenv1 "kitchen/proto/gen/kitchen/v1"
	kitchenv1connect "kitchen/proto/gen/kitchen/v1/kitchenv1connect"

	"connectrpc.com/connect"
)

var _ kitchenv1connect.KitchenServiceHandler = &Server{}

type Server struct {
	manager manager.Manager
}

func NewServer(manager manager.Manager) *Server {
	return &Server{manager: manager}
}

func (s *Server) CreatePost(ctx context.Context, req *connect.Request[kitchenv1.CreatePostRequest]) (*connect.Response[kitchenv1.CreatePostResponse], error) {
	resp, err := s.manager.CreatePost(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

func (s *Server) GetPost(ctx context.Context, req *connect.Request[kitchenv1.GetPostRequest]) (*connect.Response[kitchenv1.GetPostResponse], error) {
	resp, err := s.manager.GetPost(ctx, req.Msg)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}
