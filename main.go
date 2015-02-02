//go:generate ./gen-include
package main

import (
	"fmt"
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
	categoryTmpl := newTemplate("category", baseTemplText, contentIs("category"))
	articleTmpl := newTemplate("article", baseTemplText, contentIs("article"))
	homepageTmpl := newTemplate("homepage", baseTemplText, contentIs("homepage"))
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
			content, fi := readAllAndStat(path.Join(srcDir, cat.Name, am.Name+".html"))
			modTime := fi.ModTime()
			if modTime.After(lastModified) {
				lastModified = modTime
			}

			a := article{am, cat.Name, string(content), rfc3339Time(modTime)}

			// Generate article
			ad := &articleDot{base, a}
			executeToFile(articleTmpl, ad, path.Join(catDir, am.Name+".html"))

			articles = insertNewArticle(articles, a, narticles)
		}
	}
	// Generate homepage
	homepage.baseDot = base
	homepage.Articles = articlesToDots(homepage.baseDot, articles, conf.IndexPosts)
	executeToFile(homepageTmpl, homepage, path.Join(dstDir, "index.html"))

	// Generate feed
	feedArticles := articles[:min(len(articles), conf.FeedPosts)]
	fd := feedDot{base, feedArticles, rfc3339Time(lastModified)}
	executeToFile(feedTmpl, fd, path.Join(dstDir, "feed.atom"))
}
