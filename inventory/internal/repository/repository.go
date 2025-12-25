package repository

import (
	"context"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, UUID string) (*model.Part, error)
	ListParts(ctx context.Context, filters *model.PartsFilter) ([]*model.Part, error)
	CreatePart(ctx context.Context, part *model.Part) error
}
