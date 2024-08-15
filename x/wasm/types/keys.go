package types

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the contract module
	ModuleName = "wasm"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// TStoreKey is the string transient store representation
	TStoreKey = "transient_" + ModuleName

	// QuerierRoute is the querier route for the wasm module
	QuerierRoute = ModuleName

	// RouterKey is the msg router key for the wasm module
	RouterKey = ModuleName
)

var (
	CodeKeyPrefix                                  = []byte{0x01}
	ContractKeyPrefix                              = []byte{0x02}
	ContractStorePrefix                            = []byte{0x03}
	SequenceKeyPrefix                              = []byte{0x04}
	ContractCodeHistoryElementPrefix               = []byte{0x05}
	ContractByCodeIDAndCreatedSecondaryIndexPrefix = []byte{0x06}
	PinnedCodeIndexPrefix                          = []byte{0x07}
	TXCounterPrefix                                = []byte{0x08}
	ContractsByCreatorPrefix                       = []byte{0x09}
	ParamsKey                                      = []byte{0x10}
	AsyncAckKeyPrefix                              = []byte{0x11}
	GaslessContractIndexPrefix                     = []byte{0x0a}

	KeySequenceCodeID     = append(SequenceKeyPrefix, []byte("lastCodeId")...)
	KeySequenceInstanceID = append(SequenceKeyPrefix, []byte("lastContractId")...)
)

// GetCodeKey constructs the key for retreiving the ID for the WASM code
func GetCodeKey(codeID uint64) []byte {
	contractIDBz := sdk.Uint64ToBigEndian(codeID)
	return append(CodeKeyPrefix, contractIDBz...)
}

// GetContractAddressKey returns the key for the WASM contract instance
func GetContractAddressKey(addr sdk.AccAddress) []byte {
	return append(ContractKeyPrefix, addr...)
}

// GetContractsByCreatorPrefix returns the contracts by creator prefix for the WASM contract instance
func GetContractsByCreatorPrefix(addr sdk.AccAddress) []byte {
	bz := address.MustLengthPrefix(addr)
	return append(ContractsByCreatorPrefix, bz...)
}

// GetContractStorePrefix returns the store prefix for the WASM contract instance
func GetContractStorePrefix(addr sdk.AccAddress) []byte {
	return append(ContractStorePrefix, addr...)
}

// GetAsyncPacketKey returns the key for a packet that is acknowledged asynchronously
func GetAsyncPacketKey(destChannel string, sequence uint64) []byte {
	// key is a concatenation of length-prefixed destination channel and sequence
	channel := []byte(destChannel)
	channelLen := make([]byte, 4)
	binary.BigEndian.PutUint32(channelLen, uint32(len(channel)))
	seq := make([]byte, 8)
	binary.BigEndian.PutUint64(seq, sequence)

	return append(append(channelLen, channel...), seq...)
}

// GetAsyncAckStorePrefix returns the store prefix for packets that are acknowledged asynchronously
func GetAsyncAckStorePrefix(portID string) []byte {
	return append(AsyncAckKeyPrefix, portID...)
}

// GetContractByCreatedSecondaryIndexKey returns the key for the secondary index:
// `<prefix><codeID><created/last-migrated><contractAddr>`
func GetContractByCreatedSecondaryIndexKey(contractAddr sdk.AccAddress, c ContractCodeHistoryEntry) []byte {
	prefix := GetContractByCodeIDSecondaryIndexPrefix(c.CodeID)
	prefixLen := len(prefix)
	contractAddrLen := len(contractAddr)
	r := make([]byte, prefixLen+AbsoluteTxPositionLen+contractAddrLen)
	copy(r[0:], prefix)
	copy(r[prefixLen:], c.Updated.Bytes())
	copy(r[prefixLen+AbsoluteTxPositionLen:], contractAddr)
	return r
}

// GetContractByCodeIDSecondaryIndexPrefix returns the prefix for the second index: `<prefix><codeID>`
func GetContractByCodeIDSecondaryIndexPrefix(codeID uint64) []byte {
	prefixLen := len(ContractByCodeIDAndCreatedSecondaryIndexPrefix)
	const codeIDLen = 8
	r := make([]byte, prefixLen+codeIDLen)
	copy(r[0:], ContractByCodeIDAndCreatedSecondaryIndexPrefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(codeID))
	return r
}

// GetContractByCreatorSecondaryIndexKey returns the key for the second index: `<prefix><creatorAddress length><created time><creatorAddress><contractAddr>`
func GetContractByCreatorSecondaryIndexKey(bz, position []byte, contractAddr sdk.AccAddress) []byte {
	prefixBytes := GetContractsByCreatorPrefix(bz)
	lenPrefixBytes := len(prefixBytes)
	r := make([]byte, lenPrefixBytes+AbsoluteTxPositionLen+len(contractAddr))

	copy(r[:lenPrefixBytes], prefixBytes)
	copy(r[lenPrefixBytes:lenPrefixBytes+AbsoluteTxPositionLen], position)
	copy(r[lenPrefixBytes+AbsoluteTxPositionLen:], contractAddr)

	return r
}

// GetContractCodeHistoryElementKey returns the key a contract code history entry: `<prefix><contractAddr><position>`
func GetContractCodeHistoryElementKey(contractAddr sdk.AccAddress, pos uint64) []byte {
	prefix := GetContractCodeHistoryElementPrefix(contractAddr)
	prefixLen := len(prefix)
	r := make([]byte, prefixLen+8)
	copy(r[0:], prefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(pos))
	return r
}

// GetContractCodeHistoryElementPrefix returns the key prefix for a contract code history entry: `<prefix><contractAddr>`
func GetContractCodeHistoryElementPrefix(contractAddr sdk.AccAddress) []byte {
	prefixLen := len(ContractCodeHistoryElementPrefix)
	contractAddrLen := len(contractAddr)
	r := make([]byte, prefixLen+contractAddrLen)
	copy(r[0:], ContractCodeHistoryElementPrefix)
	copy(r[prefixLen:], contractAddr)
	return r
}

// GetPinnedCodeIndexPrefix returns the key prefix for a code id pinned into the wasmvm cache
func GetPinnedCodeIndexPrefix(codeID uint64) []byte {
	prefixLen := len(PinnedCodeIndexPrefix)
	r := make([]byte, prefixLen+8)
	copy(r[0:], PinnedCodeIndexPrefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(codeID))
	return r
}

// GetGaslessContractIndexPrefix returns the key prefix for a gasless contract into the wasmvm cache
func GetGaslessContractIndexPrefix(contractAddr sdk.AccAddress) []byte {
	prefixLen := len(PinnedCodeIndexPrefix)
	contractAddrLen := len(contractAddr)
	r := make([]byte, prefixLen+contractAddrLen)
	copy(r[0:], GaslessContractIndexPrefix)
	copy(r[prefixLen:], contractAddr)
	return r
}

// ParsePinnedCodeIndex converts the serialized code ID back.
func ParsePinnedCodeIndex(s []byte) uint64 {
	return sdk.BigEndianToUint64(s)
}
