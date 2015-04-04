package main

import (
	"encoding/json"
	"strings"
)

// constants in map form
var (
	Meals = map[int]string{
		1: "Breakfast",
		2: "Lunch",
		4: "Dinner",
	}
	MealsToInt = map[string]int{
		"Breakfast": 1,
		"Lunch":     2,
		"Dinner":    4,
	}
	Seasons = map[int]string{
		1: "Spring",
		2: "Summer",
		4: "Winter",
		8: "Fall",
	}
	SeasonsToInt = map[string]int{
		"Spring": 1,
		"Summer": 2,
		"Winter": 4,
		"Fall":   8,
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

// ParseIngredients returns a list of ingredients from a string with
// delimiter ;
func ParseIngredients(ingredients string) []string {
	things := strings.Split(ingredients, ";")
	for k, v := range things {
		things[k] = strings.TrimSpace(v)
	}
	return things
}

// ParseMealtype returns a map of whether or not the number
// represents a particular mealtype.
// This map is a map from meal to bool (t/f)
func ParseMealtype(mealtype int) map[string]bool {
	result := make(map[string]bool)
	for meal, name := range Meals {
		if meal&mealtype > 0 {
			result[name] = true
		} else {
			result[name] = false
		}
	}
	return result
}

// ParseSeason returns a map of whether or not the number
// represents a particular season.
// The map maps seasons (string) to bool
func ParseSeason(season int) map[string]bool {
	result := make(map[string]bool)
	for s, name := range Seasons {
		if s&season > 0 {
			result[name] = true
		} else {
			result[name] = false
		}
	}
	return result
}

// EncodeMealtype will encode a mealtype int
// based on the meals present.
func EncodeMealtype(meals []string) (mealtype int) {
	for _, m := range meals {
		mealtype += MealsToInt[m]
	}
	return
}

// EncodeSeason will encode a season int
// based on the seasons present.
func EncodeSeason(seasons []string) (season int) {
	for _, s := range seasons {
		season += SeasonsToInt[s]
	}
	return
}
