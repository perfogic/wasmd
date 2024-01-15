package legacy

import (
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterInterfaces(registry types.InterfaceRegistry) {
	// support legacy cosmwasm
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgStoreCode{},
		&MsgInstantiateContract{},
		&MsgExecuteContract{},
		&MsgMigrateContract{},
		&MsgUpdateAdmin{},
		&MsgClearAdmin{},
	)

	// // support legacy proposal querying
	// registry.RegisterImplementations(
	// 	(*v1beta1.Content)(nil),
	// 	&StoreCodeProposal{},
	// 	&InstantiateContractProposal{},
	// 	&MigrateContractProposal{},
	// 	&UpdateAdminProposal{},
	// 	&ClearAdminProposal{},
	// 	&PinCodesProposal{},
	// 	&UnpinCodesProposal{},
	// )
}
