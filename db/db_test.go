package db

/*import (
	"testing"
	"github.com/vahriin/MT/model"
	"time"
	"fmt"
)

var config string = "user=vahriin dbname=MT_DB sslmode=disable"

/*func TestPassUser(t *testing.T) {
	db, err := InitDB(config)
	if err != nil {
		t.Fatal(err)
	}

	gleb := model.PassUser{
		Id: 0,
		Nick: "Gleb",
		PassHash: []byte("abcdef"),
	}

	andrew := model.PassUser{
		Id: 0,
		Nick: "Andrew",
		PassHash: []byte("ffreafdf"),
	}

	err = db.AddPassUser(&gleb)
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddPassUser(&andrew)
	if err != nil {
		t.Fatal(err)
	}

	gleb0, err := db.GetPassUserByEmail("Gleb")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gleb0.Id)

	andrew0, err := db.GetPassUserByEmail("Andrew")
	if err != nil {
		t.Fatal(err)
	}

	pass, err := db.GetPassUserByEmail("Empty")
	if err != nil || pass != nil {
		t.Log(err)
	}

	err = db.DeletePassUser(gleb0)
	if err != nil {
		t.Fatal(err)
	}

	pass, err = db.GetPassUserByEmail("Gleb")
	if err != nil {
		t.Log(err)
	}

	err = db.DeletePassUser(andrew0)
}*/

/*func initUser(db AppDB) {

	gleb := model.PassUser{
		Id: 0,
		Nick: "Gleb",
		Email: "eeee@gmail.com",
		PassHash: []byte("abcdef"),
	}

	andrew := model.PassUser{
		Id: 0,
		Nick: "Andrew",
		Email: "trueprogrammer@gmail.com",
		PassHash: []byte("ffreafdf"),
	}

	alex := model.PassUser{
		Id: 0,
		Nick: "Alex",
		Email: "truemath@yandex.ru",
		PassHash: []byte("telwqhfd"),
	}

	pavel := model.PassUser{
		Id: 0,
		Nick: "Pavel",
		Email: "Ilikegod@gmail.ru",
		PassHash: []byte("aefdjkaf"),
	}

	_ = db.AddPassUser(&gleb)
	_ = db.AddPassUser(&andrew)
	_ = db.AddPassUser(&alex)
	_ = db.AddPassUser(&pavel)
}


func TestTransactions(t *testing.T) {
	db, err := InitDB(config)
	if err != nil {
		t.Fatal(err)
	}
	initUser(db)

	gleb, _ := db.GetUserByNick("Gleb")
	andrew, _ := db.GetUserByNick("Andrew")
	pavel, _ := db.GetUserByNick("Pavel")
	alex, _ := db.GetUserByNick("Alex")

	transact0 := model.Transaction{
		Transaction: model.Transaction{
			Id: 0,
			Date: time.Now(),
			Source: *gleb,
			Sum: 10000,
			Matter: "Chocolate",
			Comment: "Very taste",
		},
		Targets: []model.User{*andrew, *pavel},
	}




	transact1 := model.Transaction{
		Transaction: model.Transaction{
			Id: 0,
			Date: time.Now(),
			Source: *andrew,
			Sum: 30000,
			Matter: "Chips",
			Comment: "Kracks",
		},
		Targets: []model.User{*gleb, *pavel, *alex},
	}

	err = db.AddTransaction(&transact0, []int{1,1,1})
	if err != nil {
		t.Fatal(err)
	}

	err = db.AddTransaction(&transact1, []int{1,1,1,1})
	if err != nil {
		t.Fatal(err)
	}

	transact01, err := db.GetTransactionsBySource(gleb)
	if err != nil {
		t.Fatal(err)
	}

	transact11, err := db.GetTransactionsBySource(andrew)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(transact01)
	fmt.Println(transact11)

	_ = db.DeleteTransaction(&transact01[0])
	_ = db.DeleteTransaction(&transact11[0])
}
*/
