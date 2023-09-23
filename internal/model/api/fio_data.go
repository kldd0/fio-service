package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kldd0/fio-service/internal/model/domain_models"
)

type FioApi interface {
	GetAge(name string) (age int, err error)
	GetGender(name string) (gender string, err error)
	GetNationality(name string) (nationality string, err error)

	FillModel(model *domain_models.FioStruct) error
}

type FioAPIClient struct {
}

func (api FioAPIClient) GetAge(name string) (age int, err error) {
	const op = "model.api.fio_data.GetAge"

	client := http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", "https://api.agify.io/?name="+name, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: failed constructing request: %w", op, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%s: failed making request: %w", op, err)
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return 0, fmt.Errorf("%s: failed unmarshalling response: %w", op, err)
	}

	return int(result["age"].(float64)), nil
}

func (api FioAPIClient) GetGender(name string) (gender string, err error) {
	const op = "model.api.fio_data.GetGender"

	client := http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", "https://api.genderize.io/?name="+name, nil)
	if err != nil {
		return "", fmt.Errorf("%s: failed constructing request: %w", op, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: failed making request: %w", op, err)
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("%s: failed unmarshalling response: %w", op, err)
	}

	return result["gender"].(string), nil
}

func (api FioAPIClient) GetNationality(name string) (nationality string, err error) {
	const op = "model.api.fio_data.GetGender"

	client := http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest("GET", "https://api.nationalize.io/?name="+name, nil)
	if err != nil {
		return "", fmt.Errorf("%s: failed constructing request: %w", op, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: failed making request: %w", op, err)
	}

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("%s: failed unmarshalling response: %w", op, err)
	}

	nationality = result["country"].([]interface{})[0].(map[string]interface{})["country_id"].(string)
	return
}

func (api FioAPIClient) FillModel(model *domain_models.FioStruct) error {
	const op = "model.api.fio_data.FillModel"

	age, err := api.GetAge(model.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	gender, err := api.GetGender(model.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	nationality, err := api.GetNationality(model.Name)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	model.Age = age
	model.Gender = gender
	model.Nationality = nationality

	return nil
}
