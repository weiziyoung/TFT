package utils

func Contain(container []interface{}, object interface{}) bool{
	for _,candidate:=range container{
		if candidate == object{
			return true
		}
	}
	return false
}

func Deduplicate(list []int)[]int{
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