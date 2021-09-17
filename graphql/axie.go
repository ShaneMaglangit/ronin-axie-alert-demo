package graphql

import (
	"github.com/shanemaglangit/agp"
)

type AxieRes struct {
	Data map[string]Axie `json:"data"`
}

type Axie struct {
	ID         string    `json:"id"`
	Class      agp.Class `json:"class"`
	Stage      int       `json:"stage"`
	MatronId   int       `json:"matronId"`
	SireId     int       `json:"sireId"`
	Pureness   int       `json:"pureness"`
	NumMystic  int       `json:"numMystic"`
	BreedCount int       `json:"breedCount"`
	Stats      struct {
		HP     int `json:"hp"`
		Speed  int `json:"speed"`
		Skill  int `json:"skill"`
		Morale int `json:"morale"`
	}
	BirthDate int64  `json:"birthDate"`
	Genes     string `json:"genes"`
}
