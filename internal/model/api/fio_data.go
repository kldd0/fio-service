package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FIOApiClient struct {
}

func (api FIOApiClient) GetAge(name string) (age int, err error) {
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

func (api FIOApiClient) GetGender(name string) (gender string, err error) {
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

func (api FIOApiClient) GetNationality(name string) (nationality string, err error) {
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
