package local

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/files"

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
}

// NewContracter _
func NewContracter(ctx context.Context, dealer models.IContracterDealer, registry models.IContracterRegistry, tempPath files.Pather) (*ContracterInstance, error) {
	publisher, err := dealer.AllocatePublisherAuthority("local")

	if err != nil {
		return nil, errors.Wrap(err, "Allocating publisher authority")
	}

	return &ContracterInstance{
		ctx:       ctx,
		tempPath:  tempPath,
		dealer:    dealer,
		publisher: publisher,
		registry:  registry,
		wg:        &sync.WaitGroup{},
	}, nil
}
