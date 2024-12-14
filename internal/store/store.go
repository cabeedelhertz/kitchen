package store

import (
	"context"

	kitchenv1 "kitchen/proto/gen/kitchen/v1"
)

type Store interface {
	CreatePost(ctx context.Context, req *kitchenv1.CreatePostRequest) (*kitchenv1.CreatePostResponse, error)
	GetPost(ctx context.Context, id string) (*kitchenv1.Post, error)
}
