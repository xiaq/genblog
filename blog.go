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
	IndexPosts int
	FeedPosts  int
	RootURL    string
	L10N       l10nConf
}

// l10Conf represents the L10N section of the blog configuration.
type l10nConf struct {
	AllArticles string
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
	Category  string
	Timestamp string
}

// article represents an article, including all information that is needed to
// render it.
type article struct {
	articleMeta
	Category     string
	Content      string
	LastModified rfc3339Time
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

func articlesToDots(b *baseDot, as []article, n int) []homepageArticleDot {
	ads := make([]homepageArticleDot, min(n, len(as)))
	for i, a := range as {
		if i == n {
			break
		}
		ads[i] = homepageArticleDot(articleDot{b, a})
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
	for i := range conf.Articles {
		conf.Articles[i].Category = cat
	}
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
