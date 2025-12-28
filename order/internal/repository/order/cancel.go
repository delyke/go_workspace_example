package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/delyke/go_workspace_example/order/internal/model"
)

func (r *repository) Cancel(ctx context.Context, uuid string, status model.OrderStatus) error {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("order_status", status).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query for cancel: %v", err)
		return err
	}
	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to cancel order: %v", err)
		return err
	}
	log.Printf("cancelled orders count: %d", res.RowsAffected())
	return nil
}
