package helpers

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-unidecode"
)

type TextNormalize struct{}

func NewTextNormalize() *TextNormalize {
	return &TextNormalize{}
}

func (tp *TextNormalize) Normalize(text string) string {
	text = unidecode.Unidecode(text)

	text = strings.ToLower(text)

	reg, _ := regexp.Compile(`[^a-z0-9\s]+`)
	text = reg.ReplaceAllString(text, "")

	whiteSpaceReg, _ := regexp.Compile(`\s+`)
	text = whiteSpaceReg.ReplaceAllString(text, " ")

	text = strings.TrimSpace(text)

	return text
}
