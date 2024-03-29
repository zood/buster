package resources

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"
)

type Resources struct {
	DevMode       bool
	path          string
	htmlTemplates *template.Template
	posts         postSlice
	postsByID     map[int]Post
	postHtml      map[int][]byte
}

func New(resourcesPath string) (*Resources, error) {
	r := &Resources{
		path: resourcesPath,
	}

	if err := r.loadAll(); err != nil {
		return nil, err
	}

	return r, nil
}

// ExecuteTemplate ...
func (r *Resources) ExecuteTemplate(tmplName string, w io.Writer, data map[string]interface{}) {
	r.ExecuteTemplateCode(tmplName, w, data, http.StatusOK)
}

// ExecuteTemplateCode ...
func (r *Resources) ExecuteTemplateCode(tmplName string, w io.Writer, data map[string]interface{}, httpCode int) {
	if r.DevMode {
		if err := r.loadAll(); err != nil {
			log.Printf("Error reloading templates: %v", err)
			return
		}
	}

	if data == nil {
		data = map[string]interface{}{}
	}
	data["currentYear"] = strconv.Itoa(time.Now().Year())
	if err := r.htmlTemplates.ExecuteTemplate(w, tmplName, data); err != nil {
		log.Printf("Error rendering template '%s': %v", tmplName, err)
	}
}

func (r *Resources) loadPosts() error {
	err := r.parsePostsManifest()
	if err != nil {
		return fmt.Errorf("parsing posts manifest: %v", err)
	}
	sort.Sort(r.posts)
	r.postsByID = make(map[int]Post)
	r.postHtml = make(map[int][]byte)
	for _, p := range r.posts {
		r.postsByID[p.ID] = p
		// load the html of the post
		htmlPath := filepath.Join(r.path, "posts", fmt.Sprintf("%d", p.ID), "index.html")
		htmlData, err := os.ReadFile(htmlPath)
		if err != nil {
			return fmt.Errorf("loading html for post %d: %v", p.ID, err)
		}
		r.postHtml[p.ID] = htmlData
	}

	return nil
}

func (r *Resources) loadTemplates() error {
	tmplsPath := filepath.Join(r.path, "templates")
	entries, err := os.ReadDir(tmplsPath)
	if err != nil {
		return fmt.Errorf("reading templates directory: %v", err)
	}

	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		paths = append(paths, filepath.Join(tmplsPath, entry.Name()))
	}

	tmpls, err := template.New("").ParseFiles(paths...)
	if err != nil {
		return fmt.Errorf("parsing template: %v", err)
	}
	r.htmlTemplates = tmpls

	return nil
}

func (r *Resources) parsePostsManifest() error {
	manifestPath := filepath.Join(r.path, "posts", "manifest.json")
	file, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("opening posts manifest: %v", err)
	}
	defer file.Close()

	r.posts = postSlice{}
	err = json.NewDecoder(file).Decode(&r.posts)
	if err != nil {
		return fmt.Errorf("parsing posts manifest: %v", err)
	}

	return nil
}

func (r *Resources) PostAssetPath(postID int, asset string) string {
	return filepath.Join(r.path, "posts", strconv.Itoa(postID), asset)
}

func (r *Resources) PostBody(id int) template.HTML {
	if r.DevMode {
		if err := r.loadAll(); err != nil {
			log.Printf("Unable to reload posts: %v", err)
			return template.HTML(fmt.Sprintf("<error reloading posts: %v>", err))
		}
	}

	htmlBytes := r.postHtml[id]
	if htmlBytes == nil {
		log.Printf("ERROR: body for post %d requested and not found", id)
		return "<not found>"
	}

	return template.HTML(string(htmlBytes))
}

func (r *Resources) PostById(id int) *Post {
	if r.DevMode {
		if err := r.loadAll(); err != nil {
			log.Printf("Unable to reload resources: %v", err)
			return nil
		}
	}

	post, ok := r.postsByID[id]
	if !ok {
		return nil
	}
	return &post
}

func (r *Resources) Posts(limit, offset uint) []Post {
	if r.DevMode {
		if err := r.loadAll(); err != nil {
			log.Printf("Unable to reload posts: %v", err)
			return nil
		}
	}

	if int(offset) >= len(r.posts) {
		return nil
	}
	start := offset
	end := int(offset + limit)
	if end > len(r.posts) {
		end = len(r.posts)
	}

	return r.posts[start:end]
}

func (r *Resources) loadAll() error {
	err := r.loadTemplates()
	if err != nil {
		return fmt.Errorf("loading templates: %v", err)
	}

	err = r.loadPosts()
	if err != nil {
		return fmt.Errorf("loading posts: %v", err)
	}

	return nil
}
