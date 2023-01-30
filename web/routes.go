package web

import (
	// "strings"
	"bridge-allowance/config"
	grpcClient "bridge-allowance/pkg/grpc/client"
	handler "bridge-allowance/web/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"go.uber.org/zap"
)

func InitRoutes(config *config.Config, logger *zap.SugaredLogger, route *gin.Engine) {
	grpcClientManager := grpcClient.NewGrpcClientManager(config, logger)
	// http := utils.NewHttpRequest(logger)
	webHandler := handler.NewHandler(config, logger, grpcClientManager)
	groupRoute := route.Group("/v2")

	groupRoute.GET("/getAllowance", webHandler.AllowanceHandler)

	//swagger api
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
