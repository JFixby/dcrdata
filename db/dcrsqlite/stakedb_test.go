package dcrsqlite

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/decred/dcrd/blockchain/stake"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/database"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrdata/testutil"
	"github.com/decred/dcrdata/txhelpers"
)

func TestStakeDB(t *testing.T) {
	testutil.BindCurrentTestSetup(t)

	tag := "blocks_0-5000"
	tdb := testutil.LoadTestDataHandler(tag)
	firstBlockIndex := tdb.GetFirstBlockNumber()
	lastBlockIndex := tdb.GetLastBlockNumber()

	//load blocks
	testBlockchain := testutil.PreLoadBlocks(tag, firstBlockIndex, lastBlockIndex)

	params := &chaincfg.MainNetParams
	bestNode := produceGenesisNode(params)

	// test
	for i := firstBlockIndex; i <= lastBlockIndex; i++ {
		index := int64(i)
		block := testBlockchain[index]
		bestNode = connectBlock(block, bestNode, testBlockchain, params)
	}

}
func produceGenesisNode(params *chaincfg.Params) *stake.Node {
	testName := testutil.TestName()
	testutil.ResetTempFolder(&testName)
	target := filepath.Join(testName, testutil.DefaultDBFileName)
	dbPath := testutil.FilePathInsideTempDir(target)

	testDb, err := database.Create("ffldb", dbPath, params.Net)
	if err != nil {
		testutil.ReportTestIsNotAbleToTest("error creating db: %v", err)
	}

	var bestNode *stake.Node
	err = testDb.Update(func(dbTx database.Tx) error {
		var errLocal error
		bestNode, errLocal = stake.InitDatabaseState(dbTx, params)
		return errLocal
	})
	if err != nil {
		testutil.ReportTestIsNotAbleToTest(err.Error())
	}
	return bestNode
}

// See func (db *StakeDatabase) ConnectBlock(block *dcrutil.Block) error method in stakedb.go
func connectBlock(
	block *dcrutil.Block, BestNode *stake.Node, blockCache map[int64]*dcrutil.Block, params *chaincfg.Params) *stake.Node {

	//block.Bytes() // serialize block

	testutil.Log("ConnectBlock:", block.MsgBlock().Header.Height)

	height := block.Height()
	maturingHeight := height - int64(params.TicketMaturity)

	var maturingTickets []chainhash.Hash
	if maturingHeight >= 0 {
		maturingBlock := blockCache[maturingHeight]
		maturingTickets, _ = txhelpers.TicketsInBlock(maturingBlock)
	}

	revokedTickets := txhelpers.RevokedTicketsInBlock(block)
	votedTickets := txhelpers.TicketsSpentInBlock(block)

	hB, err := block.BlockHeaderBytes()
	if err != nil {
		err := fmt.Errorf("unable to serialize block header: %v", err)
		testutil.ReportTestFailed("err", err)
	}
	bestNode, err := BestNode.ConnectNode(stake.CalcHash256PRNGIV(hB), votedTickets, revokedTickets, maturingTickets)
	if err != nil {
		testutil.ReportTestFailed("err", err)
	}
	if bestNode == nil {
		err := fmt.Errorf("failed to ConnectNode at BestNode")
		testutil.ReportTestFailed("err", err)
	}
	return bestNode

}
