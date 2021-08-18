package cos

import (
	"context"
	"fmt"
	"os"

	"github.com/Jsoneft/cos_service/client"
)

//go:generate mockgen -destination=./mock/cos_mock.go -package=mock . COS

// COS cache object service interface.
type COS interface {
	// Serve Start to Serve.
	Serve(ctx context.Context) error
}

// Impl is the default implement.
type Impl struct {
	client.Client
}

var New = func() COS {
	return &Impl{client.New()}
}

// Serve Start to Serve.
func (i *Impl) Serve(ctx context.Context) error {
	args := os.Args[1:]
	urls, err := i.Client.Upload(ctx, args)
	if err != nil {
		return err
	}
	// 打印到标准输出里
	for _, url := range urls {
		fmt.Println(url)
	}
	return nil
}
