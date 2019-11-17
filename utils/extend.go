package utils

// Contain 相当于 python 的 in
func Contain(container []int, object int) bool{
	for _,candidate:=range container{
		if candidate == object{
			return true
		}
	}
	return false
}

// DeduplicateAndFilter 去重并且过滤比当前英雄位置先的英雄
func DeduplicateAndFilter(list []int, node int)[]int{
	var lastOne int
	var resultList []int
	lastOne = 65535
	for _, item:=range list{
		if 	item > node && item != lastOne{
			resultList = append(resultList, item)
			lastOne = item
		}
	}
	return resultList
}

// Deduplicate 去重
func Deduplicate(list []int) []int{
	var lastOne int
	var resultList []int
	lastOne = 65535
	for _, item:=range list{
		if 	item != lastOne{
			resultList = append(resultList, item)
			lastOne = item
		}
	}
	return resultList
}

// IsEqual 判断两个切片是否相等
func IsEqual(list1 []int, list2 []int) bool{
	if len(list1)!=len(list2){
		return false
	}
	for i:=0;i<len(list1);i++{
		if list1[i]!=list2[i]{
			return false
		}
	}
	return true
}