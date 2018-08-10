package main

import (
	"bytes"
	"github.com/davecgh/go-spew/spew"
	"github.com/decred/dcrdata/testutil"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestBlock(t *testing.T) {

	targetFolder := "testdata"
	blockFileName := testutil.BlockFilename(4095)
	blockFileName = filepath.Join(targetFolder, blockFileName)

	bytesA, err := ioutil.ReadFile(blockFileName)

	block, err := testutil.ReadBlock(blockFileName)
	if err != nil {
		testutil.Log(" failed", err)
	}

	bytesB, err := block.Bytes()

	if err != nil {
		testutil.Log(" failed", err)
	}

	if !bytes.Equal(bytesB, bytesA) {
		t.Fatalf("TestBlock: block does not appear valid - "+
			"got %v, want %v", spew.Sdump(bytesB),
			spew.Sdump(bytesA))
	}

}
