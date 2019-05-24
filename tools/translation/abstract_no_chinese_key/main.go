package main

import (
	"bufio"
	"github.com/pineal-niwan/busybox/tools/translation/parse"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"unicode"
)

func main() {
	cmd := &cli.App{
		Name:    "提取需要中文翻译的文档中没有翻译的情况",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "input",
				Usage: "输入的文本文件",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "输出的excel文件",
			},
		},
		Action: abstractNoChinese,
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

func containChineseString(val string) bool {
	r := []rune(val)
	for _, v := range r {
		if unicode.Is(unicode.Han, v) {
			return true
		}
	}
	return false
}

//提取不同
func abstractNoChinese(c *cli.Context) error {
	var noChineseKeyList []parse.TranslatePairDiff

	input := c.String("input")
	inputPairList, err := getKVFormFile(input)
	if err != nil {
		return err
	}

	for _, k := range inputPairList {
		if !containChineseString(k.Val) {
			noChineseKeyList = append(noChineseKeyList, parse.TranslatePairDiff{
				TranslatePair: k,
			})
		}
	}

	excelFile := xlsx.NewFile()
	sheet, err := excelFile.AddSheet("sheet1")
	if err != nil {
		return err
	}
	row := sheet.AddRow()
	cell := row.AddCell()
	cell.SetString("没有翻译成中文的条目")
	for _, addKey := range noChineseKeyList {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(addKey.Key)
		cell = row.AddCell()
		cell.SetString(addKey.Val)
	}

	err = excelFile.Save(c.String("output"))
	return err
}
