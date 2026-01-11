package integration

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

func (env *TestEnvironment) GetExistingPartInfo() *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          "55555555-5555-5555-5555-555555555555",
		Name:          "Wing Module R",
		Description:   "Right aerodynamic wing",
		Price:         210000,
		StockQuantity: 7,
		Category:      inventoryV1.Category_CATEGORY_WING,
		Tags:          []string{"wing", "right"},
	}
}

// ClearPartsCollection - удаляет все записи из коллекции деталей
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}
	return nil
}
