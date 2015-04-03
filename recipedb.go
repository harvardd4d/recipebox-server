package main

import (
	"container/list"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// RecipeDB represents a recipe database. Wraps a sqlx.DB.
type RecipeDB struct {
	DB *sqlx.DB
}

// GetRecipe gets a Recipe based on its id.
func (recipeDB *RecipeDB) GetRecipe(id int) (recipe *Recipe, err error) {
	row := recipeDB.DB.QueryRowx("SELECT * FROM recipes WHERE id=$1", id)
	recipe = new(Recipe)
	err = row.StructScan(recipe)
	return
}

// EditRecipe takes an edited recipe and inserts in into the database
func (recipeDB *RecipeDB) EditRecipe(recipe *Recipe) (err error) {
	// 8 things, TODO insert picture
	update := `UPDATE recipes SET ` +
		`name=$2,description=$3,cuisine=$4,mealtype=$5,` +
		`season=$6, ingredientlist=$7, instructions=$8 WHERE id=$1`
	_, err = recipeDB.DB.Exec(update, recipe.ID, recipe.Name,
		recipe.Description, recipe.Cuisine, recipe.Mealtype, recipe.Season,
		recipe.Ingredientlist, recipe.Instructions)
	return err
}

// GetRecipesStrict gets a Recipe based on a strict search
func (recipeDB *RecipeDB) GetRecipesStrict(name string, cuisine,
	mealtype, season int) (recipes *list.List, err error) {

	fmt.Printf("Getting %v %v %v %v", name, cuisine, mealtype, season)
	rows, err := recipeDB.DB.Queryx(`SELECT * `+
		`FROM recipes WHERE name LIKE $1 AND `+
		`cuisine=$2 AND `+
		`mealtype=$3 AND `+
		`season=$4`,
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
