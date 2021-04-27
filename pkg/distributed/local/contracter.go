package local

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/subchen/go-trylock/v2"
	"github.com/wailorman/fftb/pkg/chwg"
	"github.com/wailorman/fftb/pkg/ctxlog"
	"github.com/wailorman/fftb/pkg/files"

	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// DefaultSegmentDuration in seconds
const DefaultSegmentDuration = 30

// LockOrderTimeout _
const LockOrderTimeout = time.Duration(10 * time.Second)

// ContracterInstance _
type ContracterInstance struct {
	ctx            context.Context
	tempPath       files.Pather
	dealer         models.IContracterDealer
	publisher      models.IAuthor
	registry       models.IContracterRegistry
	wg             *sync.WaitGroup
	logger         logrus.FieldLogger
	storageClient  models.IStorageClient
	orderQueueLock trylock.TryLocker
	orderMutator   models.IOrderMutator
}

// NewContracter _
func NewContracter(
	ctx context.Context,
	dealer models.IContracterDealer,
	registry models.IContracterRegistry,
	storageClient models.IStorageClient,
	orderMutator models.IOrderMutator,
	tempPath files.Pather) (*ContracterInstance, error) {

	publisher, err := dealer.AllocatePublisherAuthority(ctx, "local")

	if err != nil {
		return nil, errors.Wrap(err, "Allocating publisher authority")
	}

	var logger logrus.FieldLogger
	if logger = ctxlog.FromContext(ctx, dlog.PrefixContracter); logger == nil {
		logger = ctxlog.New(dlog.PrefixContracter)
	}

	return &ContracterInstance{
		ctx:            ctx,
		tempPath:       tempPath,
		dealer:         dealer,
		publisher:      publisher,
		registry:       registry,
		wg:             &sync.WaitGroup{},
		logger:         logger,
		storageClient:  storageClient,
		orderQueueLock: trylock.New(),
		orderMutator:   orderMutator,
	}, nil
}

// GetAllOrders _
func (contracter *ContracterInstance) GetAllOrders(ctx context.Context) ([]models.IOrder, error) {
	return contracter.SearchAllOrders(ctx, models.EmptyOrderFilters())
}

// SearchAllOrders _
func (contracter *ContracterInstance) SearchAllOrders(ctx context.Context, search models.IOrderSearchCriteria) ([]models.IOrder, error) {
	return contracter.registry.SearchAllOrders(ctx, func(order models.IOrder) bool { return search.Select(order) })
}

// GetAllSegments _
func (contracter *ContracterInstance) GetAllSegments(ctx context.Context) ([]models.ISegment, error) {
	return contracter.SearchAllSegments(ctx, models.EmptySegmentFilters())
}

// SearchAllSegments _
func (contracter *ContracterInstance) SearchAllSegments(ctx context.Context, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	allOrders, err := contracter.GetAllOrders(ctx)

	if err != nil {
		return nil, err
	}

	allSegments := make([]models.ISegment, 0)

	for _, order := range allOrders {
		if err := ctx.Err(); err != nil {
			return nil, ctx.Err()
		}

		segments, err := contracter.dealer.GetSegmentsByOrderID(ctx, contracter.publisher, order.GetID(), search)

		if err != nil {
			return nil, errors.Wrapf(err, "Getting segments by order id `%s`", order.GetID())
		}

		allSegments = append(allSegments, segments...)
	}

	return allSegments, nil
}

// GetSegmentsByOrderID _
func (contracter *ContracterInstance) GetSegmentsByOrderID(ctx context.Context, orderID string) ([]models.ISegment, error) {
	return contracter.dealer.GetSegmentsByOrderID(ctx, contracter.publisher, orderID, models.EmptySegmentFilters())
}

// SearchSegmentsByOrderID _
func (contracter *ContracterInstance) SearchSegmentsByOrderID(ctx context.Context, orderID string, search models.ISegmentSearchCriteria) ([]models.ISegment, error) {
	return contracter.dealer.GetSegmentsByOrderID(ctx, contracter.publisher, orderID, search)
}

// GetSegmentByID _
func (contracter *ContracterInstance) GetSegmentByID(ctx context.Context, id string) (models.ISegment, error) {
	return contracter.dealer.GetSegmentByID(ctx, contracter.publisher, id)
}

// CancelOrderByID _
func (contracter *ContracterInstance) CancelOrderByID(ctx context.Context, orderID string, reason string) error {
	contracter.wg.Add(1)
	defer contracter.wg.Done()

	order, err := contracter.registry.FindOrderByID(ctx, orderID)

	if err != nil {
		return errors.Wrapf(err, "Getting order by id `%s`", orderID)
	}

	logger := dlog.WithOrder(contracter.logger, order)

	segments, err := contracter.GetSegmentsByOrderID(ctx, orderID)

	if err != nil {
		return errors.Wrap(err, "Getting order segments")
	}

	err = contracter.orderMutator.CancelOrder(
		NewDelaerMutatorProxy(
			ctx,
			logger,
			contracter.dealer,
			contracter.publisher,
			WithDelaerMutatorProxyIgnoreErrors()),
		order,
		segments,
		reason,
	)

	if err != nil {
		return errors.Wrap(err, "Cancel order")
	}

	err = contracter.registry.PersistOrder(ctx, order)

	if err != nil {
		return errors.Wrap(err, "Persisting order to registry")
	}

	return nil
}

// FailOrderByID _
func (contracter *ContracterInstance) FailOrderByID(ctx context.Context, orderID string, reportedErr error) error {

	order, err := contracter.registry.FindOrderByID(ctx, orderID)

	if err != nil {
		return errors.Wrapf(err, "Getting order by id `%s`", orderID)
	}

	logger := dlog.WithOrder(contracter.logger, order)

	segments, err := contracter.GetSegmentsByOrderID(ctx, orderID)

	if err != nil {
		return errors.Wrap(err, "Getting order segments")
	}

	err = contracter.orderMutator.FailOrder(
		NewDelaerMutatorProxy(
			ctx,
			logger,
			contracter.dealer,
			contracter.publisher,
			WithDelaerMutatorProxyIgnoreErrors()),
		order,
		segments,
		reportedErr,
	)

	if err != nil {
		return errors.Wrap(err, "Fail order")
	}

	err = contracter.registry.PersistOrder(ctx, order)

	if err != nil {
		return errors.Wrap(err, "Persist order")
	}

	return nil
}

// PickOrderFromQueue _
func (contracter *ContracterInstance) PickOrderFromQueue(ctx context.Context) (models.IOrder, error) {
	if !contracter.orderQueueLock.TryLockTimeout(LockSegmentTimeout) {
		return nil, models.ErrLockTimeout
	}

	defer contracter.orderQueueLock.Unlock()

	return contracter.registry.SearchOrder(ctx, func(order models.IOrder) bool {
		return order.GetCanPublish()
	})
}

// ObserveOrders _
func (contracter *ContracterInstance) ObserveOrders(ctx context.Context, wg chwg.WaitGrouper) {
	go func() {
		var err error

		contracter.logger.Debug("Orders observer started")

		wg.Add(1)
		defer wg.Done()

		ticker := time.NewTicker(PollingInterval)
		unretryableLogger := ctxlog.WithPrefix(contracter.logger, dlog.PrefixContracter+".orders_observer.unretryable")

		for {
			select {
			case <-ctx.Done():
				return

			case <-ticker.C:
				err = cancelUnretryableOrders(ctx,
					unretryableLogger,
					contracter,
					contracter.registry)

				if err != nil {
					unretryableLogger.WithError(err).
						Warn("Failed to cancel unretryable orders")
				}
			}
		}
	}()
}

// OrderSearcher _
type OrderSearcher interface {
	SearchOrder(fctx context.Context, check func(models.IOrder) bool) (models.IOrder, error)
}

// OrderCanceller _
type OrderCanceller interface {
	CancelOrderByID(ctx context.Context, orderID string, reason string) error
}

// OrderSegmentsGetter _
type OrderSegmentsGetter interface {
	GetSegmentsByOrderID(ctx context.Context, orderID string) ([]models.ISegment, error)
}

func cancelUnretryableOrders(
	ctx context.Context,
	logger logrus.FieldLogger,
	orderCanceller OrderCanceller,
	orderSearcher OrderSearcher) error {

	_, searchErr := orderSearcher.SearchOrder(ctx, func(order models.IOrder) bool {
		// TODO: lock order

		if !order.GetCanRetry() && order.GetState() != models.OrderStateCancelled {
			cErr := orderCanceller.CancelOrderByID(ctx, order.GetID(), models.CancellationReasonFailed)

			if cErr != nil {
				dlog.WithOrder(logger, order).
					WithError(cErr).
					Warn("Failed to cancel order")
			}
		}

		return false
	})

	if searchErr != nil && !errors.Is(searchErr, models.ErrNotFound) {
		return searchErr
	}

	return nil
}

func cancelOrdersWithCancelledSegments(
	ctx context.Context,
	logger logrus.FieldLogger,
	orderCanceller OrderCanceller,
	segmentsGetter OrderSegmentsGetter,
	orderSearcher OrderSearcher) error {

	_, searchErr := orderSearcher.SearchOrder(ctx, func(order models.IOrder) bool {
		// TODO: lock order

		if order.GetState() != models.OrderStateCancelled {
			segments, segmentsErr := segmentsGetter.GetSegmentsByOrderID(ctx, order.GetID())

			if segmentsErr != nil {
				dlog.WithOrder(logger, order).
					WithError(segmentsErr).
					Warn("Failed to get order segments")

				return false
			}

			for _, segment := range segments {
				if segment.GetState() == models.SegmentStateCancelled {
					cancellationErr := orderCanceller.CancelOrderByID(ctx, order.GetID(), models.CancellationReasonFailed)

					if cancellationErr != nil {
						dlog.WithOrder(logger, order).
							WithError(cancellationErr).
							Warn("Failed to cancel order")
					}

					return false
				}
			}
		}

		return false
	})

	if searchErr != nil && !errors.Is(searchErr, models.ErrNotFound) {
		return searchErr
	}

	return nil
}
