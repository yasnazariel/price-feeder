package client

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	cryptocodec "github.com/cosmos/evm/crypto/codec"
	kiiparams "github.com/kiichain/kiichain/v2/app/params"
)

var (
	encodingConfig kiiparams.EncodingConfig
)

func init() {
	encodingConfig = kiiparams.MakeEncodingConfig()

	// Register cosmos-sdk interfaces
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// Register the pubkey EVM interface
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
}
