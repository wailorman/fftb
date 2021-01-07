package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/wailorman/fftb/pkg/distributed/ukvs"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// Segment _
type Segment struct {
	ObjectType                 string     `json:"object_type"`
	ID                         string     `json:"id"`
	OrderID                    string     `json:"order_id"`
	Kind                       string     `json:"kind"`
	InputStorageClaimIdentity  string     `json:"input_storage_claim_identity"`
	OutputStorageClaimIdentity string     `json:"output_storage_claim_identity"`
	Payload                    string     `json:"payload"`
	LockedUntil                *time.Time `json:"locked_until"`
	LockedBy                   string     `json:"locked_by"`
	State                      string     `json:"state"`
	// CreatedAt                  *time.Time
	// UpdatedAt                  *time.Time
}

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// FindSegmentByID _
func (r *Instance) FindSegmentByID(id string) (models.ISegment, error) {
	result, err := r.store.Get(fmt.Sprintf("v1/segments/%s", id))

	if err != nil {
		if errors.Is(err, ukvs.ErrNotFound) {
			return nil, models.ErrNotFound
		}

		return nil, errors.Wrap(err, "Accessing store for segment")
	}

	dbSeg := &Segment{}
	err = unmarshalObject(result, SegmentObjectType, dbSeg)

	if err != nil {
		return nil, err
	}

	segment, err := fromDbSegment(dbSeg)

	if err != nil {
		return nil, err
	}

	return segment, nil
}

// FindNotLockedSegment _
func (r *Instance) FindNotLockedSegment() (models.ISegment, error) {
	// TODO: receive cancellation context
	ctx, cancel := context.WithCancel(context.TODO())
	results, failures := r.store.FindAll(ctx, "v1/segments/*")

	for {
		select {
		case err := <-failures:
			if err != nil {
				cancel()
				return nil, err
			}

		case data := <-results:
			fmt.Printf("FindNotLockedSegment data: %#v\n", data)
			fmt.Printf("FindNotLockedSegment string(data): %#v\n", string(data))
			// fmt.Printf("FindNotLockedSegment ok: %#v\n", ok)

			if len(data) == 0 {
				cancel()
				return nil, models.ErrNotFound
			}

			dbSeg := &Segment{}

			err := unmarshalObject(data, SegmentObjectType, dbSeg)

			if err != nil {
				// TODO: log unmarshaling errors
				log.Println(err)
				continue
			}

			if dbSeg.LockedUntil == nil || time.Now().After(*dbSeg.LockedUntil) {
				modSeg, err := fromDbSegment(dbSeg)

				if err != nil {
					// TODO: log unmarshaling errors
					log.Println(err)
					continue
				}

				if modSeg.GetState() != models.SegmentPublishedState {
					// TODO: log debug info about segment if not published
					continue
				}

				cancel()
				return modSeg, nil
			}

			// TODO: log
			continue
		default:
			continue
		}
	}
}

// LockSegmentByID _
func (r *Instance) LockSegmentByID(segmentID string, lockedBy string) error {
	if lockedBy == "" {
		return models.ErrMissingLockAuthor
	}

	result, err := r.store.Get(fmt.Sprintf("v1/segments/%s", segmentID))

	if err != nil {
		if errors.Is(err, ukvs.ErrNotFound) {
			return models.ErrNotFound
		}

		return errors.Wrap(err, "Accessing store for segment")
	}

	dbSeg := &Segment{}
	err = unmarshalObject(result, SegmentObjectType, dbSeg)

	if err != nil {
		return err
	}

	lockedUntil := time.Now().Add(LockSegmentTimeout)
	dbSeg.LockedUntil = &lockedUntil
	dbSeg.LockedBy = lockedBy

	data, err := marshalObject(dbSeg)

	if err != nil {
		return errors.Wrap(err, "Marshaling db segment for store")
	}

	err = r.store.Set(fmt.Sprintf("v1/segments/%s", dbSeg.ID), data)

	if err != nil {
		return errors.Wrap(err, "Persisting segment to store")
	}

	return nil
}

// FindSegmentsByOrderID _
func (r *Instance) FindSegmentsByOrderID(orderID string) ([]models.ISegment, error) {
	modSegs := make([]models.ISegment, 0)

	ctx, cancel := context.WithCancel(context.TODO())
	results, failures := r.store.FindAll(ctx, "v1/segments/*")

	for {
		select {
		case data := <-results:
			if data == nil || len(data) == 0 {
				cancel()
				return nil, models.ErrNotFound
			}

			dbSeg := &Segment{}
			err := unmarshalObject(data, SegmentObjectType, dbSeg)

			if err != nil {
				// TODO: log unmarshaling errors
				continue
			}

			if dbSeg.OrderID != orderID {
				continue
			}

			modSeg, err := fromDbSegment(dbSeg)

			if err != nil {
				// TODO: log unmarshaling errors
				continue
			}

			modSegs = append(modSegs, modSeg)

		case err := <-failures:
			cancel()
			return nil, errors.Wrap(err, "Failed to list segments")
		}
	}
}

// PersistSegment _
func (r *Instance) PersistSegment(segment models.ISegment) error {
	if segment == nil {
		return models.ErrMissingSegment
	}

	dbSeg, err := toDbSegment(segment)

	if err != nil {
		return errors.Wrap(err, "Converting segment model to registry format")
	}

	data, err := marshalObject(dbSeg)

	if err != nil {
		return errors.Wrap(err, "Marshaling db segment for store")
	}

	err = r.store.Set(fmt.Sprintf("v1/segments/%s", segment.GetID()), data)

	if err != nil {
		return errors.Wrap(err, "Persisting segment to store")
	}

	return nil
}

func toDbSegment(segment models.ISegment) (*Segment, error) {
	var err error

	dbSegment := &Segment{ID: segment.GetID()}

	if segment.GetType() != models.ConvertV1Type {
		return dbSegment, models.ErrUnknownSegmentType
	}

	dbSegment.Payload, err = segment.GetPayload()

	if err != nil {
		return nil, errors.Wrap(err, "Getting segment payload")
	}

	dbSegment.OrderID = segment.GetOrderID()
	dbSegment.ObjectType = SegmentObjectType
	dbSegment.State = segment.GetState()
	dbSegment.Kind = segment.GetType()
	dbSegment.InputStorageClaimIdentity = segment.GetInputStorageClaimIdentity()
	dbSegment.OutputStorageClaimIdentity = segment.GetOutputStorageClaimIdentity()

	return dbSegment, nil
}

func fromDbSegment(dbSeg *Segment) (models.ISegment, error) {
	if dbSeg.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownSegmentType
	}

	modSeg := &models.ConvertSegment{}

	// fmt.Printf("dbSeg.Payload: %#v\n", dbSeg.Payload)

	err := json.Unmarshal([]byte(dbSeg.Payload), modSeg)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling payload")
	}

	modSeg.Identity = dbSeg.ID
	modSeg.Type = dbSeg.Kind
	modSeg.OrderIdentity = dbSeg.OrderID
	modSeg.State = dbSeg.State
	modSeg.InputStorageClaimIdentity = dbSeg.InputStorageClaimIdentity
	modSeg.OutputStorageClaimIdentity = dbSeg.OutputStorageClaimIdentity

	return modSeg, nil
}
