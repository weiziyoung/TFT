package globals

import (
	"TFT/models"
	"encoding/json"
	"github.com/schollz/progressbar"
	"os"
)

var (
	traitList    models.TraitList
	ChampionList models.ChampionList
	OneTraitChampionNameList []string
	TraitDict    map[string]models.Trait
	ChampionDict map[string]models.Champion
	TranslateDict map[string]string
	Bar          *progressbar.ProgressBar
	Counter      int64
	GainLevel 	 float32
)

func init() {
	// Init traits file
	traitsFile, err := os.Open("data/traits.json")
	if err != nil {
		panic("Open traits file fail!")
	}
	Counter = 0
	jsonParser := json.NewDecoder(traitsFile)
	if err = jsonParser.Decode(&traitList); err != nil {
		panic(err.Error())
	}
	TraitDict = make(map[string]models.Trait)
	for _, trait := range traitList {
		TraitDict[trait.Name] = models.Trait{
			BonusNum:  trait.BonusNum,
			Scope:     trait.Scope,
			Champions: trait.Champions,
			Strength: trait.Strength,
		}
	}
	GainLevel = 1.2

	// Init champion file
	championFile, err := os.Open("data/champions.json")
	if err != nil {
		panic("Open champion file fail!")
	}
	jsonParser = json.NewDecoder(championFile)
	if err = jsonParser.Decode(&ChampionList); err != nil {
		panic(err.Error())
	}

	translateFile, err := os.Open("data/language.json")
	if err != nil {
		panic("Open language file fail")
	}
	jsonParser = json.NewDecoder(translateFile)
	if err = jsonParser.Decode(&TranslateDict); err!=nil {
		panic(err.Error())
	}
	ChampionDict = make(map[string]models.Champion)
	for _, champion := range ChampionList {
		ChampionDict[champion.Name] = models.Champion{
			Avatar: champion.Avatar,
			Price:  champion.Price,
			Origin: champion.Origin,
			Class:  champion.Class,
		}
	}
	for _, trait := range traitList{
		if trait.BonusNum[0] == 1{
			OneTraitChampionNameList = append(OneTraitChampionNameList, trait.Champions...)
		}
	}
}
