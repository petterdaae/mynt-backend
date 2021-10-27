package spendings

import (
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RawSpending struct {
	CategoryID *int64 `json:"category_id"`
	Amount     int64  `json:"amount"`
	ParentID   *int64 `json:"parent_id"`
}

type Spending struct {
	CategoryID     *int64 `json:"category_id"`
	Amount         int64  `json:"amount"`
	PositiveAmount int64  `json:"positive_amount"`
	NegativeAmount int64  `json:"negative_amount"`
}

type Category struct {
	ID       *int64 `json:"id"`
	ParentID *int64 `json:"parent_id"`
}

type Result struct {
	Spendings []Spending
}

func List(c *gin.Context) {
	database, _ := c.MustGet("database").(*utils.Database)
	sub := c.GetString("sub")

	rows, err := database.Query(
		`SELECT id, parent_id
		FROM categories
		WHERE user_id = $1
		AND DELETED is NULL`,
		sub,
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var category Category
		err = rows.Scan(&category.ID, &category.ParentID)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		categories = append(categories, category)
	}

	rows, err = database.Query(
		`SELECT tc.category_id, SUM(tc.amount), c.parent_id
		FROM transactions_to_categories AS tc, transactions AS t, categories AS c
		WHERE tc.transaction_id = t.id
		AND tc.category_id = c.id
		AND t.user_id = $1
		AND tc.user_id = $1
		AND t.accounting_date >= $2
		AND t.accounting_date <= $3
		GROUP BY tc.category_id, c.id`,
		sub,
		c.Query("from_date")+"T00:00:00",
		c.Query("to_date")+"T00:00:00",
	)
	if err != nil {
		utils.InternalServerError(c, err)
		return
	}
	defer rows.Close()

	rawSpendings := []RawSpending{}
	for rows.Next() {
		var spending RawSpending
		err := rows.Scan(&spending.CategoryID, &spending.Amount, &spending.ParentID)
		if err != nil {
			utils.InternalServerError(c, err)
			return
		}
		rawSpendings = append(rawSpendings, spending)
	}

	result := Result{}
	groupSpendings(nil, &rawSpendings, &categories, &result)

	c.JSON(http.StatusOK, result.Spendings)
}

func groupSpendings(
	categoryID *int64,
	rawSpendings *[]RawSpending,
	categories *[]Category,
	result *Result,
) (total, positive, negative int64) {
	spending := Spending{
		CategoryID: categoryID,
		Amount:     0,
	}

	for _, rawSpending := range *rawSpendings {
		if categoryID != nil && *rawSpending.CategoryID == *categoryID {
			spending.Amount += rawSpending.Amount
			if rawSpending.Amount > 0 {
				spending.PositiveAmount += rawSpending.Amount
			} else {
				spending.NegativeAmount += rawSpending.Amount
			}
		}
	}

	for _, category := range *categories {
		if categoryID == nil && category.ParentID == nil {
			a, p, n := groupSpendings(category.ID, rawSpendings, categories, result)
			spending.Amount += a
			spending.PositiveAmount += p
			spending.NegativeAmount += n
		} else if categoryID != nil && category.ParentID != nil && *category.ParentID == *categoryID {
			a, p, n := groupSpendings(category.ID, rawSpendings, categories, result)
			spending.Amount += a
			spending.PositiveAmount += p
			spending.NegativeAmount += n
		}
	}

	result.Spendings = append(result.Spendings, spending)

	return spending.Amount, spending.PositiveAmount, spending.NegativeAmount
}
