package parse

import (
	"regexp"
	"strings"
)

var (
	keyReg      *regexp.Regexp
	valueReg    *regexp.Regexp
	quoteReg    *regexp.Regexp
	resourceReg *regexp.Regexp
)

type TranslatePair struct {
	Key string
	Val string
}

type TranslatePairDiff struct {
	TranslatePair
	OldVal string
}

func ParseLine(txt string) TranslatePair {
	var key, value string

	txt = strings.TrimSpace(txt)
	if txt == "" {
		return TranslatePair{}
	} else {
		key = KeyReg().FindString(txt)
		key = QuoteReg().FindString(key)
		value = ValueReg().FindString(txt)
		value = QuoteReg().FindString(value)
		key = strings.Trim(key, `"`)
		value = strings.Trim(value, `"`)

		return TranslatePair{
			Key: key,
			Val: value,
		}
	}
}

//key对应的正则表达式
func KeyReg() *regexp.Regexp {
	if keyReg != nil {
		return keyReg
	} else {
		keyReg = regexp.MustCompile(`".+".*=`)
		return keyReg
	}
}

//value对应的正则表达式
func ValueReg() *regexp.Regexp {
	if valueReg != nil {
		return valueReg
	} else {
		valueReg = regexp.MustCompile(`=.*".+"`)
		return valueReg
	}
}

//引号对应的正则表达式
func QuoteReg() *regexp.Regexp {
	if quoteReg != nil {
		return quoteReg
	} else {
		quoteReg = regexp.MustCompile(`".+"`)
		return quoteReg
	}
}

//资源对应的正则表达式
func ResourceReg() *regexp.Regexp {
	if resourceReg != nil {
		return resourceReg
	} else {
		resourceReg = regexp.MustCompile(`LangPackString.langPackString\({1}[^\(]+\){1}`)
		return resourceReg
	}
}
