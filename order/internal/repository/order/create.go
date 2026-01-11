package order

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/order/internal/repository/converter"
	repoModel "github.com/delyke/go_workspace_example/order/internal/repository/model"
	"github.com/delyke/go_workspace_example/platform/pkg/logger"
)

func (r *repository) Create(ctx context.Context, order *model.Order) (*model.Order, error) {
	repoOrder := converter.OrderToRepoModel(order)

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_uuids", "total_price", "order_status").
		Values(repoOrder.UserUUID, repoOrder.PartUuids, repoOrder.TotalPrice, repoOrder.OrderStatus).
		Suffix("RETURNING *")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return nil, err
	}

	var ouuid uuid.UUID
	var userUUID uuid.UUID
	var partUUIDs []uuid.UUID
	var totalPrice float64
	var transactionUUID *uuid.UUID
	var orderStatus string
	var paymentMethod *string
	var createdAt time.Time
	var updatedAt *time.Time
	err = r.pool.QueryRow(ctx, query, args...).Scan(&ouuid, &userUUID, &partUUIDs, &totalPrice, &transactionUUID, &orderStatus, &paymentMethod, &createdAt, &updatedAt)
	if err != nil {
		logger.Error(ctx, "failed to insert order", zap.Error(err))
		return nil, err
	}

	insertedOrder := &repoModel.Order{
		UUID:          ouuid,
		UserUUID:      userUUID,
		PartUuids:     partUUIDs,
		TotalPrice:    totalPrice,
		OrderStatus:   orderStatus,
		PaymentMethod: paymentMethod,
	}

	return converter.RepoToOrderModel(insertedOrder)
}
