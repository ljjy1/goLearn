package redis

/**
 * @Description: 初始化redis连接
 */
import (
	"context"
	"fmt"
	"homework4/config"
	"homework4/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

// 初始化redis配置
func InitRedis() {
	dsn := config.Cfg.Redis.Host + ":" + fmt.Sprintf("%d", config.Cfg.Redis.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     dsn,
		Password: config.Cfg.Redis.Password,
		DB:       config.Cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		logger.AppLog.Fatal("redis连接失败", logger.WrapMeta(err)...)
	}
	logger.AppLog.Info("redis连接成功")
}
