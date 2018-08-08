package dcrsqlite

import (
	"testing"

	"github.com/decred/dcrdata/testutil"
)

func TestEmptyDBGetBestBlockHeight(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := InitTestDB(DBPathForTest())
	testEmptyDBGetBestBlockHeight(db)
}

// Empty DB, should return -1
func testEmptyDBGetBestBlockHeight(db *DB) {
	h := db.GetBestBlockHeight()
	if h != -1 {
		testutil.ReportTestFailed(
			"db.GetBestBlockHeight() is %v,"+
				" should be -1",
			h)
	}
}

func TestGetBestBlockHeight(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := GetDB("synced_up_to_260241")
	testGetBestBlockHeight(db)
}

func testGetBestBlockHeight(tdb *TestDB) {
	db := tdb.DB()
	exh := tdb.GetExpectedBestBlockHeight()
	h := db.GetBestBlockHeight()
	if exh != h {
		testutil.ReportTestFailed(
			"db.GetBestBlockHeight() is %v,"+
				" should be %v",
			h,
			exh)
	}
}
