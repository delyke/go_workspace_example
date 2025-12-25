package v1

import (
	"context"

	clientConverter "github.com/delyke/go_workspace_example/order/internal/client/converter"
	"github.com/delyke/go_workspace_example/order/internal/model"
	generatedInventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	resp, err := c.generatedClient.ListParts(
		ctx,
		&generatedInventoryV1.ListPartsRequest{
			Filter: clientConverter.PartsFilterToProto(filter),
		},
	)
	if err != nil {
		return nil, err
	}
	return clientConverter.ProtoPartsToModel(resp.Parts), nil
}
