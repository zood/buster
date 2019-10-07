package l10n

import (
	"html/template"
	"log"

	"github.com/russross/blackfriday/v2"
	"golang.org/x/text/language"
)

var markdownCache = map[StringAsset]template.HTML{}

func Markdown(tag language.Tag, asset StringAsset) template.HTML {
	if cache, ok := markdownCache[asset]; ok {
		return cache
	}

	s, ok := enStrings[asset]
	if !ok {
		log.Printf("WARNING: no entry found for string asset: %d", asset)
		return "&lt;undefined&gt;"
	}

	output := template.HTML(string(blackfriday.Run([]byte(s))))
	markdownCache[asset] = output

	return output
}
