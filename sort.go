package main

import "sort"

type articleMetas []articleMeta

func (si articleMetas) Len() int           { return len(si) }
func (si articleMetas) Less(i, j int) bool { return si[i].Timestamp > si[j].Timestamp }
func (si articleMetas) Swap(i, j int)      { si[i], si[j] = si[j], si[i] }

func sortArticleMetas(ams []articleMeta) {
	sort.Sort(articleMetas(ams))
}
