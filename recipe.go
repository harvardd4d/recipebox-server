package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// Recipe represents a recipe.
type Recipe struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Cuisine        int    `json:"cuisine"`
	Mealtype       int    `json:"mealtype"`
	Season         int    `json:"season"`
	Ingredientlist string `json:"ingredientlist"`
	Instructions   string `json:"instructions"`
	Picture        []byte `json:"picture"`
}

// RecipeDB represents a recipe database. Wraps a sqlx.DB.
type RecipeDB struct {
	DB *sqlx.DB
}

// ToJSON turns a Recipe into a JSON string
func (recipe *Recipe) ToJSON() (result string) {
	resultBytes, err := json.Marshal(recipe)
	if err != nil {
		result = "Chosen recipe cannot be formatted into json form"
	} else {
		result = string(resultBytes)
	}
	return
}

// GetRecipe gets a Recipe based on its id.
func (recipeDB *RecipeDB) GetRecipe(id int) (recipe *Recipe, err error) {
	row := recipeDB.DB.QueryRowx("SELECT * FROM recipes WHERE id=$1", id)
	recipe = new(Recipe)
	err = row.StructScan(recipe)
	return
}

// GetRecipesStrict gets a Recipe based on a strict search
func (recipeDB *RecipeDB) GetRecipesStrict(name string, cuisine,
	mealtype, season int) (recipes *list.List, err error) {

	fmt.Printf("Getting %v %v %v %v", name, cuisine, mealtype, season)
	rows, err := recipeDB.DB.Queryx("SELECT * "+
		"FROM recipes WHERE name LIKE $1 AND "+
		"cuisine=$2 AND "+
		"mealtype=$3 AND "+
		"season=$4",
		name, cuisine, mealtype, season)

	if err != nil {
		fmt.Printf("[WARNING] in GetRecipesStrict: %s", err.Error())
	}

	recipes = list.New()

	for rows.Next() {
		recipe := new(Recipe)
		err = rows.StructScan(recipe)
		if err == nil {
			recipes.PushBack(recipe)
		} else {
			fmt.Printf("[WARNING] StructScan: %s", err.Error())
		}
	}
	return
}

// GetRecipesLoose gets a Recipe based on a loose search.
func (recipeDB *RecipeDB) GetRecipesLoose(name string, cuisine,
	mealtype, season int) (recipes *list.List, err error) {

	recipes, err = recipeDB.GetRecipesStrict("%"+name+"%", cuisine, mealtype, season)
	return
}
