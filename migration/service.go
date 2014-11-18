package migration

import (
	"github.com/elwinar/rambler/migration/driver"
)

type Service interface {
	driver.Driver
}
