package v1

import (
	def "github.com/delyke/go_workspace_example/order/internal/client/grpc"
	generatedInventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient generatedInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient generatedInventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
