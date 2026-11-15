package l10n

import (
	"html/template"
	"log"
	"reflect"
	"sync"

	"github.com/russross/blackfriday/v2"
	"golang.org/x/text/language"
)

var (
	markdownMu    sync.RWMutex
	markdownCache = map[StringAsset]template.HTML{}
)

func Markdown(tag language.Tag, asset StringAsset) template.HTML {
	markdownMu.RLock()
	cache, ok := markdownCache[asset]
	markdownMu.RUnlock()
	if ok {
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

	markdownMu.Lock()
	markdownCache[asset] = output
	markdownMu.Unlock()

	return output
}
