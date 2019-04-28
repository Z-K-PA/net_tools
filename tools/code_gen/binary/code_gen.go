package binary

import (
	"errors"
	"github.com/go-yaml/yaml"
	"github.com/pineal-niwan/busybox/util"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"os"
	"text/template"
)

func GenCode(c *cli.Context, logger *zap.Logger, val interface{}) error {
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

	err = util.UnMarshalFile2Object(yaml.Unmarshal, inFileName, val)
	if err != nil {
		return err
	}

	_, err = os.Stat(outFileName)
	if !os.IsNotExist(err) {
		rmErr := os.Remove(outFileName)
		if rmErr != nil {
			logger.Error("rm file fail",
				zap.Error(rmErr))
		}
	}

	outFile, err := os.OpenFile(outFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	t, err := template.New("go").Funcs(FuncHash).Parse(string(templateBuf))
	if err != nil {
		return err
	}

	err = t.Execute(outFile, val)
	if err == nil {
		logger.Info("generate object code completed.")
	}

	return err
}
