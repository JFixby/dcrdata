package dcrsqlite

import (
	"testing"

	"github.com/decred/dcrdata/testutil"
)

func TestEmptyDBGetStakeInfoHeight(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := InitTestDB(DBPathForTest())
	testEmptyDBGetStakeInfoHeight(db)
}

func testEmptyDBGetStakeInfoHeight(db *DB) {
	endHeight, err := db.GetStakeInfoHeight()
	if err != nil {
		testutil.ReportTestFailed(
			"GetStakeInfoHeight() failed: %v",
			err)
	}
	if endHeight != -1 {
		testutil.ReportTestFailed(
			"GetStakeInfoHeight() failed: endHeight=%v,"+
				" should be -1",
			endHeight)
	}
}

func TestGetStakeInfoHeight(t *testing.T) {
	testutil.BindCurrentTestSetup(t)
	db := GetDB("synced_up_to_260241")
	testGetStakeInfoHeight(db)
}

func testGetStakeInfoHeight(tdb *TestDB) {
	db := tdb.DB()
	exh := tdb.GetExpectedStakeInfoHeight()
	h, err := db.GetStakeInfoHeight()
	if err != nil {
		testutil.ReportTestFailed(
			"GetStakeInfoHeight() failed: %v",
			err)
	}
	if exh != h {
		testutil.ReportTestFailed(
			"db.GetStakeInfoHeight() is %v,"+
				" should be %v",
			h,
			exh)
	}
}
