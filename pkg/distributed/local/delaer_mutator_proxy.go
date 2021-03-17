package local

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/distributed/dlog"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// DealerSegmentCanceller _
type DealerSegmentCanceller interface {
	CancelSegment(ctx context.Context, publisher models.IAuthor, id string, reason string) error
}

// DelaerMutatorProxy _
type DelaerMutatorProxy struct {
	ctx          context.Context
	logger       logrus.FieldLogger
	author       models.IAuthor
	dealer       DealerSegmentCanceller
	ignoreErrors bool
}

// DelaerMutatorProxyOption _
type DelaerMutatorProxyOption func(*DelaerMutatorProxy)

// NewDelaerMutatorProxy _
func NewDelaerMutatorProxy(
	ctx context.Context,
	logger logrus.FieldLogger,
	dealer DealerSegmentCanceller,
	author models.IAuthor,
	options ...DelaerMutatorProxyOption) *DelaerMutatorProxy {

	mp := &DelaerMutatorProxy{
		ctx:    ctx,
		logger: logger,
		author: author,
		dealer: dealer,
	}

	for _, option := range options {
		option(mp)
	}

	return mp
}

// CancelSegment _
func (dp *DelaerMutatorProxy) CancelSegment(segment models.ISegment, reason string) error {
	err := dp.dealer.CancelSegment(dp.ctx, dp.author, segment.GetID(), reason)

	if err != nil {
		if dp.ignoreErrors {
			dlog.WithSegment(dp.logger, segment).
				Warn("Failed to cancel segment")
		} else {
			return err
		}
	}

	return nil
}

// WithDelaerMutatorProxyIgnoreErrors _
func WithDelaerMutatorProxyIgnoreErrors() DelaerMutatorProxyOption {
	return func(mp *DelaerMutatorProxy) {
		mp.ignoreErrors = true
	}
}
