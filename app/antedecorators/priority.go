package antedecorators

import sdk "github.com/cosmos/cosmos-sdk/types"

type PriorityDecorator struct{}

func NewPriorityDecorator() PriorityDecorator {
	return PriorityDecorator{}
}

// Assigns higher priority to certain types of transactions including oracle
func (pd PriorityDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Cap priority to MAXINT64 - 1000
	// Use higher priorities for tiers including oracle tx's
	return next(ctx, tx, simulate)
}
