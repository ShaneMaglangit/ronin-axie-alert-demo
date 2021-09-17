package firebase

import (
	"github.com/shanemaglangit/agp"
	"time"
)

type Sale struct {
	TxHash      string    `json:"txHash" firestore:"txHash"`
	ID          int       `json:"id" firestore:"id"`
	StartPrice  float64   `json:"startPrice" firestore:"startPrice"`
	EndPrice    float64   `json:"endPrice" firestore:"endPrice"`
	Timestamp   time.Time `json:"timestamp" firestore:"timestamp"`
	Class       agp.Class `json:"class" firestore:"class"`
	Purity      int       `json:"purity" firestore:"purity"`
	BreedCount  int       `json:"breedCount" firestore:"breedCount"`
	NumMystic   int       `json:"numMystic" firestore:"numMystic"`
	HP          int       `json:"hp" firestore:"hp"`
	Speed       int       `json:"speed" firestore:"speed"`
	Skill       int       `json:"skill" firestore:"skill"`
	Morale      int       `json:"morale" firestore:"morale"`
	EyesD       string    `json:"eyesD" firestore:"eyesD"`
	EyesR1      string    `json:"eyesR1" firestore:"eyesR1"`
	EyesR2      string    `json:"eyesR2" firestore:"eyesR2"`
	EarsD       string    `json:"earsD" firestore:"earsD"`
	EarsR1      string    `json:"earsR1" firestore:"earsR1"`
	EarsR2      string    `json:"earsR2" firestore:"earsR2"`
	HornD       string    `json:"hornD" firestore:"hornD"`
	HornR1      string    `json:"hornR1" firestore:"hornR1"`
	HornR2      string    `json:"hornR2" firestore:"hornR2"`
	MouthD      string    `json:"mouthD" firestore:"mouthD"`
	MouthR1     string    `json:"mouthR1" firestore:"mouthR1"`
	MouthR2     string    `json:"mouthR2" firestore:"mouthR2"`
	BackD       string    `json:"backD" firestore:"backD"`
	BackR1      string    `json:"backR1" firestore:"backR1"`
	BackR2      string    `json:"backR2" firestore:"backR2"`
	TailD       string    `json:"tailD" firestore:"tailD"`
	TailR1      string    `json:"tailR1" firestore:"tailR1"`
	TailR2      string    `json:"tailR2" firestore:"tailR2"`
	GeneQuality float64   `json:"geneQuality" firestore:"geneQuality"`
}
