package storage

import (
	"context"
	"net/url"

	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
	"github.com/wailorman/fftb/pkg/files"
)

// ErrUnknownClaimType typically returned when claim's url scheme is not from you passed to UniversalClient constructor
var ErrUnknownClaimType = errors.New("UNKNOWN_CLAIM_TYPE")

// ClientsMap is map type with specific storage clients
type ClientsMap map[string]models.IStorageClient

// UniversalClient implements IStorageClient interface, providing universal storage claim for multiple protocols
type UniversalClient struct {
	clientsMap ClientsMap
}

// NewUniversalClient returns a new UniversalClient instance
func NewUniversalClient(clientsMap ClientsMap) *UniversalClient {
	return &UniversalClient{clientsMap: clientsMap}
}

// BuildStorageClaimByURL implements IStorageClient's same method
func (uc *UniversalClient) BuildStorageClaimByURL(claimURL string) (models.IStorageClaim, error) {
	u, err := url.Parse(claimURL)

	if err != nil {
		return nil, errors.Wrapf(err, "Parsing storage claim url (`%s`)", claimURL)
	}

	claimType := u.Scheme

	if uc.clientsMap[claimType] == nil {
		return nil, errors.Wrapf(ErrUnknownClaimType, "(received claim type: `%s` from url `%s`)", claimType, claimURL)
	}

	storageClient := uc.clientsMap[claimType]

	return storageClient.BuildStorageClaimByURL(claimURL)
}

// RemoveLocalCopy implements IStorageClient's same method
func (uc *UniversalClient) RemoveLocalCopy(ctx context.Context, sc models.IStorageClaim) error {
	claimType := sc.GetType()

	if uc.clientsMap[claimType] == nil {
		return errors.Wrapf(ErrUnknownClaimType, "(received claim type: `%s`)", claimType)
	}

	storageClient := uc.clientsMap[claimType]

	return storageClient.RemoveLocalCopy(ctx, sc)
}

// MakeLocalCopy implements IStorageClient's same method
func (uc *UniversalClient) MakeLocalCopy(ctx context.Context, sc models.IStorageClaim, p chan models.IProgress) (files.Filer, error) {
	claimType := sc.GetType()

	if uc.clientsMap[claimType] == nil {
		return nil, errors.Wrapf(ErrUnknownClaimType, "(received claim type: `%s`)", claimType)
	}

	storageClient := uc.clientsMap[claimType]

	return storageClient.MakeLocalCopy(ctx, sc, p)
}

// MoveFileToStorageClaim implements IStorageClient's same method
func (uc *UniversalClient) MoveFileToStorageClaim(ctx context.Context, file files.Filer, sc models.IStorageClaim, p chan models.IProgress) error {
	claimType := sc.GetType()

	if uc.clientsMap[claimType] == nil {
		return errors.Wrapf(ErrUnknownClaimType, "(received claim type: `%s`)", claimType)
	}

	storageClient := uc.clientsMap[claimType]

	return storageClient.MoveFileToStorageClaim(ctx, file, sc, p)
}
