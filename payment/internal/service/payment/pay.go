package payment

import (
	"context"

	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context, _, _, _ string) (string, error) {
	txId := uuid.NewString()
	return txId, nil
}
