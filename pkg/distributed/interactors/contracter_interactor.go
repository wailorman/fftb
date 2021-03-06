package interactors

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// ContracterInteractor _
type ContracterInteractor struct {
	contracter models.IContracter
}

// NewContracterInteractor _
func NewContracterInteractor(contracter models.IContracter) *ContracterInteractor {
	return &ContracterInteractor{
		contracter: contracter,
	}
}

// OrdersListItem _
type OrdersListItem struct {
	ID string `json:"id"`

	// InputFile is local path or name of local contracter file
	InputFile string `json:"input_file"`

	// OutputFile is local path or name of local contracter file
	OutputFile string `json:"output_file"`

	// Progress is float value 0..1
	Progress float64 `json:"progress"`

	State string `json:"state"`
}

// GetAllOrders _
func (ci *ContracterInteractor) GetAllOrders(ctx context.Context, search models.IOrderSearchCriteria) ([]*OrdersListItem, error) {
	orders, err := ci.contracter.GetAllOrders(ctx, search)

	if err != nil {
		return nil, errors.Wrap(err, "Getting all orders from contracter")
	}

	allOrders := make([]*OrdersListItem, 0)

	for _, order := range orders {
		segments, err := ci.contracter.GetSegmentsByOrderID(ctx, order.GetID(), models.EmptySegmentFilters())

		if err != nil {
			return nil, errors.Wrapf(err, "Gettings segments by order id `%s`", order.GetID())
		}

		allOrders = append(allOrders, &OrdersListItem{
			ID:         order.GetID(),
			InputFile:  order.GetInputFile().Name(),
			OutputFile: order.GetOutputFile().Name(),
			Progress:   order.CalculateProgress(segments),
			State:      order.GetState(),
		})
	}

	return allOrders, nil
}

// OrderItem _
type OrderItem struct {
	ID string `json:"id"`

	// InputFile is local path or name of local contracter file
	InputFile string `json:"input_file"`

	// OutputFile is local path or name of local contracter file
	OutputFile string `json:"output_file"`

	// Progress is float value 0..1
	Progress float64 `json:"progress"`

	State         string `json:"state"`
	SegmentsCount int    `json:"segments_count"`
}

// GetOrderByID _
func (ci *ContracterInteractor) GetOrderByID(ctx context.Context, id string) (*OrderItem, error) {
	orders, err := ci.contracter.GetAllOrders(ctx, models.EmptyOrderFilters())

	if err != nil {
		return nil, errors.Wrap(err, "Getting all orders from contracter")
	}

	order := searchOrderByID(id, orders)

	if order == nil {
		return nil, models.ErrNotFound
	}

	segments, err := ci.contracter.GetSegmentsByOrderID(ctx, order.GetID(), models.EmptySegmentFilters())

	if err != nil {
		return nil, errors.Wrapf(err, "Getting segments by order id `%s`", order.GetID())
	}

	return &OrderItem{
		ID:            order.GetID(),
		InputFile:     order.GetInputFile().Name(),
		OutputFile:    order.GetOutputFile().Name(),
		Progress:      order.CalculateProgress(segments),
		State:         order.GetState(),
		SegmentsCount: len(segments),
	}, nil
}

// SegmentsListItem _
type SegmentsListItem struct {
	ID        string `json:"id"`
	State     string `json:"state"`
	Performer string `json:"performer"`
}

// GetSegmentsByOrderID _
func (ci *ContracterInteractor) GetSegmentsByOrderID(ctx context.Context, id string, search models.ISegmentSearchCriteria) ([]*SegmentsListItem, error) {
	var segments []models.ISegment

	if id != "" {
		orders, err := ci.contracter.GetAllOrders(ctx, models.EmptyOrderFilters())

		if err != nil {
			return nil, errors.Wrap(err, "Getting all orders from contracter")
		}

		order := searchOrderByID(id, orders)

		if order == nil {
			return nil, models.ErrNotFound
		}

		segments, err = ci.contracter.GetSegmentsByOrderID(ctx, order.GetID(), search)

		if err != nil {
			return nil, errors.Wrapf(err, "Getting segments by order id `%s`", order.GetID())
		}
	} else {
		var err error
		segments, err = ci.contracter.GetAllSegments(ctx, search)

		if err != nil {
			return nil, errors.Wrap(err, "Getting all segments")
		}
	}

	segmentsData := make([]*SegmentsListItem, 0)

	for _, segment := range segments {
		var performerName string

		if segment.GetPerformer() != nil {
			performerName = segment.GetPerformer().GetName()
		}

		segmentsData = append(segmentsData, &SegmentsListItem{
			ID:        segment.GetID(),
			State:     segment.GetCurrentState(),
			Performer: performerName,
		})
	}

	return segmentsData, nil
}

// SegmentItem _
type SegmentItem struct {
	ID        string `json:"id"`
	State     string `json:"state"`
	Performer string `json:"performer"`
}

// GetSegmentByID _
func (ci *ContracterInteractor) GetSegmentByID(ctx context.Context, id string) (*SegmentItem, error) {
	segments, err := ci.contracter.GetAllSegments(ctx, models.EmptySegmentFilters())

	if err != nil {
		return nil, errors.Wrap(err, "Getting all segments")
	}

	segment := searchSegmentByID(id, segments)

	if segment == nil {
		return nil, models.ErrNotFound
	}

	var performerName string

	if segment.GetPerformer() != nil {
		performerName = segment.GetPerformer().GetName()
	}

	return &SegmentItem{
		ID:        segment.GetID(),
		State:     segment.GetCurrentState(),
		Performer: performerName,
	}, nil
}

// CancelOrderByID _
func (ci *ContracterInteractor) CancelOrderByID(ctx context.Context, id string) error {
	orders, err := ci.contracter.GetAllOrders(ctx, models.EmptyOrderFilters())

	if err != nil {
		return errors.Wrap(err, "Getting all orders from contracter")
	}

	order := searchOrderByID(id, orders)

	if order == nil {
		return models.ErrNotFound
	}

	err = ci.contracter.CancelOrderByID(ctx, order.GetID())

	if err != nil {
		return errors.Wrap(err, "Cancelling order")
	}

	return nil
}

func searchOrderByID(passedID string, allObjects []models.IOrder) models.IOrder {
	for _, object := range allObjects {
		if strings.Index(object.GetID(), passedID) == 0 {
			return object
		}
	}

	return nil
}

func searchSegmentByID(passedID string, allObjects []models.ISegment) models.ISegment {
	for _, object := range allObjects {
		if strings.Index(object.GetID(), passedID) == 0 {
			return object
		}
	}

	return nil
}
