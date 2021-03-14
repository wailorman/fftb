package local

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// FinishedOrdersPollingInterval _
const FinishedOrdersPollingInterval = time.Duration(5 * time.Second)

// ContracterConcatWorker _
type ContracterConcatWorker struct {
	ctx        context.Context
	contracter *ContracterInstance
	closed     chan struct{}
	logger     logrus.FieldLogger
}

// NewContracterConcatWorker _
func NewContracterConcatWorker(ctx context.Context, contracter *ContracterInstance) *ContracterConcatWorker {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixContracterConcatWorker); logger == nil {
		logger = ctxlog.New(dlog.PrefixContracterConcatWorker)
	}

	return &ContracterConcatWorker{
		ctx:        ctx,
		contracter: contracter,
		closed:     make(chan struct{}),
		logger:     logger,
	}
}

// Start _
func (pC *ContracterConcatWorker) Start() {
	go func() {
		ticker := time.NewTicker(QueuedSegmentsPollingInterval)

		for {
			select {
			case <-pC.ctx.Done():
				close(pC.closed)
				return

			case <-ticker.C:
				finishedOrder, err := pC.contracter.PickOrderForConcat(pC.ctx)

				if err != nil {
					if errors.Is(err, models.ErrNotFound) {
						pC.logger.Debug("Concatenatable orders not found")
						continue
					}

					pC.logger.WithError(err).
						Warn("Failed to pick finished order")
					continue
				}

				logger := pC.logger.WithField(dlog.KeyOrderID, finishedOrder.GetID())

				logger.WithField(dlog.KeyOrderID, finishedOrder.GetID()).
					Info("Found finished order for concatenation")

				err = pC.contracter.ConcatOrder(pC.ctx, finishedOrder)

				if err != nil {
					logger.WithError(err).
						Warn("Failed to concat finished order")

					failErr := pC.contracter.FailOrderByID(pC.ctx, finishedOrder.GetID(), err)

					if err != nil {
						logger.WithError(failErr).
							Warn("Failed to report failed order")
					}
				}

				continue
			}
		}
	}()
}

// Closed _
func (pC *ContracterConcatWorker) Closed() <-chan struct{} {
	return pC.closed
}
