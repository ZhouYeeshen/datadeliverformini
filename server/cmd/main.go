package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"business-report-system/internal/config"
	"business-report-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := config.NewDB(cfg.DB)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	if err := config.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("数据库迁移完成")

	if err := config.SeedDefaultAdmin(db); err != nil {
		log.Printf("默认管理员创建失败: %v", err)
	}

	rdb, err := config.NewRedis(cfg.Redis)
	if err != nil {
		log.Printf("Redis 连接失败 (非致命): %v", err)
	} else {
		defer rdb.Close()
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Logger(), gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "time": time.Now().Unix()})
	})

	r.Static("/uploads", "./uploads")

	api := r.Group("/api/v1")
	svcs := registerRoutes(api, db, cfg)
	startScheduler(svcs)

	// Serve admin panel SPA
	r.Static("/assets", "./admin-dist/assets")
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Let API/upload/health/pos routes return their own 404
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/uploads") ||
			path == "/health" || strings.HasPrefix(path, "/pos") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.File("./admin-dist/index.html")
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("服务启动，端口: %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务关闭异常: %v", err)
	}
	log.Println("服务已关闭")
}
