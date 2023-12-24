package database

import (
	"blog-server/types"
)

func GetCategories() ([]types.Category, error) {
    var categories []types.Category
	rows, err := db.Query("SELECT name, parent FROM categories")
	if err != nil { return categories, err; }

	for rows.Next() {
		var category types.Category
		err := rows.Scan(&category.Name, &category.Parent)
		if err != nil { return categories, err }
		categories = append(categories, category)
	}
	return categories, nil
}
