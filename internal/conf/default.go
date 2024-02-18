package conf

import (
	"errors"
	"fmt"
	"os"
)

const defaultConfig = `
addr: 127.0.0.1:8090

routes:
- path: GET /articles/{id}
  body: |
    {
      "id": "{id}",
      "title": "Beans are now free!"
    }
`

func GenerateDefault() error {
	_, err := findDefaultPaths()
	if err != nil && !errors.Is(err, ErrDefaultPathNotFound) {
		return err
	} else if err == nil {
		return ErrDefaultConfigAlreadyExists
	}

	f, err := os.Create(defaultPaths[0])
	if err != nil {
		return fmt.Errorf("create default config: %w", err)
	}

	if _, err = f.WriteString(defaultConfig); err != nil {
		return fmt.Errorf("write default config: %w", err)
	}

	return nil
}

func findDefaultPaths() (string, error) {
	for _, p := range defaultPaths {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", ErrDefaultPathNotFound
}
