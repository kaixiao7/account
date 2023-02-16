package account

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kaixiao7/account/internal/account/store"
	"kaixiao7/account/internal/pkg/middleware"
	"kaixiao7/account/internal/pkg/token"
	"kaixiao7/account/internal/pkg/validatetrans"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"kaixiao7/account/internal/pkg/log"
)

const (
	// 默认配置文件名称
	defaultConfigName = "account.yaml"

	// 应用名称
	appName = "account"
)

var cfgFile string

func NewAccountServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          appName,
		Short:        "account server",
		Long:         "",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// 初始化日志配置
			log.Init(logOptions())
			defer log.Sync()

			return run()
		},
	}

	cobra.OnInitialize(initConfig)

	// 将命令行传递的配置文件持久化到变量cfgFile
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	return cmd
}

func run() error {
	// 初始化数据库配置
	db, err := store.Init(&store.DBOption{
		Dsn:                   viper.GetString("db.dsn"),
		MaxIdleConnections:    viper.GetInt("db.max-idle-connections"),
		MaxOpenConnections:    viper.GetInt("db.max-open-connections"),
		MaxConnectionLifeTime: viper.GetInt("db.max-connection-life-time"),
	})
	if err != nil {
		return err
	}
	defer store.Close()

	// 初始化jwt相关设置
	token.Init(viper.GetString("jwt.secret"),
		viper.GetInt64("jwt.token-access-expire"),
		viper.GetInt64("jwt.token-refresh-expire"))

	// 初始化gin参数验证
	if err := validatetrans.Init(); err != nil {
		return err
	}

	// 设置gin
	gin.SetMode(viper.GetString("run_mode"))
	g := gin.New()
	loadRouter(g, middleware.RequestId(), middleware.SqlDB(db))

	server := &http.Server{
		Addr:    viper.GetString("port"),
		Handler: g,
	}

	// 在goroutine中启动server，将不会阻塞下面的优雅关闭逻辑
	go func() {
		log.Infof("HTTP server ListenAndServe on port: %s", viper.GetString("port"))
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Errorf("HTTP server ListenAndServe: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("HTTP server Shutdown: %v", err)
		return err
	}

	log.Info("Server exiting\n")

	return nil
}

// 读取配置文件
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// 添加搜索路径
		viper.AddConfigPath(".")
		viper.SetConfigName(defaultConfigName)
	}

	viper.SetConfigType("yaml")

	// 读取配置文件内容
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
}

func logOptions() *log.Options {
	return &log.Options{
		Level:            viper.GetString("log.level"),
		Format:           viper.GetString("log.format"),
		EnableColor:      viper.GetBool("log.enable-color"),
		Path:             viper.GetString("log.path"),
		Filename:         viper.GetString("log.filename"),
		MaxSize:          viper.GetInt("log.max-size"),
		MaxAge:           viper.GetInt("log.max-age"),
		MaxBackups:       viper.GetInt("log.max-backups"),
		Compress:         viper.GetBool("log.compress"),
		EnableCaller:     viper.GetBool("log.enable-caller"),
		EnableStacktrace: viper.GetBool("log.enable-stacktrace"),
		EnableStdout:     viper.GetBool("log.enable-stdout"),
	}
}
