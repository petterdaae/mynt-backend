package synchronize

import (
	"backend/internal/utils"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var shops = []string{
	"KIWI",
	"REMA 1000",
	"Sammen Kantine",
	"Fjordkraft",
	"Telenor",
	"IN2BAR",
	"CIRCLE K NESTTUN",
	"Vipps fra Petter",
	"EXTRA Bergen",
	"Bien Snackbar",
}

func ResetDemoAccount(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)

	check(database.Exec("DELETE from accounts WHERE user_id = 'demo'"))
	check(database.Exec("DELETE from transactions WHERE user_id = 'demo'"))
	check(database.Exec("DELETE from categories WHERE user_id = 'demo'"))
	check(database.Exec("DELETE from categorizations WHERE user_id = 'demo'"))

	minDate := time.Date(2021, 8, 0, 0, 0, 0, 0, time.UTC).Unix()
	maxDate := time.Date(2021, 12, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := maxDate - minDate

	check(database.Exec("INSERT INTO accounts " +
		"(id, user_id, name, account_number, available, external_id, balance) " +
		"VALUES ('0', 'demo', 'Card', '123', 10000, 'lol', 10000)"))

	for i := 0; i < 150; i++ {
		amount := -rand.Int31n(50000) + 10000
		text := shops[rand.Intn(len(shops))]
		date := time.Unix(rand.Int63n(delta)+minDate, 0).Format("2006-01-02") + "T00:00:00"
		check(database.Exec(
			"INSERT INTO transactions (id, user_id, account_id, accounting_date, interest_date, amount, text, external_id) "+
				"VALUES ($1, 'demo', '0', $2, $3, $4, $5, '123')",
			i, date, date, amount, text,
		))
	}
	c.Status(http.StatusOK)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
