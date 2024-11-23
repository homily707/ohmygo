package util

import (
	"errors"
	"time"

	myerrors "github.com/homily707/ohmygo/pkg/errors"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
)

var threeTimesBackoff = wait.Backoff{
	Steps:    3,
	Duration: 500 * time.Millisecond,
	Factor:   5.0,
	Jitter:   0.1,
}

func RetryNoMatterWhat(f func() error) error {
	return retry.OnError(threeTimesBackoff, func(err error) bool { return true }, f)
}

func Retry(f func() error) error {
	return retry.OnError(threeTimesBackoff, func(err error) bool {
		return errors.Is(err, myerrors.ErrRetriable)
	}, f)
}

func RetryWithBackoff(backoff wait.Backoff, f func() error) error {
	return retry.OnError(backoff, func(err error) bool {
		return errors.Is(err, myerrors.ErrRetriable)
	}, f)
}
