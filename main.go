//go:generate ./gen-include
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/BurntSushi/toml"
)

type blogConf struct {
	Title      string
	Categories []categoryMeta
	Recent int
	RecentTitle string
	CategoryListTitle string
}

type categoryConf struct {
	Articles []articleMeta
}

type baseDot struct {
	BlogTitle  string
	RecentTitle string
	CategoryListTitle string
	Categories []categoryMeta

	CategoryMap map[string]string
	CSS         string
}

type categoryMeta struct {
	Name  string
	Title string
}

func newBaseDot(t, rt, clt string, cs []categoryMeta) *baseDot {
	b := &baseDot{t, rt, clt, cs, make(map[string]string), css}
	for _, m := range cs {
		b.CategoryMap[m.Name] = m.Title
	}
	return b
}

func (b *baseDot) IsCategory() bool {
	return false
}

func (b *baseDot) IsHomepage() bool {
	return false
}

func (b *baseDot) Root() string {
	return ".."
}

type articleMeta struct {
	Name      string
	Title     string
	Category  string
	Timestamp string
}

type articleDot struct {
	baseDot
	articleMeta
	Category string
	Content  string
}

type categoryDot struct {
	baseDot
	Category string
	Articles []articleMeta
}

func (c *categoryDot) IsCategory() bool {
	return true
}

const recentArticles = 5

type homepageDot struct {
	articleDot
	Articles []articleMeta
}

func (h *homepageDot) insertNewArticle(am articleMeta, nmax int) {
	var i int
	for i = len(h.Articles); i > 0; i-- {
		if h.Articles[i-1].Timestamp < am.Timestamp {
			break
		}
	}
	if i == len(h.Articles) {
		if i < nmax {
			h.Articles = append(h.Articles, am)
		}
		return
	}

	if len(h.Articles) < nmax {
		h.Articles = append(h.Articles, articleMeta{})
	}

	copy(h.Articles[i+1:], h.Articles[i:])
	h.Articles[i] = am
}

func (h *homepageDot) IsHomepage() bool {
	return true
}

func (h *homepageDot) Root() string {
	return "."
}

func decodeFile(fname string, v interface{}) {
	_, err := toml.DecodeFile(fname, v)
	if err != nil {
		log.Fatalf("when reading %v: %v", fname, err)
	}
}

func openForWrite(fname string) (*os.File, error) {
	return os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func readCategoryConf(cat, fname string) *categoryConf {
	cf := &categoryConf{}
	decodeFile(fname, &cf)
	for i := range cf.Articles {
		cf.Articles[i].Category = cat
	}
	return cf
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: genblog <src dir> <dst dir>")
		os.Exit(1)
	}
	srcDir := os.Args[1]
	dstDir := os.Args[2]

	bf := &blogConf{}
	decodeFile(path.Join(srcDir, "index.toml"), &bf)
	bd := newBaseDot(bf.Title, bf.RecentTitle, bf.CategoryListTitle, bf.Categories)

	categoryTmpl := template.New("category")
	template.Must(categoryTmpl.Parse(baseTemplateText))
	template.Must(categoryTmpl.Parse(categoryTemplateText))

	articleTmpl := template.New("article")
	template.Must(articleTmpl.Parse(baseTemplateText))
	template.Must(articleTmpl.Parse(articleTemplateText))

	homepage := &homepageDot{}

	for _, cm := range bf.Categories {
		cf := readCategoryConf(cm.Name, path.Join(srcDir, cm.Name, "index.toml"))

		cDir := path.Join(dstDir, cm.Name)
		// Create directory
		err := os.MkdirAll(cDir, 0755)
		if err != nil {
			log.Fatalf("failed to create output dir for %v: %v\n", cm.Name, err)
		}

		// Generate index
		file, err := openForWrite(path.Join(cDir, "index.html"))
		if err != nil {
			log.Fatalf("failed to open index of %v: %v\n", cm.Name, err)
		}
		defer file.Close()
		cd := &categoryDot{*bd, cm.Name, cf.Articles}
		err = categoryTmpl.Execute(file, cd)
		if err != nil {
			log.Fatalf("failed to render index of %v: %v\n", cm.Name, err)
		}

		// Generate articles
		for _, am := range cf.Articles {
			// Read article
			content, err := ioutil.ReadFile(path.Join(srcDir, cm.Name, am.Name+".html"))
			if err != nil {
				log.Fatalf("failed to read article %v/%v: %v\n", cm.Name, am.Name, err)
			}

			// Generate article
			file, err := openForWrite(path.Join(cDir, am.Name+".html"))
			if err != nil {
				log.Fatalf("failed to open article %v/%v: %v\n", cm.Name, am.Name, err)
			}
			defer file.Close()
			ad := &articleDot{*bd, am, cm.Name, string(content)}
			err = articleTmpl.Execute(file, ad)
			if err != nil {
				log.Fatalf("failed to render article %v/%v: %v\n", cm.Name, am.Name, err)
			}

			homepage.insertNewArticle(ad.articleMeta, bf.Recent)

			if homepage.Timestamp < ad.Timestamp {
				homepage.articleDot = *ad
			}
		}
	}
	// Generate homepage
	file, err := openForWrite(path.Join(dstDir, "index.html"))
	err = articleTmpl.Execute(file, homepage)
	if err != nil {
		log.Fatalf("failed to render homepage: %v\n", err)
	}
}
