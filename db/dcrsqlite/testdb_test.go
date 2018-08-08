package dcrsqlite

import (
	"github.com/decred/dcrdata/testutil"
	"path/filepath"
	"os"
)

const (
	TestDataFolderName     = "dcrdata-testdata"
	TestDBFileName         = "test-data.db"
	TestDescriptorFileName = "test-description.json"
)

type TestDB struct {
	tag      string
	db       *DB
	expected TestDescriptorJson
}

type TestDescriptorJson struct {
	TestName                                string `json:"test_name"`
	ExpectedBestBlockHeight                 int64  `json:"expected_block_height"`
	ExpectedStakeInfoHeight                 int64  `json:"expected_stake_info_height"`
	ExpectedBestBlockHash                   string `json:"expected_best_block_hash"`
	ExpectedMD5OfRetrieveBlockFeeInfoString string `json:"expected_md5_of_retrieve_block_fee_info_string"`
	FirstBlockNumber                        uint64 `json:"first_block_number"`
	LastBlockNumber                         uint64 `json:"last_block_number"`
}

func (tdb *TestDB) DB() *DB {
	if tdb.db != nil {
		return tdb.db
	}
	testDBFile := PathToTestDBFile(tdb.tag)
	testutil.Log("           reading", testDBFile)
	tdb.db = InitTestDB(testDBFile)
	return tdb.db
}

func (db *TestDB) GetExpectedBestBlockHeight() int64 {
	return db.expected.ExpectedBestBlockHeight
}

func (db *TestDB) GetExpectedStakeInfoHeight() int64 {
	return db.expected.ExpectedStakeInfoHeight
}

func (db *TestDB) GetExpectedBestBlockHash() string {
	return db.expected.ExpectedBestBlockHash
}

func (db *TestDB) GetExpectedMD5OfRetrieveBlockFeeInfoString() string {
	return db.expected.ExpectedMD5OfRetrieveBlockFeeInfoString
}

func (db *TestDB) GetFirstBlockNumber() uint64 {
	return db.expected.FirstBlockNumber
}

func (db *TestDB) GetLastBlockNumber() uint64 {
	return db.expected.LastBlockNumber
}

var testDBs = make(map[string]*TestDB)

func GetDB(tag string) *TestDB {
	tdb := testDBs[tag]
	if tdb == nil {
		tdb = loadTDB(tag)
		testDBs[tag] = tdb
	}
	return tdb
}

func loadTDB(tag string) *TestDB {
	var tdb TestDB
	testDescriptorFile := PathToTestDescriptorFile(tag)

	testutil.Log("reading:")
	testutil.Log("testDescriptorFile", testDescriptorFile)

	// uncomment this to produce example descriptor file
	//ProduceExampleTestDescriptor(PathToTestDescriptorFile(tag))

	testDescriptor := ReadTestDescriptorJson(testDescriptorFile)
	tdb = TestDB{
		expected: testDescriptor,
		tag:      tag,
	}

	return &tdb
}
func parseTestDescriptorJson(jsonString string) TestDescriptorJson {
	desc := TestDescriptorJson{}
	testutil.FromJson(jsonString, &desc)
	return desc
}
func ReadTestDescriptorJson(targetFile string) TestDescriptorJson {
	jsonString := testutil.ReadFileToString(targetFile)
	desc := parseTestDescriptorJson(jsonString)
	return desc
}

func PathToTestDataFolder(tag string) string {
	return testutil.FullPathToFile(
		filepath.Join(TestDataFolderName, tag))
}

func PathToTestDBFile(tag string) string {
	return testutil.FullPathToFile(
		filepath.Join(PathToTestDataFolder(tag),
			TestDBFileName))
}

func PathToTestDescriptorFile(tag string) string {
	return testutil.FullPathToFile(
		filepath.Join(PathToTestDataFolder(tag),
			TestDescriptorFileName))
}

func (tdb *TestDB) PathToTestFile(fileName string) string {
	return testutil.FullPathToFile(
		filepath.Join(PathToTestDataFolder(tdb.tag),
			fileName))
}
func (tdb *TestDB) EmptyDB() {
	targetFile := PathToTestDBFile(tdb.tag)
	testutil.Log("            delete", targetFile)

	err := os.Remove(targetFile)

	if err != nil {
		testutil.ReportTestIsNotAbleToTest(
			"Failed to delete DB file %v,\n"+
				"%v",
			targetFile,
			err)
	}
}

func ProduceExampleTestDescriptor(targetFile string) {
	exampleDescriptor := TestDescriptorJson{
		TestName:                                "exampleTest",
		ExpectedBestBlockHeight:                 -1,
		ExpectedStakeInfoHeight:                 -1,
		ExpectedBestBlockHash:                   "",
		ExpectedMD5OfRetrieveBlockFeeInfoString: "d41d8cd98f00b204e9800998ecf8427e",
	}

	jsonString := testutil.ToJson(exampleDescriptor)

	testutil.Log("writing", targetFile)
	testutil.Log(jsonString)
	testutil.WriteStringToFile(targetFile, jsonString)
}
