package converters

import (
	"github.com/pkg/errors"

	"github.com/twitchtv/twirp"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// ToRPCError converts internal error to RPC format
func ToRPCError(err error) error {
	return twirp.NewError(twirp.Unknown, err.Error())
}

// FromRPCError converts RPC error to internal format
func FromRPCError(err error) error {
	if twirpError, ok := errors.Cause(err).(twirp.Error); ok {
		switch twirpError.Code() {
		case twirp.NotFound:
			return errors.Wrap(models.ErrNotFound, err.Error())
		}
	}

	return err
}
