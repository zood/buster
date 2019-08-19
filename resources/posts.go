package resources

// Post contains all the information about a blog/news item, except the actual body
type Post struct {
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
	Intro string `json:"intro"`
}

type postSlice []Post

func (ps postSlice) Len() int {
	return len(ps)
}

func (ps postSlice) Less(i, j int) bool {
	return ps[i].ID > ps[j].ID
}

func (ps postSlice) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// func loadPosts(resourcesPath string) ([]Post, error) {
// 	manifestPath := filepath.Join(resourcesPath, "posts", "manifest.json")
// 	file, err := os.Open(manifestPath)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to open posts manifest")
// 	}
// 	defer file.Close()

// 	var posts []Post
// 	err = json.NewDecoder(file).Decode(&posts)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to parse posts manifest")
// 	}

// 	return posts, nil
// }
