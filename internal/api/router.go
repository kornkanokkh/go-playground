// internal/api/router.go
package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-playground/internal/config"
	"go-playground/internal/logger"
	"go.uber.org/zap" // For zap.Field
)

// Router struct จะเก็บ Echo instance และ dependencies ที่จำเป็น
type Router struct {
	EchoInstance *echo.Echo
	logger       logger.Logger
	config       *config.Config // หรือจะส่งแค่ AppConfig เข้ามาก็ได้
	//dbRepo       *repositories.GormDatabase // เปลี่ยนชื่อจาก dbSql เป็น dbRepo เพื่อให้ชัดเจน
}

// NewRouter สร้างและคืนค่า Router instance พร้อมตั้งค่าพื้นฐานและ routes
func NewRouter(e *echo.Echo, cfg *config.Config, appLogger logger.Logger) *Router {
	r := &Router{
		EchoInstance: e,
		logger:       appLogger,
		config:       cfg,
		//dbRepo:       dbRepo,
	}

	r.setupMiddleware() // ตั้งค่า Middleware ทั่วไป
	r.setupRoutes()     // ตั้งค่า Routes

	return r
}

// setupMiddleware ตั้งค่า Global Middleware สำหรับ Echo
func (r *Router) setupMiddleware() {
	e := r.EchoInstance
	appLogger := r.logger // ใช้ logger จาก Router instance

	e.Use(middleware.Recover()) // Recovers from panics anywhere in the chain and handles the error
	e.Use(middleware.CORS())    // Enables Cross-Origin Resource Sharing (CORS)

	// Custom Logger for Echo middleware (to use Zap instead of default log)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			// Only log start/end in dev for less noise in production
			start := r.config.App.Env == "development"

			if start {
				appLogger.Debug("Request started",
					zap.String("method", req.Method),
					zap.String("uri", req.RequestURI),
					zap.String("remote_ip", c.RealIP()))
			}

			err := next(c)

			if start {
				appLogger.Debug("Request finished",
					zap.String("method", req.Method),
					zap.String("uri", req.RequestURI),
					zap.Int("status", res.Status),
					zap.Int64("response_size", res.Size),
					zap.Error(err))
			} else {
				// Log only essential info for production, especially errors or high status codes
				if err != nil || res.Status >= 400 {
					appLogger.Warn("Request finished with status or error",
						zap.String("method", req.Method),
						zap.String("uri", req.RequestURI),
						zap.Int("status", res.Status),
						zap.Error(err))
				}
			}
			return err
		}
	})

}

// setupRoutes กำหนดเส้นทาง API ทั้งหมด
func (r *Router) setupRoutes() {
	e := r.EchoInstance
	appLogger := r.logger // ใช้ logger จาก Router instance
	//dbRepo := r.dbRepo    // ใช้ dbRepo จาก Router instance

	// Health Check endpoint
	e.GET("/health", r.healthCheckHandler(appLogger))
	e.GET("/hello/:name", r.helloHandler(appLogger))

	apiGroup := e.Group("/api/v1")

	apiGroup.GET("/items", r.getItemsHandler(appLogger))
	apiGroup.POST("/items", r.createItemHandler(appLogger))
}
