package dcrsqlite

import (
	"github.com/decred/dcrdata/rpcutils"
	"math"
)

func ResyncDB(
	baseDB wiredDB, quit chan struct{}, blockGetter rpcutils.BlockGetter) (
	int64, error) {
	h, err := baseDB.resyncDB(quit, blockGetter, int64(math.MaxInt32))
	return h, err
}
