package gin

import (
	"net/http"
	"strconv"

	"assignment/internal/infra/gin/middleware"

	"github.com/gin-contrib/sessions"
	sessredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Config setting http config
type Config struct {
	Debug                  bool          `yaml:"debug" mapstructure:"debug"`
	Address                string        `yaml:"address" mapstructure:"address"`
	AppID                  string        `yaml:"app_id" mapstructure:"app_id"`
	DisableBodyDump        bool          `yaml:"disable_body_dump" mapstructure:"disable_body_dump"`
	DisablePprof           bool          `yaml:"disable_pprof" mapstructure:"disable_pprof"`
	BodyDumpIgnoreURLPath  []string      `yaml:"body_dump_ignore_url_path" mapstructure:"body_dump_ignore_url_path"`
	CORSAllowOrigins       []string      `yaml:"cors_allow_origins" mapstructure:"cors_allow_origins"`
	AccessLogIgnoreURLPath []string      `yaml:"access_log_ignore_url_path" mapstructure:"access_log_ignore_url_path"`
	Session                SessionConfig `yaml:"session" mapstructure:"session"`
}

type SessionConfig struct {
	Addr string `yaml:"addr" mapstructure:"addr"`
	Port int    `yaml:"port" mapstructure:"port"`
}

// New create new engine for handler to register
func New(cfg *Config) *gin.Engine {

	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	var store sessredis.Store
	if cfg.Session.Addr != "" {
		addr := cfg.Session.Addr + ":" + strconv.Itoa(cfg.Session.Port)
		store, _ = sessredis.NewStore(10, "tcp", addr, "", []byte("secret_key"))
		store.Options(sessions.Options{MaxAge: 60 * 60 * 24})
	}

	e := gin.New()
	e.HandleMethodNotAllowed = true
	e.NoMethod(NotFoundHandler)
	e.NoRoute(NotFoundHandler)
	e.Use(middleware.RequestID)
	e.Use(middleware.AccessLog)
	e.Use(middleware.BodyDump())
	e.Use(middleware.Recovery())
	if cfg.Session.Addr != "" {
		e.Use(sessions.Sessions("session", store))
	}
	return e
}

// NotFoundHandler responds not found response.
func NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, "page not found")
}
