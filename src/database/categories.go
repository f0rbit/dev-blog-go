package database

import (
	"blog-server/types"
	"errors"
)

func GetCategories(user *types.User) ([]types.Category, error) {
    if user == nil {
        return nil, errors.New("Invalid user reference")
    }
    var categories []types.Category
	rows, err := db.Query("SELECT name, parent FROM categories WHERE owner_id = ?", user.ID)
	if err != nil { return categories, err; }

	for rows.Next() {
		var category types.Category
        category.OwnerID = user.ID
		err := rows.Scan(&category.Name, &category.Parent)
		if err != nil { return categories, err }
		categories = append(categories, category)
	}
	return categories, nil
}

func ConstructCategoryGraph(categories []types.Category, root string) types.CategoryNode {
    var node = types.CategoryNode{
        Name: root,
        Children: make([]types.CategoryNode, 0),
        OwnerID: categories[0].OwnerID,
    }
	for _, cat := range categories {
		if cat.Parent == node.Name {
            node.Children = append(node.Children, ConstructCategoryGraph(categories, cat.Name))
		}
	}
	return node;
}

func CreateCategory(category types.Category) error {
    _, err := db.Exec(`INSERT INTO categories (owner_id, name, parent) VALUES (?, ?, ?)`, category.OwnerID, category.Name, category.Parent);
    return err;
}
