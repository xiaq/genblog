//go:generate ./gen-include
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

var (
	printDefaultTemplate = flag.Bool("print-default-template", false, "Print default template")
	printDefaultCSS      = flag.Bool("print-default-css", false, "Print default CSS")
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
	// Parse flags.
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: genblog [options] <src dir> <dst dir>")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	switch {
	case *printDefaultTemplate:
		fmt.Print(defaultTemplate)
		return
	case *printDefaultCSS:
		fmt.Print(defaultCSS)
		return
	}
	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}
	srcDir, dstDir := args[0], args[1]

	// Read blog configuration.
	conf := &blogConf{}
	decodeFile(path.Join(srcDir, "index.toml"), conf)
	genFeed, genSitemap := true, true
	if conf.RootURL == "" {
		fmt.Fprintln(os.Stderr, "No rootURL specified, generation of feed and sitemap disabled.")
		genFeed, genSitemap = false, false
	}
	template := defaultTemplate
	if conf.Template != "" {
		template = readAll(path.Join(srcDir, conf.Template))
	}
	css := defaultCSS
	if conf.CSS != "" {
		css = readAll(path.Join(srcDir, conf.CSS))
	}

	// Initialize templates. They are all initialized from the same source code,
	// plus a snippet to fix the "content" reference.
	categoryTmpl := newTemplate("category", "..", template, contentIs("category"))
	articleTmpl := newTemplate("article", "..", template, contentIs("article"))
	homepageTmpl := newTemplate("homepage", ".", template, contentIs("article"))
	feedTmpl := newTemplate("feed", ".", feedTemplText)

	// Base for the {{ . }} object used in all templates.
	base := newBaseDot(conf, css)

	// Up to conf.FeedPosts recent posts, used in the feed.
	recents := recentArticles{nil, conf.FeedPosts}
	// Last modified time of the newest post, used in the feed.
	var lastModified time.Time

	// Whether the "all" category has been requested.
	hasAllCategory := false
	// Meta of all articles, used to generate the index of the "all", if if is
	// requested.
	allArticleMetas := []articleMeta{}

	// Paths of all generated URLs, relative to the destination directory,
	// always without "index.html". Used to generate the sitemap.
	allPaths := []string{""}

	// Render a category index.
	renderCategoryIndex := func(name, prelude string, articles []articleMeta) {
		// Add category index to the sitemap, without "/index.html"
		allPaths = append(allPaths, name)
		// Create directory
		catDir := path.Join(dstDir, name)
		err := os.MkdirAll(catDir, 0755)
		if err != nil {
			log.Fatal(err)
		}

		// Generate index
		cd := &categoryDot{base, name, prelude, articles}
		executeToFile(categoryTmpl, cd, path.Join(catDir, "index.html"))
	}

	for _, cat := range conf.Categories {
		if cat.Name == "all" {
			// The "all" category has been requested. It is a pseudo-category in
			// that it doesn't need to have any associated category
			// configuration file. We cannot render the category index now
			// because we haven't seen all articles yet. Render it later.
			hasAllCategory = true
			continue
		}

		catConf := readCategoryConf(cat.Name, path.Join(srcDir, cat.Name, "index.toml"))
		sortArticleMetas(catConf.Articles)

		var prelude string
		if catConf.Prelude != "" {
			prelude = readAll(
				path.Join(srcDir, cat.Name, catConf.Prelude+".html"))
		}
		renderCategoryIndex(cat.Name, prelude, catConf.Articles)

		// Generate articles
		for _, am := range catConf.Articles {
			// Add article URL to sitemap.
			p := path.Join(cat.Name, am.Name+".html")
			allPaths = append(allPaths, p)
			// Read article
			content, fi := readAllAndStat(path.Join(srcDir, p))
			modTime := fi.ModTime()
			if modTime.After(lastModified) {
				lastModified = modTime
			}

			a := article{am, false, cat.Name, content, rfc3339Time(modTime)}

			// Generate article page.
			ad := &articleDot{base, a}
			executeToFile(articleTmpl, ad, path.Join(dstDir, p))

			allArticleMetas = append(allArticleMetas, a.articleMeta)
			recents.insert(a)
		}
	}

	// Generate "all category"
	if hasAllCategory {
		sortArticleMetas(allArticleMetas)
		renderCategoryIndex("all", "", allArticleMetas)
	}

	// Generate index page. XXX(xiaq): duplicated code with generating ordinary
	// article pages.
	content, fi := readAllAndStat(path.Join(srcDir, conf.Index.Name+".html"))
	modTime := fi.ModTime()
	a := article{conf.Index, true, "homepage", content, rfc3339Time(modTime)}
	ad := &articleDot{base, a}
	executeToFile(homepageTmpl, ad, path.Join(dstDir, "index.html"))

	// Generate feed
	if genFeed {
		feedArticles := recents.articles
		fd := feedDot{base, feedArticles, rfc3339Time(lastModified)}
		executeToFile(feedTmpl, fd, path.Join(dstDir, "feed.atom"))
	}

	if genSitemap {
		file, err := openForWrite(path.Join(dstDir, "sitemap.txt"))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		for _, p := range allPaths {
			fmt.Fprintf(file, "%s/%s\n", conf.RootURL, p)
		}
	}
}
