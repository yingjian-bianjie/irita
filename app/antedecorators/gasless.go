package antedecorators

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkacltypes "github.com/cosmos/cosmos-sdk/types/accesscontrol"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type GaslessDecorator struct {
	wrapped []sdk.AnteFullDecorator
}

func (gd GaslessDecorator) AnteDeps(txDeps []sdkacltypes.AccessOperation, tx sdk.Tx, next sdk.AnteDepGenerator) (newTxDeps []sdkacltypes.AccessOperation, err error) {
	deps := []sdkacltypes.AccessOperation{}
	terminatorDeps := func(txDeps []sdkacltypes.AccessOperation, _ sdk.Tx) ([]sdkacltypes.AccessOperation, error) {
		return txDeps, nil
	}
	for _, depGen := range gd.wrapped {
		deps, _ = depGen.AnteDeps(deps, tx, terminatorDeps)
	}
	return next(append(txDeps, deps...), tx)
}

func NewGaslessDecorator(wrapped []sdk.AnteFullDecorator) GaslessDecorator {
	return GaslessDecorator{wrapped: wrapped}
}

func (gd GaslessDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	originalGasMeter := ctx.GasMeter()
	// eagerly set infinite gas meter so that queries performed by isTxGasless will not incur gas cost
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	gas := feeTx.GetGas()
	// If non-zero gas limit is provided by the TX, we then consider it exempt from the gasless TX, and then prioritize it accordingly
	if err != nil {
		return ctx, err
	}
	if gas > 0 {
		ctx = ctx.WithGasMeter(originalGasMeter)
		// if not gasless, then we use the wrappers

		// AnteHandle always takes a `next` so we need a no-op to execute only one handler at a time
		terminatorHandler := func(ctx sdk.Context, _ sdk.Tx, _ bool) (sdk.Context, error) {
			return ctx, nil
		}
		// iterating instead of recursing the handler for readability
		for _, handler := range gd.wrapped {
			ctx, err = handler.AnteHandle(ctx, tx, simulate, terminatorHandler)
			if err != nil {
				return ctx, err
			}
		}
		return next(ctx, tx, simulate)
	}

	return next(ctx, tx, simulate)
}
