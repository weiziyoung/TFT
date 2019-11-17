package main

import (
	"TFT/globals"
	"TFT/search"
	"TFT/utils"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

func start(teamSize int){
	start := time.Now()
	result := search.TraitBasedGraphSearch(globals.ChampionList, int(teamSize))
	elapsed := time.Since(start)
	fmt.Println("人口:",teamSize,
		"共耗时:", elapsed.Seconds(), "s",
		"平均每秒遍历组合数:", float64(globals.Counter)/elapsed.Seconds(),
		"共搜索结点个数:", globals.Counter)

	// 排序
	sort.Slice(result, func(i, j int)bool{
		return result[i].TotalStrength > result[j].TotalStrength
	})
	utils.Translate(&result)
	jsonString, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	}
	dstFile, err := os.Create(globals.Global.OutPutPath+"champions_comb"+strconv.Itoa(teamSize)+".json")
	if err != nil {
		panic("打开文件句柄失败")
	}
	_, err = dstFile.WriteString(string(jsonString))
	if err != nil {
		panic("写入文件失败")
	}
}

func main() {
	start(7)
}