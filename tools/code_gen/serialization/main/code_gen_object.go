package main

import (
	"errors"
	"github.com/go-yaml/yaml"
	"github.com/pineal-niwan/busybox/tools/code_gen/serialization"
	"github.com/pineal-niwan/busybox/util"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"log"
	"os"
	"text/template"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("不能初始化logger %+v", err)
	}
	defer logger.Sync()

	runWithLogger := func(c *cli.Context) error {
		return generatewithLogger(c, logger)
	}

	app := cli.App{
		Name:    "二进制序列化工具",
		Usage:   "用于将自定义的结构体做序列化与反序列化的代码生成",
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

func generatewithLogger(c *cli.Context, logger *zap.Logger) error {
	var packageDef serialization.Package

	templateFileName := c.String("template")
	inFileName := c.String("in")
	outFileName := c.String("out")

	if templateFileName == "" {
		return errors.New("not specific template file path")
	}
	if inFileName == "" {
		return errors.New("not specific input file path")
	}
	if outFileName == "" {
		return errors.New("not specific output file path")
	}

	templateBuf, err := util.ReadFile2Buffer(templateFileName)
	if err != nil {
		return err
	}

	err = util.UnMarshalFile2Object(yaml.Unmarshal, inFileName, &packageDef)
	if err != nil {
		return err
	}

	rmErr := os.Remove(outFileName)
	if rmErr != nil {
		logger.Error("rm file fail",
			zap.Error(rmErr))
	}

	outFile, err := os.OpenFile(outFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	t, err := template.New("go").Funcs(serialization.FuncHash).Parse(string(templateBuf))
	if err != nil {
		return err
	}

	err = t.Execute(outFile, packageDef)
	if err == nil {
		logger.Info("generate object code completed.")
	}

	return err
}
