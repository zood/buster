package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func blogArchivesHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("blog-archive.html", w, map[string]interface{}{
		"title": "Blog Archive | Zood",
	})
}

func blogHomeHandler(w http.ResponseWriter, r *http.Request) {
	rsrcs := resourcesFromContext(r.Context())
	rsrcs.ExecuteTemplate("blog-home.html", w, map[string]interface{}{
		"posts": rsrcs.Posts(5, 0),
		"title": "Blog | Zood",
	})
}

func blogPostHandler(w http.ResponseWriter, r *http.Request) {
	/*
	 We're flexible in the URLs we will accept, but we should redirect
	 the user to the canonical URL if that's not what they visited.
	 Specifically, we need to make sure the slug is present
	 e.g. /blog/{{id}} should go to /blog/{{id}}/{{slug}} for readability
	*/
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["id"])
	if err != nil {
		notFoundHandler(w, r)
		return
	}

	rsrcs := resourcesFromContext(r.Context())
	post := rsrcs.PostById(postID)
	if post == nil {
		// unknown post id
		notFoundHandler(w, r)
		return
	}

	slug := vars["slug"]
	if slug == "" {
		http.Redirect(w, r, fmt.Sprintf("/blog/%d/%s", post.ID, post.Slug), http.StatusMovedPermanently)
		return
	}

	if slug != post.Slug {
		// They're probably requesting a post asset
		assetPath := rsrcs.PostAssetPath(post.ID, slug)
		http.ServeFile(w, r, assetPath)
		return
	}

	// Just serve the post html
	rsrcs.ExecuteTemplate("post-single.html", w, map[string]interface{}{
		// Convert the post body into a template.HTML so it doesn't get escaped
		"body":  template.HTML(rsrcs.PostBody(postID)),
		"title": fmt.Sprintf("%s | Zood Blog", post.Title),
	})
}