package tmdb

import (
	"encoding/json"
)

const configURL = "https://api.themoviedb.org/3/configuration?api_key=ae802ff2638e8a186add7079dda29e03"

type Config struct {
	Images struct {
		BaseURL       string   `json:"baseURL,omitempty"`
		SecureBaseUrl string   `json:"secure_base_url,omitempty"`
		BackdropSizes []string `json:"backdrop_sizes,omitempty"`
		LogoSizes     []string `json:"logo_sizes,omitempty"`
		PosterSizes   []string `json:"poster_sizes,omitempty"`
		ProfileSizes  []string `json:"profile_sizes,omitempty"`
		StillSizes    []string `json:"still_sizes,omitempty"`
	} `json:"images,omitempty"`
	ChangeKeys []string `json:"change_keys,omitempty"`
}

func NewConfig(data []byte) (*Config, error) {
	cfg := Config{}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) String() string {
	res, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(res)
}

func GetConfig(getter HttpGetter) (*Config, error) {
	body, err := getBody(configURL, getter)
	if err != nil {
		return nil, err
	}
	return NewConfig(body)
}
