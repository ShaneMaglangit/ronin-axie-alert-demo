package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
	"log"
)

type Firebase struct {
	db *db.Client
}

const RowCap = 1000000

func New() *Firebase {
	conf := firebase.Config{DatabaseURL: "..."}
	opt := option.WithCredentialsFile("key.json")
	app, err := firebase.NewApp(context.Background(), &conf, opt)
	if err != nil {
		log.Fatal(err)
	}
	db, err := app.Database(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return &Firebase{db: db}
}

func (fb *Firebase) SaveSalesMultiple(sales []Sale) {
	salesRef := fb.db.NewRef("sales")
	for _, sale := range sales {
		err := salesRef.Child(sale.TxHash).Set(context.Background(), sale)
		if err != nil {
			log.Fatal(err)
		}
	}
	//fb.LimitSalesRecord()
}

func (fb *Firebase) LimitSalesRecord() {
	var hashes map[string]bool
	salesRef := fb.db.NewRef("sales")
	_ = salesRef.GetShallow(context.Background(), &hashes)
	if len(hashes) > RowCap {
		fb.ReduceExcessSalesRecord(len(hashes) - RowCap)
	}
}

func (fb *Firebase) ReduceExcessSalesRecord(n int) {
	salesRef := fb.db.NewRef("sales")
	res, err := salesRef.OrderByChild("timestamp").GetOrdered(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < n; i++ {
		var sale Sale
		if err := res[i].Unmarshal(&sale); err != nil {
			log.Fatalln("Error unmarshalling result:", err)
		}
		fb.DeleteSaleRecord(sale.TxHash)
	}
}

func (fb *Firebase) DeleteSaleRecord(id string) {
	err := fb.db.NewRef("sales").Child(id).Delete(context.Background())
	if err != nil {
		log.Fatalln("Error while deleting excess", err)
	}
}
