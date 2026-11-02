package feed

import (
	"encoding/xml"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"zood.xyz/buster/resources"
)

const siteURL = "https://www.zood.xyz"

// parsedFeed is the decoding counterpart of rss. Its tags use namespace URLs
// instead of the literal prefixes in rss, because encoding/xml resolves
// prefixes to namespace URLs when decoding.
type parsedFeed struct {
	Version string `xml:"version,attr"`
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Language    string `xml:"language"`
		// A field tag without a namespace matches elements in any
		// namespace, so Links holds both <link> and <atom:link>.
		Links []struct {
			XMLName xml.Name
			Href    string `xml:"href,attr"`
			Rel     string `xml:"rel,attr"`
			Value   string `xml:",chardata"`
		} `xml:"link"`
		LastBuildDate string `xml:"lastBuildDate"`
		Items         []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
			GUID  struct {
				IsPermaLink string `xml:"isPermaLink,attr"`
				Value       string `xml:",chardata"`
			} `xml:"guid"`
			PubDate     string `xml:"pubDate"`
			Description string `xml:"description"`
			Content     string `xml:"http://purl.org/rss/1.0/modules/content/ encoded"`
		} `xml:"item"`
	} `xml:"channel"`
}

func renderAndParse(t *testing.T, posts []resources.PostWithBody) parsedFeed {
	t.Helper()
	data, err := RSS(Channel{
		SiteURL:     siteURL,
		Title:       "Zood Blog",
		Description: "News & announcements",
		Language:    "en",
	}, posts)
	require.NoError(t, err)

	var feed parsedFeed
	require.NoError(t, xml.Unmarshal(data, &feed))
	return feed
}

func TestRSS(t *testing.T) {
	posts := []resources.PostWithBody{
		{
			Post: resources.Post{ID: 7, Slug: "second", Title: "Second <post>", Intro: "Intro & more", Date: 1751609101},
			Body: "<h1>Second post</h1>\n\n" +
				`<p>See <a href="/blog/6/first">the first</a> and <img src="shot.png" /></p>` + "\n<p>]]> survives</p>",
		},
		{
			Post: resources.Post{ID: 6, Slug: "first", Title: "First", Intro: "Hi", Date: 1607446800},
			Body: `<p>Hello</p>`,
		},
	}
	feed := renderAndParse(t, posts)

	require.Equal(t, "2.0", feed.Version)
	ch := feed.Channel
	require.Equal(t, "Zood Blog", ch.Title)
	require.Equal(t, "News & announcements", ch.Description)
	require.Equal(t, "en", ch.Language)
	require.Equal(t, "Fri, 04 Jul 2025 06:05:01 +0000", ch.LastBuildDate)
	require.Len(t, ch.Links, 2)
	require.Equal(t, "", ch.Links[0].XMLName.Space)
	require.Equal(t, siteURL+"/blog", ch.Links[0].Value)
	require.Equal(t, "http://www.w3.org/2005/Atom", ch.Links[1].XMLName.Space)
	require.Equal(t, siteURL+"/blog/rss.xml", ch.Links[1].Href)
	require.Equal(t, "self", ch.Links[1].Rel)

	require.Len(t, ch.Items, 2)
	item := ch.Items[0]
	require.Equal(t, "Second <post>", item.Title)
	require.Equal(t, siteURL+"/blog/7/second", item.Link)
	require.Equal(t, siteURL+"/blog/7", item.GUID.Value)
	require.Equal(t, "false", item.GUID.IsPermaLink)
	require.Equal(t, "Intro & more", item.Description)
	pubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
	require.NoError(t, err)
	require.Equal(t, int64(1751609101), pubDate.Unix())
	require.Equal(t,
		`<p>See <a href="https://www.zood.xyz/blog/6/first">the first</a> and <img src="https://www.zood.xyz/blog/7/shot.png"/></p>`+
			"\n<p>]]> survives</p>",
		item.Content)

	require.Equal(t, siteURL+"/blog/6/first", ch.Items[1].Link)
}

func TestRSSWithNoPosts(t *testing.T) {
	feed := renderAndParse(t, nil)
	require.Empty(t, feed.Channel.Items)
	require.Empty(t, feed.Channel.LastBuildDate)
}

func TestAbsolutizeURLs(t *testing.T) {
	base, err := url.Parse("https://www.zood.xyz/blog/2/some-post")
	require.NoError(t, err)

	tests := []struct {
		name, in, out string
	}{
		{"relative", `<a href="pic.png">x</a>`, `<a href="https://www.zood.xyz/blog/2/pic.png">x</a>`},
		{"root relative", `<a href="/about">x</a>`, `<a href="https://www.zood.xyz/about">x</a>`},
		{"fragment", `<a href="#notes">x</a>`, `<a href="https://www.zood.xyz/blog/2/some-post#notes">x</a>`},
		{"protocol relative", `<img src="//example.com/a.png">`, `<img src="https://example.com/a.png">`},
		{"self-closing", `<img src="a.png" />`, `<img src="https://www.zood.xyz/blog/2/a.png"/>`},
		{"absolute untouched", `<a  HREF="https://Example.com/a/../b?q=1&amp;r=2">x</a>`, `<a  HREF="https://Example.com/a/../b?q=1&amp;r=2">x</a>`},
		{"mailto untouched", `<a href="mailto:hello@zood.xyz">x</a>`, `<a href="mailto:hello@zood.xyz">x</a>`},
		{"other attributes kept", `<img class="wide" src="a.png" alt="A &amp; B">`, `<img class="wide" src="https://www.zood.xyz/blog/2/a.png" alt="A &amp; B">`},
		{"text untouched", "<p>It&rsquo;s &amp; ok “quoted” href=\"x\"</p>\n", "<p>It&rsquo;s &amp; ok “quoted” href=\"x\"</p>\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := absolutizeURLs(tc.in, base)
			require.NoError(t, err)
			require.Equal(t, tc.out, out)
		})
	}
}

func TestStripLeadingH1(t *testing.T) {
	tests := []struct {
		name, in, out string
	}{
		{"leading h1", "<h1>Title</h1>\n\n<p>Body</p>\n", "<p>Body</p>\n"},
		{"markup inside h1", `<h1>A <em>big</em> <a href="/x">day</a></h1><p>Body</p>`, "<p>Body</p>"},
		{"whitespace and comment first", "\n  <!-- draft -->\n<h1>Title</h1>\n<p>Body</p>", "<p>Body</p>"},
		{"only the first h1", "<h1>One</h1>\n<h1>Two</h1>", "<h1>Two</h1>"},
		{"starts with paragraph", "<p>Body</p>\n<h1>Later</h1>", "<p>Body</p>\n<h1>Later</h1>"},
		{"text before h1", "Hello <h1>Title</h1>", "Hello <h1>Title</h1>"},
		{"wrapped h1", "<div><h1>Title</h1></div>", "<div><h1>Title</h1></div>"},
		{"h2 first", "<h2>Title</h2><p>Body</p>", "<h2>Title</h2><p>Body</p>"},
		{"unclosed h1", "<h1>Title", "<h1>Title"},
		{"empty", "", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.out, stripLeadingH1(tc.in))
		})
	}
}

// TestRSSFromAssets checks the feed built from the posts in ../assets. It
// fails if a post's href or src is still relative in the feed, or if a post
// still begins with its h1.
func TestRSSFromAssets(t *testing.T) {
	rsrcs, err := resources.New("../assets")
	require.NoError(t, err)
	posts := rsrcs.PostsWithBodies(1000, 0)
	require.NotEmpty(t, posts)

	feed := renderAndParse(t, posts)
	require.Len(t, feed.Channel.Items, len(posts))

	attrRE := regexp.MustCompile(`\b(?:href|src)="([^"]*)"`)
	for i, item := range feed.Channel.Items {
		require.Equal(t, posts[i].Title, item.Title)
		// Posts that aren't in the manifest, such as the sample in posts/5,
		// must not appear in the feed.
		require.NotEqual(t, siteURL+"/blog/5", item.GUID.Value)
		require.False(t, strings.HasPrefix(item.Content, "<h1"), "post %s still opens with its h1", item.Link)
		for _, m := range attrRE.FindAllStringSubmatch(item.Content, -1) {
			u, err := url.Parse(m[1])
			require.NoError(t, err)
			require.True(t, u.IsAbs(), "post %s has a relative link: %s", item.Link, m[1])
		}
	}
}
