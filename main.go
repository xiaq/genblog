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
	categoryTmpl := newTemplate("category", "..", baseTemplText, contentIs("category"))
	articleTmpl := newTemplate("article", "..", baseTemplText, contentIs("article"))
	homepageTmpl := newTemplate("homepage", ".", baseTemplText, contentIs("article"))
	feedTmpl := newTemplate("feed", ".", feedTemplText)

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "Usage: genblog <src dir> <dst dir>")
		os.Exit(1)
	}
	srcDir := os.Args[1]
	dstDir := os.Args[2]

	conf := &blogConf{}
	decodeFile(path.Join(srcDir, "index.toml"), conf)
	base := newBaseDot(conf)

	recents := recentArticles{nil, conf.FeedPosts}

	var lastModified time.Time

	allArticleMetas := []articleMeta{}

	renderCategoryIndex := func(name, prelude string, articles []articleMeta) string {
		// Create directory
		catDir := path.Join(dstDir, name)
		err := os.MkdirAll(catDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		// Generate index
		cd := &categoryDot{base, name, prelude, articles}
		executeToFile(categoryTmpl, cd, path.Join(catDir, "index.html"))

		return catDir
	}

	hasAllCategory := false

	for _, cat := range conf.Categories {
		if cat.Name == "all" {
			hasAllCategory = true
			continue
		}
		catConf := readCategoryConf(cat.Name, path.Join(srcDir, cat.Name, "index.toml"))
		sortArticleMetas(catConf.Articles)

		var prelude string
		if catConf.Prelude != "" {
			bytes, _ := readAllAndStat(path.Join(srcDir, cat.Name, catConf.Prelude+".html"))
			prelude = string(bytes)
		}
		catDir := renderCategoryIndex(cat.Name, prelude, catConf.Articles)

		// Generate articles
		for _, am := range catConf.Articles {
			// Read article
			content, fi := readAllAndStat(path.Join(srcDir, cat.Name, am.Name+".html"))
			modTime := fi.ModTime()
			if modTime.After(lastModified) {
				lastModified = modTime
			}

			a := article{am, false, cat.Name, string(content), rfc3339Time(modTime)}

			// Generate article page.
			ad := &articleDot{base, a}
			executeToFile(articleTmpl, ad, path.Join(catDir, am.Name+".html"))

			allArticleMetas = append(allArticleMetas, a.articleMeta)
			recents.insert(a)
		}
	}
	// Generate "all category"
	if hasAllCategory {
		sortArticleMetas(allArticleMetas)
		renderCategoryIndex("all", "", allArticleMetas)
	}

	// Generate index page. XXX(xiaq): duplicated code.
	content, fi := readAllAndStat(path.Join(srcDir, conf.Index.Name+".html"))
	modTime := fi.ModTime()
	a := article{conf.Index, true, "homepage", string(content), rfc3339Time(modTime)}
	ad := &articleDot{base, a}
	executeToFile(homepageTmpl, ad, path.Join(dstDir, "index.html"))

	// Generate feed
	feedArticles := recents.articles
	fd := feedDot{base, feedArticles, rfc3339Time(lastModified)}
	executeToFile(feedTmpl, fd, path.Join(dstDir, "feed.atom"))
}
