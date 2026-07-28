package assembly

import "time"

type service struct {
	buildTime time.Duration
}

func NewService(buildTime time.Duration) Assembler {
	return &service{
		buildTime: buildTime,
	}
}
