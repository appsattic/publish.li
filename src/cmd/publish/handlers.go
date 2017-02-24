// --------------------------------------------------------------------------------------------------------------------
//
// This file is part of https://github.com/appsattic/publish.li
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.
//
// --------------------------------------------------------------------------------------------------------------------

package main

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/Machiel/slugify"
	"github.com/boltdb/bolt"
	"github.com/russross/blackfriday"
)

var baseUrl string
var tmpl *template.Template

func init() {
	baseUrl = os.Getenv("BASE_URL")

	tmpl1, err := template.New("").Delims("[[", "]]").ParseGlob("./templates/*.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl = tmpl1
}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	err := tmpl.ExecuteTemplate(w, templateName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func apiPut(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")
		website := r.FormValue("website")
		twitter := r.FormValue("twitter")
		github := r.FormValue("github")
		facebook := r.FormValue("facebook")
		instagram := r.FormValue("instagram")

		// validation

		// check that the title has something in it (other than whitespace)
		slug := slugify.Slugify(title)
		if slug == "" {
			sendError(w, "Provide a title")
			return
		}
		if website != "" {
			u, err := url.ParseRequestURI(website)
			if err != nil {
				sendError(w, "Invalid website URL")
				return
			}
			website = u.String()
		}
		if twitter != "" {
			if !isValidTwitterHandle(twitter) {
				sendError(w, "Invalid Twitter Handle. Only letters, numbers, and underscore allowed.")
				return
			}
		}
		if github != "" {
			if !isValidGitHubHandle(github) {
				sendError(w, "Invalid GitHub Handle. Only letters, numbers, and dash allowed.")
				return
			}
		}
		if facebook != "" {
			if !isValidFacebookHandle(facebook) {
				sendError(w, "Invalid Facebook Handle. Only letters, numbers, and dot allowed.")
				return
			}
		}
		if instagram != "" {
			if !isValidInstagramHandle(instagram) {
				sendError(w, "Invalid Instagram Handle. Only letters and numbers allowed.")
				return
			}
		}

		// fill in the other fields to save this page
		now := time.Now()
		page := Page{
			Id:        randStr(16),
			Name:      slug + "-" + randStr(8),
			Title:     title,
			Author:    author,
			Website:   website,
			Content:   content,
			Twitter:   twitter,
			Facebook:  facebook,
			GitHub:    github,
			Instagram: instagram,
			Inserted:  now,
			Updated:   now,
		}

		// and finally, create the HTML
		html := blackfriday.MarkdownCommon([]byte(page.Content))
		page.Html = template.HTML(html)

		errIns := storePutPage(db, page)
		if errIns != nil {
			http.Error(w, errIns.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Ok      bool              `json:"ok"`
			Msg     string            `json:"msg"`
			Payload map[string]string `json:"payload"`
		}{
			Ok:      true,
			Msg:     "Saved",
			Payload: make(map[string]string),
		}
		data.Payload["id"] = page.Id
		data.Payload["name"] = page.Name

		fmt.Printf("data=%#v\n", data)

		sendJson(w, data)
	}
}

func apiPost(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.FormValue("id")
		name := r.FormValue("name")
		title := r.FormValue("title")
		content := r.FormValue("content")
		author := r.FormValue("author")
		website := r.FormValue("website")
		twitter := r.FormValue("twitter")
		github := r.FormValue("github")
		facebook := r.FormValue("facebook")
		instagram := r.FormValue("instagram")

		// using the page.Name, retrieve this page then check it's Id is correct
		existPage, errGet := storeGetPage(db, name)
		if errGet != nil {
			log.Printf("Error: %v\n", errGet)
			sendError(w, "Internal Error. Please try again later.")
			return
		}

		if existPage == nil {
			sendError(w, "This page name does not exist.")
			return
		}

		// check that this page has this Id
		if existPage.Id != id {
			sendError(w, "Permission denied.")
			return
		}

		// validation

		// check that the title has something in it (other than whitespace)
		slug := slugify.Slugify(title)
		if slug == "" {
			sendError(w, "Provide a title")
			return
		}
		if website != "" {
			u, err := url.ParseRequestURI(website)
			if err != nil {
				sendError(w, "Invalid website URL")
				return
			}
			website = u.String()
		}
		if twitter != "" {
			if !isValidTwitterHandle(twitter) {
				sendError(w, "Invalid Twitter Handle. Only letters, numbers, and underscore allowed.")
				return
			}
		}
		if github != "" {
			if !isValidGitHubHandle(github) {
				sendError(w, "Invalid GitHub Handle. Only letters, numbers, and dash allowed.")
				return
			}
		}
		if facebook != "" {
			if !isValidFacebookHandle(facebook) {
				sendError(w, "Invalid Facebook Handle. Only letters, numbers, and dot allowed.")
				return
			}
		}
		if instagram != "" {
			if !isValidInstagramHandle(instagram) {
				sendError(w, "Invalid Instagram Handle. Only letters and numbers allowed.")
				return
			}
		}

		// We don't trust what is in the incoming params, but we know `existPage` is fine, so we'll just update a
		// the fields there to then re-save.
		now := time.Now()
		existPage.Title = title
		existPage.Content = content
		existPage.Author = author
		existPage.Website = website
		existPage.Twitter = twitter
		existPage.Facebook = facebook
		existPage.GitHub = github
		existPage.Instagram = instagram
		existPage.Updated = now

		// and finally, create the HTML
		html := blackfriday.MarkdownCommon([]byte(content))
		existPage.Html = template.HTML(html)

		errIns := storePutPage(db, *existPage)
		if errIns != nil {
			http.Error(w, errIns.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Ok      bool              `json:"ok"`
			Msg     string            `json:"msg"`
			Payload map[string]string `json:"payload"`
		}{
			Ok:      true,
			Msg:     "Saved",
			Payload: make(map[string]string),
		}
		data.Payload["id"] = existPage.Id
		data.Payload["name"] = existPage.Name

		fmt.Printf("data=%#v\n", data)

		sendJson(w, data)
	}
}

func apiGet(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get this Id from the incoming params
		id := r.FormValue("id")
		log.Printf("looking up id=%s\n", id)

		// retrieve this page
		page, errGet := storeGetPageUsingId(db, id)
		if errGet != nil {
			log.Printf("Error: %v\n", errGet)
			sendError(w, "Internal Error. Please try again later.")
			return
		}

		if page == nil {
			sendError(w, "This page Id does not exist.")
			return
		}

		data := struct {
			Ok      bool   `json:"ok"`
			Msg     string `json:"msg"`
			Payload *Page  `json:"payload"`
		}{
			Ok:      true,
			Msg:     "Saved",
			Payload: page,
		}

		sendJson(w, data)
	}
}

func apiHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	insPost := apiPut(db)
	savePost := apiPost(db)
	getPost := apiGet(db)

	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path != "/api" {
			http.NotFoundHandler().ServeHTTP(w, r)
			return
		}

		if r.Method == "PUT" {
			insPost(w, r)
			return
		}

		if r.Method == "POST" {
			savePost(w, r)
			return
		}

		if r.Method == "GET" {
			getPost(w, r)
			return
		}
	}
}

func servePage(w http.ResponseWriter, r *http.Request, db *bolt.DB) {
	// everything else
	name := r.URL.Path[1:]
	log.Printf("Page=%q\n", html.EscapeString(name))

	page, errPage := storeGetPage(db, name)
	if errPage != nil {
		http.Error(w, errPage.Error(), http.StatusInternalServerError)
		return
	}

	if page == nil {
		log.Printf("Not Found : %s\n", name)
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	// serve the page
	data := struct {
		Layout string
		Page   *Page
	}{
		"page",
		page,
	}
	render(w, "page.html", data)
}

func sitemap(w http.ResponseWriter, r *http.Request, baseUrl string, db *bolt.DB) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, baseUrl+"/\n")

	// loop through all the pages
	err := storeIteratePages(db, func(k, v []byte) error {
		fmt.Fprintf(w, "%s/%s\n", baseUrl, string(k))
		return nil
	})

	if err != nil {
		log.Printf("Error: %v", err)
	}
}

func homeHandler(db *bolt.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		log.Printf("path=%s\n", path)

		// simple routing instead of using a complicated router (until we need one)
		if path == "/" {
			// serve the page
			data := struct {
				Layout string
			}{
				"home",
			}
			render(w, "home.html", data)

		} else if path == "/favicon.ico" {
			http.ServeFile(w, r, "static/favicon.ico")

		} else if path == "/robots.txt" {
			http.ServeFile(w, r, "static/robots.txt")

		} else if path == "/sitemap.txt" {
			sitemap(w, r, baseUrl, db)

		} else {
			servePage(w, r, db)
		}
	}
}
