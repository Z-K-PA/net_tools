package main

import (
	"bufio"
	"github.com/pineal-niwan/busybox/tools/translation/parse"
	"github.com/tealeg/xlsx"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"strings"
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
			&cli.StringSliceFlag{
				Name:  "skipWords",
				Usage: "跳过比较的文字",
			},
			&cli.StringFlag{
				Name:  "output",
				Usage: "输出的excel文件",
			},
		},
		Action: abstractDiff,
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
func abstractDiff(c *cli.Context) error {
	var addedKeyList []parse.TranslatePairDiff
	var deletedKeyList []parse.TranslatePairDiff
	var differentKeyList []parse.TranslatePairDiff

	input1 := c.String("input1")
	input2 := c.String("input2")

	skipWordList := c.StringSlice("skipWords")
	log.Printf("skip word list:%+v\n", skipWordList)

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
			k1Val := k1.Val
			k2Val := k2.Val

			if k1Val != k2Val {
				k1Val = strings.Replace(k1Val, " ", "", -1)
				k2Val = strings.Replace(k2Val, " ", "", -1)

				for _, skipWord := range skipWordList {
					k1Val = strings.Replace(k1Val, skipWord, "", -1)
					k2Val = strings.Replace(k2Val, skipWord, "", -1)
				}

				if k1Val != k2Val {
					differentKeyList = append(differentKeyList, parse.TranslatePairDiff{
						TranslatePair: k2,
						OldVal:        k1.Val,
					})
				}
			}
		}
	}

	for _, k1 := range inputPairList1 {
		_, ok := m2[k1.Key]
		if !ok {
			//input1有input2没有
			deletedKeyList = append(deletedKeyList, parse.TranslatePairDiff{
				TranslatePair: k1,
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
	cell.SetString("新增的条目")
	for _, addKey := range addedKeyList {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(addKey.Key)
		cell = row.AddCell()
		cell.SetString(addKey.Val)
	}
	row = sheet.AddRow()
	row = sheet.AddRow()
	row = sheet.AddRow()
	row = sheet.AddRow()
	row = sheet.AddRow()

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetString("更改的条目")
	for _, diffKey := range differentKeyList {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(diffKey.Key)
		cell = row.AddCell()
		cell.SetString(diffKey.Val)
		cell = row.AddCell()
		cell.SetString(diffKey.OldVal)
	}
	row = sheet.AddRow()
	row = sheet.AddRow()
	row = sheet.AddRow()
	row = sheet.AddRow()
	row = sheet.AddRow()

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetString("删除的条目")
	for _, deletedKey := range deletedKeyList {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.SetString(deletedKey.Key)
		cell = row.AddCell()
		cell.SetString(deletedKey.Val)
	}

	err = excelFile.Save(c.String("output"))

	return err
}
