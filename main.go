//go:generate ./gen-include
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var (
	printDefaultTemplate = flag.Bool("print-default-template", false, "Print default template")
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
	flag.Parse()
	args := flag.Args()

	if *printDefaultTemplate {
		fmt.Print(baseTemplText)
		return
	}

	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: genblog <src dir> <dst dir>")
		os.Exit(1)
	}
	srcDir, dstDir := args[0], args[1]

	conf := &blogConf{}
	decodeFile(path.Join(srcDir, "index.toml"), conf)

	template := baseTemplText
	if conf.Template != "" {
		bytes, err := ioutil.ReadFile(path.Join(srcDir, conf.Template))
		if err != nil {
			log.Fatal(err)
		}
		template = string(bytes)
	}

	categoryTmpl := newTemplate("category", "..", template, contentIs("category"))
	articleTmpl := newTemplate("article", "..", template, contentIs("article"))
	homepageTmpl := newTemplate("homepage", ".", template, contentIs("article"))
	feedTmpl := newTemplate("feed", ".", feedTemplText)

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
