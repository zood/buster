// Package feed renders the blog's RSS 2.0 feed.
package feed

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
	"zood.xyz/buster/resources"
)

// Path is the URL path at which the feed is served.
const Path = "/blog/rss.xml"

// A Channel holds the metadata that applies to the whole feed.
type Channel struct {
	// SiteURL is the scheme and host prefixed to every URL in the feed.
	// It must not end in a slash.
	SiteURL     string
	Title       string
	Description string
	Language    string
}

// rss is the root element of the feed.
//
// encoding/xml does not support marshaling namespace prefixes, so the xmlns
// attributes here and the prefixed element names in rssChannel and item are
// written out as literal names.
type rss struct {
	XMLName   xml.Name   `xml:"rss"`
	Version   string     `xml:"version,attr"`
	AtomNS    string     `xml:"xmlns:atom,attr"`
	ContentNS string     `xml:"xmlns:content,attr"`
	Channel   rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title         string   `xml:"title"`
	Link          string   `xml:"link"`
	Description   string   `xml:"description"`
	Language      string   `xml:"language"`
	LastBuildDate string   `xml:"lastBuildDate,omitempty"`
	AtomLink      atomLink `xml:"atom:link"`
	Items         []item   `xml:"item"`
}

type atomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        guid   `xml:"guid"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	Content     cdata  `xml:"content:encoded"`
}

type guid struct {
	IsPermaLink bool   `xml:"isPermaLink,attr"`
	Value       string `xml:",chardata"`
}

// cdata marshals its value as a CDATA section, so that post HTML stays
// legible in the feed. Marshaled as character data, it would have every
// newline and markup character escaped.
type cdata struct {
	Value string `xml:",cdata"`
}

// RSS returns an RSS 2.0 document with an item for each of posts, in order,
// so posts should be sorted newest first.
//
// Each item contains its post's full body with two changes: a leading h1 is
// removed, because feed readers display the item title, and relative URLs are
// resolved against the post's URL, because feed readers display the body
// apart from the page it was written for.
func RSS(ch Channel, posts []resources.PostWithBody) ([]byte, error) {
	feed := rss{
		Version:   "2.0",
		AtomNS:    "http://www.w3.org/2005/Atom",
		ContentNS: "http://purl.org/rss/1.0/modules/content/",
		Channel: rssChannel{
			Title:       ch.Title,
			Link:        ch.SiteURL + "/blog",
			Description: ch.Description,
			Language:    ch.Language,
			AtomLink: atomLink{
				Href: ch.SiteURL + Path,
				Rel:  "self",
				Type: "application/rss+xml",
			},
		},
	}

	var newest int64
	for _, p := range posts {
		postURL := fmt.Sprintf("%s/blog/%d/%s", ch.SiteURL, p.ID, p.Slug)
		base, err := url.Parse(postURL)
		if err != nil {
			return nil, fmt.Errorf("parsing url of post %d: %v", p.ID, err)
		}
		body, err := absolutizeURLs(stripLeadingH1(string(p.Body)), base)
		if err != nil {
			return nil, fmt.Errorf("rewriting urls in post %d: %v", p.ID, err)
		}

		feed.Channel.Items = append(feed.Channel.Items, item{
			Title: p.Title,
			Link:  postURL,
			// The GUID omits the slug so that renaming a post doesn't change
			// its GUID, which feed readers would display as a new post.
			GUID:        guid{Value: fmt.Sprintf("%s/blog/%d", ch.SiteURL, p.ID)},
			PubDate:     rfc822(p.Date),
			Description: p.Intro,
			Content:     cdata{Value: body},
		})
		newest = max(newest, p.Date)
	}
	// The newest post's date is used instead of the current time so that the
	// output is identical across calls until the posts change.
	if newest != 0 {
		feed.Channel.LastBuildDate = rfc822(newest)
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc := xml.NewEncoder(&buf)
	enc.Indent("", "  ")
	if err := enc.Encode(feed); err != nil {
		return nil, err
	}
	buf.WriteString("\n")

	return buf.Bytes(), nil
}

func rfc822(unix int64) string {
	return time.Unix(unix, 0).UTC().Format(time.RFC1123Z)
}

// stripLeadingH1 returns body without its leading h1 element, any whitespace
// and comments before it, and the whitespace after it. If body doesn't begin
// with a complete h1 element, stripLeadingH1 returns it unchanged.
func stripLeadingH1(body string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(body))
	offset := 0
	inH1 := false
	for {
		typ := tokenizer.Next()
		offset += len(tokenizer.Raw())
		switch typ {
		case html.ErrorToken:
			return body
		case html.CommentToken:
		case html.TextToken:
			if !inH1 && strings.TrimSpace(string(tokenizer.Raw())) != "" {
				return body
			}
		case html.StartTagToken:
			if inH1 {
				continue
			}
			if name, _ := tokenizer.TagName(); string(name) != "h1" {
				return body
			}
			inH1 = true
		case html.EndTagToken:
			if !inH1 {
				return body
			}
			if name, _ := tokenizer.TagName(); string(name) == "h1" {
				return strings.TrimLeft(body[offset:], " \t\r\n")
			}
		default:
			if !inH1 {
				return body
			}
		}
	}
}

// absolutizeURLs returns body with every relative href and src attribute
// value resolved against base. A tag containing such an attribute is
// re-serialized, which can change its formatting; the rest of body is
// returned unchanged.
func absolutizeURLs(body string, base *url.URL) (string, error) {
	var out strings.Builder
	tokenizer := html.NewTokenizer(strings.NewReader(body))
	for {
		typ := tokenizer.Next()
		switch typ {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return out.String(), nil
			}
			return "", tokenizer.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			raw := string(tokenizer.Raw())
			tok := tokenizer.Token()
			changed := false
			for i, attr := range tok.Attr {
				if attr.Namespace != "" || (attr.Key != "href" && attr.Key != "src") {
					continue
				}
				ref, err := url.Parse(strings.TrimSpace(attr.Val))
				if err != nil || ref.IsAbs() {
					continue
				}
				tok.Attr[i].Val = base.ResolveReference(ref).String()
				changed = true
			}
			if changed {
				out.WriteString(tok.String())
			} else {
				out.WriteString(raw)
			}
		default:
			out.Write(tokenizer.Raw())
		}
	}
}
