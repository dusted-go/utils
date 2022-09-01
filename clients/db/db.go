package db

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dusted-go/utils/fault"
)

// Entity represents a single entity stored in Google Cloud Datastore.
type Entity interface {
	Kind() string
	ID() string
}

// Client exposes read and write operations to Google Cloud Datastore.
type Client[T Entity] struct {
	client    *datastore.Client
	namespace string
}

// NewClient creates a new instance of Client.
func NewClient[T Entity](ctx context.Context, projectID string, namespace string) (*Client[T], error) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil,
			fault.SystemWrap(err, "db", "NewClient", "failed to create Google Cloud Datastore client")
	}
	return &Client[T]{
		client:    client,
		namespace: namespace}, nil
}

// Upsert creates a new or updates an existing entity in GCP Datastore.
func (c *Client[T]) Upsert(ctx context.Context, entity T) error {

	key := datastore.NameKey(entity.Kind(), entity.ID(), nil)
	key.Namespace = c.namespace
	if _, err := c.client.Put(ctx, key, entity); err != nil {
		return fault.SystemWrap(err, "db", "PutEntity", "error writing to Google Cloud Datastore")
	}
	return nil
}

// Insert creates a new entity in GCP Datastore or fails with an error.
func (c *Client[T]) Insert(ctx context.Context, entity T) (alreadyExists bool, err error) {

	key := datastore.NameKey(entity.Kind(), entity.ID(), nil)
	key.Namespace = c.namespace
	insert := datastore.NewInsert(key, entity)
	_, dbErr := c.client.Mutate(ctx, insert)
	if dbErr != nil {
		// Get the underlying GRPC error if it's been wrapped as a MultiError
		// nolint: errorlint
		if multiErr, ok := dbErr.(datastore.MultiError); ok {
			dbErr = multiErr[0]
		}
		if status.Code(dbErr) == codes.AlreadyExists {
			return true, nil
		}

		return false, fault.SystemWrap(dbErr, "db", "Insert", "error writing to Google Cloud Datastore")
	}
	return false, nil
}

// Get loads the single entity which matches the kind and key of the given object.
// The function will return false if the entity cannot be found or an error has occurred.
func (c *Client[T]) Get(ctx context.Context, kind, id string) (*T, error) {
	key := datastore.NameKey(kind, id, nil)
	key.Namespace = c.namespace
	var entity T
	if err := c.client.Get(ctx, key, &entity); err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		}
		return nil, fault.SystemWrap(err, "db", "Get", "error reading from Google Cloud Datastore")
	}
	return &entity, nil
}

// Query finds all entities which match the given query.
func (c *Client[T]) Query(
	ctx context.Context,
	query *datastore.Query) ([]T, error) {
	q := query.Namespace(c.namespace)
	var entities []T
	if _, err := c.client.GetAll(ctx, q, &entities); err != nil {
		return nil, fault.SystemWrap(err, "db", "Query", "error reading from Google Cloud Datastore")
	}
	return entities, nil
}

// Count returns the total count of items resulting from a given query.
func (c *Client[T]) Count(
	ctx context.Context,
	query *datastore.Query) (int, error) {
	q := query.Namespace(c.namespace).KeysOnly()

	keys, err := c.client.GetAll(ctx, q, nil)
	if err != nil {
		return -1,
			fault.SystemWrap(err, "db", "Count", "error reading from Google Cloud Datastore")
	}

	return len(keys), nil
}
