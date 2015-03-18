package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// The Recipe structure
type Recipe struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Cuisine        int    `json:"cuisine"`
	Mealtype       int    `json:"mealtype"`
	Season         int    `json:"season"`
	Ingredientlist string `json:"ingredientlist"`
	Instructions   string `json:"instructions"`
	Picture        []byte `json:"picture"`
}

// Turns a Recipe into a JSON string
func (self *Recipe) ToJSON() (result string) {
	resultBytes, err := json.Marshal(self)
	if err != nil {
		result = "Chosen recipe cannot be formatted into json form"
	} else {
		result = string(resultBytes)
	}
	return
}

// Recipe Database struct.  Wraps a sqlx.DB
type RecipeDatabase struct {
	DB *sqlx.DB
}

// Get a Recipe based on its id.
func (this *RecipeDatabase) GetRecipe(id int) (recipe *Recipe, err error) {
	row := this.DB.QueryRowx("SELECT * FROM recipes WHERE id=?", id)
	recipe = new(Recipe)
	err = row.StructScan(recipe)
	return
}

// Get a Recipe based on a strict search
func (this *RecipeDatabase) GetRecipesStrict(name string, cuisine,
	mealtype, season int) (recipes *list.List, err error) {

	fmt.Printf("Getting %v %v %v %v", name, cuisine, mealtype, season)
	rows, err := this.DB.Queryx("SELECT * "+
		"FROM recipes WHERE name LIKE ? AND "+
		"cuisine=? AND "+
		"mealtype=? AND "+
		"season=?",
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

// Get a Recipe based on a loose search.
func (this *RecipeDatabase) GetRecipesLoose(name string, cuisine,
	mealtype, season int) (recipes *list.List, err error) {

	recipes, err = this.GetRecipesStrict("%"+name+"%", cuisine, mealtype, season)
	return
}
