package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
	"time"
)

// This file contains functions and types for rendering the blog.

type baseDot struct {
	BlogTitle  string
	Author     string
	RootURL    string
	L10N       *l10nConf
	Categories []categoryMeta

	CategoryMap map[string]string
	CSS         string
}

func newBaseDot(bc *blogConf) *baseDot {
	b := &baseDot{bc.Title, bc.Author, bc.RootURL, &bc.L10N,
		bc.Categories, make(map[string]string), css}
	for _, m := range bc.Categories {
		b.CategoryMap[m.Name] = m.Title
	}
	return b
}

type articleDot struct {
	*baseDot
	article
}

type categoryDot struct {
	*baseDot
	Category string
	Articles []articleMeta
}

type homepageArticleDot articleDot

type homepageDot struct {
	*baseDot
	Articles []homepageArticleDot
}

type feedDot struct {
	*baseDot
	Articles     []article
	LastModified rfc3339Time
}

type rfc3339Time time.Time

func (t rfc3339Time) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func contentIs(what string) string {
	return fmt.Sprintf(
		`{{ define "content" }} {{ template "%s-content" . }} {{ end }}`,
		what)
}

func newTemplate(name, root string, sources ...string) *template.Template {
	t := template.New(name).Funcs(template.FuncMap(map[string]interface{}{
		"is": func(s string) bool { return s == name },
		"homepageURL": func() string {
			return root + "/index.html"
		},
		"categoryURL": func(cat string) string {
			return root + "/" + cat + "/index.html"
		},
		"articleURL": func(cat, article string) string {
			return root + "/" + cat + "/" + article + ".html"
		},
	}))
	for _, source := range sources {
		template.Must(t.Parse(source))
	}
	return t
}

func openForWrite(fname string) (*os.File, error) {
	return os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func executeToFile(t *template.Template, data interface{}, fname string) {
	file, err := openForWrite(fname)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	err = t.Execute(file, data)
	if err != nil {
		log.Fatalln(err)
	}
}
