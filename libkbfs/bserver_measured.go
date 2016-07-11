// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package libkbfs

import (
	metrics "github.com/rcrowley/go-metrics"
	"golang.org/x/net/context"
)

// BlockServerMeasured delegates to another BlockServer instance but
// also keeps track of stats.
type BlockServerMeasured struct {
	delegate                    IFCERFTBlockServer
	getTimer                    metrics.Timer
	putTimer                    metrics.Timer
	addBlockReferenceTimer      metrics.Timer
	removeBlockReferenceTimer   metrics.Timer
	archiveBlockReferencesTimer metrics.Timer
}

var _ IFCERFTBlockServer = BlockServerMeasured{}

// NewBlockServerMeasured creates and returns a new
// BlockServerMeasured instance with the given delegate and registry.
func NewBlockServerMeasured(delegate IFCERFTBlockServer, r metrics.Registry) BlockServerMeasured {
	getTimer := metrics.GetOrRegisterTimer("BlockServer.Get", r)
	putTimer := metrics.GetOrRegisterTimer("BlockServer.Put", r)
	addBlockReferenceTimer := metrics.GetOrRegisterTimer("BlockServer.AddBlockReference", r)
	removeBlockReferenceTimer := metrics.GetOrRegisterTimer("BlockServer.RemoveBlockReference", r)
	archiveBlockReferencesTimer := metrics.GetOrRegisterTimer("BlockServer.ArchiveBlockReferences", r)
	return BlockServerMeasured{
		delegate:                    delegate,
		getTimer:                    getTimer,
		putTimer:                    putTimer,
		addBlockReferenceTimer:      addBlockReferenceTimer,
		removeBlockReferenceTimer:   removeBlockReferenceTimer,
		archiveBlockReferencesTimer: archiveBlockReferencesTimer,
	}
}

// Get implements the BlockServer interface for BlockServerMeasured.
func (b BlockServerMeasured) Get(ctx context.Context, id IFCERFTBlockID, tlfID IFCERFTTlfID, context IFCERFTBlockContext) (
	buf []byte, serverHalf IFCERFTBlockCryptKeyServerHalf, err error) {
	b.getTimer.Time(func() {
		buf, serverHalf, err = b.delegate.Get(ctx, id, tlfID, context)
	})
	return buf, serverHalf, err
}

// Put implements the BlockServer interface for BlockServerMeasured.
func (b BlockServerMeasured) Put(ctx context.Context, id IFCERFTBlockID, tlfID IFCERFTTlfID, context IFCERFTBlockContext, buf []byte,
	serverHalf IFCERFTBlockCryptKeyServerHalf) (err error) {
	b.putTimer.Time(func() {
		err = b.delegate.Put(ctx, id, tlfID, context, buf, serverHalf)
	})
	return err
}

// AddBlockReference implements the BlockServer interface for
// BlockServerMeasured.
func (b BlockServerMeasured) AddBlockReference(ctx context.Context, id IFCERFTBlockID, tlfID IFCERFTTlfID, context IFCERFTBlockContext) (err error) {
	b.addBlockReferenceTimer.Time(func() {
		err = b.delegate.AddBlockReference(ctx, id, tlfID, context)
	})
	return err
}

// RemoveBlockReference implements the BlockServer interface for
// BlockServerMeasured.
func (b BlockServerMeasured) RemoveBlockReference(ctx context.Context,
	tlfID IFCERFTTlfID, contexts map[IFCERFTBlockID][]IFCERFTBlockContext) (
	liveCounts map[IFCERFTBlockID]int, err error) {
	b.removeBlockReferenceTimer.Time(func() {
		liveCounts, err = b.delegate.RemoveBlockReference(ctx, tlfID, contexts)
	})
	return liveCounts, err
}

// ArchiveBlockReferences implements the BlockServer interface for
// BlockServerRemote
func (b BlockServerMeasured) ArchiveBlockReferences(ctx context.Context,
	tlfID IFCERFTTlfID, contexts map[IFCERFTBlockID][]IFCERFTBlockContext) (err error) {
	b.archiveBlockReferencesTimer.Time(func() {
		err = b.delegate.ArchiveBlockReferences(ctx, tlfID, contexts)
	})
	return err
}

// Shutdown implements the BlockServer interface for
// BlockServerMeasured.
func (b BlockServerMeasured) Shutdown() {
	b.delegate.Shutdown()
}

// RefreshAuthToken implements the BlockServer interface for
// BlockServerMeasured.
func (b BlockServerMeasured) RefreshAuthToken(ctx context.Context) {
	b.delegate.RefreshAuthToken(ctx)
}

// GetUserQuotaInfo implements the BlockServer interface for BlockServerMeasured
func (b BlockServerMeasured) GetUserQuotaInfo(ctx context.Context) (info *IFCERFTUserQuotaInfo, err error) {
	return b.delegate.GetUserQuotaInfo(ctx)
}
