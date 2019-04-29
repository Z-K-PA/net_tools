package main

import (
	"github.com/pineal-niwan/busybox/tools/code_gen/binary"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"log"
	"os"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("不能初始化logger %+v", err)
	}
	defer logger.Sync()

	runWithLogger := func(c *cli.Context) error {
		return generateMsgAPICode(c, logger)
	}

	app := cli.App{
		Name:    "二进制消息API封装工具",
		Usage:   "用于将定义的消息封装成API的代码生成",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "template",
				Usage: "模版文件",
			},
			&cli.StringFlag{
				Name:  "in",
				Usage: "输入文件",
			},
			&cli.StringFlag{
				Name:  "out",
				Usage: "输出文件",
			},
		},
		Action: runWithLogger,
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Error(
			"app error",
			zap.Error(err),
		)
	}
}

func generateMsgAPICode(c *cli.Context, logger *zap.Logger) error {
	var apiDef binary.API
	err := binary.GenCode(c, logger, &apiDef)
	return err
}
