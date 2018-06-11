package tmdb

import (
	"encoding/json"
	"fmt"
)

const seasonURLFmt = "https://api.themoviedb.org/3/tv/%d/season/%d?api_key=ae802ff2638e8a186add7079dda29e03&language=en-US"

type Season struct {
	ID_      string `json:"_id,omitempty"`
	AirDate  string `json:"air_date,omitempty"`
	Episodes []struct {
		AirDate string `json:"air_date,omitempty"`
		Crew    []struct {
			ID          int    `json:"id,omitempty"`
			CreditID    string `json:"credit_id,omitempty"`
			Name        string `json:"name,omitempty"`
			Department  string `json:"department,omitempty"`
			Job         string `json:"job,omitempty"`
			ProfilePath string `json:"profile_path,omitempty"`
		} `json:"crew,omitempty"`
		EpisodeNumber int `json:"episode_number,omitempty"`
		GuestStars    []struct {
			ID          int    `json:"id,omitempty"`
			Name        string `json:"name,omitempty"`
			CreditID    string `json:"credit_id,omitempty"`
			Character   string `json:"character,omitempty"`
			Order       int    `json:"order,omitempty"`
			ProfilePath string `json:"profile_path,omitempty"`
		} `json:"guest_stars,omitempty"`
		Name           string  `json:"name,omitempty"`
		Overview       string  `json:"overview,omitempty"`
		ID             int     `json:"id,omitempty"`
		ProductionCode string  `json:"production_code,omitempty"`
		SeasonNumber   int     `json:"season_number,omitempty"`
		StillPath      string  `json:"still_path,omitempty"`
		VoteAverage    float64 `json:"vote_average,omitempty"`
		VoteCount      int     `json:"vote_count,omitempty"`
	} `json:"episodes,omitempty"`
	Name         string `json:"name,omitempty"`
	Overview     string `json:"overview,omitempty"`
	ID           int    `json:"id,omitempty"`
	PosterPath   string `json:"poster_path,omitempty"`
	SeasonNumber int    `json:"season_number,omitempty"`
}

func NewSeason(data []byte) (*Season, error) {
	season := Season{}
	if err := json.Unmarshal(data, &season); err != nil {
		return nil, err
	}
	return &season, nil
}

func GetSeason(showID, seasonNumber int, getter HttpGetter) (*Season, error) {
	url := fmt.Sprintf(seasonURLFmt, showID, seasonNumber)
	body, err := getBody(url, getter)
	if err != nil {
		return nil, err
	}
	return NewSeason(body)
}

func (s *Season) String() string {
	res, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(res)
}
