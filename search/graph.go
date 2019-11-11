package search

import (
	"TFT/globals"
	"TFT/models"
	"TFT/utils"
	"sort"
)

type Graph map[int][]int


func GenerateGraph(championList models.ChampionList) Graph{
	graph := make(Graph)
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
				if index > no{
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
		children = utils.Deduplicate(children)
		graph[no] = children
	}
	return graph
}