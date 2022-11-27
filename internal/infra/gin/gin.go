package gin

import (
	"net/http"

	"shorten/internal/infra/gin/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Config setting http config
type Config struct {
	Debug                  bool     `yaml:"debug" mapstructure:"debug"`
	Address                string   `yaml:"address" mapstructure:"address"`
	AppID                  string   `yaml:"app_id" mapstructure:"app_id"`
	DisableAccessLog       bool     `yaml:"disable_access_log" mapstructure:"disable_access_log"`
	DisableBodyDump        bool     `yaml:"disable_body_dump" mapstructure:"disable_body_dump"`
	DisablePprof           bool     `yaml:"disable_pprof" mapstructure:"disable_pprof"`
	BodyDumpIgnoreURLPath  []string `yaml:"body_dump_ignore_url_path" mapstructure:"body_dump_ignore_url_path"`
	CORSAllowOrigins       []string `yaml:"cors_allow_origins" mapstructure:"cors_allow_origins"`
	AccessLogIgnoreURLPath []string `yaml:"access_log_ignore_url_path" mapstructure:"access_log_ignore_url_path"`
}

// New create new engine for handler to register
func New(cfg Config) *gin.Engine {

	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	e := gin.New()
	if !cfg.DisablePprof {
		pprof.Register(e)
	}
	e.HandleMethodNotAllowed = true
	e.NoMethod(NotFoundHandler)
	e.NoRoute(NotFoundHandler)
	e.Use(middleware.RequestID)
	e.Use(middleware.AccessLog(cfg.DisableAccessLog))
	e.Use(middleware.BodyDump(cfg.DisableBodyDump))
	e.Use(middleware.Recovery())
	return e
}

// NotFoundHandler responds not found response.
func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "page not found"})
}
