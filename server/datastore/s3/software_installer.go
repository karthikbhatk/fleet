package s3

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fleetdm/fleet/v4/server/config"
	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
)

// SoftwareInstallerStore implements the fleet.SoftwareInstallerStore to store
// and retrieve software installers from S3.
type SoftwareInstallerStore struct {
	*s3store
}

// NewSoftwareInstallerStore creates a new instance with the given S3 config.
func NewSoftwareInstallerStore(config config.S3Config) (*SoftwareInstallerStore, error) {
	s3store, err := newS3store(config)
	if err != nil {
		return nil, err
	}
	return &SoftwareInstallerStore{s3store}, nil
}

// Get retrieves the requested software installer from S3.
func (i *SoftwareInstallerStore) Get(ctx context.Context, installerID string) (io.ReadCloser, int64, error) {
	key := i.keyForInstaller(installerID)

	req, err := i.s3client.GetObject(&s3.GetObjectInput{Bucket: &i.bucket, Key: &key})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey, s3.ErrCodeNoSuchBucket, "NotFound":
				return nil, int64(0), installerNotFoundError{}
			}
		}
		return nil, int64(0), ctxerr.Wrap(ctx, err, "retrieving installer from store")
	}
	return req.Body, *req.ContentLength, nil
}

// Put uploads a software installer to S3.
func (i *SoftwareInstallerStore) Put(ctx context.Context, installerID string, content io.ReadSeeker) error {
	key := i.keyForInstaller(installerID)
	_, err := i.s3client.PutObject(&s3.PutObjectInput{
		Bucket: &i.bucket,
		Body:   content,
		Key:    &key,
	})
	return err
}

// Exists checks if a software installer exists in the S3 bucket for the ID.
func (i *SoftwareInstallerStore) Exists(ctx context.Context, installerID string) (bool, error) {
	key := i.keyForInstaller(installerID)

	_, err := i.s3client.HeadObject(&s3.HeadObjectInput{Bucket: &i.bucket, Key: &key})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey, s3.ErrCodeNoSuchBucket, "NotFound":
				return false, nil
			}
		}
		return false, ctxerr.Wrap(ctx, err, "checking existence on file store")
	}
	return true, nil
}

// keyForInstaller builds an S3 key to search for the software installer.
func (i *SoftwareInstallerStore) keyForInstaller(installerID string) string {
	file := fmt.Sprintf("%s.%s", executable, installer.Kind)
	dir := ""
	if installer.Desktop {
		dir = desktopPath
	}
	return path.Join(i.prefix, installer.EnrollSecret, dir, file)
}
