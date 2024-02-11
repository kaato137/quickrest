package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

type RequestRecorder struct {
	recordPath string
	files      map[string]*os.File
}

func NewRequestRecorder(recordPath string) *RequestRecorder {
	return &RequestRecorder{
		recordPath: recordPath,
		files:      make(map[string]*os.File),
	}
}

func (rec *RequestRecorder) Record(name string, r *http.Request) error {
	f, err := rec.openOrCreateFile(name)
	if err != nil {
		return err
	}

	if err := rec.writeRequest(f, r); err != nil {
		return fmt.Errorf("write req: %w", err)
	}

	return nil
}

func (rec *RequestRecorder) Close() error {
	var compoundErr error
	for i := range rec.files {
		if err := rec.files[i].Close(); err != nil {
			compoundErr = errors.Join(compoundErr, err)
		}
	}

	return compoundErr
}

func (rec *RequestRecorder) writeRequest(f *os.File, r *http.Request) error {
	if _, err := fmt.Fprintf(f, "%s %s\n", r.Method, r.URL.String()); err != nil {
		return err
	}

	if _, err := io.Copy(f, r.Body); err != nil {
		return err
	}

	if _, err := f.WriteString("\n\n"); err != nil {
		return err
	}

	return nil
}

func (rec *RequestRecorder) openOrCreateFile(name string) (*os.File, error) {
	fullPath := path.Join(rec.recordPath, name)

	file, exists := rec.files[fullPath]
	if exists {
		return file, nil
	}

	file, err := rec.createFileWithPath(fullPath)
	if err != nil {
		return nil, err
	}

	rec.files[fullPath] = file

	return file, nil
}

func (rec *RequestRecorder) createFileWithPath(name string) (*os.File, error) {
	if _, err := os.Stat(name); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		fileDir := path.Dir(name)
		if err := os.MkdirAll(fileDir, os.FileMode(0700)); err != nil {
			return nil, err
		}
	}

	file, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	return file, nil
}
