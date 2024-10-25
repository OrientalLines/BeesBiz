package tikv

import (
	"context"
	"fmt"

	"github.com/pingcap/kvproto/pkg/kvrpcpb"
	"github.com/tikv/client-go/v2/rawkv"
	"go.uber.org/zap"
)

type TiKV struct {
	client *rawkv.Client
}

func New(pdEndpoints []string) (*TiKV, error) {
	client, err := rawkv.NewClientWithOpts(context.Background(), pdEndpoints, rawkv.WithAPIVersion(kvrpcpb.APIVersion_V2),
		rawkv.WithKeyspace("a"))
	if err != nil {
		return nil, fmt.Errorf("error connecting to TiKV: %w", err)
	}

	zap.L().Info("Successfully connected to TiKV")
	return &TiKV{client: client}, nil
}

func (tk *TiKV) Close() error {
	if tk.client != nil {
		return tk.client.Close()
	}
	return nil
}

func (tk *TiKV) GetClient() *rawkv.Client {
	return tk.client
}
