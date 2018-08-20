// Copyright (c) 2017, The dcrdata developers
// See LICENSE for details.

package explorer

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type contextKey int

// AllowFutureBlocks signals if instead of an error, for any block numbers that
// are ahead of the current height, we should display a page which says
// the block is not available yet.
var AllowFutureBlocks = true

const (
	ctxSearch     contextKey = iota
	ctxBlockIndex
	ctxBlockHash
	ctxTxHash
	ctxAddress
	ctxAgendaId
)

func (exp *explorerUI) BlockHashPathOrIndexCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		blockHashOrIndex := chi.URLParam(r, "blockhash")

		// Is it a block index? Let's check.
		blockIndex := tryExtractBlockIndex(blockHashOrIndex)
		if blockIndex != -1 {
			// Yes it is the block index.
			// Serve block by the index.
			ctx := context.WithValue(r.Context(), ctxBlockIndex, blockIndex)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// This is not a block index.

		// Is it a block hash? Let's check.
		_, err := exp.blockData.GetBlockHeight(blockHashOrIndex)
		if err != nil {
			// No it is not a valid block hash, also it is not a valid index. Report.
			log.Errorf("GetBlockHeight(%s) failed: %v", blockHashOrIndex, err)
			exp.StatusPage(w, defaultErrorCode, "could not find that block", NotFoundStatusType)
			return
		}

		// Yes it is the block hash.
		hash := blockHashOrIndex

		// Serve block by the hash.
		ctx := context.WithValue(r.Context(), ctxBlockHash, hash)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func tryExtractBlockIndex(blockHashOrHeight string) int64 {
	height, err := strconv.ParseInt(blockHashOrHeight, 10, 0)
	if err != nil { // this is not height
		return -1
	}
	return height
}

func getBlockHashCtx(r *http.Request) string {
	hash, ok := r.Context().Value(ctxBlockHash).(string)
	if !ok {
		log.Trace("Block Hash not set")
		return ""
	}
	return hash
}

func getBlockHeightCtx(r *http.Request) int64 {
	idx, ok := r.Context().Value(ctxBlockIndex).(int64)
	if !ok {
		log.Trace("Block Height not set")
		return -1
	}
	return idx
}

func getTxIDCtx(r *http.Request) string {
	hash, ok := r.Context().Value(ctxTxHash).(string)
	if !ok {
		log.Trace("Txid not set")
		return ""
	}
	return hash
}

func getAgendaIDCtx(r *http.Request) string {
	hash, ok := r.Context().Value(ctxAgendaId).(string)
	if !ok {
		log.Trace("Agendaid not set")
		return ""
	}
	return hash
}

// TransactionHashCtx embeds "txid" into the request context
func TransactionHashCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		txid := chi.URLParam(r, "txid")
		ctx := context.WithValue(r.Context(), ctxTxHash, txid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AddressPathCtx embeds "address" into the request context
func AddressPathCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		ctx := context.WithValue(r.Context(), ctxAddress, address)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AgendaPathCtx embeds "agendaid" into the request context
func AgendaPathCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		agendaid := chi.URLParam(r, "agendaid")
		ctx := context.WithValue(r.Context(), ctxAgendaId, agendaid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
