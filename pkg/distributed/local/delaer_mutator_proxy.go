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

// DealerMutatorProxy _
type DealerMutatorProxy struct {
	ctx          context.Context
	logger       logrus.FieldLogger
	author       models.IAuthor
	dealer       DealerSegmentCanceller
	ignoreErrors bool
}

// DealerMutatorProxyOption _
type DealerMutatorProxyOption func(*DealerMutatorProxy)

// NewDealerMutatorProxy _
func NewDealerMutatorProxy(
	ctx context.Context,
	logger logrus.FieldLogger,
	dealer DealerSegmentCanceller,
	author models.IAuthor,
	options ...DealerMutatorProxyOption) *DealerMutatorProxy {

	mp := &DealerMutatorProxy{
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
func (dp *DealerMutatorProxy) CancelSegment(segment models.ISegment, reason string) error {
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

// WithDealerMutatorProxyIgnoreErrors _
func WithDealerMutatorProxyIgnoreErrors() DealerMutatorProxyOption {
	return func(mp *DealerMutatorProxy) {
		mp.ignoreErrors = true
	}
}
