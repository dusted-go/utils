package db

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dusted-go/fault/fault"
)

// Repo exposes read and write operations to Google Cloud Datastore.
type Repo[T any] struct {
	client    *datastore.Client
	namespace string
	kind      string
}

// NewRepo creates a new instance of Repo.
func NewRepo[T any](
	ctx context.Context,
	projectID,
	namespace,
	kind string,
) (
	*Repo[T],
	error,
) {
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil,
			fault.SystemWrap(err, "db", "NewRepo", "failed to create Google Cloud Datastore Repo")
	}
	return &Repo[T]{
		client:    client,
		namespace: namespace,
		kind:      kind}, nil
}

// NewQuery creates a new *datastore.Query for the current kind.
func (r *Repo[T]) NewQuery() *datastore.Query {
	return datastore.NewQuery(r.kind)
}

// Upsert creates a new entity or updates an existing one in GCP Datastore.
// The entity should be a struct pointer.
func (r *Repo[T]) Upsert(ctx context.Context, key string, entity *T) error {
	k := datastore.NameKey(r.kind, key, nil)
	k.Namespace = r.namespace
	if _, err := r.client.Put(ctx, k, entity); err != nil {
		return fault.SystemWrap(err, "db", "Upsert", "error writing to Google Cloud Datastore")
	}
	return nil
}

// Insert creates a new entity in GCP Datastore or fails with an error.
// The entity should be a struct pointer.
func (r *Repo[T]) Insert(ctx context.Context, key string, entity *T) (duplicate bool, err error) {
	k := datastore.NameKey(r.kind, key, nil)
	k.Namespace = r.namespace
	insert := datastore.NewInsert(k, entity)
	_, dbErr := r.client.Mutate(ctx, insert)
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
// The function will return nil if the entity cannot be found or an error has occurred.
func (r *Repo[T]) Get(ctx context.Context, key string) (*T, error) {
	k := datastore.NameKey(r.kind, key, nil)
	k.Namespace = r.namespace
	var entity T
	if err := r.client.Get(ctx, k, &entity); err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		}
		return nil, fault.SystemWrap(err, "db", "Get", "error reading from Google Cloud Datastore")
	}
	return &entity, nil
}

// Query finds all entities which match the given query.
func (r *Repo[T]) Query(
	ctx context.Context,
	query *datastore.Query) ([]*T, error) {
	q := query.Namespace(r.namespace)
	var entities []*T
	if _, err := r.client.GetAll(ctx, q, &entities); err != nil {
		return nil, fault.SystemWrap(err, "db", "Query", "error reading from Google Cloud Datastore")
	}
	return entities, nil
}

// Count returns the total count of items resulting from a given query.
func (r *Repo[T]) Count(
	ctx context.Context,
	query *datastore.Query) (int, error) {
	q := query.Namespace(r.namespace).KeysOnly()

	keys, err := r.client.GetAll(ctx, q, nil)
	if err != nil {
		return -1,
			fault.SystemWrap(err, "db", "Count", "error reading from Google Cloud Datastore")
	}

	return len(keys), nil
}

func (r *Repo[T]) Delete(ctx context.Context, key string) error {
	k := datastore.NameKey(r.kind, key, nil)
	k.Namespace = r.namespace
	if err := r.client.Delete(ctx, k); err != nil {
		return fault.SystemWrap(err, "db", "Delete", "error deleting entity in Google Cloud Datastore")
	}
	return nil
}
