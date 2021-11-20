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

// Service exposes read and write operations to Google Cloud Datastore.
type Service struct {
	client    *datastore.Client
	namespace string
}

// NewService creates a new instance of Service.
func NewService(projectID string, namespace string) (*Service, error) {

	client, err := datastore.NewClient(context.Background(), projectID)
	if err != nil {
		return nil,
			fault.SystemWrap("db", "NewService", "failed to create Google Cloud Datastore client", err)
	}
	return &Service{
		client:    client,
		namespace: namespace}, nil
}

// Upsert creates a new or updates an existing entity in GCP Datastore.
func (svc *Service) Upsert(ctx context.Context, e Entity) error {

	key := datastore.NameKey(e.Kind(), e.ID(), nil)
	key.Namespace = svc.namespace
	if _, err := svc.client.Put(ctx, key, e); err != nil {
		return fault.SystemWrap("db", "PutEntity", "error writing to Google Cloud Datastore", err)
	}
	return nil
}

// Insert creates a new entity in GCP Datastore or fails with an error.
func (svc *Service) Insert(ctx context.Context, e Entity) (alreadyExists bool, err error) {

	key := datastore.NameKey(e.Kind(), e.ID(), nil)
	key.Namespace = svc.namespace
	insert := datastore.NewInsert(key, e)
	_, dbErr := svc.client.Mutate(ctx, insert)
	if dbErr != nil {
		// Get the underlying GRPC error if it's been wrapped as a MultiError
		// nolint: errorlint
		if multiErr, ok := dbErr.(datastore.MultiError); ok {
			dbErr = multiErr[0]
		}
		if status.Code(dbErr) == codes.AlreadyExists {
			return true, nil
		}

		return false, fault.SystemWrap("db", "InsertEntity", "error writing to Google Cloud Datastore", dbErr)
	}
	return false, nil
}

// Get loads the single entity which matches the kind and key of the given object.
func (svc *Service) Get(ctx context.Context, e Entity) error {

	key := datastore.NameKey(e.Kind(), e.ID(), nil)
	key.Namespace = svc.namespace
	if err := svc.client.Get(ctx, key, e); err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil
		}

		return fault.SystemWrap("db", "GetEntity", "error reading from Google Cloud Datastore", err)
	}

	return nil
}

// Query finds all entities which match the given query.
func (svc *Service) Query(
	ctx context.Context,
	query *datastore.Query,
	entities interface{}) error {
	q := query.Namespace(svc.namespace)

	if _, err := svc.client.GetAll(ctx, q, entities); err != nil {
		return fault.SystemWrap("db", "QueryEntities", "error reading from Google Cloud Datastore", err)
	}

	return nil
}

// Count returns the total count of items resulting from a given query.
func (svc *Service) Count(
	ctx context.Context,
	query *datastore.Query) (int, error) {
	q := query.Namespace(svc.namespace).KeysOnly()

	keys, err := svc.client.GetAll(ctx, q, nil)
	if err != nil {
		return -1,
			fault.SystemWrap("db", "Count", "error reading from Google Cloud Datastore", err)
	}

	return len(keys), nil
}
