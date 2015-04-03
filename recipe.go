package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	Breakfast = 1 << iota
	Lunch     = 1 << iota
	Dinner    = 1 << iota
)

const (
	Spring = 1 << iota
	Summer = 1 << iota
	Fall   = 1 << iota
	Winter = 1 << iota
)

var (
	Meals = map[int]string{
		Breakfast: "Breakfast",
		Lunch:     "Lunch",
		Dinner:    "Dinner",
	}
	Seasons = map[int]string{
		Spring: "Spring",
		Summer: "Summer",
		Winter: "Winter",
		Fall:   "Fall",
	}
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

// ParseIngredients returns a list of ingredients from a string with
// delimiter ;
func ParseIngredients(ingredients string) []string {
	things := strings.Split(ingredients, ";")
	for k, v := range things {
		things[k] = strings.TrimSpace(v)
	}
	return things
}

// ParseMealtype returns a list of mealtypes (strings) based on meal
func ParseMealtype(mealtype int) []string {
	return nil
}

// ParseSeason returns a list of seasons(string) based on season
func ParseSeason(season int) []string {
	return nil
}

// MealIs returns whether or not the mealtype is a particular meal
func MealIs(meal string, mealtype int) bool {
	switch meal {
	default:
		return false
	case "Breakfast":
		return mealtype&Breakfast > 0
	case "Lunch":
		return mealtype&Lunch > 0
	case "Dinner":
		return mealtype&Dinner > 0
	}
}

// SeasonIs returns whether or not the season is a particular season
func SeasonIs(season string, seasontype int) bool {
	switch season {
	default:
		return false
	case "Spring":
		return seasontype&Spring > 0
	case "Summer":
		return seasontype&Summer > 0
	case "Fall":
		return seasontype&Fall > 0
	case "Winter":
		return seasontype&Winter > 0
	}
}
