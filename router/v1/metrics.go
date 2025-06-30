package v1

import "github.com/cosmos/cosmos-sdk/telemetry"

// Metrics defines the interface for gathering telemetry metrics.
type Metrics interface {
	Gather(format string) (telemetry.GatherResponse, error)
}
