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
	ctx           context.Context
	tempPath      files.Pather
	dealer        models.IContracterDealer
	publisher     models.IAuthor
	registry      models.IContracterRegistry
	wg            *sync.WaitGroup
	logger        logrus.FieldLogger
	storageClient models.IStorageClient
}

// NewContracter _
func NewContracter(
	ctx context.Context,
	dealer models.IContracterDealer,
	registry models.IContracterRegistry,
	storageClient models.IStorageClient,
	tempPath files.Pather) (*ContracterInstance, error) {

	publisher, err := dealer.AllocatePublisherAuthority("local")

	if err != nil {
		return nil, errors.Wrap(err, "Allocating publisher authority")
	}

	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixContracter); logger == nil {
		logger = ctxlog.New(dlog.PrefixContracter)
	}

	return &ContracterInstance{
		ctx:           ctx,
		tempPath:      tempPath,
		dealer:        dealer,
		publisher:     publisher,
		registry:      registry,
		wg:            &sync.WaitGroup{},
		logger:        logger,
		storageClient: storageClient,
	}, nil
}

// GetAllOrders _
func (contracter *ContracterInstance) GetAllOrders(ctx context.Context, search models.IOrderSearchCriteria) ([]models.IOrder, error) {
	return contracter.registry.SearchAllOrders(ctx, func(order models.IOrder) bool { return search.Select(order) })
}

// GetAllSegments _
func (contracter *ContracterInstance) GetAllSegments(ctx context.Context, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	allOrders, err := contracter.GetAllOrders(ctx, models.EmptyOrderFilters())

	if err != nil {
		return nil, err
	}

	allSegments := make([]models.ISegment, 0)

	for _, order := range allOrders {
		if err := ctx.Err(); err != nil {
			return nil, ctx.Err()
		}

		segments, err := contracter.dealer.GetSegmentsByOrderID(ctx, order.GetID(), search)

		if err != nil {
			return nil, errors.Wrapf(err, "Getting segments by order id `%s`", order.GetID())
		}

		allSegments = append(allSegments, segments...)
	}

	return allSegments, nil
}

// GetSegmentsByOrderID _
func (contracter *ContracterInstance) GetSegmentsByOrderID(ctx context.Context, orderID string, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	return contracter.dealer.GetSegmentsByOrderID(ctx, orderID, search)
}

// GetSegmentByID _
func (contracter *ContracterInstance) GetSegmentByID(id string) (models.ISegment, error) {
	return contracter.dealer.GetSegmentByID(id)
}

// CancelOrderByID _
func (contracter *ContracterInstance) CancelOrderByID(ctx context.Context, orderID string) error {
	logger := contracter.logger.WithField(dlog.KeyOrderID, orderID)

	order, err := contracter.registry.FindOrderByID(orderID)

	if err != nil {
		return errors.Wrapf(err, "Getting order by id `%s`", orderID)
	}

	convOrder, ok := order.(*models.ConvertOrder)

	if !ok {
		return errors.Wrapf(models.ErrUnknownOrderType, "Received order of type `%s`", order.GetType())
	}

	convOrder.State = models.OrderStateCancelled

	for _, segmentID := range convOrder.SegmentIDs {
		if ctx.Err() != nil {
			break
		}

		err = contracter.dealer.CancelSegment(contracter.publisher, segmentID)

		if err != nil {
			logger.Println()
			logger.WithField(dlog.KeySegmentID, segmentID).
				WithError(err).
				Warn("Problem with cancelling segment via dealer")
		}
	}

	if ctx.Err() != nil {
		for _, segmentID := range convOrder.SegmentIDs {
			err = contracter.dealer.RepublishSegment(contracter.publisher, segmentID)

			logger.WithField(dlog.KeySegmentID, segmentID).
				WithError(err).
				Warn("Problem with republishing segment via dealer")
		}

		return nil
	}

	err = contracter.registry.PersistOrder(convOrder)

	if err != nil {
		return errors.Wrap(err, "Persisting order to registry")
	}

	return nil
}
