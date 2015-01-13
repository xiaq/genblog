package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

// This file contains functions and types for parsing and manipulating the
// in-memory representation of the blog.

type blogConf struct {
	Title      string
	Author     string
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

type categoryMeta struct {
	Name  string
	Title string
}

type categoryConf struct {
	Articles []articleMeta
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

func readCategoryConf(cat, fname string) *categoryConf {
	conf := &categoryConf{}
	decodeFile(fname, conf)
	for i := range conf.Articles {
		conf.Articles[i].Category = cat
	}
	return conf
}
