package storage

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
	"github.com/dusted-go/utils/fault"
)

type Service struct {
	client *storage.Client
}

func NewService(ctx context.Context) (*Service, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil,
			fault.SystemWrap(err, "storage", "NewService", "failed to create Google Cloud Storage client")
	}
	return &Service{client: client}, nil
}

func (s *Service) PutFile(
	ctx context.Context,
	bucketName string,
	fileName string,
	file multipart.File,
	mimeType string,
	cacheControl string,
	aclEntity string,
	aclRole string) error {

	obj := s.client.Bucket(bucketName).Object(fileName)
	_, err := obj.Attrs(ctx)

	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			w := obj.NewWriter(ctx)
			w.CacheControl = cacheControl
			w.ContentType = mimeType

			if _, err := io.Copy(w, file); err != nil {
				return fault.SystemWrap(err, "storage", "PutFile", "failed to write file to Google Cloud Storage bucket")
			}

			if err := w.Close(); err != nil {
				return fault.SystemWrap(err, "storage", "PutFile", "failed to close *storage.Writer after writing file to Google Cloud Storage bucket")
			}

			if err := obj.ACL().Set(
				ctx,
				storage.ACLEntity(aclEntity),
				storage.ACLRole(aclRole)); err != nil {
				return fault.SystemWrap(err, "storage", "PutFile", "failed to set ACL on Google Cloud Storage object")
			}
			return nil
		}

		return fault.SystemWrap(err, "storage", "PutFile", "failed to check if file exists in Google Cloud Storage")
	}

	return nil
}
