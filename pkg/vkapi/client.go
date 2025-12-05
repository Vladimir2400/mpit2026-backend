package vkapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	apiVersion = "5.131"
	apiBaseURL = "https://api.vk.com/method"
)

// Client represents VK API client
type Client struct {
	httpClient *http.Client
}

// NewClient creates new VK API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VKUserInfo represents VK user information
type VKUserInfo struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Photo200  string `json:"photo_200"`
	Photo400  string `json:"photo_400_orig"`
	City      *struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"city"`
	BirthDate string `json:"bdate"`
	Sex       int    `json:"sex"` // 0=not specified, 1=female, 2=male
}

// VKAPIResponse represents standard VK API response
type VKAPIResponse struct {
	Response []VKUserInfo `json:"response"`
	Error    *struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	} `json:"error"`
}

// GetUserInfo fetches user information from VK API
func (c *Client) GetUserInfo(accessToken string, userID int) (*VKUserInfo, error) {
	params := url.Values{}
	params.Set("user_ids", fmt.Sprintf("%d", userID))
	params.Set("fields", "photo_200,photo_400_orig,city,bdate,sex")
	params.Set("access_token", accessToken)
	params.Set("v", apiVersion)

	apiURL := fmt.Sprintf("%s/users.get?%s", apiBaseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call VK API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("VK API returned status %d", resp.StatusCode)
	}

	var apiResp VKAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("failed to decode VK API response: %w", err)
	}

	if apiResp.Error != nil {
		return nil, fmt.Errorf("VK API error %d: %s", apiResp.Error.ErrorCode, apiResp.Error.ErrorMsg)
	}

	if len(apiResp.Response) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &apiResp.Response[0], nil
}

// GetUserWall fetches user's wall posts
func (c *Client) GetUserWall(accessToken string, userID int, count int) ([]map[string]interface{}, error) {
	params := url.Values{}
	params.Set("owner_id", fmt.Sprintf("%d", userID))
	params.Set("count", fmt.Sprintf("%d", count))
	params.Set("access_token", accessToken)
	params.Set("v", apiVersion)

	apiURL := fmt.Sprintf("%s/wall.get?%s", apiBaseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call VK API: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode VK API response: %w", err)
	}

	if errData, ok := result["error"]; ok {
		return nil, fmt.Errorf("VK API error: %v", errData)
	}

	responseData, ok := result["response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	items, ok := responseData["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid items format")
	}

	var posts []map[string]interface{}
	for _, item := range items {
		if post, ok := item.(map[string]interface{}); ok {
			posts = append(posts, post)
		}
	}

	return posts, nil
}

// GetUserGroups fetches user's groups
func (c *Client) GetUserGroups(accessToken string, userID int) ([]map[string]interface{}, error) {
	params := url.Values{}
	params.Set("user_id", fmt.Sprintf("%d", userID))
	params.Set("extended", "1")
	params.Set("access_token", accessToken)
	params.Set("v", apiVersion)

	apiURL := fmt.Sprintf("%s/groups.get?%s", apiBaseURL, params.Encode())

	resp, err := c.httpClient.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to call VK API: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode VK API response: %w", err)
	}

	if errData, ok := result["error"]; ok {
		return nil, fmt.Errorf("VK API error: %v", errData)
	}

	responseData, ok := result["response"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format")
	}

	items, ok := responseData["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid items format")
	}

	var groups []map[string]interface{}
	for _, item := range items {
		if group, ok := item.(map[string]interface{}); ok {
			groups = append(groups, group)
		}
	}

	return groups, nil
}
