package search

import (
	"TFT/evaluate"
	"TFT/globals"
	"TFT/models"
	"container/heap"
	"fmt"
	"sort"
)

type Graph map[int][]int

var startPoint int
var lim int
var graph Graph
var Result = make(models.ComboMetricHeap, 0, 128)

func getSlice(start int, end int) []int {
	rangeList := make([]int, end-start)
	for i := start; i < end; i++ {
		rangeList[i-start] = i
	}
	return rangeList
}

func GraphSearch(championList models.ChampionList, teamSize int) models.ComboMetricHeap {
	// 初始化
	graph = make(Graph)
	lim = teamSize
	for no, _ := range championList {
		graph[no] = getSlice(no+1, len(championList))
	}
	heap.Init(&Result)
	startPoint = len(championList)
	graph[startPoint] = getSlice(0, len(championList)-teamSize+1)
	traverse(championList, startPoint, make([]int, 0, teamSize))
	return Result
}

// 基于羁绊图的图搜索
func TraitBasedGraphSearch(championList models.ChampionList, teamSize int) models.ComboMetricHeap {
	graph = make(Graph)
	lim = teamSize
	// 建立键为英雄名字，值为它在列表中所在位置的表
	positionMap := make(map[string]int)
	for index, champion := range championList {
		positionMap[champion.Name] = index
	}
	for no, champion := range championList {
		// children 排序
		children := make([]int, 0, 30)
		// 加入相同职业的英雄
		classes := champion.Class
		for _, class := range classes {
			sameClassChampions := globals.TraitDict[class].Champions
			for _, champion := range sameClassChampions {
				index := positionMap[champion]
				if index > no {
					children = append(children, index)
				}
			}
		}
		// 加入相同种族的英雄
		origins := champion.Origin
		for _, origin := range origins {
			sameOriginChampions := globals.TraitDict[origin].Champions
			for _, champion := range sameOriginChampions {
				index := positionMap[champion]
				if index > no {
					children = append(children, index)
				}
			}
		}
		// 加入1羁绊的英雄
		for _, championName := range globals.OneTraitChampionNameList {
			index := positionMap[championName]
			if index > no {
				children = append(children, index)
			}
		}
		// 对index从小到大排序
		sort.Ints(children)
		graph[no] = children
	}
	heap.Init(&Result)
	startPoint = len(championList)
	graph[startPoint] = getSlice(0, len(championList)-teamSize+1)
	traverse(championList, startPoint, make([]int, 0, teamSize))
	return Result
}

// 遍历
func traverse(championList models.ChampionList, node int, selected []int) {
	if node != startPoint {
		selected = append(selected, node)
	}
	if len(selected) == lim {
		combo := make(models.ChampionList, lim)
		for index, no := range selected {
			unit := championList[no]
			combo[index] = unit
		}
		globals.Counter += 1
		if globals.Counter%100000 == 0 {
			fmt.Println(globals.Counter)
		}
		metric := evaluate.Evaluate(combo)
		if metric.TraitNum > 6{
			heap.Push(&Result, metric)
		}
		// 超过最大就pop
		if len(Result) == cap(Result) {
			Result.Pop()
		}
		return
	}
	for _, child := range graph[node] {
		copySelected := make([]int, len(selected), lim)
		copy(copySelected, selected)
		traverse(championList, child, copySelected)
	}
}
