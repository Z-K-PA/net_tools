package main

import (
	"bufio"
	"fmt"
	"github.com/pineal-niwan/busybox/tools/translation/parse"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	cmd := &cli.App{
		Name:    "提取需要翻译的文档中的相同key",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "base",
				Usage: "基准的文本文件",
			},
			&cli.StringFlag{
				Name:  "input",
				Usage: "输入的文本文件",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "输出的文本文件",
			},
		},
		Action: abstractCommon,
	}
	err := cmd.Run(os.Args)
	if err != nil {
		log.Printf("运行失败:%+v\n", err)
	}
}

//获取键值对
func getKVFormFile(input string) ([]parse.TranslatePair, error) {
	var keyPairList []parse.TranslatePair
	f, err := os.Open(input)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		txt, _, err := r.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		pair := parse.ParseLine(string(txt))
		if pair.Key != "" {
			keyPairList = append(keyPairList, pair)
		}
	}
	return keyPairList, nil
}

//提取不同
func abstractCommon(c *cli.Context) error {
	var commonKeyList []parse.TranslatePairDiff

	base := c.String("base")
	input := c.String("input")
	output := c.String("output")

	basePairList, err := getKVFormFile(base)
	if err != nil {
		return err
	}
	inputPairList, err := getKVFormFile(input)
	if err != nil {
		return err
	}

	m2 := make(map[string]parse.TranslatePairDiff)
	for _, k2 := range inputPairList {
		m2[k2.Key] = parse.TranslatePairDiff{
			TranslatePair: k2,
		}
	}

	for _, k1 := range basePairList {
		k2, ok := m2[k1.Key]
		if ok {
			//两个文件都有
			commonKeyList = append(commonKeyList, k2)
		}
	}

	var outputStr strings.Builder
	for _, kv := range commonKeyList {
		fmt.Fprintf(&outputStr, "\"%s\" = \"%s\"\n", kv.Key, kv.Val)
	}
	err = ioutil.WriteFile(output, []byte(outputStr.String()), 0644)

	return err
}
