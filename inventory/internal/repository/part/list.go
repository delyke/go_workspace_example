package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/delyke/go_workspace_example/inventory/internal/model"
	"github.com/delyke/go_workspace_example/inventory/internal/repository/converter"
	repoModel "github.com/delyke/go_workspace_example/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filters *model.PartsFilter) ([]*model.Part, error) {
	var qFilter bson.M
	if filters != nil {
		qFilter = buildFilters(filters)
	}

	cursor, err := r.collection.Find(ctx, qFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := cursor.Close(ctx); err != nil {
			log.Printf("error cursor.Closing: %v", cerr)
		}
	}()

	var parts []repoModel.Part

	err = cursor.All(ctx, &parts)
	if err != nil {
		return nil, err
	}
	return converter.RepoListPartsToModel(parts), nil
}

func buildFilters(filters *model.PartsFilter) bson.M {
	filter := bson.M{}
	and := bson.A{}

	if len(filters.UUIDs) > 0 {
		and = append(and, bson.M{"uuid": bson.M{"$in": filters.UUIDs}})
	}

	if len(filters.Names) > 0 {
		and = append(and, bson.M{"name": bson.M{"$in": filters.Names}})
	}

	if len(filters.Categories) > 0 {
		cats := make([]int, len(filters.Categories))
		for _, cat := range filters.Categories {
			cats = append(cats, int(cat))
		}
		and = append(and, bson.M{"category": bson.M{"$in": cats}})
	}

	if len(filters.ManufacturerCountries) > 0 {
		and = append(and, bson.M{"manufacturer.country": bson.M{"$in": filters.ManufacturerCountries}})
	}

	if len(filters.Tags) > 0 {
		and = append(and, bson.M{"tags": bson.M{"$in": filters.Tags}})
	}

	if len(and) == 0 {
		return filter
	}

	if len(and) == 1 {
		for k, v := range and[0].(bson.M) {
			filter[k] = v
		}
		return filter
	}

	filter["$and"] = and
	return filter
}
