//go:build !darwin && !linux

package pdf

import "errors"

func execLookPath(name string) (string, error) {
	return "", errors.New("unsupported platform")
}
