package utils

import (
	"TFT/globals"
	"TFT/models"
)

func Translate(combos *models.ComboMetricHeap){
	for comboNo, combo := range *combos{
		for nameNo, name := range combo.Combo{
			newName := globals.TranslateDict[name]
			(*combos)[comboNo].Combo[nameNo] = newName
		}
		for traitName, traitNum := range combo.TraitDetail{
			newTraitName := globals.TranslateDict[traitName]
			if newTraitName == ""{
				continue
			}
			(*combos)[comboNo].TraitDetail[newTraitName] = traitNum
			delete((*combos)[comboNo].TraitDetail, traitName)
		}
	}
}