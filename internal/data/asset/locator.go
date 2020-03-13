package asset

import (
	"io"
)

type Locator interface {
	Open(uri string) (io.ReadCloser, error)
}
