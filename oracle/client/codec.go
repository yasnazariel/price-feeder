package client

import (
	kiiparams "github.com/kiichain/kiichain/v2/app/params"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	cryptocodec "github.com/cosmos/evm/crypto/codec"
)

var encodingConfig kiiparams.EncodingConfig

func init() {
	encodingConfig = kiiparams.MakeEncodingConfig()

	// Register cosmos-sdk interfaces
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// Register the pubkey EVM interface
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
}
