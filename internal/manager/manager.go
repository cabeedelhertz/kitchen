package manager

import (
	"context"

	kitchenv1 "kitchen/proto/gen/kitchen/v1"
	"kitchen/internal/store"
)

type Manager struct {
	store store.Store
}

func NewManager(store store.Store) *Manager {
	return &Manager{store: store}
}

func (m *Manager) CreatePost(ctx context.Context, req *kitchenv1.CreatePostRequest) (*kitchenv1.CreatePostResponse, error) {
	return m.store.CreatePost(ctx, req)
}

func (m *Manager) GetPost(ctx context.Context, req *kitchenv1.GetPostRequest) (*kitchenv1.GetPostResponse, error) {
	post, err := m.store.GetPost(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &kitchenv1.GetPostResponse{Post: post}, nil
}
