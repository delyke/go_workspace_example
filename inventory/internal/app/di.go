package app

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	inventoryApi "github.com/delyke/go_workspace_example/inventory/internal/api/inventory/v1"
	"github.com/delyke/go_workspace_example/inventory/internal/config"
	"github.com/delyke/go_workspace_example/inventory/internal/repository"
	partRepository "github.com/delyke/go_workspace_example/inventory/internal/repository/part"
	"github.com/delyke/go_workspace_example/inventory/internal/service"
	partService "github.com/delyke/go_workspace_example/inventory/internal/service/part"
	"github.com/delyke/go_workspace_example/platform/pkg/closer"
	"github.com/delyke/go_workspace_example/platform/pkg/logger"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	partRepo       repository.PartRepository
	partService    service.PartService
	inventoryV1Api inventoryV1.InventoryServiceServer
	mongoDBClient  *mongo.Client
	mongoDBHandle  *mongo.Database
}

func NewDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PartRepo(ctx context.Context) repository.PartRepository {
	if d.partRepo == nil {
		//nolint:contextcheck
		d.partRepo = partRepository.NewRepository(d.MongoDBHandle(ctx))
	}
	return d.partRepo
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			logger.Error(ctx, "Error connecting to mongo db", zap.Error(err))
		}

		closer.AddNamed("Mongo DB", func(ctx context.Context) error {
			cerr := mongoClient.Disconnect(ctx)
			if cerr != nil {
				logger.Error(ctx, "Error disconnecting from mongo db", zap.Error(cerr))
				return cerr
			}
			return nil
		})

		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			logger.Error(ctx, "Error pinging mongo db", zap.Error(err))
			return nil
		}
		d.mongoDBClient = mongoClient
	}
	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database("inventory")
	}
	return d.mongoDBHandle
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepo(ctx))
	}
	return d.partService
}

func (d *diContainer) InventoryV1Api(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1Api == nil {
		d.inventoryV1Api = inventoryApi.NewApi(d.PartService(ctx))
	}
	return d.inventoryV1Api
}
