package search

import (
	"TFT/evaluate"
	"TFT/globals"
	"TFT/models"
	"container/heap"
	"fmt"
)

// 遍历
func Traverse(championList models.ChampionList,graph Graph, node int, selected []int) {
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
		heap.Push(&Result, metric)

		// 超过最大就pop
		if len(Result) == 1024{
			heap.Remove(&Result, 0)
		}
		return
	}
	for _, child := range graph[node] {
		copySelected := make([]int, len(selected), lim)
		copy(copySelected, selected)
		Traverse(championList, graph,  child, copySelected)
	}
}
