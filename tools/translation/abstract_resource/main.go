package main

import (
	"fmt"
	"github.com/pineal-niwan/busybox/tools/translation/parse"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	cmd := &cli.App{
		Name:    "提取翻译资源文档",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "输入的文本文件",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "输出的文本文件",
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
	output := c.String("output")

	var outputStr strings.Builder


	txt, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	kvList := parse.ResourceReg().FindAllString(string(txt), -1)

	braceReg := regexp.MustCompile(`\(.+\)`)
	kReg := regexp.MustCompile(`key[^\,]+`)
	vReg := regexp.MustCompile(`value[^\)]+`)


	for _, kv := range kvList {
		kv = strings.Replace(kv, "\n", "\\n", -1)
		kv = braceReg.FindString(kv)

		k := kReg.FindString(kv)
		k = strings.Replace(k, "key:", "", -1)
		k = strings.TrimSpace(k)

		v := vReg.FindString(kv)
		v = strings.Replace(v, "value:", "", -1)
		v = strings.Replace(v, `"`, `\"`, -1)
		v = strings.TrimSpace(v)

		log.Printf(`"%s" = "%s"`, k, v)

		fmt.Fprintf(&outputStr, "\"%s\" = \"%s\"\n", k, v)
	}

	err = ioutil.WriteFile(output, []byte(outputStr.String()), 0644)

	return err
}
