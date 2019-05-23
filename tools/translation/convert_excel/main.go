package main

import (
	"bufio"
	"github.com/pineal-niwan/busybox/tools/translation/parse"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
)

func main() {
	cmd := &cli.App{
		Name:    "提取需要翻译的文档",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "输入的文本文件",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "输入的excel文件",
			},
		},
		Action: convertFile,
	}
	err := cmd.Run(os.Args)
	if err != nil {
		log.Printf("运行失败:%+v\n", err)
	}
}

func convertFile(c *cli.Context) error {
	input := c.String("input")
	f, err := os.Open(input)
	if err != nil {
		return err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		txt, _, err := r.ReadLine()
		if err == io.EOF {
			return nil
		}
		pair := parse.ParseLine(string(txt))
		log.Printf("k:%+v -- v:%+v\n", pair.Key, pair.Val)
	}

}
