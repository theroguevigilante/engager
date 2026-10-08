package service

import (
	"context"

	"github.com/adaptive-learning/backend/internal/models"
)

func (s *Service) IngestSignals(ctx context.Context, attemptID string, batch models.TelemetryBatch) error {
	return s.repo.SaveSignalBatch(ctx, attemptID, batch)
}
