package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (r *repository) Pay(
	ctx context.Context,
	uuid string,
	method model.PaymentMethod,
	txUUID string,
	status model.OrderStatus,
) (*model.Order, error) {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("transaction_uuid", txUUID).
		Set("order_status", status).
		Set("payment_method", method).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		return nil, err
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	if res.RowsAffected() == 0 {
		return nil, model.ErrOrderNotFound
	}

	var updatedOrder *model.Order
	updatedOrder, err = r.Get(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return updatedOrder, nil
}
