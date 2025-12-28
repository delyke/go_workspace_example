package part

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
)

func (r *repository) CreatePart(ctx context.Context, part *model.Part) error {
	rModel := *converter.ModelPartToRepo(part)
	if rModel.CreatedAt.IsZero() {
		rModel.CreatedAt = lo.ToPtr(time.Now())
	}

	_, err := r.collection.InsertOne(ctx, rModel)
	if err != nil {
		return err
	}

	return nil
}
