package search

import (
	"TFT/models"
	"container/heap"
)

var startPoint int
var lim int
var Result = make(models.ComboMetricHeap, 0, 1024)

// getSlice 与python range(x,y)相同效果
func getSlice(start int, end int) []int {
	rangeList := make([]int, end-start)
	for i := start; i < end; i++ {
		rangeList[i-start] = i
	}
	return rangeList
}


// TraitBasedGraphSearch 基于羁绊图的图搜索
func TraitBasedGraphSearch(championList models.ChampionList, teamSize int) models.ComboMetricHeap {
	graph := GenerateGraph(championList)
	lim = teamSize

	heap.Init(&Result)
	startPoint = len(championList)
	graph[startPoint] = getSlice(0, len(championList)-teamSize+1)
	Traverse(championList, graph, startPoint, make([]int, 0, teamSize))

	return Result
}
