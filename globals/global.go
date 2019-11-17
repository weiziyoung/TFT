package globals

import (
	"TFT/models"
	"encoding/json"
	"github.com/go-yaml/yaml"
	"github.com/schollz/progressbar"
	"io/ioutil"
	"os"
)

var (
	traitList                models.TraitList
	ChampionList             models.ChampionList
	OneTraitChampionNameList []string
	TraitDict                map[string]models.Trait
	ChampionDict             map[string]models.Champion
	TranslateDict            map[string]string
	Bar                      *progressbar.ProgressBar
	Counter                  int64
	Global                   GlobalStruct
)

type GlobalStruct struct {
	TraitPath    string  `yaml:"trait_path"`
	ChampionPath string  `yaml:"champion_path"`
	LanguagePath string  `yaml:"language_path"`
	OutPutPath   string  `yaml:"output_path"`
	GainLevel    float64 `yaml:"gain_level"`
	MaximumHeap  int     `yaml:"maximum_heap"`
	Evaluation   string  `yaml:"evaluation"`
}

func init() {
	// 读取配置文件
	content, _ := ioutil.ReadFile("config/app.yaml")
	err := yaml.Unmarshal(content, &Global)
	// Init traits file
	traitsFile, err := os.Open(Global.TraitPath)
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
			Strength:  trait.Strength,
		}
	}
	// Init champion file
	championFile, err := os.Open(Global.ChampionPath)
	if err != nil {
		panic("Open champion file fail!")
	}
	jsonParser = json.NewDecoder(championFile)
	if err = jsonParser.Decode(&ChampionList); err != nil {
		panic(err.Error())
	}

	translateFile, err := os.Open(Global.LanguagePath)
	if err != nil {
		panic("Open language file fail")
	}
	jsonParser = json.NewDecoder(translateFile)
	if err = jsonParser.Decode(&TranslateDict); err != nil {
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
	for _, trait := range traitList {
		if trait.BonusNum[0] == 1 {
			OneTraitChampionNameList = append(OneTraitChampionNameList, trait.Champions...)
		}
	}
}
