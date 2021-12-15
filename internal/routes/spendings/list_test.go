package spendings

import "testing"

func TestGroupSpendings(t *testing.T) {
	var categoryID int64 = 60
	var childCategoryID int64 = 65
	rawSpendings := []RawSpending{
		{
			CategoryID: &childCategoryID,
			ParentID:   &categoryID,
			Amount:     100,
		},
	}
	cateogries := []Category{
		{
			ID:       &categoryID,
			ParentID: nil,
		},
		{
			ID:       &childCategoryID,
			ParentID: &categoryID,
		},
	}
	result := Result{}
	groupSpendings(nil, &rawSpendings, &cateogries, &result)
	if result.Spendings[0].Amount != 100 {
		t.Error()
	}
}

func TestGroupSpendingsTopLevelAmounts(t *testing.T) {
	var categoryID int64 = 60
	rawSpendings := []RawSpending{
		{
			CategoryID: &categoryID,
			ParentID:   nil,
			Amount:     100,
		},
		{
			CategoryID: &categoryID,
			ParentID:   nil,
			Amount:     -100,
		},
	}
	cateogries := []Category{
		{
			ID:       &categoryID,
			ParentID: nil,
		},
	}
	result := Result{}
	groupSpendings(nil, &rawSpendings, &cateogries, &result)
	if result.Spendings[0].Amount != 0 || result.Spendings[0].PositiveAmount != 100 || result.Spendings[0].NegativeAmount != -100 {
		t.Error()
	}
}
