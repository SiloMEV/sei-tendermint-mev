package types

import (
	"crypto/sha256"
	"errors"
	"fmt"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

// ErrTxInCache is returned to the client if we saw tx earlier
var ErrTxInCache = errors.New("tx already exists in cache")

// TxKey is the fixed length array key used as an index.
type TxKey [sha256.Size]byte

// ToProto converts Data to protobuf
func (txKey *TxKey) ToProto() *tmproto.TxKey {
	tp := new(tmproto.TxKey)

	txBzs := make([]byte, len(txKey))
	if len(txKey) > 0 {
		for i := range txKey {
			txBzs[i] = txKey[i]
		}
		tp.TxKey = txBzs
	}

	return tp
}

// TxKeyFromProto takes a protobuf representation of TxKey &
// returns the native type.
func TxKeyFromProto(dp *tmproto.TxKey) (TxKey, error) {
	if dp == nil {
		return TxKey{}, errors.New("nil data")
	}
	var txBzs [sha256.Size]byte
	for i := range dp.TxKey {
		txBzs[i] = dp.TxKey[i]
	}

	return txBzs, nil
}

func TxKeysListFromProto(dps []*tmproto.TxKey) ([]TxKey, error) {
	var txKeys []TxKey
	for _, txKey := range dps {
		txKey, err := TxKeyFromProto(txKey)
		if err != nil {
			return nil, err
		}
		txKeys = append(txKeys, txKey)
	}
	return txKeys, nil
}

// ErrWrongHeight means the tx is asking to be in a height that doesn't match the current auction
type ErrWrongHeight struct {
	desiredHeight        int
	currentAuctionHeight int
}

func (e ErrWrongHeight) Error() string {
	return fmt.Sprintf("Tx submitted for wrong height, asked for %d, but current auction height is %d", e.desiredHeight, e.currentAuctionHeight)
}

// ErrBundleFull means the tx is trying to enter a bundle that has already reached its limit
type ErrBundleFull struct {
	bundleId     int64
	bundleHeight int64
}

func (e ErrBundleFull) Error() string {
	return fmt.Sprintf("Tx submitted but bundle is full, for bundleId %d with bundle size %d", e.bundleId, e.bundleHeight)
}

// ErrTxMalformedForBundle is a general malformed error for specific cases
type ErrTxMalformedForBundle struct {
	bundleId     int64
	bundleSize   int64
	bundleHeight int64
	bundleOrder  int64
}

func (e ErrTxMalformedForBundle) Error() string {
	return fmt.Sprintf("Tx submitted but malformed with respect to bundling, for bundleId %d, at height %d, with bundleSize %d, and bundleOrder %d", e.bundleId, e.bundleHeight, e.bundleSize, e.bundleOrder)
}

// ErrTxTooLarge defines an error when a transaction is too big to be sent in a
// message to other peers.
type ErrTxTooLarge struct {
	Max    int
	Actual int
}

func (e ErrTxTooLarge) Error() string {
	return fmt.Sprintf("Tx too large. Max size is %d, but got %d", e.Max, e.Actual)
}

// ErrMempoolIsFull defines an error where Tendermint and the application cannot
// handle that much load.
type ErrMempoolIsFull struct {
	NumTxs      int
	MaxTxs      int
	TxsBytes    int64
	MaxTxsBytes int64
}

func (e ErrMempoolIsFull) Error() string {
	return fmt.Sprintf(
		"mempool is full: number of txs %d (max: %d), total txs bytes %d (max: %d)",
		e.NumTxs,
		e.MaxTxs,
		e.TxsBytes,
		e.MaxTxsBytes,
	)
}

// ErrMempoolPendingIsFull defines an error where there are too many pending transactions
// not processed yet
type ErrMempoolPendingIsFull struct {
	NumTxs      int
	MaxTxs      int
	TxsBytes    int64
	MaxTxsBytes int64
}

func (e ErrMempoolPendingIsFull) Error() string {
	return fmt.Sprintf(
		"mempool pending set is full: number of txs %d (max: %d), total txs bytes %d (max: %d)",
		e.NumTxs,
		e.MaxTxs,
		e.TxsBytes,
		e.MaxTxsBytes,
	)
}

// ErrPreCheck defines an error where a transaction fails a pre-check.
type ErrPreCheck struct {
	Reason error
}

func (e ErrPreCheck) Error() string {
	return e.Reason.Error()
}

// IsPreCheckError returns true if err is due to pre check failure.
func IsPreCheckError(err error) bool {
	return errors.As(err, &ErrPreCheck{})
}
