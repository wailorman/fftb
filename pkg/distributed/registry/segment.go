package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/ukvs"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/media/convert"
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
	Publisher                  string     `json:"publisher"`
}

// ConvertSegmentPayload _
type ConvertSegmentPayload struct {
	Params convert.Params `json:"params"`
	Muxer  string         `json:"muxer"`
}

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// FreeSegmentTimeout _
const FreeSegmentTimeout = time.Duration(20 * time.Second)

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
	if !r.freeSegmentLock.TryLockTimeout(FreeSegmentTimeout) {
		return nil, models.ErrTimeoutReached
	}

	ctx, cancel := context.WithCancel(r.ctx)
	results, failures := r.store.FindAll(ctx, "v1/segments/*")

	for {
		select {
		case err := <-failures:
			if err != nil {
				cancel()
				return nil, err
			}

		case data := <-results:
			if len(data) == 0 {
				cancel()
				return nil, models.ErrNotFound
			}

			dbSeg := &Segment{}

			err := unmarshalObject(data, SegmentObjectType, dbSeg)

			if err != nil {
				r.logger.WithError(err).
					Trace("Unmarshalling error")

				continue
			}

			if dbSeg.LockedUntil == nil || time.Now().After(*dbSeg.LockedUntil) {
				modSeg, err := fromDbSegment(dbSeg)

				if err != nil {
					r.logger.WithField(dlog.KeySegmentID, dbSeg.ID).
						WithError(err).
						Trace("Serializing from db segment model")

					continue
				}

				if modSeg.GetState() != models.SegmentPublishedState {
					r.logger.WithField(dlog.KeySegmentID, modSeg.GetID()).
						WithField(dlog.KeySegmentState, modSeg.GetState()).
						Trace("Segment is not published")

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
func (r *Instance) LockSegmentByID(segmentID string, lockedBy models.IAuthor) error {
	if lockedBy == nil || lockedBy.GetName() == "" {
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
	dbSeg.LockedBy = lockedBy.GetName()

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

// UnlockSegmentByID _
func (r *Instance) UnlockSegmentByID(segmentID string) error {
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

	dbSeg.LockedBy = ""
	dbSeg.LockedUntil = nil

	data, err := marshalObject(dbSeg)

	if err != nil {
		return errors.Wrap(err, "Marshaling db segment for store")
	}

	err = r.store.Set(fmt.Sprintf("v1/segments/%s", dbSeg.ID), data)

	return err
}

func toDbSegment(segment models.ISegment) (*Segment, error) {
	var err error

	dbSegment := &Segment{ID: segment.GetID()}

	if segment.GetType() != models.ConvertV1Type {
		return dbSegment, models.ErrUnknownSegmentType
	}

	dbSegment.Payload, err = serializeSegmentPayload(segment)

	if err != nil {
		return nil, errors.Wrap(err, "Serializing segment payload")
	}

	dbSegment.OrderID = segment.GetOrderID()
	dbSegment.ObjectType = SegmentObjectType
	dbSegment.State = segment.GetState()
	dbSegment.Kind = segment.GetType()
	dbSegment.InputStorageClaimIdentity = segment.GetInputStorageClaimIdentity()
	dbSegment.OutputStorageClaimIdentity = segment.GetOutputStorageClaimIdentity()
	dbSegment.LockedUntil = segment.GetLockedUntil()

	if lockedBy := segment.GetLockedBy(); lockedBy != nil {
		dbSegment.LockedBy = lockedBy.GetName()
	}

	if publisher := segment.GetPublisher(); publisher != nil {
		dbSegment.Publisher = publisher.GetName()
	}

	return dbSegment, nil
}

func fromDbSegment(dbSeg *Segment) (models.ISegment, error) {
	if dbSeg.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownSegmentType
	}

	modSeg := &models.ConvertSegment{}

	err := deserializeSegmentPayload(dbSeg, modSeg)

	if err != nil {
		return nil, errors.Wrap(err, "Deserializing segment payload")
	}

	modSeg.Identity = dbSeg.ID
	modSeg.Type = dbSeg.Kind
	modSeg.OrderIdentity = dbSeg.OrderID
	modSeg.State = dbSeg.State
	modSeg.InputStorageClaimIdentity = dbSeg.InputStorageClaimIdentity
	modSeg.OutputStorageClaimIdentity = dbSeg.OutputStorageClaimIdentity

	modSeg.LockedUntil = dbSeg.LockedUntil

	if dbSeg.LockedBy != "" {
		modSeg.LockedBy = &models.Author{Name: dbSeg.LockedBy}
	}

	if dbSeg.Publisher != "" {
		modSeg.Publisher = &models.Author{Name: dbSeg.Publisher}
	}

	return modSeg, nil
}

func serializeSegmentPayload(modSeg models.ISegment) (string, error) {
	convSeg, ok := modSeg.(*models.ConvertSegment)

	if !ok {
		return "", models.ErrUnknownSegmentType
	}

	payload := &ConvertSegmentPayload{
		Params: convSeg.Params,
		Muxer:  convSeg.Muxer,
	}

	bPayload, err := json.Marshal(payload)

	if err != nil {
		return "", err
	}

	return string(bPayload), nil
}

func deserializeSegmentPayload(dbSeg *Segment, modSeg models.ISegment) error {
	if dbSeg.Kind != models.ConvertV1Type {
		return models.ErrUnknownSegmentType
	}

	convSeg, ok := modSeg.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownSegmentType
	}

	convPayload := &ConvertSegmentPayload{}

	err := json.Unmarshal([]byte(dbSeg.Payload), convPayload)

	if err != nil {
		return err
	}

	convSeg.Params = convPayload.Params
	convSeg.Muxer = convPayload.Muxer

	return nil
}
