package descriptor

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func Identical[T any](res1, res2 T) (bool, error) {
	buffer1, err := yaml.Marshal(res1)
	if err != nil {
		return false, err
	}
	buffer2, err := yaml.Marshal(res2)
	if err != nil {
		return false, err
	}

	return bytes.Equal(buffer1, buffer2), nil
}
