package dcrsqlite

import (
	"testing"

	"github.com/decred/dcrdata/testutil"
)

func TestEmptyDBGetBestBlockHash(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := InitTestDB(DBPathForTest())
	testEmptyDBGetBestBlockHash(db)
}

func testEmptyDBGetBestBlockHash(db *DB) {
	str := db.GetBestBlockHash()
	if str != "" {
		// Open question: Should it really be the empty string?
		// Maybe error instead to avoid confusion?
		testutil.ReportTestFailed(
			"GetBestBlockHash() failed: %v",
			str)
	}
}

func TestGetBestBlockHash(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := GetDB("synced_up_to_260241")
	testGetBestBlockHash(db)
}

func testGetBestBlockHash(tdb *TestDB) {
	db := tdb.DB()
	exh := tdb.GetExpectedBestBlockHash()
	h := db.GetBestBlockHash()
	if exh != h {
		testutil.ReportTestFailed(
			"db.GetBestBlockHash() is %v,"+
				" should be %v",
			h,
			exh)
	}
}
