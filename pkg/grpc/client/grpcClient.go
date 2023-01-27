package client

import (
	"bridge-allowance/config"
	"bridge-allowance/pkg/grpc/proto/pb"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
	"time"
)

type GrpcClientManager struct {
	config *config.Config
	log    *zap.SugaredLogger
}

func NewGrpcClientManager(config *config.Config, log *zap.SugaredLogger) *GrpcClientManager {
	return &GrpcClientManager{
		config: config,
		log:    log,
	}
}

func (g *GrpcClientManager) MapGrpcClient() map[string]pb.UnifrontClient {
	grpcClientMap := make(map[string]pb.UnifrontClient)
	nonEVMServerAddr := g.grpcConn(g.config.NONEVM_GRPC_ENDPOINT, g.config.NONEVM_DATADOG_SERVICE)
	cosmosServerAddr := g.grpcConn(g.config.COSMOS_GRPC_ENDPOINT, g.config.COSMOS_DATADOG_SERVICE)
	evmServerAddr := g.grpcConn(g.config.EVM_GRPC_ENDPOINT, g.config.EVM_DATADOG_SERVICE)
	bridgeServerAddr := g.grpcConn(g.config.BRIDGE_GRPC_ENDPOINT, g.config.BRIDGE_DATADOG_SERVICE)

	grpcClientMap["cosmos"] = cosmosServerAddr
	grpcClientMap["nonevm"] = nonEVMServerAddr
	grpcClientMap["cosmos_network"] = cosmosServerAddr

	grpcClientMap["evm"] = evmServerAddr
	grpcClientMap["bridge"] = bridgeServerAddr
	return grpcClientMap
}

func (g *GrpcClientManager) grpcConn(endpoint string, grpcTraceService string) pb.UnifrontClient {
	if endpoint == g.config.NONEVM_GRPC_ENDPOINT {
		opts := []grpc_retry.CallOption{
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(5 * time.Second)),
			grpc_retry.WithCodes(codes.FailedPrecondition),
			grpc_retry.WithMax(5),
		}
		// Create the client interceptor using the grpc trace package.
		si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(grpcTraceService))
		ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(grpcTraceService))
		conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStreamInterceptor(si),
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(ui, grpc_retry.UnaryClientInterceptor(opts...))))
		if err != nil {
			g.log.Error(err)
		}
		client := pb.NewUnifrontClient(conn)
		return client
	}
	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(grpcTraceService))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(grpcTraceService))
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithStreamInterceptor(si),
		grpc.WithUnaryInterceptor(ui))
	if err != nil {
		g.log.Error(err)
	}
	client := pb.NewUnifrontClient(conn)
	return client
}
