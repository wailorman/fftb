package local

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/files"

	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// DefaultSegmentSize _
const DefaultSegmentSize = 10

// ContracterInstance _
type ContracterInstance struct {
	ctx       context.Context
	tempPath  files.Pather
	dealer    models.IContracterDealer
	publisher models.IAuthor
	registry  models.IContracterRegistry
	wg        *sync.WaitGroup
	logger    logrus.FieldLogger
}

// NewContracter _
func NewContracter(ctx context.Context, dealer models.IContracterDealer, registry models.IContracterRegistry, tempPath files.Pather) (*ContracterInstance, error) {
	publisher, err := dealer.AllocatePublisherAuthority("local")

	if err != nil {
		return nil, errors.Wrap(err, "Allocating publisher authority")
	}

	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixContracter); logger == nil {
		logger = ctxlog.New(dlog.PrefixContracter)
	}

	return &ContracterInstance{
		ctx:       ctx,
		tempPath:  tempPath,
		dealer:    dealer,
		publisher: publisher,
		registry:  registry,
		wg:        &sync.WaitGroup{},
		logger:    logger,
	}, nil
}

// GetAllOrders _
func (contracter *ContracterInstance) GetAllOrders(ctx context.Context) ([]models.IOrder, error) {
	return contracter.registry.SearchAllOrders(ctx, func(models.IOrder) bool { return true })
}

// GetAllSegments _
func (contracter *ContracterInstance) GetAllSegments(ctx context.Context) ([]models.ISegment, error) {
	allOrders, err := contracter.GetAllOrders(ctx)

	if err != nil {
		return nil, err
	}

	allSegments := make([]models.ISegment, 0)

	for _, order := range allOrders {
		if err := ctx.Err(); err != nil {
			return nil, ctx.Err()
		}

		segments, err := contracter.dealer.GetSegmentsByOrderID(ctx, order.GetID())

		if err != nil {
			return nil, errors.Wrapf(err, "Getting segments by order id `%s`", order.GetID())
		}

		allSegments = append(allSegments, segments...)
	}

	return allSegments, nil
}

// GetSegmentsByOrderID _
func (contracter *ContracterInstance) GetSegmentsByOrderID(ctx context.Context, orderID string) ([]models.ISegment, error) {
	return contracter.dealer.GetSegmentsByOrderID(ctx, orderID)
}

// GetSegmentByID _
func (contracter *ContracterInstance) GetSegmentByID(id string) (models.ISegment, error) {
	return contracter.dealer.GetSegmentByID(id)
}
