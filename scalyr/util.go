package scalyr

import (
	"encoding/json"
	"fmt"
)

func TransformType[K []any, V any](input K) ([]V, error) {
	items := make([]V, 0, len(input))

	for _, value := range input {
		var item V

		data := []byte(fmt.Sprint(value))

		if err := json.Unmarshal(data, &item); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
