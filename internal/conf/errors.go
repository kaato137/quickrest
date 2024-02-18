package conf

import "errors"

var (
	ErrConfigFileNotExist         = errors.New("file not exist")
	ErrDefaultPathNotFound        = errors.New("defaut path not found")
	ErrDefaultConfigAlreadyExists = errors.New("default config already exists")
)
