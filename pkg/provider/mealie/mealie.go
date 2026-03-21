package mealie

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
}

func NewClient(baseURL, apiToken string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiToken:   apiToken,
		httpClient: &http.Client{},
	}
}

type PaginatedResponse struct {
	Page       int             `json:"page"`
	PerPage    int             `json:"perPage"`
	Total      int             `json:"total"`
	TotalPages int             `json:"totalPages"`
	Items      []Recipe `json:"items"`
}

type Recipe struct {
	ID                 *string          `json:"id"`
	UserID             string           `json:"userId"`
	HouseholdID        string           `json:"householdId"`
	GroupID            string           `json:"groupId"`
	Name               *string          `json:"name"`
	Slug               string           `json:"slug"`
	Image              interface{}      `json:"image"`
	RecipeServings     float64          `json:"recipeServings"`
	RecipeYieldQuantity float64         `json:"recipeYieldQuantity"`
	RecipeYield        *string          `json:"recipeYield"`
	TotalTime          *string          `json:"totalTime"`
	PrepTime           *string          `json:"prepTime"`
	CookTime           *string          `json:"cookTime"`
	PerformTime        *string          `json:"performTime"`
	Description        *string          `json:"description"`
	RecipeCategory     []RecipeCategory `json:"recipeCategory"`
	Tags               []RecipeTag      `json:"tags"`
	Tools              []RecipeTool     `json:"tools"`
	Rating             *float64         `json:"rating"`
	OrgURL             *string          `json:"orgURL"`
	DateAdded          *string          `json:"dateAdded"`
	DateUpdated        *string          `json:"dateUpdated"`
	CreatedAt          *string          `json:"createdAt"`
	UpdatedAt          *string          `json:"updatedAt"`
	LastMade           *string          `json:"lastMade"`
}

type RecipeCategory struct {
	ID      *string `json:"id"`
	GroupID *string `json:"groupId"`
	Name    string  `json:"name"`
	Slug    string  `json:"slug"`
}

type RecipeTag struct {
	ID      *string `json:"id"`
	GroupID *string `json:"groupId"`
	Name    string  `json:"name"`
	Slug    string  `json:"slug"`
}

type RecipeTool struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (c *Client) GetAllRecipes() ([]Recipe, error) {
	var allRecipes []Recipe
	page := 1

	for {
		log.Debugf("Fetching recipes page %d", page)
		resp, err := c.getRecipesPage(page)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch recipes page %d: %w", page, err)
		}

		// if there's only one page resp.TotalPages is 0 so we explicitly set it here for log output
		// otherwise it says "fetched page 1/0" which is confusing.
		if resp.TotalPages == 0 {
			resp.TotalPages = 1
		}

		allRecipes = append(allRecipes, resp.Items...)
		log.Debugf("Fetched page %d/%d (%d recipes)...", resp.Page, resp.TotalPages, len(resp.Items))

		if page >= resp.TotalPages {
			break
		}
		page++
	}

	log.Debugf("Fetched %d total recipes", len(allRecipes))
	return allRecipes, nil
}

func (c *Client) getRecipesPage(page int) (*PaginatedResponse, error) {
	url := fmt.Sprintf("%s/api/recipes?page=%d&perPage=50", c.baseURL, page)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var paginatedResp PaginatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&paginatedResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paginatedResp, nil
}
