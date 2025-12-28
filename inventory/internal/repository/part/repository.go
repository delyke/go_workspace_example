package part

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/delyke/go_workspace_example/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		panic(err)
	}

	r := &repository{
		collection: collection,
	}

	initErr := r.initIfEmpty(ctx)
	if initErr != nil {
		log.Printf("InitIfEmpty error: %v", initErr)
	}

	return r
}

func (r *repository) initIfEmpty(ctx context.Context) error {
	countParts, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Printf("error counting parts: %v", err)
		return err
	}
	if countParts == 0 {
		err = r.Init(ctx)
		if err != nil {
			log.Printf("init parts err: %v", err)
			return err
		}
	}
	return nil
}
