package rpcv1

import (
	"github.com/tsubasa597/ASoulCnkiBackend/pkg/config"
	"github.com/tsubasa597/BILIBILI-HELPER/rpc/service"
	"google.golang.org/grpc"
)

var (
	connent *grpc.ClientConn
)

func Setup() error {
	var err error
	if connent, err = grpc.Dial(config.RPCPath, grpc.WithInsecure()); err != nil {
		return err
	}

	dynamicClient = service.NewDynamicClient(connent)
	commentClient = service.NewCommentClient(connent)
	return nil
}

func Stop() error {
	return connent.Close()
}
