package evaluate

import (
	"TFT/globals"
	"TFT/models"
	"TFT/utils"
	"math"
)

// Evaluate 评估当前组合的羁绊数量、单位收益羁绊总数、羁绊强度
func Evaluate(combo []models.ChampionDict) models.ComboMetric {
	var traitDetail = make(map[string]int)

	comboName := make([]string, 0, len(combo))
	traitNum := 0
	totalTraitNum := 0
	totalTraitStrength := float32(0.0)

	// 初始化英雄强度向量
	unitsStrength := make([]float64, len(combo), len(combo))
	traitChampionsDict := make(map[string][]int)
	for index, unit := range combo {
		comboName = append(comboName, unit.Name)
		unitStrength := math.Pow(float64(globals.GainLevel), float64(unit.Price-1))
		unitsStrength[index] = unitStrength
		for _, origin := range unit.Origin {
			traitChampionsDict[origin] = append(traitChampionsDict[origin], index)
		}
		for _, class := range unit.Class {
			traitChampionsDict[class] = append(traitChampionsDict[class], index)
		}
	}

	for trait, champions := range traitChampionsDict {
		num := len(champions)
		bonusRequirement := globals.TraitDict[trait].BonusNum
		var bonusLevel = len(bonusRequirement)
		for index, requirement := range bonusRequirement {
			if requirement > num {
				bonusLevel = index
				break
			}
		}

		// 忍者只有在1只和4只时触发，其他不触发
		if trait == "ninja" && 1 < num && num < 4 {
			bonusLevel = 0
		}
		if bonusLevel > 0 {
			traitDetail[trait] = bonusRequirement[bonusLevel-1]
			bonusScope := globals.TraitDict[trait].Scope[bonusLevel-1]
			traitNum += bonusLevel
			bonusStrength := globals.TraitDict[trait].Strength[bonusLevel-1]
			benefitedNum := 0
			switch bonusScope {
			case 1:
				{
					benefitedNum = 1 // 单体Buff，例如 机器人、浪人、三贵族、双帝国
					for _, champion := range champions {
						unitsStrength[champion] *= float64(bonusStrength)
					}
				}
			case 2:
				{
					benefitedNum = num // 对同一种族的Buff，大多数羁绊都是这种
					for _, champion := range champions {
						unitsStrength[champion] *= float64(bonusStrength)
					}
				}
			case 3:
				{
					benefitedNum = len(combo) // 群体Buff，如骑士、六贵族、四帝国
					for index, _ := range unitsStrength {
						unitsStrength[index] *= float64(bonusStrength)
					}
				}
			case 4:
				{
					benefitedNum = len(combo) - 2 // 护卫Buff，比较特殊，除护卫本身外，其他均能吃到buff
					for index, _ := range unitsStrength {
						isGuard := false
						for _, champion := range champions {
							if index == champion {
								isGuard = true
								break
							}
						}
						if !isGuard {
							unitsStrength[index] *= float64(bonusStrength)
						}
					}
				}
			}
			totalTraitNum += bonusLevel * benefitedNum
			totalTraitStrength += float32(benefitedNum) * bonusStrength
		}
	}
	metric := models.ComboMetric{
		Combo:              comboName,
		TraitNum:           traitNum,
		TotalTraitNum:      totalTraitNum,
		TraitDetail:        traitDetail,
		TotalTraitStrength: totalTraitStrength,
		TotalStrength:      utils.Sum(unitsStrength),
	}
	return metric
}
