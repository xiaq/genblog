package main

import "sort"

type articleMetas []articleMeta

func (a articleMetas) Len() int           { return len(a) }
func (a articleMetas) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }
func (a articleMetas) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortArticleMetas(a []articleMeta) { sort.Sort(articleMetas(a)) }

type articles []article

func (a articles) Len() int           { return len(a) }
func (a articles) Less(i, j int) bool { return a[i].Timestamp > a[j].Timestamp }
func (a articles) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func sortArticles(a articles) { sort.Sort(articles(a)) }
