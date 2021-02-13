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
