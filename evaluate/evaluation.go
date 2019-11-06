package evaluate

import (
	"TFT/globals"
	"TFT/models"
)

// Evaluate 评估当前组合的羁绊数量、单位收益羁绊总数、羁绊强度
func Evaluate(combo []models.ChampionDict) models.ComboMetric {
	var nowTraitDict = make(map[string]int)
	var traitDetail = make(map[string]int)
	for key := range globals.TraitDict {
		nowTraitDict[key] = 0
	}
	comboName := make([]string, 0, len(combo))
	traitNum := 0
	totalTraitNum := 0
	totalStrength := float32(0.0)

	for _, unit := range combo {
		comboName = append(comboName, unit.Name)
		for _, origin := range unit.Origin {
			nowTraitDict[origin] += 1
		}
		for _, class := range unit.Class {
			nowTraitDict[class] += 1
		}
	}
	for trait, num := range nowTraitDict {
		if num > 0 {
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
					}
				case 2:
					{
						benefitedNum = num // 对同一种族的Buff，大多数羁绊都是这种
					}
				case 3:
					{
						benefitedNum = len(combo) // 群体Buff，如骑士、六贵族、四帝国
					}
				case 4:
					{
						benefitedNum = len(combo) - 2 // 护卫Buff，比较特殊，除护卫本身外，其他均能吃到buff
					}
				}
				totalTraitNum += bonusLevel * benefitedNum
				totalStrength += float32(bonusLevel) * float32(benefitedNum) * bonusStrength
			}
		}
	}
	metric := models.ComboMetric{
		Combo:         comboName,
		TraitNum:      traitNum,
		TotalTraitNum: totalTraitNum,
		TraitDetail: traitDetail,
		TotalStrength: totalStrength,
	}
	return metric
}
