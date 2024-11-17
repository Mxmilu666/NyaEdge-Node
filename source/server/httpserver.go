package server

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"nyaedge-node/source"
	"nyaedge-node/source/zaplogger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 请求数据结构
type RequestData struct {
	Nonce     string `json:"nonce"`
	Timestamp string `json:"timestamp"`
	Hash      string `json:"hash"`
}

func StartServer(config *source.Config, logger *zap.Logger) error {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(zaplogger.ZapLogger(logger))

	api := r.Group("/api")
	{
		nodeapi := api.Group("/node")
		{
			nodeapi.GET("/ping", AuthMiddleware(config.Center.NodeSecret), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"status":  "success",
					"message": "pong",
				})
			})
		}
	}

	address := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)

	logger.Info(fmt.Sprintf("Server is running at http://%s", address))
	return r.Run(address)
}

// AuthMiddleware 验证中间件
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data RequestData

		if err := c.BindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "invalid request", "error": err.Error()})
			c.Abort()
			return
		}

		nonce := data.Nonce
		timestamp := data.Timestamp
		hashValue := data.Hash

		// 验证时间戳是否在允许的时间窗口内
		requestTime, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "invalid timestamp", "error": err.Error()})
			c.Abort()
			return
		}

		if time.Since(requestTime) > 5*time.Minute {
			c.JSON(http.StatusForbidden, gin.H{"status": "request expired"})
			c.Abort()
			return
		}

		// 服务器使用相同的密钥和时间戳生成哈希值并进行对比
		expectedHash := createHash(nonce, secretKey, timestamp)
		if expectedHash != hashValue {
			c.JSON(http.StatusBadRequest, gin.H{"status": "failure", "message": "Invalid nonce or hash"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// createHash 使用密钥加盐生成哈希值
func createHash(nonce, secretKey, timestamp string) string {
	saltedNonce := nonce + secretKey + timestamp
	hash := sha256.New()
	hash.Write([]byte(saltedNonce))
	return hex.EncodeToString(hash.Sum(nil))
}
