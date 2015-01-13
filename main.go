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

func newBaseDot(bc *blogConf) *baseDot {
	b := &baseDot{bc.Title, bc.RootURL, &bc.L10N,
		bc.Categories, make(map[string]string), css}
	for _, m := range bc.Categories {
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
	*baseDot
	article
}

type categoryDot struct {
	*baseDot
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
	*baseDot
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
		log.Fatalln(err)
	}
}

func openForWrite(fname string) (*os.File, error) {
	return os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
}

func readCategoryConf(cat, fname string) *categoryConf {
	conf := &categoryConf{}
	decodeFile(fname, conf)
	for i := range conf.Articles {
		conf.Articles[i].Category = cat
	}
	return conf
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func newTemplate(name string, sources ...string) *template.Template {
	t := template.New(name)
	for _, source := range sources {
		template.Must(t.Parse(source))
	}
	return t
}

func main() {
	categoryTmpl := newTemplate("category", baseTemplText, categoryTemplText)
	articleTmpl := newTemplate("article", baseTemplText, articleTemplText)
	feedTmpl := newTemplate("feed", feedTemplText)

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: genblog <src dir> <dst dir>")
		os.Exit(1)
	}
	srcDir := os.Args[1]
	dstDir := os.Args[2]

	conf := &blogConf{}
	decodeFile(path.Join(srcDir, "index.toml"), conf)
	base := newBaseDot(conf)

	homepage := &homepageDot{}

	narticles := max(conf.IndexPosts, conf.FeedPosts)
	articles := make([]article, 0, narticles)

	var lastModified time.Time

	for _, cat := range conf.Categories {
		catConf := readCategoryConf(cat.Name, path.Join(srcDir, cat.Name, "index.toml"))

		catDir := path.Join(dstDir, cat.Name)
		// Create directory
		err := os.MkdirAll(catDir, 0755)
		if err != nil {
			log.Fatalln(err)
		}

		// Generate index
		file, err := openForWrite(path.Join(catDir, "index.html"))
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()
		cd := &categoryDot{base, cat.Name, catConf.Articles}
		err = categoryTmpl.Execute(file, cd)
		if err != nil {
			log.Fatalln(err)
		}

		// Generate articles
		for _, am := range catConf.Articles {
			// Read article
			fname := path.Join(srcDir, cat.Name, am.Name+".html")
			file, err := os.Open(fname)
			if err != nil {
				log.Fatalln(err)
			}
			content, err := ioutil.ReadAll(file)
			if err != nil {
				log.Fatalln(err)
			}
			fi, err := file.Stat()
			if err != nil {
				log.Fatalln(err)
			}
			modTime := fi.ModTime()
			if modTime.After(lastModified) {
				lastModified = modTime
			}

			a := article{am, cat.Name, string(content), modTime}

			// Generate article
			file, err = openForWrite(path.Join(catDir, am.Name+".html"))
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
			ad := &articleDot{base, a}
			err = articleTmpl.Execute(file, ad)
			if err != nil {
				log.Fatalln(err)
			}

			articles = insertNewArticle(articles, a, narticles)

			if homepage.Timestamp < am.Timestamp {
				homepage.articleDot = *ad
			}
		}
	}
	// Generate homepage
	homepage.Articles = articlesToMetas(articles, conf.IndexPosts)
	file, err := openForWrite(path.Join(dstDir, "index.html"))
	err = articleTmpl.Execute(file, homepage)
	if err != nil {
		log.Fatalln(err)
	}

	// Generate feed
	feedArticles := articles
	if len(articles) > conf.FeedPosts {
		feedArticles = articles[:conf.FeedPosts]
	}
	feed := feedDot{base, feedArticles, rfc3339Time(lastModified)}
	file, err = openForWrite(path.Join(dstDir, "feed.atom"))
	err = feedTmpl.Execute(file, feed)
	if err != nil {
		log.Fatalln(err)
	}
}
