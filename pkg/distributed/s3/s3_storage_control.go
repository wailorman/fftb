package s3

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/pkg/errors"
	"github.com/wailorman/fftb/pkg/distributed/models"
)

// SignedURLDuration _
var SignedURLDuration = time.Duration(24 * time.Hour)

// StorageControl _
type StorageControl struct {
	s3Client *s3.S3
	bucket   string
}

// NewStorageControl _
func NewStorageControl() *StorageControl {
	// host := os.Getenv("FFTB_S3_HOST")
	host := "https://s3.hermes.ha.wailorman.ru"
	// key := os.Getenv("FFTB_S3_KEY")
	key := "NAUBWE2PZUNCEHLH2I3P"
	// secret := os.Getenv("FFTB_S3_SECRET")
	secret := "EUWVAfkRoA8KE/JrrsDW3/gmehNgj29Lm8iFPckB"
	// bucket := os.Getenv("FFTB_S3_BUCKET")
	bucket := "fftb-slow"

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:         aws.String(host),
		Region:           aws.String("us-east-1"),
		S3ForcePathStyle: aws.Bool(true),
	}

	newSession := session.New(s3Config)
	s3Client := s3.New(newSession)

	return &StorageControl{
		s3Client: s3Client,
		bucket:   bucket,
	}
}

// AllocateStorageClaim _
func (sc *StorageControl) AllocateStorageClaim(ctx context.Context, identity string) (models.IStorageClaim, error) {
	req, _ := sc.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(identity),
	})

	urlStr, err := req.Presign(SignedURLDuration)

	if err != nil {
		return nil, errors.Wrap(err, "Signing url")
	}

	fmt.Printf("urlStr: %#v\n", urlStr)

	claim := &StorageClaim{
		id:   identity,
		url:  urlStr,
		size: 0,
	}

	return claim, nil
}

// BuildStorageClaim _
func (sc *StorageControl) BuildStorageClaim(identity string) (models.IStorageClaim, error) {

	size, err := sc.GetStorageClaimSize(context.Background(), identity)

	if err != nil {
		return nil, errors.Wrap(err, "Calculating storage claim size")
	}

	url, err := sc.GetURLForStorageClaim(context.Background(), identity)

	if err != nil {
		return nil, errors.Wrap(err, "Building URL for storage claim")
	}

	claim := &StorageClaim{
		id:   identity,
		url:  url,
		size: size,
	}

	return claim, nil
}

// PurgeStorageClaim _
func (sc *StorageControl) PurgeStorageClaim(ctx context.Context, claim models.IStorageClaim) error {
	_, err := sc.s3Client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(claim.GetID()),
	})

	if err != nil {
		return castAmzErr(err)
	}

	return nil
}

// GetURLForStorageClaim _
func (sc *StorageControl) GetURLForStorageClaim(ctx context.Context, identity string) (string, error) {
	req, _ := sc.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(identity),
	})

	urlStr, err := req.Presign(SignedURLDuration)

	if err != nil {
		return "", errors.Wrap(castAmzErr(err), "Signing url")
	}

	return urlStr, nil
}

// GetStorageClaimSize _
func (sc *StorageControl) GetStorageClaimSize(ctx context.Context, identity string) (int, error) {
	fmt.Printf("identity: %#v\n", identity)

	result, err := sc.s3Client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(sc.bucket),
		Key:    aws.String(identity),
	})

	if err != nil {
		return 0, errors.Wrap(castAmzErr(err), "Performing S3 HEAD request")
	}

	return int(aws.Int64Value(result.ContentLength)), nil
}

func castAmzErr(err error) error {
	if err != nil {
		if aErr, ok := err.(awserr.Error); ok {
			if aErr.Code() == s3.ErrCodeNoSuchKey {
				return errors.Wrapf(models.ErrNotFound, "No such S3 key (`%s`)", aErr.Error())
			}
		}
	}

	return err
}
