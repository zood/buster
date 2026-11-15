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
	DevMode bool
	path    string
	loaded  *snapshot
}

func New(resourcesPath string) (*Resources, error) {
	s, err := newSnapshot(resourcesPath)
	if err != nil {
		return nil, err
	}

	return &Resources{path: resourcesPath, loaded: s}, nil
}

// current returns the snapshot to read from. In dev mode, it returns a new
// snapshot loaded from disk, so edits to the resources take effect on the
// next call.
func (r *Resources) current() (*snapshot, error) {
	if r.DevMode {
		// Storing the new snapshot in r.loaded would race with concurrent
		// reads of r.loaded, which take no lock.
		return newSnapshot(r.path)
	}
	return r.loaded, nil
}

// ExecuteTemplate ...
func (r *Resources) ExecuteTemplate(tmplName string, w io.Writer, data map[string]interface{}) {
	r.ExecuteTemplateCode(tmplName, w, data, http.StatusOK)
}

// ExecuteTemplateCode ...
func (r *Resources) ExecuteTemplateCode(tmplName string, w io.Writer, data map[string]interface{}, httpCode int) {
	s, err := r.current()
	if err != nil {
		log.Printf("Error reloading templates: %v", err)
		return
	}

	if data == nil {
		data = map[string]interface{}{}
	}
	data["currentYear"] = strconv.Itoa(time.Now().Year())
	if err := s.htmlTemplates.ExecuteTemplate(w, tmplName, data); err != nil {
		log.Printf("Error rendering template '%s': %v", tmplName, err)
	}
}

func (r *Resources) PostAssetPath(postID int, asset string) string {
	return filepath.Join(r.path, "posts", strconv.Itoa(postID), asset)
}

func (r *Resources) PostBody(id int) template.HTML {
	s, err := r.current()
	if err != nil {
		log.Printf("Unable to reload posts: %v", err)
		return template.HTML(fmt.Sprintf("<error reloading posts: %v>", err))
	}

	htmlBytes := s.postHtml[id]
	if htmlBytes == nil {
		log.Printf("ERROR: body for post %d requested and not found", id)
		return "<not found>"
	}

	return template.HTML(string(htmlBytes))
}

func (r *Resources) PostById(id int) *Post {
	s, err := r.current()
	if err != nil {
		log.Printf("Unable to reload resources: %v", err)
		return nil
	}

	post, ok := s.postsByID[id]
	if !ok {
		return nil
	}
	return &post
}

func (r *Resources) Posts(limit, offset uint) []Post {
	s, err := r.current()
	if err != nil {
		log.Printf("Unable to reload posts: %v", err)
		return nil
	}

	return s.page(limit, offset)
}

// PostsWithBodies returns the same posts as Posts, along with their bodies.
// In dev mode, it reloads the resources once, whereas calling PostBody for
// each post reloads them on every call.
func (r *Resources) PostsWithBodies(limit, offset uint) []PostWithBody {
	s, err := r.current()
	if err != nil {
		log.Printf("Unable to reload posts: %v", err)
		return nil
	}

	posts := s.page(limit, offset)
	result := make([]PostWithBody, 0, len(posts))
	for _, p := range posts {
		result = append(result, PostWithBody{Post: p, Body: template.HTML(s.postHtml[p.ID])})
	}
	return result
}

// A snapshot holds the templates and posts read from a resources directory.
// It is not modified after newSnapshot returns.
type snapshot struct {
	htmlTemplates *template.Template
	posts         postSlice
	postsByID     map[int]Post
	postHtml      map[int][]byte
}

func newSnapshot(path string) (*snapshot, error) {
	s := &snapshot{}

	err := s.loadTemplates(path)
	if err != nil {
		return nil, fmt.Errorf("loading templates: %v", err)
	}

	err = s.loadPosts(path)
	if err != nil {
		return nil, fmt.Errorf("loading posts: %v", err)
	}

	return s, nil
}

func (s *snapshot) loadPosts(path string) error {
	err := s.parsePostsManifest(path)
	if err != nil {
		return fmt.Errorf("parsing posts manifest: %v", err)
	}
	sort.Sort(s.posts)
	s.postsByID = make(map[int]Post)
	s.postHtml = make(map[int][]byte)
	for _, p := range s.posts {
		s.postsByID[p.ID] = p
		// load the html of the post
		htmlPath := filepath.Join(path, "posts", fmt.Sprintf("%d", p.ID), "index.html")
		htmlData, err := os.ReadFile(htmlPath)
		if err != nil {
			return fmt.Errorf("loading html for post %d: %v", p.ID, err)
		}
		s.postHtml[p.ID] = htmlData
	}

	return nil
}

func (s *snapshot) loadTemplates(path string) error {
	tmplsPath := filepath.Join(path, "templates")
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
	s.htmlTemplates = tmpls

	return nil
}

// page returns up to limit posts, newest first, skipping the first offset.
func (s *snapshot) page(limit, offset uint) []Post {
	if int(offset) >= len(s.posts) {
		return nil
	}
	start := offset
	end := int(offset + limit)
	if end > len(s.posts) {
		end = len(s.posts)
	}

	return s.posts[start:end]
}

func (s *snapshot) parsePostsManifest(path string) error {
	manifestPath := filepath.Join(path, "posts", "manifest.json")
	file, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("opening posts manifest: %v", err)
	}
	defer file.Close()

	s.posts = postSlice{}
	err = json.NewDecoder(file).Decode(&s.posts)
	if err != nil {
		return fmt.Errorf("parsing posts manifest: %v", err)
	}

	return nil
}
