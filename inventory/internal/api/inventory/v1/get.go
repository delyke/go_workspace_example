package v1

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/inventory/internal/converter"
	"github.com/delyke/go_workspace_example/inventory/internal/model"
	inventoryV1 "github.com/delyke/go_workspace_example/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.partService.GetPart(ctx, req.GetUuid())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(codes.NotFound, "no parts found with uuid: %s", req.GetUuid())
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve part: %s", err)
	}
	return &inventoryV1.GetPartResponse{
		Part: converter.ModelPartToProto(part),
	}, nil
}
