package db

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
)

const (
	testKind = "unitTestFoo"
)

type Foo struct {
	Bar string
}

type Bar struct {
	Foo string
}

func Test_CRUD_Foo(t *testing.T) {
	skipTestsVar := "SKIP_INTEGRATION_TESTS"
	skipTests := os.Getenv(skipTestsVar)
	if len(skipTests) > 0 {
		t.Skipf("Skipping database tests because %s was set.", skipTestsVar)
	}

	ctx := context.Background()
	repo, err := NewRepo[Foo](ctx, "dusted-codes", "unittest", testKind)
	if err != nil {
		t.Error(err.Error())
	}
	foo := Foo{
		Bar: "yolo",
	}
	key := uuid.NewString()
	_, err = repo.Insert(ctx, key, &foo)
	if err != nil {
		t.Error(err.Error())
	}

	foo2, err := repo.Get(ctx, key)
	if err != nil {
		t.Error(err.Error())
	}
	if foo.Bar != foo2.Bar {
		t.Error("something went wrong")
	}

	foo2.Bar = "yada yada"
	err = repo.Upsert(ctx, key, foo2)
	if err != nil {
		t.Error(err.Error())
	}

	foos, err := repo.Query(ctx, repo.NewQuery().FilterField("Bar", "=", "yada yada"))
	if err != nil {
		t.Error(err.Error())
	}
	if len(foos) != 1 && foos[0].Bar != "yada yada" {
		t.Error("Something went wrong")
	}
	err = repo.Delete(ctx, key)
	if err != nil {
		t.Error(err.Error())
	}
	foo3, err := repo.Get(ctx, key)
	if err != nil {
		t.Error(err.Error())
	}
	if foo3 != nil {
		t.Error("Delete operation failed")
	}
}

func Test_CRUD_Bar(t *testing.T) {
	skipTestsVar := "SKIP_INTEGRATION_TESTS"
	skipTests := os.Getenv(skipTestsVar)
	if len(skipTests) > 0 {
		t.Skipf("Skipping database tests because %s was set.", skipTestsVar)
	}

	ctx := context.Background()
	repo, err := NewRepo[Bar](ctx, "dusted-codes", "unittest", testKind)
	if err != nil {
		t.Error(err.Error())
	}
	bar := &Bar{
		Foo: "yolo",
	}
	key := uuid.NewString()
	_, err = repo.Insert(ctx, key, bar)
	if err != nil {
		t.Error(err.Error())
	}

	bar2, err := repo.Get(ctx, key)
	if err != nil {
		t.Error(err.Error())
	}
	if bar.Foo != bar2.Foo {
		t.Error("something went wrong")
	}
}
