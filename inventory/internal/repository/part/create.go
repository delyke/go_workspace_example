package part

import (
	"context"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
)

func (r *repository) CreatePart(_ context.Context, part *model.Part) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.parts[part.Uuid]; exists {
		return model.ErrPartAlreadyExists
	}
	r.parts[part.Uuid] = converter.ModelPartToRepo(part)
	return nil
}
