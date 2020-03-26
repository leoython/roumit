package roumit

import (
	"context"
	"testing"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

func TestMap(t *testing.T) {
	k := [10]int{}
	Map(10, 2, func(i int) error {
		k[i] = i + 1
		return nil
	})
	for i := 0; i < 10; i++ {
		if k[i] != i+1 {
			t.Error("id", i, "not working")
		}
	}
}

func TestMapWithTimeout(t *testing.T) {
	k := [10]int{}
	err := MapWithTimeout(10, 2, time.Millisecond*50, func(ctx context.Context, i int) error {
		if i%2 == 0 {
			time.Sleep(time.Millisecond * 100)
		}
		if ctx.Err() == nil {
			k[i] = i + 1
		}
		return nil
	})
	if err == nil {
		t.Error("No timeout error returned")
	}
	merr, ok := (err).(*multierror.Error)
	if !ok {
		t.Error("return error is not of type multierror.Error")
	}
	if len(merr.Errors) != 5 {
		t.Error("Timed out error num not correct, expected 5, actual", len(merr.Errors))
	}
	for i := 0; i < 10; i++ {
		if i%2 == 0 && k[i] != 0 || i%2 == 1 && k[i] != i+1 {
			t.Error("Id", i, "not working")
		}
	}
}
