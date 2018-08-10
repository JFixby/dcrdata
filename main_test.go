package main

import (
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/rpcclient"
	"github.com/decred/dcrdata/blockdata"
	"github.com/decred/dcrdata/db/agendadb"
	"github.com/decred/dcrdata/db/dcrpg"
	"github.com/decred/dcrdata/db/dcrsqlite"
	"github.com/decred/dcrdata/explorer"
	"github.com/decred/dcrdata/mempool"
	notify "github.com/decred/dcrdata/notification"
	"github.com/decred/dcrdata/rpcutils"
	"github.com/decred/dcrdata/testutil"
	"github.com/decred/dcrdata/version"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func TestMainn(t *testing.T) {
	testutil.BindCurrentTestSetup(t)

	// Parse the configuration file, and setup logger.
	var cfggood, _ = loadConfig()
	testutil.Log("loadConfig", cfggood)

	var cfg = &defaultConfig
	cfg.DcrdUser = "6vML++Uu+Wlt7U3LPTTlsiEr6cs="
	cfg.DcrdPass = "RXXfqiCG3zsaMb9wCrQMqe/51F8="
	cfg.DcrdCert = "D:\\PICFIGHT\\dev-s\\rpc.cert"
	cfg.DcrdServ = "localhost:9109"
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	cfg.HomeDir = filepath.Join("test.data", "SyncTest", timestamp)
	cfg.LogDir = filepath.Join(cfg.HomeDir, defaultLogDirname)
	cfg.DataDir = filepath.Join(cfg.HomeDir, defaultDataDirname)

	// Start with version info
	ver := &version.Ver

	testutil.Log("       cfg", cfg)

	// Connect to dcrd RPC server using websockets

	// Set up the notification handler to deliver blocks through a channel.

	// Daemon client connection
	ntfnHandlers, _ := notify.MakeNodeNtfnHandlers()
	dcrdClient, err := proxyConnect(cfg, ntfnHandlers)

	if err != nil || dcrdClient == nil {
		testutil.ReportTestIsNotAbleToTest("Connection to dcrd failed: %v", err)
	}
	smartClient := rpcutils.NewBlockGate(dcrdClient, 10)

	// Sqlite output
	dbPath := filepath.Join(cfg.DataDir, cfg.DBFileName)
	dbInfo := dcrsqlite.DBInfo{FileName: dbPath}
	baseDB, cleanupDB, err := dcrsqlite.InitWiredDB(&dbInfo,
		notify.NtfnChans.UpdateStatusDBHeight, dcrdClient, activeChain, cfg.DataDir)
	defer cleanupDB()
	if err != nil {
		testutil.ReportTestIsNotAbleToTest("Unable to initialize SQLite database: %v", err)
	}
	log.Infof("SQLite DB successfully opened: %s", cfg.DBFileName)
	defer baseDB.Close()

	// PostgreSQL
	var auxDB *dcrpg.ChainDB

	// Ctrl-C to shut down.
	// Nothing should be sent the quit channel.  It should only be closed.
	quit := make(chan struct{})
	// Only accept a single CTRL+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	_, _, err = dcrdClient.GetBestBlock()
	if err != nil {
		testutil.ReportTestFailed("Unable to get block from node: %v", err)
	}

	// AgendaDB upgrade check
	if err = agendadb.CheckForUpdates(dcrdClient); err != nil {
		testutil.ReportTestFailed("agendadb upgrade failed: %v", err)
	}

	// Block data collector. Needs a StakeDatabase too.
	collector := blockdata.NewCollector(dcrdClient, activeChain, baseDB.GetStakeDB())
	if collector == nil {
		testutil.ReportTestFailed("Failed to create block data collector")
	}

	// Build a slice of each required saver type for each data source
	var blockDataSavers []blockdata.BlockDataSaver
	var mempoolSavers []mempool.MempoolDataSaver

	blockDataSavers = append(blockDataSavers, auxDB)

	// For example, dumping all mempool fees with a custom saver
	if cfg.DumpAllMPTix {
		log.Debugf("Dumping all mempool tickets to file in %s.\n", cfg.OutFolder)
		mempoolFeeDumper := mempool.NewMempoolFeeDumper(cfg.OutFolder, "mempool-fees")
		mempoolSavers = append(mempoolSavers, mempoolFeeDumper)
	}

	blockDataSavers = append(blockDataSavers, &baseDB)
	mempoolSavers = append(mempoolSavers, baseDB.MPC)

	// Create the explorer system
	explore := explorer.New(&baseDB, auxDB, cfg.UseRealIP, ver.String(), !cfg.NoDevPrefetch)
	if explore == nil {
		testutil.ReportTestFailed("failed to create new explorer (templates missing?)")
	}
	explore.UseSIGToReloadTemplates()
	defer explore.StopWebsocketHub()
	defer explore.StopMempoolMonitor(notify.NtfnChans.ExpNewTxChan)

	blockDataSavers = append(blockDataSavers, explore)

	// Synchronization between DBs via rpcutils.BlockGate

	// stakedb (in baseDB) connects blocks *after* ChainDB retrieves them, but
	// it has to get a notification channel first to receive them. The BlockGate
	// will provide this for blocks after fetchHeight.
	//baseDB.SyncDBAsync(make(chan dbtypes.SyncResult), quit, smartClient, int64(math.MaxInt32))

	//var blockgetter rpcutils.BlockGetter = newBlockLoader(smartClient)

	//_, err = dcrsqlite.ResyncDB(baseDB, quit, smartClient)
	_, err = dcrsqlite.ResyncDB(baseDB, quit, smartClient)
	if err != nil {
		testutil.Log("&v", err)
		testutil.ReportTestFailed("&v", err)
		os.Exit(-1)
	}

}

type ProxyClient struct {
	client *rpcclient.Client
}

func proxyConnect(cfg *config, ntfnHandlers *rpcclient.NotificationHandlers) (*rpcclient.Client, error) {
	client, _, z := connectNodeRPC(cfg, ntfnHandlers)

	//proxy := ProxyClient{
	//	client: client,
	//}
	return client, z
}

func newBlockLoader(g *rpcutils.BlockGate) *BlockLoader {
	//x := BlockLoader{
	//	g:             g,
	//	height:        -1,
	//	fetchToHeight: -1,
	//	hashAtHeight:  make(map[int64]chainhash.Hash),
	//	blockWithHash: make(map[chainhash.Hash]*dcrutil.Block),
	//	heightWaiters: make(map[int64][]chan chainhash.Hash),
	//	hashWaiters:   make(map[chainhash.Hash][]chan int64),
	//}
	return nil
}

type BlockLoader struct {
	g             *rpcutils.BlockGate
	height        int64
	fetchToHeight int64
	hashAtHeight  map[int64]chainhash.Hash
	blockWithHash map[chainhash.Hash]*dcrutil.Block
	heightWaiters map[int64][]chan chainhash.Hash
	hashWaiters   map[chainhash.Hash][]chan int64
}

func (g *BlockLoader) NodeHeight() (int64, error) {
	r, e := g.g.NodeHeight()
	return r, e
}
func (g *BlockLoader) BestBlockHeight() int64 {
	h := g.g.BestBlockHeight()
	return h
}
func (g *BlockLoader) BestBlockHash() (chainhash.Hash, int64, error) {
	h, x, y := g.g.BestBlockHash()
	return h, x, y
}
func (g *BlockLoader) BestBlock() (*dcrutil.Block, error) {
	r, e := g.g.BestBlock()
	return r, e
}
func (g *BlockLoader) Block(hash chainhash.Hash) (*dcrutil.Block, error) {
	r, e := g.g.Block(hash)
	return r, e
}
func (g *BlockLoader) WaitForHeight(height int64) chan chainhash.Hash {
	c := g.g.WaitForHeight(height)
	return c
}
func (g *BlockLoader) WaitForHash(hash chainhash.Hash) chan int64 {
	c := g.g.WaitForHash(hash)
	return c
}
