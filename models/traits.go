package models

type Trait struct{
	BonusNum []int `json:"bonus_num"`
	Scope []int `json:"scope"`
	Champions []string `json:"champions"`
	Strength []float32 `json:"strength"`
}

type TraitDict struct {
	Name string `json:"name"`
	*Trait
}

type TraitList []TraitDict