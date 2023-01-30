package web

import (
	"strings"
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
	groupRoute.Use(func(ctx *gin.Context) {
		ClientManager(ctx, logger, config)
	})

	utilsGroupRoute := groupRoute.Group("/utils")
	utilsGroupRoute.GET("/getAllowance", webHandler.AllowanceHandler)

	//swagger api
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func ClientManager(ctx *gin.Context, logger *zap.SugaredLogger, config *config.Config) {
	// Get a specific header
	userAgent := ctx.GetHeader("User-Agent")
	//url := ctx.Request.URL.String()

	// Loading client codes from default_config and log the data
	keys := config.ClientCodes
	apiKey := ctx.GetHeader("x-api-key")
	// Checking if we have a valid API key client
	value, ok := keys[strings.ToLower(apiKey)]
	if !ok {
		value = "Mobile"
	}
	logger.Infof("Client: %s", value)
	logger.Infof("User-Agent: %v", userAgent)

}
