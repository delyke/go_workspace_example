package part

import (
	"context"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filters *model.PartsFilter) ([]*model.Part, error) {
	parts, err := s.partRepository.ListParts(ctx, filters)
	if err != nil {
		return nil, err
	}
	return parts, nil
}
