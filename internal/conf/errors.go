package conf

import "errors"

var (
	ErrDefaultPathNotFound        = errors.New("defaut path not found")
	ErrDefaultConfigAlreadyExists = errors.New("default config already exists")
)
