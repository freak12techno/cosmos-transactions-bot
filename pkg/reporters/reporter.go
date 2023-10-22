package reporters

import (
	"main/pkg/types"
)

type Reporter interface {
	Init()
	Name() string
	Enabled() bool
	Send(report types.Report) error
}
