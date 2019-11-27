package l10n

import (
	"html/template"
	"log"
	"reflect"

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

	var output template.HTML
	switch t := s.(type) {
	case string:
		output = template.HTML(blackfriday.Run([]byte(t)))
	case template.HTML:
		output = template.HTML(blackfriday.Run([]byte(t)))
	default:
		log.Printf("ERROR: unsupported type (%v) for asset %d", reflect.TypeOf(s), asset)
		return "&lt;undefined&gt;"
	}
	markdownCache[asset] = output

	return output
}
