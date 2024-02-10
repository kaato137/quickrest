package filewatch

import (
	"bytes"
	"context"
	"crypto/md5"
	"io"
	"os"
	"time"
)

const defaultInterval = 1 * time.Second

type watcher struct {
	path     string
	inverval time.Duration
	checksum []byte

	onChangeFn func() error
	onErrorFn  func(error) (ignore bool)
}

func WatchFilePath(path string) *watcher {
	return &watcher{path: path, inverval: defaultInterval}
}

func (w *watcher) WithInterval(val time.Duration) *watcher {
	w.inverval = val

	return w
}

func (w *watcher) OnChange(cb func() error) *watcher {
	w.onChangeFn = cb

	return w
}

func (w *watcher) OnError(cb func(err error) bool) *watcher {
	w.onErrorFn = cb

	return w
}

func (w *watcher) Run(ctx context.Context) (err error) {
	w.checksum, err = checksumForPath(w.path)
	if err != nil {
		return err
	}

	go w.loop(ctx)

	return nil
}

func (w *watcher) loop(ctx context.Context) {
	ticker := time.NewTicker(w.inverval)
	for {
		newChecksum, err := checksumForPath(w.path)
		if err != nil {
			if w.onErrorFn != nil {
				if ignore := w.onErrorFn(err); !ignore {
					return
				}
			}
		}

		if !bytes.EqualFold(w.checksum, newChecksum) {
			if w.onChangeFn != nil {
				if err := w.onChangeFn(); err != nil {
					if ignore := w.onErrorFn(err); !ignore {
						return
					}
				}
			}

			w.checksum = newChecksum
		}

		select {
		case <-ticker.C:
			// continue
		case <-ctx.Done():
			return
		}
	}
}

func checksumForPath(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return md5.New().Sum(bytes), nil
}
