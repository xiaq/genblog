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

func renderCategoryIndex() {
}

func main() {
	categoryTmpl := newTemplate("category", "..", baseTemplText, contentIs("category"))
	articleTmpl := newTemplate("article", "..", baseTemplText, contentIs("article"))
	homepageTmpl := newTemplate("homepage", ".", baseTemplText, contentIs("homepage"))
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

	recentArticles := make([]article, 0, max(conf.IndexPosts, conf.FeedPosts))

	var lastModified time.Time

	allArticleMetas := []articleMeta{}

	renderCategoryIndex := func(name string, articles []articleMeta) string {
		// Create directory
		catDir := path.Join(dstDir, name)
		err := os.MkdirAll(catDir, 0755)
		if err != nil {
			log.Fatalln(err)
		}

		// Generate index
		cd := &categoryDot{base, name, articles}
		executeToFile(categoryTmpl, cd, path.Join(catDir, "index.html"))

		return catDir
	}

	hasAllCategory := false

	for _, cat := range conf.Categories {
		if cat.Name == "-" {
			hasAllCategory = true
			continue
		}
		catConf := readCategoryConf(cat.Name, path.Join(srcDir, cat.Name, "index.toml"))
		sortArticleMetas(catConf.Articles)
		catDir := renderCategoryIndex(cat.Name, catConf.Articles)

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

			allArticleMetas = append(allArticleMetas, a.articleMeta)
			recentArticles = insertNewArticle(recentArticles, a, cap(recentArticles))
		}
	}
	// Generate "all category"
	if hasAllCategory {
		sortArticleMetas(allArticleMetas)
		renderCategoryIndex("-", allArticleMetas)
	}

	// Generate homepage
	homepage := &homepageDot{
		base, articlesToDots(base, recentArticles, conf.IndexPosts),
	}
	executeToFile(homepageTmpl, homepage, path.Join(dstDir, "index.html"))

	// Generate feed
	feedArticles := recentArticles[:min(len(recentArticles), conf.FeedPosts)]
	fd := feedDot{base, feedArticles, rfc3339Time(lastModified)}
	executeToFile(feedTmpl, fd, path.Join(dstDir, "feed.atom"))
}
