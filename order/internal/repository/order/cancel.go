package order

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"

	"github.com/delyke/go_workspace_example/order/internal/model"
	"github.com/delyke/go_workspace_example/platform/pkg/logger"
)

func (r *repository) Cancel(ctx context.Context, uuid string, status model.OrderStatus) error {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("order_status", status).
		Where(sq.Eq{"uuid": uuid})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		logger.Error(ctx, "failed to build query for cancel", zap.Error(err))
		return err
	}
	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "failed to cancel order", zap.Error(err))
		return err
	}
	logger.Debug(ctx, "cancelled orders count", zap.Int64("orders count", res.RowsAffected()))
	return nil
}
