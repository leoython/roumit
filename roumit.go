package roumit

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/leoython/roumit/internal"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func Map(length, workNums int, f func(i int) error) error {
	return MapWithTimeout(length, workNums, time.Hour*24*365, func(ctx context.Context, i int) error {
		err := f(i)
		return err
	})
}

func MapWithTimeout(length, workNums int, timeout time.Duration, f func(ctx context.Context, i int) error) error {
	var allerror error
	ch := make(chan int)
	spinlock := internal.NewSpinLock()
	go func() {
		defer close(ch)
		for i := 0; i < length; i++ {
			ch <- i
		}
	}()
	errGroup := errgroup.Group{}
	for i := 0; i < workNums; i++ {
		errGroup.Go(func() error {
			for idx := range ch {
				done := make(chan struct{})
				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				defer cancel()
				go func() {
					err := f(ctx, idx)
					if err != nil {
						spinlock.Lock()
						allerror = multierror.Append(err)
						spinlock.Unlock()
					}
					done <- struct{}{}
				}()
				select {
				case <-done:
				case <-ctx.Done():
					allerror = multierror.Append(allerror, errors.Errorf("timeout at %d", idx))
				}
			}
			return nil
		})
	}
	errGroup.Wait()
	return allerror
}
