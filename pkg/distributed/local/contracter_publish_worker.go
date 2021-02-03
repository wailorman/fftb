package local

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
)

// QueuedSegmentsPollingInterval _
const QueuedSegmentsPollingInterval = time.Duration(5 * time.Second)

// MaxQueuedSegmentsCount _
const MaxQueuedSegmentsCount = 15

// ContracterPublishWorker _
type ContracterPublishWorker struct {
	ctx        context.Context
	contracter *ContracterInstance
	closed     chan struct{}
	logger     logrus.FieldLogger
}

// NewContracterPublishWorker _
func NewContracterPublishWorker(ctx context.Context, contracter *ContracterInstance) *ContracterPublishWorker {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixContracterPublishWorker); logger == nil {
		logger = ctxlog.New(dlog.PrefixContracterPublishWorker)
	}

	return &ContracterPublishWorker{
		ctx:        ctx,
		contracter: contracter,
		closed:     make(chan struct{}),
		logger:     logger,
	}
}

// Start _
func (pW *ContracterPublishWorker) Start() {
	go func() {
		ticker := time.NewTicker(QueuedSegmentsPollingInterval)

		for {
			select {
			case <-pW.ctx.Done():
				close(pW.closed)
				return
			case <-ticker.C:
				queuedSegmentsCount, err := pW.contracter.dealer.GetQueuedSegmentsCount(pW.ctx)

				if err != nil {
					pW.logger.WithError(err).
						Warn("Failed to count queued segments")
					continue
				}

				if queuedSegmentsCount < MaxQueuedSegmentsCount {
					queuedOrder, err := pW.contracter.registry.PickOrderFromQueue(pW.ctx)

					if err != nil {
						pW.logger.WithError(err).
							Warn("Failed to pick new order from queue")
						continue
					}

					err = pW.contracter.PublishOrder(pW.ctx, queuedOrder)

					if err != nil {
						pW.logger.WithError(err).
							Warn("Failed to publish new order")
						continue
					}
				}
			}
		}
	}()
}

// Closed _
func (pW *ContracterPublishWorker) Closed() <-chan struct{} {
	return pW.closed
}
