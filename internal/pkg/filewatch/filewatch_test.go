package filewatch

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const defaultTimeout = 30 * time.Second

func TestWatchFilePath(t *testing.T) {
	t.Run("Should detect change in file", func(t *testing.T) {
		file, clean := createTempFile(t)
		defer clean()

		var isChangeDetected bool
		isChangeDetectedFn := func() bool {
			return isChangeDetected
		}

		err := WatchFilePath(file.Name()).
			OnChange(func() error {
				isChangeDetected = true
				return nil
			}).
			Run(context.Background())

		require.NoError(t, err)

		changeFile(t, file)

		require.Eventually(t,
			isChangeDetectedFn,
			defaultTimeout,
			time.Second,
			"file change should be detected",
		)
	})

	t.Run("Call OnError callback when file was deleted", func(t *testing.T) {
		file, close := createTempFile(t)
		defer close()

		var isErrorOccurred bool
		var theErrorThatOccurred error
		isErrorOccurredFn := func() bool {
			return isErrorOccurred
		}

		err := WatchFilePath(file.Name()).
			OnError(func(err error) bool {
				isErrorOccurred = true
				theErrorThatOccurred = err

				return false
			}).
			Run(context.Background())

		require.NoError(t, err)

		deleteFile(t, file)

		require.Eventually(t,
			isErrorOccurredFn,
			defaultTimeout,
			time.Second,
			"error should occur",
		)

		require.ErrorIs(t, theErrorThatOccurred, os.ErrNotExist)
	})

	t.Run("Inverval option is working", func(t *testing.T) {
		file, close := createTempFile(t)
		defer close()

		theInterval := time.Hour

		err := WatchFilePath(file.Name()).
			WithInterval(theInterval).
			OnChange(func() error {
				require.FailNow(t,
					"change of file should never be detected "+
						"given such long check interval")
				return nil
			}).
			Run(context.Background())

		require.NoError(t, err)

		leeway := 2 * time.Second
		time.Sleep(defaultInterval + leeway)
	})
}

func createTempFile(t *testing.T) (*os.File, func()) {
	t.Helper()

	file, err := os.CreateTemp("", "temp.*.txt")
	require.NoError(t, err)
	return file, func() {
		_ = file.Close()
		_ = os.Remove(file.Name())
	}
}

func changeFile(t *testing.T, f *os.File) {
	t.Helper()

	_, err := f.WriteString("1\n")
	require.NoError(t, err)
}

func deleteFile(t *testing.T, f *os.File) {
	t.Helper()

	err := os.Remove(f.Name())
	require.NoError(t, err)
}
