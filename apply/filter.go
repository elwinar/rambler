package apply

import (
	"fmt"
)

// Filter return a slice containing the migrations version to apply (those not applied
// already)
func Filter(available, applied []uint64) ([]uint64, error) {
	var i, j int = 0, 0

	for i < len(available) && j < len(applied) {
		if available[i] < applied[j] {
			return nil, fmt.Errorf(errOutOfOrder, available[i])
		}

		if available[i] > applied[j] {
			return nil, fmt.Errorf(errMissingMigration, applied[j])
		}

		i++
		j++
	}

	if j != len(applied) {
		return nil, fmt.Errorf(errMissingMigration, applied[j])
	}

	return available[i:], nil
}
