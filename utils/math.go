package utils

import "math/big"

func Sum(nums[]float64) float64{
	result := float64(0)
	for _, num := range nums{
		result += num
	}
	return result
}

func Factorial(num int) *big.Int {
	result := big.NewInt(int64(1))
	result = result.MulRange(1, int64(num))
	return result
}