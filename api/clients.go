package api

import (
	"fmt"
	"strings"
	"net/url"

	"resty.dev/v3"

	"cwl/config"
	"cwl/types"
)



type ClashClient struct {
	*resty.Client
}

func NewClashClient(cfg *config.ClashConfig) *ClashClient {

	return &ClashClient{
		resty.New().
			SetBaseURL(cfg.BaseURL).
			SetHeader("Authorization", "Bearer " + cfg.ApiKey).
			SetHeader("Content-Type", "application/json"),
	}
}



func (c *ClashClient) verifyAccount(tag string, token string) (bool, error) {

	req := c.R().SetBody(map[string]string{"token": token})

	encodedTag := url.PathEscape("#" + tag)

	urlFragment := fmt.Sprintf("/players/%s/verifytoken", encodedTag)
	resp, err := req.Post(urlFragment)


	if err != nil {
		fmt.Println("Could not send request")
		return false, err
	}

	if resp.StatusCode() == 500 && strings.Contains(resp.String(), "unknownException") {
		fmt.Println("Gibberish Tag")
		return false, fmt.Errorf("[!] Tag can't exist: Invalid Tag Format")
	}

	if resp.StatusCode() != 200 {
		return false, fmt.Errorf("[!] API returned status %d: %s", resp.StatusCode(), resp.String())
	}




	if strings.Contains(resp.String(), `"status":"ok"`) {
		fmt.Println("Verified!")
		return true, nil

	} else if strings.Contains(resp.String(), `"status":"invalid"`) {
		fmt.Println("Invalid Token")
		return false, nil

	} else {
		fmt.Println("Unknown Error / Malformed Response")
		return false, fmt.Errorf("unexpected API response")
	}
}




func (c *ClashClient) getAccount(tag string) (*types.ClashData, error) {

	urlFragment := fmt.Sprintf("/players/%%23%s", tag)
	account := &types.ClashData{Tag: tag}

	req := c.R().SetResult(account)
	resp, err := req.Get(urlFragment)

	if err != nil {
		fmt.Println("Could not send request")
		return nil, err
	}

	if resp.StatusCode() != 200 {
		fmt.Println("API error")
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode(), resp.String())
	}

	fmt.Println(*account)

	return account, err
}












type DiscordClient struct {
	*resty.Client
}