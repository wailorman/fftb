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
	Position                   int        `json:"position"`

	RetriesCount int        `json:"retries_count"`
	RetryAt      *time.Time `json:"retry_at"`
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

// SearchTimeout _
const SearchTimeout = time.Duration(10 * time.Second)

// FindSegmentByID _
func (r *Instance) FindSegmentByID(id string) (models.ISegment, error) {
	result, err := r.store.Get(fmt.Sprintf("v1/segments/%s", id))

	if err != nil {
		if errors.Is(err, ukvs.ErrNotFound) {
			return nil, models.ErrNotFound
		}

		return nil, errors.Wrap(err, "Accessing store for segment")
	}

	return unmarshalSegmentModel(result)
}

func (r *Instance) searchSegments(fctx context.Context, multiple bool, check func(models.ISegment) bool) ([]models.ISegment, error) {
	ffctx, ffcancel := context.WithCancel(fctx)
	defer ffcancel()

	results, failures := r.store.FindAll(ffctx, "v1/segments/*")
	segments := make([]models.ISegment, 0)

	for {
		select {
		case <-r.ctx.Done():
			return segments, nil

		case <-fctx.Done():
			return segments, nil

		case err := <-failures:
			if err != nil {
				return nil, errors.Wrap(err, "Searching for free order")
			}

		case res, ok := <-results:
			if !ok {
				return segments, nil
			}

			modSegment, err := unmarshalSegmentModel(res)

			if err != nil {
				r.logger.WithError(err).
					WithField(dlog.KeyStorePayload, string(res)).
					Warn("Unmarshalling order model from store")

				continue
			}

			if check(modSegment) {
				segments = append(segments, modSegment)

				if !multiple {
					return segments, nil
				}
			}

		case <-time.After(SearchTimeout):
			return nil, models.ErrTimeoutReached
		}
	}
}

// SearchSegment _
func (r *Instance) SearchSegment(fctx context.Context, check func(models.ISegment) bool) (models.ISegment, error) {
	segments, err := r.searchSegments(fctx, false, check)

	if err != nil {
		return nil, err
	}

	if len(segments) == 0 {
		return nil, models.ErrNotFound
	}

	return segments[0], nil
}

// SearchAllSegments _
func (r *Instance) SearchAllSegments(fctx context.Context, check func(models.ISegment) bool) ([]models.ISegment, error) {
	segments, err := r.searchSegments(fctx, true, check)

	if err != nil {
		return nil, err
	}

	return segments, nil
}

// FindNotLockedSegment _
func (r *Instance) FindNotLockedSegment(fctx context.Context) (models.ISegment, error) {
	if !r.freeSegmentLock.TryLockTimeout(FreeSegmentTimeout) {
		return nil, models.ErrTimeoutReached
	}

	defer r.freeSegmentLock.Unlock()

	return r.SearchSegment(fctx, func(segment models.ISegment) bool {
		return !segment.GetIsLocked() && segment.GetState() == models.SegmentStatePublished
	})
}

// FindSegmentsByOrderID _
func (r *Instance) FindSegmentsByOrderID(fctx context.Context, orderID string) ([]models.ISegment, error) {
	return r.SearchAllSegments(fctx, func(segment models.ISegment) bool {
		return segment.GetOrderID() == orderID
	})
}

// PersistSegment _
func (r *Instance) PersistSegment(modSegment models.ISegment) error {
	segmentBefore, _ := r.FindSegmentByID(modSegment.GetID())

	dlog.WithSegment(r.logger, modSegment).
		WithField("after", dlog.JSON(modSegment)).
		WithField("before", dlog.JSON(segmentBefore)).
		Trace("Persisting segment")

	if modSegment == nil {
		return models.ErrMissingSegment
	}

	if validationErr := modSegment.Validate(); validationErr != nil {
		return validationErr
	}

	data, err := marshalSegmentModel(modSegment)

	if err != nil {
		return errors.Wrap(err, "Marshaling db segment for store")
	}

	err = r.store.Set(fmt.Sprintf("v1/segments/%s", modSegment.GetID()), data)

	if err != nil {
		return errors.Wrap(err, "Persisting segment to store")
	}

	return nil
}

// LockSegmentByID _
func (r *Instance) LockSegmentByID(segmentID string, lockedBy models.IAuthor) error {
	panic("deprecated")
}

// UnlockSegmentByID _
func (r *Instance) UnlockSegmentByID(segmentID string) error {
	panic("deprecated")
}

func unmarshalSegmentModel(data []byte) (models.ISegment, error) {
	dbSegment := &Segment{}
	err := dbSegment.unmarshal(data)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling")
	}

	return dbSegment.toModel()
}

func marshalSegmentModel(modSegment models.ISegment) ([]byte, error) {
	dbSegment := &Segment{}
	err := dbSegment.fromModel(modSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Converting from model")
	}

	return dbSegment.marshal()
}

func (dbSegment *Segment) unmarshal(data []byte) error {
	return unmarshalObject(data, ObjectTypeSegment, dbSegment)
}

func (dbSegment *Segment) marshal() ([]byte, error) {
	return marshalObject(dbSegment)
}

func (dbSegment *Segment) toModel() (models.ISegment, error) {
	if dbSegment.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownSegmentType
	}

	modSeg := &models.ConvertSegment{}

	err := deserializeSegmentPayload(dbSegment, modSeg)

	if err != nil {
		return nil, errors.Wrap(err, "Deserializing segment payload")
	}

	modSeg.Identity = dbSegment.ID
	modSeg.Type = dbSegment.Kind
	modSeg.OrderIdentity = dbSegment.OrderID
	modSeg.State = dbSegment.State
	modSeg.InputStorageClaimIdentity = dbSegment.InputStorageClaimIdentity
	modSeg.OutputStorageClaimIdentity = dbSegment.OutputStorageClaimIdentity
	modSeg.Position = dbSegment.Position
	modSeg.RetriesCount = dbSegment.RetriesCount
	modSeg.RetryAt = dbSegment.RetryAt

	modSeg.LockedUntil = dbSegment.LockedUntil

	if dbSegment.LockedBy != "" {
		modSeg.LockedBy = &models.Author{Name: dbSegment.LockedBy}
	}

	if dbSegment.Publisher != "" {
		modSeg.Publisher = &models.Author{Name: dbSegment.Publisher}
	}

	return modSeg, nil
}

func (dbSegment *Segment) fromModel(modSegment models.ISegment) error {
	var err error

	if modSegment.GetType() != models.ConvertV1Type {
		return models.ErrUnknownSegmentType
	}

	err = serializeSegmentPayload(modSegment, dbSegment)

	if err != nil {
		return errors.Wrap(err, "Serializing segment payload")
	}

	dbSegment.ID = modSegment.GetID()
	dbSegment.OrderID = modSegment.GetOrderID()
	dbSegment.ObjectType = ObjectTypeSegment
	dbSegment.State = modSegment.GetState()
	dbSegment.Kind = modSegment.GetType()
	dbSegment.InputStorageClaimIdentity = modSegment.GetInputStorageClaimIdentity()
	dbSegment.OutputStorageClaimIdentity = modSegment.GetOutputStorageClaimIdentity()
	dbSegment.LockedUntil = modSegment.GetLockedUntil()
	dbSegment.Position = modSegment.GetPosition()
	dbSegment.RetriesCount = modSegment.GetRetriesCount()
	dbSegment.RetryAt = modSegment.GetRetryAt()

	if lockedBy := modSegment.GetLockedBy(); lockedBy != nil {
		dbSegment.LockedBy = lockedBy.GetName()
	}

	if publisher := modSegment.GetPublisher(); publisher != nil {
		dbSegment.Publisher = publisher.GetName()
	}

	return nil
}

func serializeSegmentPayload(modSeg models.ISegment, dbSeg *Segment) error {
	convSeg, ok := modSeg.(*models.ConvertSegment)

	if !ok {
		return models.ErrUnknownSegmentType
	}

	payload := &ConvertSegmentPayload{
		Params: convSeg.Params,
		Muxer:  convSeg.Muxer,
	}

	bPayload, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	dbSeg.Payload = string(bPayload)

	return nil
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
