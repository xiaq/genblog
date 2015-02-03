package main

import "sort"

type sortInterface []articleMeta

func (si sortInterface) Len() int {
	return len(si)
}

func (si sortInterface) Less(i, j int) bool {
	return si[i].Timestamp > si[j].Timestamp
}

func (si sortInterface) Swap(i, j int) {
	si[i], si[j] = si[j], si[i]
}

func sortArticleMetas(ams []articleMeta) {
	sort.Sort(sortInterface(ams))
}
