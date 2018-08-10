package dcrsqlite

import (
	"bytes"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
	"github.com/decred/dcrdata/testutil"
)

func testBlock(block *dcrutil.Block) {
	testutil.SaveBlockToFile(block, "testdata")

	if 1 == 1 {
		return
	}

	data, err := block.MsgBlock().Bytes()
	if err != nil {
		testutil.Log(" failed", err)
	}

	var msgBlock wire.MsgBlock
	err = msgBlock.Deserialize(bytes.NewReader(data))
	if err != nil {
		testutil.Log(" failed", err)
	}
	blockR := dcrutil.NewBlock(&msgBlock)

	compare(block, blockR)
}
func compare(a *dcrutil.Block, b *dcrutil.Block) {
	a.STransactions()
}
