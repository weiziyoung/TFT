package search

import (
	"TFT/globals"
	"TFT/models"
	"container/heap"
)

var lim int
var Result = make(models.ComboMetricHeap, 0, globals.Global.MaximumHeap)

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
	startPoints := getSlice(0, len(championList)-teamSize + 1)
	for _,startNode := range startPoints{
		Traverse(championList, graph, startNode, make([]int, 0, teamSize), make([]int, 0, 57))
	}

	return Result
}

