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
