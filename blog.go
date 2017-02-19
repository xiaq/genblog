package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// This file contains functions and types for parsing and manipulating the
// in-memory representation of the blog.

// blogConf represents the global blog configuration.
type blogConf struct {
	Title      string
	Author     string
	Categories []categoryMeta
	Index      articleMeta
	IndexPosts int
	FeedPosts  int
	RootURL    string
}

// categoryMeta represents the metadata of a cateogory, found in the global
// blog configuration.
type categoryMeta struct {
	Name  string
	Title string
}

// categoryConf represents the configuration of a category. Note that the
// metadata is found in the global blog configuration and not duplicated here.
type categoryConf struct {
	Articles []articleMeta
}

// articleMeta represents the metadata of an article, found in a category
// configuration.
type articleMeta struct {
	Name      string
	Title     string
	Timestamp string
}

// article represents an article, including all information that is needed to
// render it.
type article struct {
	articleMeta
	IsHomepage   bool
	Category     string
	Content      string
	LastModified rfc3339Time
}

type recentArticles struct {
	articles []article
	max      int
}

func (ra *recentArticles) insert(a article) {
	// Find a place to insert.
	var i int
	for i = len(ra.articles); i > 0; i-- {
		if ra.articles[i-1].Timestamp > a.Timestamp {
			break
		}
	}
	// If we are at the end, insert only if we haven't reached the maximum
	// number of articles.
	if i == len(ra.articles) {
		if i < ra.max {
			ra.articles = append(ra.articles, a)
		}
		return
	}
	// If not, make space and insert.
	if len(ra.articles) < ra.max {
		ra.articles = append(ra.articles, article{})
	}
	copy(ra.articles[i+1:], ra.articles[i:])
	ra.articles[i] = a
}

func articlesToDots(b *baseDot, as []article) []articleDot {
	ads := make([]articleDot, len(as))
	for i, a := range as {
		ads[i] = articleDot(articleDot{b, a})
	}
	return ads
}

// decodeFile decodes the named file in TOML into a pointer.
func decodeFile(fname string, v interface{}) {
	_, err := toml.DecodeFile(fname, v)
	if err != nil {
		log.Fatalln(err)
	}
}

// readCatetoryConf reads a category configuration file.
func readCategoryConf(cat, fname string) *categoryConf {
	conf := &categoryConf{}
	decodeFile(fname, conf)
	return conf
}

// readAllAndStat retrieves all content of the named file and its stat.
func readAllAndStat(fname string) ([]byte, os.FileInfo) {
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
	return content, fi
}
