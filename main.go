//go:generate ./gen-include
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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
		cd := &categoryDot{base, cat.Name, catConf.Articles}
		executeToFile(categoryTmpl, cd, path.Join(catDir, "index.html"))

		// Generate articles
		for _, am := range catConf.Articles {
			// Read article
			fname := path.Join(srcDir, cat.Name, am.Name+".html")
			file, err := os.Open(fname)
			if err != nil {
				log.Fatalln(err)
			}
			defer file.Close()
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

			a := article{am, cat.Name, string(content), rfc3339Time(modTime)}

			// Generate article
			ad := &articleDot{base, a}
			executeToFile(articleTmpl, ad, path.Join(catDir, am.Name+".html"))

			articles = insertNewArticle(articles, a, narticles)

			if homepage.Timestamp < am.Timestamp {
				homepage.articleDot = *ad
			}
		}
	}
	// Generate homepage
	homepage.Articles = articlesToMetas(articles, conf.IndexPosts)
	executeToFile(articleTmpl, homepage, path.Join(dstDir, "index.html"))

	// Generate feed
	feedArticles := articles[:min(len(articles), conf.FeedPosts)]
	fd := feedDot{base, feedArticles, rfc3339Time(lastModified)}
	executeToFile(feedTmpl, fd, path.Join(dstDir, "feed.atom"))
}
