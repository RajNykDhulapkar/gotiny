package rangeallocator

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
)

type RangeManagerConfig struct {
	ServiceID string
	RangeSize int64
	Region    string
}

type RangeManager struct {
	client       Client
	config       *RangeManagerConfig
	mu           sync.RWMutex
	currentRange *pb.Range
}

func NewRangeManager(client Client, config *RangeManagerConfig) *RangeManager {
	return &RangeManager{
		client: client,
		config: config,
	}
}

func (rm *RangeManager) GetNextID(ctx context.Context) (int64, error) {
	rm.mu.RLock()
	if rm.currentRange == nil {
		rm.mu.RUnlock()
		return rm.allocateNewRange(ctx)
	}

	currentID := atomic.LoadInt64(&rm.currentRange.StartId)
	if currentID >= rm.currentRange.EndId {
		rm.mu.RUnlock()
		return rm.allocateNewRange(ctx)
	}

	nextID := atomic.AddInt64(&rm.currentRange.StartId, 1)
	rm.mu.RUnlock()

	if nextID > rm.currentRange.EndId {
		rm.mu.Lock()
		if rm.currentRange.StartId > rm.currentRange.EndId {
			// Mark current range as exhausted
			_, err := rm.client.UpdateRangeStatus(ctx,
				rm.currentRange.RangeId,
				rm.config.ServiceID,
				pb.RangeStatus_RANGE_STATUS_EXHAUSTED,
			)
			if err != nil {
				rm.mu.Unlock()
				return 0, fmt.Errorf("failed to update range status: %w", err)
			}
			rm.mu.Unlock()
			return rm.allocateNewRange(ctx)
		}
		rm.mu.Unlock()
	}

	return nextID - 1, nil
}

func (rm *RangeManager) allocateNewRange(ctx context.Context) (int64, error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.currentRange != nil && rm.currentRange.StartId < rm.currentRange.EndId {
		return atomic.AddInt64(&rm.currentRange.StartId, 1) - 1, nil
	}

	var region *string
	if rm.config.Region != "" {
		region = &rm.config.Region
	}

	size := rm.config.RangeSize
	newRange, err := rm.client.AllocateRange(ctx, rm.config.ServiceID, &size, region)
	if err != nil {
		return 0, fmt.Errorf("failed to allocate new range: %w", err)
	}

	rm.currentRange = newRange
	return atomic.AddInt64(&rm.currentRange.StartId, 1) - 1, nil
}

func (rm *RangeManager) GetCurrentRange() *pb.Range {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if rm.currentRange == nil {
		return nil
	}

	rangeCopy := *rm.currentRange
	return &rangeCopy
}
