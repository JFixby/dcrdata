package dcrsqlite

import (
	"testing"

	"fmt"
	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
	apitypes "github.com/decred/dcrdata/api/types"
	"github.com/decred/dcrdata/testutil"
	"github.com/decred/dcrdata/txhelpers"
)

func TestBlockDB(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := GetDB("blocks_0-40767")
	testStoreBlocks(db)
}

func testStoreBlocks(tdb *TestDB) {

	tdb.EmptyDB()
	db := tdb.DB()
	first := tdb.GetFirstBlockNumber()
	last := tdb.GetLastBlockNumber()

	for i := first; i <= last; i++ {
		testutil.Log("block", fmt.Sprint(i))
		blockFileName := testutil.BlockFilename(i)
		targetBlockFile := tdb.PathToTestFile(blockFileName)
		block := readBlock(targetBlockFile)
		storeBlockInDB(db, block)
	}
}

func readBlock(targetBlockFile string) *dcrutil.Block {
	testutil.Log("      read", targetBlockFile)
	block, err := testutil.ReadBlock(targetBlockFile)
	if err != nil {
		testutil.ReportTestIsNotAbleToTest(
			"Failed to read block file %v",
			targetBlockFile)
	}
	return block
}

func storeBlockInDB(db *DB, block *dcrutil.Block) {

	header := block.MsgBlock().Header
	var activeChain = &chaincfg.MainNetParams
	diffRatio := txhelpers.GetDifficultyRatio(header.Bits, activeChain)

	tpi := apitypes.TicketPoolInfo{
		Size: header.PoolSize,
	}

	blockhash := block.MsgBlock().Header.BlockHash()

	blockSummary := apitypes.BlockDataBasic{
		Height:     header.Height,
		Size:       header.Size,
		Hash:       blockhash.String(),
		Difficulty: diffRatio,
		StakeDiff:  dcrutil.Amount(header.SBits).ToCoin(),
		Time:       header.Timestamp.Unix(),
		PoolInfo:   tpi,
	}

	testutil.Log("     store", header.BlockHash().String())
	db.StoreBlockSummary(&blockSummary)
}
