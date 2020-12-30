package registry

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"gorm.io/gorm"
)

// Segment _
type Segment struct {
	ID                         string
	OrderID                    string
	Kind                       string
	InputStorageClaimIdentity  string
	OutputStorageClaimIdentity string
	Payload                    string
	LockedUntil                *time.Time
	LockedBy                   string
	// CreatedAt                  *time.Time
	// UpdatedAt                  *time.Time
}

// LockSegmentTimeout _
const LockSegmentTimeout = time.Duration(10 * time.Second)

// FindSegmentByID _
func (r *SqliteRegistry) FindSegmentByID(id string) (models.ISegment, error) {
	dbSegment := &Segment{}
	result := r.gdb.First(dbSegment, fmt.Sprintf("id = '%s'", id))

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	}

	if dbSegment.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownOrderType
	}

	return fromDbSegment(dbSegment)
}

// FindNotLockedSegment _
func (r *SqliteRegistry) FindNotLockedSegment() (models.ISegment, error) {
	dbSegments := make([]*Segment, 0)

	locked := r.freeSegmentLock.TryLockTimeout(LockSegmentTimeout)

	if !locked {
		return nil, models.ErrFreeSegmentLockTimeout
	}

	defer r.freeSegmentLock.Unlock()

	r.gdb.Find(&dbSegments)

	for _, dbSeg := range dbSegments {
		seg, err := fromDbSegment(dbSeg)

		if err != nil {
			// log.Panicf("failed to parse segment from db: %s\n", err)
			log.Printf("failed to parse segment from db: %s\n", err)
			continue
		}

		if !seg.GetIsLocked() {
			return seg, nil
		}
	}

	return nil, models.ErrNotFound
}

// LockSegmentByID _
func (r *SqliteRegistry) LockSegmentByID(segmentID string, lockedBy string) error {
	if lockedBy == "" {
		return models.ErrMissingLockAuthor
	}

	result := r.gdb.Model(&Segment{ID: segmentID}).Updates(map[string]interface{}{"locked_until": time.Now(), "locked_by": lockedBy})

	if result.RowsAffected == 0 {
		return models.ErrNotFound
	}

	return nil
}

// FindSegmentsByOrderID _
func (r *SqliteRegistry) FindSegmentsByOrderID(orderID string) ([]models.ISegment, error) {
	dbSegments := make([]*Segment, 0)

	getResults := r.gdb.Where("order_id = ?", orderID).Find(&dbSegments)

	if getResults.Error != nil {
		return nil, errors.Wrap(getResults.Error, "Retrieving segments")
	}

	segments := make([]models.ISegment, 0)

	for _, dbSegment := range dbSegments {
		segment, err := fromDbSegment(dbSegment)

		if err != nil {
			return nil, errors.Wrap(err, "")
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

// PersistSegment _
func (r *SqliteRegistry) PersistSegment(segment models.ISegment) error {
	if segment.GetType() != models.ConvertV1Type {
		return models.ErrUnknownSegmentType
	}

	dbSegment, err := toDbSegment(segment)

	// getResult := r.gdb.First(&Segment{}, segment.GetID())
	getResult := r.gdb.First(&Segment{}, fmt.Sprintf("id = '%s'", segment.GetID()))

	var result *gorm.DB

	// fmt.Printf("errors.Is(getResult.Error, gorm.ErrRecordNotFound): %#v\n", errors.Is(getResult.Error, gorm.ErrRecordNotFound))

	if errors.Is(getResult.Error, gorm.ErrRecordNotFound) {
		result = r.gdb.Create(dbSegment)
	} else {
		result = r.gdb.Save(dbSegment)
	}

	// fmt.Printf("result: %#v\n", result)

	if result.Error != nil {
		return errors.Wrap(err, "Failed to persist order")
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
	dbSegment.Kind = segment.GetType()
	dbSegment.InputStorageClaimIdentity = segment.GetInputStorageClaimIdentity()
	dbSegment.OutputStorageClaimIdentity = segment.GetOutputStorageClaimIdentity()

	return dbSegment, nil
}

func fromDbSegment(dbSegment *Segment) (models.ISegment, error) {
	if dbSegment.Kind != models.ConvertV1Type {
		return nil, models.ErrUnknownSegmentType
	}

	convertSegment := &models.ConvertSegment{}

	// fmt.Printf("dbSegment.Payload: %#v\n", dbSegment.Payload)

	err := json.Unmarshal([]byte(dbSegment.Payload), convertSegment)

	if err != nil {
		return nil, errors.Wrap(err, "Unmarshaling payload")
	}

	convertSegment.Identity = dbSegment.ID
	convertSegment.Type = dbSegment.Kind
	convertSegment.OrderIdentity = dbSegment.OrderID
	convertSegment.InputStorageClaimIdentity = dbSegment.InputStorageClaimIdentity
	convertSegment.OutputStorageClaimIdentity = dbSegment.OutputStorageClaimIdentity

	return convertSegment, nil
}
