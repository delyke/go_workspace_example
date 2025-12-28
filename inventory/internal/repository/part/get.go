package part

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
	repoModel "github.com/delyke/go_workspace_example/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	var part repoModel.Part

	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, err
	}

	return lo.ToPtr(converter.RepoPartToModel(part)), nil
}
