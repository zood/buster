package rsrc

import (
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"
)

// Development ...
var Development = false

var loadedTemplates *template.Template

var rsrcsPath string

// Init ..
func Init(resourcesPath string) error {
	rsrcsPath = resourcesPath
	return loadTemplates()
}

// LoadTemplates ...
func loadTemplates() error {
	tmplsPath := filepath.Join(rsrcsPath, "templates")
	fis, err := ioutil.ReadDir(tmplsPath)
	if err != nil {
		return err
	}

	var paths []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		paths = append(paths, filepath.Join(tmplsPath, fi.Name()))
	}

	loadedTemplates, err = template.New("").ParseFiles(paths...)
	return err
}

// ExecuteTemplate ...
func ExecuteTemplate(tmplName string, w io.Writer, data map[string]interface{}) {
	if Development {
		if err := loadTemplates(); err != nil {
			log.Printf("Error reloading templates: %v", err)
			return
		}
	}

	if data == nil {
		data = map[string]interface{}{}
	}
	data["currentYear"] = strconv.Itoa(time.Now().Year())
	if err := loadedTemplates.ExecuteTemplate(w, tmplName, data); err != nil {
		log.Printf("Error rendering template '%s': %v", tmplName, err)
	}
}
