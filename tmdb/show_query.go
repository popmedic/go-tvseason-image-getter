package tmdb

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/popmedic/go-logger/log"
)

const queryURLFmt = "https://api.themoviedb.org/3/search/tv?api_key=ae802ff2638e8a186add7079dda29e03&language=en-US&query=%s&page=1"

type ShowQuery struct {
	Page    int `json:"page,omitempty"`
	Results []struct {
		PosterPath       string   `json:"poster_path,omitempty"`
		Popularity       float64  `json:"popularity,omitempty"`
		ID               int      `json:"id,omitempty"`
		BackdropPath     string   `json:"backdrop_path,omitempty"`
		VoteAverage      float64  `json:"vote_average,omitempty"`
		Overview         string   `json:"overview,omitempty"`
		FirstAirDate     string   `json:"first_air_date,omitempty"`
		OriginCountry    []string `json:"origin_country,omitempty"`
		GenreIDs         []int    `json:"genre_ids,omitempty"`
		OriginalLanguage string   `json:"original_language,omitempty"`
		VoteCount        int      `json:"vote_count,omitempty"`
		Name             string   `json:"name,omitempty"`
		OriginalName     string   `json:"original_name,omitempty"`
	} `json:"results,omitempty"`
	TotalResults int `json:"total_results,omitempty"`
	TotalPages   int `json:"total_pages,omitempty"`
}

func NewShowQuery(data []byte) (*ShowQuery, error) {
	show := ShowQuery{}
	if err := json.Unmarshal(data, &show); err != nil {
		return nil, err
	}
	return &show, nil
}

func (s *ShowQuery) String() string {
	res, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(res)
}

func QueryShows(showName string, getter HttpGetter) (*ShowQuery, error) {
	url := fmt.Sprintf(queryURLFmt, url.QueryEscape(showName))
	log.Infof("show query url = %q", url)
	body, err := getBody(url, getter)
	if err != nil {
		return nil, err
	}
	return NewShowQuery(body)
}
