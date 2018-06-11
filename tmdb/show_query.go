package tmdb

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/xrash/smetrics"
)

const queryURLFmt = "https://api.themoviedb.org/3/search/tv?api_key=ae802ff2638e8a186add7079dda29e03&language=en-US&query=%s&page=1"

type ShowQueryResult struct {
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
}

type ShowQuery struct {
	Page         int               `json:"page,omitempty"`
	Results      []ShowQueryResult `json:"results,omitempty"`
	TotalResults int               `json:"total_results,omitempty"`
	TotalPages   int               `json:"total_pages,omitempty"`
}

func NewShowQuery(data []byte) (*ShowQuery, error) {
	show := ShowQuery{}
	if err := json.Unmarshal(data, &show); err != nil {
		return nil, err
	}
	return &show, nil
}

func QueryShows(showName string, getter HttpGetter) (*ShowQuery, error) {
	url := fmt.Sprintf(queryURLFmt, url.QueryEscape(showName))
	body, err := getBody(url, getter)
	if err != nil {
		return nil, err
	}
	return NewShowQuery(body)
}

func (s *ShowQuery) GetClosestResult(name string) *ShowQueryResult {
	var bestResult = s.Results[0]
	var bestJaroDist float64
	normalizedName := normalizeName(name)
	for _, r := range s.Results {
		jaroDist := smetrics.JaroWinkler(normalizedName, normalizeName(r.Name), 0.7, 4)
		if jaroDist > bestJaroDist {
			bestJaroDist = jaroDist
			bestResult = r
		}
	}
	return &bestResult
}

func (s *ShowQuery) String() string {
	res, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(res)
}

func normalizeName(n string) string {
	return strings.TrimPrefix(
		strings.ToLower(
			strings.TrimSpace(
				n,
			),
		),
		"the ",
	)
}
