package database

import (
	"blog-server/types"
	"blog-server/utils"
	"errors"
)

func GetCategories(user *types.User) ([]types.Category, error) {
	if user == nil {
		return nil, errors.New("Invalid user reference")
	}
	var categories []types.Category
	rows, err := db.Query("SELECT name, parent FROM categories WHERE owner_id = ?", user.ID)
	if err != nil {
		return categories, err
	}

	for rows.Next() {
		var category types.Category
		category.OwnerID = user.ID
		err := rows.Scan(&category.Name, &category.Parent)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetCategory(user *types.User, name string) (types.Category, error) {
	var category types.Category
	if user == nil {
		return category, errors.New("Invalid user reference")
	}
	err := db.QueryRow("SELECT name, parent, owner_id FROM categories WHERE owner_id = ? AND name = ?", user.ID, name).Scan(&category.Name, &category.Parent, &category.OwnerID)
	if err != nil {
		return category, err
	}
	return category, nil

}

func ConstructCategoryGraph(categories []types.Category, root string) types.CategoryNode {
	var node = types.CategoryNode{
		Name:     root,
		Children: make([]types.CategoryNode, 0),
		OwnerID:  categories[0].OwnerID,
	}
	for _, cat := range categories {
		if cat.Parent == node.Name {
			node.Children = append(node.Children, ConstructCategoryGraph(categories, cat.Name))
		}
	}
	return node
}

func CreateCategory(category types.Category) error {
	_, err := db.Exec(`INSERT INTO categories (owner_id, name, parent) VALUES (?, ?, ?)`, category.OwnerID, category.Name, category.Parent)
	return err
}

func DeleteCategory(user *types.User, name string) error {
    // we want to get all the child categories of the one we're about to delete
    // as we are going to delete those as well
    if user == nil {
		return errors.New("Invalid user reference")
    }

    categories, err := GetCategories(user);
    if err != nil {
        return err;
    }

    children := utils.GetChildrenCategories(categories, name);

    cat_list := make([]string, 0)
    cat_list = append(cat_list, name);
    for _, c := range children {
        cat_list = append(cat_list, c.Name)
    }

    err = RemoveCategoryFromPosts(user, cat_list);
    if err != nil {
        return err;
    }

    query, err := db.Prepare("DELETE FROM categories WHERE name = ? AND owner_id = ?");
    if err != nil {
        return err;
    }

    for _, c := range children {
        query.Exec(c.Name, c.OwnerID);
    }

    query.Exec(name, user.ID);

    return nil;
}
