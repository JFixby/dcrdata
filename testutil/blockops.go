package testutil

import (
	"path/filepath"
	"io/ioutil"
	"bytes"
	"strconv"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/wire"
)

// BlockFilename generates file name for a Block
func BlockFilename(index uint64) string {
	return "block-" + strconv.FormatUint(index, 10) + ".bin"
}

// SaveBlockToFile writes Block to a file in the target folder
func SaveBlockToFile(block *dcrutil.Block, targetFolder string) (*string, error) {
	bytes, err := block.Bytes()
	if err != nil {
		return nil, err
	}

	index := block.MsgBlock().Header.Height
	filename := filepath.Join(targetFolder, BlockFilename(uint64(index)))
	filename, err = filepath.Abs(filename)

	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filename, bytes, 0777)

	if err != nil {
		return nil, err
	}

	return &filename, nil
}



// ReadBlock reads Block from file
func ReadBlock(file string) *dcrutil.Block {
	file, err := filepath.Abs(file)

	if err != nil {
		panic("Unable to read file: " + file)
	}

	data, err := ioutil.ReadFile(file)

	if err != nil {
		panic("Unable to read file: " + file)
	}

	var msgBlock wire.MsgBlock
	if err := msgBlock.Deserialize(bytes.NewReader(data)); err != nil {
		panic("Could not decode block: " + file)
	}
	block := dcrutil.NewBlock(&msgBlock)
	return block
}
