package models

type Champion struct {
	Avatar string   `json:"avatar"`
	Price  int      `json:"price"`
	Origin []string `json:"origin"`
	Class  []string `json:"class"`
}

type ChampionDict struct {
	Name string `json:"name"`
	*Champion
}

type ChampionList []ChampionDict
