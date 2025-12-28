package order

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (s *service) Create(ctx context.Context, userUUID string, partUUIDs []string) (string, float64, error) {
	listParts, err := s.inventoryClient.ListParts(
		ctx,
		model.PartsFilter{
			UUIDs: partUUIDs,
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return "", 0, model.ErrInventoryInternalServerError
		}

		switch st.Code() {
		case codes.Unavailable:
			return "", 0, model.ErrInventoryServiceUnavailable
		case codes.DeadlineExceeded:
			return "", 0, model.ErrInventoryServiceDeadlineExceeded
		default:
			return "", 0, model.ErrInventoryBadGateway
		}
	}

	if len(listParts) != len(partUUIDs) {
		return "", 0, model.ErrInventoryPartNotFound
	}

	var totalPrice float64
	for _, p := range listParts {
		totalPrice += p.Price
	}

	parsedUsrUUID, parserUsrUUIDErr := uuid.Parse(userUUID)
	if parserUsrUUIDErr != nil {
		return "", 0, err
	}

	var parsedPartUUIDs []uuid.UUID
	for _, pUUID := range partUUIDs {
		parsedPUUID, perr := uuid.Parse(pUUID)
		if perr != nil {
			continue
		}
		parsedPartUUIDs = append(parsedPartUUIDs, parsedPUUID)
	}

	order := &model.Order{
		UserUUID:    parsedUsrUUID,
		PartUuids:   parsedPartUUIDs,
		TotalPrice:  totalPrice,
		OrderStatus: model.OrderStatusPENDINGPAYMENT,
	}

	createdOrder, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return "", 0, err
	}
	return createdOrder.UUID.String(), createdOrder.TotalPrice, nil
}
