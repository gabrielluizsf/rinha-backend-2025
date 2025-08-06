package date

import (
	"time"
)

func Parse(str string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, str)
}
