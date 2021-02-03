package local

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// Dealer _
type Dealer struct {
	storageController models.IStorageController
	registry          models.IDealerRegistry
	freeSegmentLock   trylock.TryLocker
	logger            logrus.FieldLogger
	ctx               context.Context
}

// NewDealer _
func NewDealer(ctx context.Context, sc models.IStorageController, r models.IDealerRegistry) (*Dealer, error) {
	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, "fftb.dealer"); logger == nil {
		logger = ctxlog.New("fftb.dealer")
	}

	return &Dealer{
		storageController: sc,
		registry:          r,
		freeSegmentLock:   trylock.New(),
		logger:            logger,
		ctx:               ctx,
	}, nil
}

// // FindSegmentByID _
// func (d *Dealer) FindSegmentByID(id string) (models.ISegment, error) {
// 	panic(models.ErrNotImplemented)
// }

// // Subscription _
// func (d *Dealer) Subscription(segment models.ISegment) (models.Subscriber, error) {
// 	panic(models.ErrNotImplemented)
// }
