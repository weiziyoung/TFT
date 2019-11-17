package search

import (
	"TFT/evaluate"
	"TFT/globals"
	"TFT/models"
	"TFT/utils"
	"container/heap"
	"fmt"
	"sort"
)

// Traverse 图遍历，
// championList, 英雄列表，固定不变。 graph 羁绊图，也是固定不变。node 为当前的结点, selected 为已选择的英雄， oldChildren是父节点的children
func Traverse(championList models.ChampionList, graph Graph, node int, selected []int, oldChildren []int) {
	selected = append(selected, node)
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
		if len(Result) == globals.Global.MaximumHeap {
			heap.Remove(&Result, 0)
		}
		return
	}
	newChildren := graph[node]
	children := append(oldChildren, newChildren...)
	sort.Ints(children)
	children = utils.DeduplicateAndFilter(children, node)
	// 注意这里children在后面会因为append方法原地操作其内部元素，因此需要拷贝一下
	copyChildren := make([]int, len(children), 50)
	copy(copyChildren, children)
	for _, child := range children {
		// selected 也是同理
		copySelected := make([]int, len(selected), lim)
		copy(copySelected, selected)
		Traverse(championList, graph, child, copySelected, copyChildren)
	}
}
