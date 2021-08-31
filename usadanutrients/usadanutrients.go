package usadanutrients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type UsadaNutrients struct {
	key string
}

func New(apiKey string) *UsadaNutrients {
	return &UsadaNutrients{
		key: apiKey,
	}
}

func (n *UsadaNutrients) Nutrients(product string) (string, error) {
	id, err := n.productID(product)
	if err != nil {
		return "", err
	}

	fn, err := n.productNutrients(id)
	if err != nil {
		return "", err
	}

	return fn.String(), nil
}

type food struct {
	FDCID       int    `json:"fdcId"`
	Description string `json:"description"`
}

type foods struct {
	Food []food `json:"foods"`
}

func (n *UsadaNutrients) productID(product string) (int, error) {
	u, err := url.Parse("https://api.nal.usda.gov/fdc/v1/foods/search")
	if err != nil {
		return 0, err
	}

	query := &url.Values{}
	query.Add("api_key", n.key)
	query.Add("query", product)
	query.Add("requireAllWords", "true")
	u.RawQuery = query.Encode()

	rsp, err := http.Get(u.String())
	if err != nil {
		return 0, err
	}

	defer rsp.Body.Close()

	var f foods
	err = json.NewDecoder(rsp.Body).Decode(&f)
	if err != nil {
		return 0, err
	}

	if len(f.Food) <= 0 {
		return 0, errors.New("not found")
	}
	return f.Food[0].FDCID, nil
}

type foodNutrients struct {
	Description string     `json:"description"`
	FN          []nutrient `json:"foodNutrients"`
}

func (fn foodNutrients) String() string {
	list := []string{fmt.Sprintf("product name: %s\n", fn.Description)}
	for _, n := range fn.FN {
		list = append(list, fmt.Sprintf("%s: %f %s", n.ND.Name, n.Amount, n.ND.UnitName))
	}

	return strings.Join(list, "\n")
}

type nutrient struct {
	ND     nutrientDescription `json:"nutrient"`
	Amount float32             `json:"amount"`
}

type nutrientDescription struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	UnitName string `json:"unitName"`
}

func (n *UsadaNutrients) productNutrients(id int) (*foodNutrients, error) {
	u, err := url.Parse(fmt.Sprintf("https://api.nal.usda.gov/fdc/v1/food/%d", id))
	if err != nil {
		return nil, err
	}

	query := &url.Values{}
	query.Add("api_key", n.key)
	query.Add("nutrients", "203") // Protein
	query.Add("nutrients", "204") // Fat
	query.Add("nutrients", "205") // Carbohydrate
	query.Add("nutrients", "208") // Energy
	u.RawQuery = query.Encode()

	rsp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer rsp.Body.Close()

	var f foodNutrients
	if err := json.NewDecoder(rsp.Body).Decode(&f); err != nil {
		return nil, err
	}

	return &f, nil
}
