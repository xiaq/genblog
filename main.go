//go:generate ./gen-include
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
)

type blogConf struct {
	Title      string
	Categories []categoryMeta
	IndexPosts int
	FeedPosts  int
	RootURL    string
	L10N       l10nConf
}

type l10nConf struct {
	RecentTitle       string
	CategoryListTitle string
}

type categoryConf struct {
	Articles []articleMeta
}

type baseDot struct {
	BlogTitle  string
	RootURL    string
	L10N       *l10nConf
	Categories []categoryMeta

	CategoryMap map[string]string
	CSS         string
}

type categoryMeta struct {
	Name  string
	Title string
}

func newBaseDot(t, u string, l *l10nConf, cs []categoryMeta) *baseDot {
	b := &baseDot{t, u, l, cs, make(map[string]string), css}
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

type article struct {
	articleMeta
	Category     string
	Content      string
	LastModified time.Time
}

type articleDot struct {
	baseDot
	article
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

func (h *homepageDot) IsHomepage() bool {
	return true
}

func (h *homepageDot) Root() string {
	return "."
}

type feedDot struct {
	baseDot
	Articles     []article
	LastModified rfc3339Time
}

type rfc3339Time time.Time

func (t rfc3339Time) String() string {
	return time.Time(t).Format(time.RFC3339)
}

func insertNewArticle(as []article, a article, nmax int) []article {
	var i int
	for i = len(as); i > 0; i-- {
		if as[i-1].Timestamp > a.Timestamp {
			break
		}
	}
	if i == len(as) {
		if i < nmax {
			as = append(as, a)
		}
		return as
	}
	if len(as) < nmax {
		as = append(as, article{})
	}
	copy(as[i+1:], as[i:])
	as[i] = a
	return as
}

func articlesToMetas(as []article, n int) []articleMeta {
	ams := make([]articleMeta, len(as))
	for i, a := range as {
		if i == n {
			break
		}
		ams[i] = a.articleMeta
	}
	return ams
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
	decodeFile(fname, cf)
	for i := range cf.Articles {
		cf.Articles[i].Category = cat
	}
	return cf
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: genblog <src dir> <dst dir>")
		os.Exit(1)
	}
	srcDir := os.Args[1]
	dstDir := os.Args[2]

	bf := &blogConf{}
	decodeFile(path.Join(srcDir, "index.toml"), bf)
	bd := newBaseDot(bf.Title, bf.RootURL, &bf.L10N, bf.Categories)

	categoryTmpl := template.New("category")
	template.Must(categoryTmpl.Parse(baseTemplateText))
	template.Must(categoryTmpl.Parse(categoryTemplateText))

	articleTmpl := template.New("article")
	template.Must(articleTmpl.Parse(baseTemplateText))
	template.Must(articleTmpl.Parse(articleTemplateText))

	feedTmpl := template.New("feed")
	// template.Must(feedTmpl.Parse(baseTemplateText))
	template.Must(feedTmpl.Parse(feedTemplateText))

	homepage := &homepageDot{}

	narticles := max(bf.IndexPosts, bf.FeedPosts)
	articles := make([]article, 0, narticles)

	var lastModified time.Time

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
			fname := path.Join(srcDir, cm.Name, am.Name+".html")
			file, err := os.Open(fname)
			if err != nil {
				log.Fatalf("failed to read article %v/%v: %v\n", cm.Name, am.Name, err)
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatalf("failed to read article %v/%v: %v\n", cm.Name, am.Name, err)
			}
			fi, err := file.Stat()
			if err != nil {
				log.Fatalf("failed to stat article %v/%v: %v\n", cm.Name, am.Name, err)
			}
			modTime := fi.ModTime()
			if modTime.After(lastModified) {
				lastModified = modTime
			}

			a := article{am, cm.Name, string(content), modTime}

			// Generate article
			file, err = openForWrite(path.Join(cDir, am.Name+".html"))
			if err != nil {
				log.Fatalf("failed to open article %v/%v: %v\n", cm.Name, am.Name, err)
			}
			defer file.Close()
			ad := &articleDot{*bd, a}
			err = articleTmpl.Execute(file, ad)
			if err != nil {
				log.Fatalf("failed to render article %v/%v: %v\n", cm.Name, am.Name, err)
			}

			articles = insertNewArticle(articles, a, narticles)

			if homepage.Timestamp < am.Timestamp {
				homepage.articleDot = *ad
			}
		}
	}
	// Generate homepage
	homepage.Articles = articlesToMetas(articles, bf.IndexPosts)
	file, err := openForWrite(path.Join(dstDir, "index.html"))
	err = articleTmpl.Execute(file, homepage)
	if err != nil {
		log.Fatalf("failed to render homepage: %v\n", err)
	}

	// Generate feed
	feedArticles := articles
	if len(articles) > bf.FeedPosts {
		feedArticles = articles[:bf.FeedPosts]
	}
	feed := feedDot{*bd, feedArticles, rfc3339Time(lastModified)}
	file, err = openForWrite(path.Join(dstDir, "feed.atom"))
	err = feedTmpl.Execute(file, feed)
	if err != nil {
		log.Fatalf("failed to render feed: %v\n", err)
	}
}
