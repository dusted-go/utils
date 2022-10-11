package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"cloud.google.com/go/storage"
)

type Client struct {
	client *storage.Client
}

func NewClient(ctx context.Context) (*Client, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating Google Cloud Storage client: %w", err)
	}
	return &Client{client: client}, nil
}

func (c *Client) PutFile(
	ctx context.Context,
	bucketName string,
	fileName string,
	file multipart.File,
	mimeType string,
	cacheControl string,
	aclEntity string,
	aclRole string) error {

	obj := c.client.Bucket(bucketName).Object(fileName)
	_, err := obj.Attrs(ctx)

	if err != nil {
		if errors.Is(err, storage.ErrObjectNotExist) {
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()

			w := obj.NewWriter(ctx)
			w.CacheControl = cacheControl
			w.ContentType = mimeType

			if _, err := io.Copy(w, file); err != nil {
				return fmt.Errorf("error writing file to Google Cloud Storage bucket '%s': %w", bucketName, err)
			}

			if err := w.Close(); err != nil {
				return fmt.Errorf("error closing *storage.Writer: %w", err)
			}

			if err := obj.ACL().Set(
				ctx,
				storage.ACLEntity(aclEntity),
				storage.ACLRole(aclRole)); err != nil {
				return fmt.Errorf("error setting ACL on Google Cloud Storage object: %w", err)
			}
			return nil
		}

		return fmt.Errorf("error retrieving object's metadata from Google Cloud Storage: %w", err)
	}

	return nil
}
