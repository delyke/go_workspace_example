package v1

import (
	"github.com/delyke/go_workspace_example/inventory/internal/service"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewApi(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
