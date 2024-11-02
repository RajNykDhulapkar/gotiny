package rangeallocator

import (
	"context"
	"fmt"
	"time"

	"github.com/RajNykDhulapkar/gotiny-range-allocator/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClientConfig struct {
	Address     string        `yaml:"address"`
	DialTimeout time.Duration `yaml:"dialTimeout"`
}

type Client interface {
	AllocateRange(ctx context.Context, serviceID string, size *int64, region *string) (*pb.Range, error)
	UpdateRangeStatus(ctx context.Context, rangeID, serviceID string, status pb.RangeStatus) (*pb.Range, error)
	GetHealth(ctx context.Context) error
	Close() error
}

type clientImpl struct {
	conn   *grpc.ClientConn
	client pb.RangeAllocatorClient
}

func NewClient(cfg *ClientConfig) (*clientImpl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx,
		cfg.Address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to range allocator: %w", err)
	}

	return &clientImpl{
		conn:   conn,
		client: pb.NewRangeAllocatorClient(conn),
	}, nil
}

func (c *clientImpl) AllocateRange(ctx context.Context, serviceID string, size *int64, region *string) (*pb.Range, error) {
	req := &pb.AllocateRangeRequest{
		ServiceId: serviceID,
		Size:      size,
		Region:    region,
	}

	resp, err := c.client.AllocateRange(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("allocate range failed: %w", err)
	}

	return resp.Range, nil
}

func (c *clientImpl) UpdateRangeStatus(ctx context.Context, rangeID, serviceID string, status pb.RangeStatus) (*pb.Range, error) {
	req := &pb.UpdateRangeStatusRequest{
		RangeId:   rangeID,
		ServiceId: serviceID,
		Status:    status,
	}

	resp, err := c.client.UpdateRangeStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("update range status failed: %w", err)
	}

	return resp, nil
}

func (c *clientImpl) GetHealth(ctx context.Context) error {
	resp, err := c.client.GetHealth(ctx, &emptypb.Empty{})
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if resp.Status != pb.ServiceStatus_SERVICE_STATUS_SERVING {
		return fmt.Errorf("service not in serving state: %s", resp.Details)
	}

	return nil
}

func (c *clientImpl) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
