package models

type ComboMetric struct {
	// 英雄组合
	Combo []string `json:"combo"`
	// 队伍总羁绊数量 = 每个羁绊 * 羁绊等级
	TraitNum int `json:"trait_num"`
	// 具体羁绊
	TraitDetail map[string]int `json:"trait_detail"`
	// 总英雄收益羁绊数量 = sigma{羁绊} 羁绊范围 * 羁绊等级
	TotalTraitNum int `json:"total_trait_num"`
	// 当前阵容羁绊强度 = sigma{羁绊} 羁绊范围 *  羁绊强度
	TotalTraitStrength float32 `json:"total_trait_strength"`
	// 当前阵容强度 = sigma{英雄} 英雄强度 * 羁绊强度
	TotalStrength float64 `json:"total_strength"`
}

// 定义一个最大堆,注意golang标准库只有最小堆的实现，所以要重载一下Less方法
type ComboMetricHeap []ComboMetric

func (h ComboMetricHeap) Len() int {
	return len(h)
}

func (h ComboMetricHeap) Less(i,j int) bool {
	return h[i].TotalStrength < h[j].TotalStrength
}

func (h ComboMetricHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *ComboMetricHeap) Push(x interface{}) {
	*h = append(*h, x.(ComboMetric))
}

func (h *ComboMetricHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}



