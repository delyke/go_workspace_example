package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/order/internal/repository/converter"
	repoModel "github.com/delyke/go_workspace_example/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, id string) (*model.Order, error) {
	builderSelectOne := sq.Select(
		"uuid",
		"user_uuid",
		"part_uuids",
		"total_price",
		"transaction_uuid",
		"order_status",
		"payment_method",
	).
		From("orders").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"uuid": id}).
		Limit(1)

	query, args, err := builderSelectOne.ToSql()
	if err != nil {
		return nil, err
	}
	var order repoModel.Order

	var ouuid, userUUID uuid.UUID
	var partUUIDs []uuid.UUID
	var totalPrice float64
	var transactionUUID *uuid.UUID
	var orderStatus string
	var paymentMethod *string

	err = r.pool.QueryRow(ctx, query, args...).Scan(&ouuid, &userUUID, &partUUIDs, &totalPrice, &transactionUUID, &orderStatus, &paymentMethod)
	if err != nil {
		return nil, err
	}

	order = repoModel.Order{
		UUID:            ouuid,
		UserUUID:        userUUID,
		PartUuids:       partUUIDs,
		TotalPrice:      totalPrice,
		TransactionUUID: transactionUUID,
		OrderStatus:     orderStatus,
		PaymentMethod:   paymentMethod,
	}

	return converter.RepoToOrderModel(&order)
}
