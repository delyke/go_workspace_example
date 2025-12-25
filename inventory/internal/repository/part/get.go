package part

import (
	"context"

	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
)

func (r *repository) GetPart(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	part, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}
	return lo.ToPtr(converter.RepoPartToModel(*part)), nil
}
