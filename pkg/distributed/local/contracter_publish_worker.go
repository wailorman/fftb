package local

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// MaxQueuedSegmentsForNextOrder _
const MaxQueuedSegmentsForNextOrder = 15

// ContracterPublishWorker _
type ContracterPublishWorker struct {
	ctx          context.Context
	contracter   *ContracterInstance
	closed       chan struct{}
	logger       logrus.FieldLogger
	orderMutator models.IOrderMutator
}

// NewContracterPublishWorker _
func NewContracterPublishWorker(ctx context.Context, contracter *ContracterInstance, orderMutator models.IOrderMutator) *ContracterPublishWorker {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixContracterPublishWorker); logger == nil {
		logger = ctxlog.New(dlog.PrefixContracterPublishWorker)
	}

	return &ContracterPublishWorker{
		ctx:          ctx,
		contracter:   contracter,
		closed:       make(chan struct{}),
		logger:       logger,
		orderMutator: orderMutator,
	}
}

// Start _
func (pW *ContracterPublishWorker) Start() {
	go func() {
		for {
			if pW.ctx.Err() != nil {
				return
			}

			queuedSegmentsCount, err := pW.contracter.dealer.GetQueuedSegmentsCount(pW.ctx, pW.contracter.publisher)

			if err != nil {
				pW.logger.WithError(err).
					Warn("Failed to count queued segments")
				time.Sleep(PollingInterval)
				continue
			}

			if queuedSegmentsCount >= MaxQueuedSegmentsForNextOrder {
				time.Sleep(PollingInterval)
				continue
			}

			order, err := pW.contracter.PickOrderFromQueue(pW.ctx)

			if err != nil {
				if errors.Is(err, models.ErrNotFound) {
					pW.logger.Debug("Queued orders not found")
				} else {
					pW.logger.WithError(err).
						Warn("Failed to pick new order from queue")
				}
				time.Sleep(PollingInterval)
				continue
			}

			oLogger := pW.logger.WithField(dlog.KeyOrderID, order.GetID())

			err = pW.contracter.publishOrder(pW.ctx, order)

			if err != nil {
				oLogger.WithError(err).
					Warn("Failed to publish new order")

				failErr := pW.contracter.FailOrderByID(pW.ctx, order.GetID(), err)

				if err != nil {
					oLogger.WithError(failErr).
						Warn("Failed to report order failure")
				}
			}
		}
	}()
}

// Closed _
func (pW *ContracterPublishWorker) Closed() <-chan struct{} {
	return pW.closed
}
