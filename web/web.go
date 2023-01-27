package web

import (
	"bridge-allowance/config"
	"bridge-allowance/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	gintracer "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

type WebServer struct {
	router *gin.Engine
	config *config.Config
	logger *zap.SugaredLogger
}

func NewWebServer(config *config.Config, logger *zap.SugaredLogger) *WebServer {
	/**
	@description Init Router
	*/
	router := gin.Default()
	//Middleware to recover from panic
	router.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.SecureJSON(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}), gintracer.Middleware(config.WEB_DATADOG_SERVICE))
	return &WebServer{router: router, config: config, logger: logger}
}

func (u *WebServer) Start() {
	u.logger.Infof("Web Server Started at :%v", u.config.Web.Port)
	InitRoutes(u.config, u.logger, u.router)
	err := u.router.Run(":" + u.config.Web.Port)
	if err != nil {
		u.logger.Error(err)
		return
	}
}

func run(cmd *cobra.Command, args []string) {
	conf := config.LoadConfig("", "")
	logger := utils.SetupLogger(conf.Logger.LogLevel, conf.Logger.LogPath+conf.Web.LogFile, conf.LOG_ENCODING_FORMAT)
	webServer := NewWebServer(conf, logger)
	webServer.Start()
}

var WebCmd = &cobra.Command{
	Use:   "web",
	Short: "web",
	Long:  `Framework:  Web Server`,
	Run:   run,
}
