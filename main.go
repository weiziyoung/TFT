package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"time"

	"TFT/globals"
	"TFT/search"
)

func factorial(num int) *big.Int {
	result := big.NewInt(int64(1))
	result = result.MulRange(1, int64(num))
	return result
}

func main() {
	var teamSize int
	teamSize = 10

	start := time.Now()
	result := search.TraitBasedGraphSearch(globals.ChampionList, int(teamSize))
	elapsed := time.Since(start)
	fmt.Println("共耗时:", elapsed.Seconds(),
		"平均每秒遍历组合数:", float64(globals.Counter)/elapsed.Seconds(),
		"共搜索结点个数:", globals.Counter)

	sort.Slice(result, func(i, j int)bool{
		return result[i].TotalStrength > result[j].TotalStrength
	})

	jsonString, err := json.MarshalIndent(&result, "", "    ")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(teamSize)
	dstFile, err := os.Create("data/champions_comb"+strconv.Itoa(teamSize)+".json")
	if err != nil {
		panic("打开文件句柄失败")
	}
	_, err = dstFile.WriteString(string(jsonString))
	if err != nil {
		panic("写入文件失败")
	}
}
