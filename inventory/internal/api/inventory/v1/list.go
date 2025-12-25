package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/inventory/internal/converter"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	parts, err := a.partService.ListParts(ctx, converter.ProtoFiltersToModel(req.GetFilter()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %s", err)
	}
	return &inventoryV1.ListPartsResponse{
		Parts: converter.ModelPartListToProto(parts),
	}, nil
}
