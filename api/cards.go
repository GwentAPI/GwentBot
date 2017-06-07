package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type Pagecard struct {
	Count   int           `json:"count"`
	Results []GenericLink `json:"results"`
}

type GenericLink struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

type Card struct {
	Faction    GenericLink     `json:"faction"`
	Name       string          `json:"name"`
	Categories *[]GenericLink  `json:"categories"`
	Flavor     *string         `json:"flavor"`
	Info       *string         `json:"info"`
	Positions  *[]string       `json:"positions"`
	Strength   *int            `json:"strength"`
	Variations []VariationLink `json:"variations"`
}
type VariationLink struct {
	Href   string      `json:"href"`
	Rarity GenericLink `json:"rarity"`
}

func (c Card) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Name: ")
	buffer.WriteString(c.Name)
	buffer.WriteString("\n")
	buffer.WriteString("Faction: ")
	buffer.WriteString(c.Faction.Name)
	buffer.WriteString("\n")
	if c.Flavor != nil {
		buffer.WriteString("Flavor: ")
		buffer.WriteString(*c.Flavor)
		buffer.WriteString("\n")
	}
	if c.Info != nil {
		buffer.WriteString("Info: ")
		buffer.WriteString(*c.Info)
		buffer.WriteString("\n")
	}
	if c.Strength != nil {
		buffer.WriteString("Strength: ")
		s := strconv.Itoa(*c.Strength)
		buffer.WriteString(s)
		buffer.WriteString("\n")
	}
	if c.Positions != nil && len(*c.Positions) > 0 {
		buffer.WriteString("Positions: ")
		for _, position := range *c.Positions {
			buffer.WriteString(position)
			buffer.WriteString(" ")
		}
		buffer.WriteString("\n")
	}
	if c.Categories != nil && len(*c.Categories) > 0 {
		buffer.WriteString("Categories: ")
		for _, category := range *c.Categories {
			buffer.WriteString(category.Name)
			buffer.WriteString(" ")
		}
		buffer.WriteString("\n")
	}
	if len(c.Variations) > 0 {
		buffer.WriteString("Rarity: ")
		buffer.WriteString(c.Variations[0].Rarity.Name)
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func RequestPage(client *http.Client, query string) (Pagecard, error) {
	var searchPage Pagecard
	var buffer bytes.Buffer
	buffer.WriteString("https://api.gwentapi.com/v0/cards?name=")
	buffer.WriteString(query)

	resp, err := client.Get(buffer.String())
	buffer.Reset()
	if err != nil {
		return searchPage, errors.New("Error contacting the service.")
	}

	switch resp.StatusCode {
	case 500:
		return searchPage, errors.New("The service has encountered an error.")
	case 404:
		return searchPage, errors.New("Page not found.")
	case 200:
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&searchPage)

		if searchPage.Count == 0 {
			return searchPage, errors.New("No result for " + query)
		}
		return searchPage, nil
	default:
		return searchPage, errors.New("Unexpected response from the service.")
	}
}

func RequestCard(client *http.Client, requestURL string) (Card, error) {
	var card Card
	var buffer bytes.Buffer

	resp, err := client.Get(requestURL)
	buffer.Reset()
	if err != nil {
		return card, errors.New("Error contacting the service.")
	}

	switch resp.StatusCode {
	case 500:
		return card, errors.New("The service has encountered an error.")
	case 404:
		return card, errors.New("Page not found.")
	case 200:
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&card)

		return card, nil
	default:
		return card, errors.New("Unexpected response from the service.")
	}
}
