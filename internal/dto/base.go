package dto

import (
	"fmt"
)

func ParseID(idParam string) (int, error) {
	var id int
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return id, nil
}
