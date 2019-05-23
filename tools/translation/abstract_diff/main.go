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
		Name:    "提取需要翻译的文档中的不同",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input1",
				Usage: "输入的文本文件1",
			},
			&cli.StringFlag{
				Name:  "input2",
				Usage: "输入的文本文件2",
			},
		},
		Action: abstractDiff,
	}
	err := cmd.Run(os.Args)
	if err != nil {
		log.Printf("运行失败:%+v\n", err)
	}
}

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

func abstractDiff(c *cli.Context) error {
	var addedKeyList []parse.TranslatePairDiff
	var deletedKeyList []parse.TranslatePairDiff
	var differentKeyList []parse.TranslatePairDiff

	input1 := c.String("input1")
	input2 := c.String("input2")

	inputPairList1, err := getKVFormFile(input1)
	if err != nil {
		return err
	}
	inputPairList2, err := getKVFormFile(input2)
	if err != nil {
		return err
	}

	m1 := make(map[string]parse.TranslatePairDiff)
	m2 := make(map[string]parse.TranslatePairDiff)

	for _, k1 := range inputPairList1 {
		m1[k1.Key] = parse.TranslatePairDiff{
			TranslatePair: k1,
		}
	}

	for _, k2 := range inputPairList2 {
		m2[k2.Key] = parse.TranslatePairDiff{
			TranslatePair: k2,
		}
	}

	for _, k2 := range inputPairList2 {
		k1, ok := m1[k2.Key]
		if !ok {
			//input2有input1没有
			addedKeyList = append(addedKeyList, parse.TranslatePairDiff{
				TranslatePair: k2,
			})
		} else {
			if k1.Val != k2.Val {
				//input1和input2不同
				differentKeyList = append(differentKeyList, parse.TranslatePairDiff{
					TranslatePair: k2,
					OldVal:        k1.Val,
				})
			}
		}
	}

	for _, k1 := range inputPairList1 {
		_, ok := m2[k1.Key]
		if !ok {
			//input1有input2没有
			deletedKeyList = append(addedKeyList, parse.TranslatePairDiff{
				TranslatePair: k1,
			})
		}
	}

	log.Printf("新增的项目：%+v\n", addedKeyList)
	log.Printf("删除的项目：%+v\n", deletedKeyList)
	log.Printf("更改的项目：%+v\n", differentKeyList)

	return nil
}
