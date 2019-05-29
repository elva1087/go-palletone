// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package light

import (
	"github.com/palletone/go-palletone/common/log"
	"github.com/palletone/go-palletone/dag/modules"
	"github.com/palletone/go-palletone/ptn/downloader"
	"sync/atomic"
	"time"
)

const (
	forceSyncCycle      = 10 * time.Second // Time interval to force syncs, even if few peers are available
	minDesiredPeerCount = 5                // Amount of peers desired to start syncing

	// This is the target size for the packs of transactions sent by txsyncLoop.
	// A pack can get larger than this if a single transactions exceeds this size.
	txsyncPackSize = 100 * 1024
)

// syncer is responsible for periodically synchronising with the network, both
// downloading hashes and blocks as well as handling the announcement handler.
func (pm *ProtocolManager) syncer() {
	// Start and ensure cleanup of sync mechanisms
	pm.fetcher.Start()
	defer pm.fetcher.Stop()
	defer pm.downloader.Terminate()

	// Wait for different events to fire synchronisation operations
	forceSync := time.Tick(forceSyncCycle)
	for {
		select {
		case <-pm.newPeerCh:
			// Make sure we have peers to select from, then sync
			if pm.peers.Len() < minDesiredPeerCount {
				break
			}
			go pm.synchronise(pm.peers.BestPeer(pm.assetId), pm.assetId)

		case <-forceSync:
			// Force a sync even if not enough peers are present
			go pm.syncall() //pm.synchronise(pm.peers.BestPeer(pm.assetId), pm.assetId)

		case <-pm.noMorePeers:
			return
		}
	}
}

func (pm *ProtocolManager) syncall() {
	return
	if atomic.LoadUint32(&pm.fastSync) == 0 {
		log.Debug("Light PalletOne syncall synchronising")
		return
	}

	p := pm.peers.BestPeer(pm.assetId)
	if p == nil {
		log.Debug("Light PalletOne syncall peer is nil")
		return
	}
	headers, err := pm.downloader.FetchAllToken(p.id)
	if err != nil {
		log.Debug("Light PalletOne syncall FetchAllToken", "err", err)
	}
	log.Debug("Light PalletOne syncall FetchAllToken", "len(headers)", len(headers), "headers", headers)
}

// synchronise tries to sync up our local block chain with a remote peer.
func (pm *ProtocolManager) synchronise(peer *peer, assetId modules.AssetId) {
	// Short circuit if no peers are available
	if peer == nil {
		return
	}

	if !pm.lightSync && pm.assetId == assetId {
		log.Debug("Light PalletOne synchronise pm.assetId == assetId")
		return
	}

	if atomic.LoadUint32(&pm.fastSync) == 0 {
		log.Debug("Light PalletOne synchronising")
		return
	}
	atomic.StoreUint32(&pm.fastSync, 0)
	headhash, number := peer.HeadAndNumber(assetId)
	log.Debug("Light PalletOne ProtocolManager synchronise", "assetid", assetId, "index", number.Index)

	if err := pm.downloader.Synchronise(peer.id, headhash, number.Index, downloader.LightSync, number.AssetID); err != nil {
		log.Debug("Light PalletOne ProtocolManager synchronise", "Synchronise err:", err)
		return
	}

	if atomic.LoadUint32(&pm.fastSync) == 0 {
		log.Debug("Fast sync complete, auto disabling")
		atomic.StoreUint32(&pm.fastSync, 1)
	}

	header := pm.dag.CurrentHeader(assetId)
	if header != nil && header.Number.Index > 0 {
		go pm.BroadcastLightHeader(header)
	}
}
