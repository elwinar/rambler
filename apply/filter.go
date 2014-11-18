package apply 

import (
	"fmt"
	"github.com/elwinar/rambler/migration"
)

func Filter(available, applied []*migration.Migration) ([]*migration.Migration, error) {
	var i, j int = 0, 0
	
	for i < len(available) && j < len(applied) {
		if available[i].Version < applied[j].Version {
			return nil, fmt.Errorf("out of order migration %d", available[i].Version)
		}
		
		if available[i].Version > applied[j].Version {
			return nil, fmt.Errorf("missing migration %d", applied[j].Version)
		}
		
		i++
		j++
	}
	
	if j != len(applied) {
		return nil, fmt.Errorf("missing migration %d", applied[j].Version)
	}
	
	return available[i:], nil
}
