package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"github.com/popmedic/go-logger/log"
	"github.com/popmedic/go-tvseason-image-getter/tmdb"
	"github.com/xrash/smetrics"
)

var showName = flag.String("show", "", "the show name ** required **")
var seasonNumber = flag.Int("season", -1, "the season number")
var out = flag.String("out", "", "file to output, defaults to \"Season <season>-SD.jpg\", or \"<show>-SD.jpg\" if no season is given.")
var height = flag.Int("h", 0, "set the height of the downloaded poster")
var width = flag.Int("w", 0, "set the width of the downloaded poster")

func download(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

func isShow() bool {
	return *seasonNumber < 0
}

func main() {
	flag.Parse()

	if len(*showName) == 0 {
		flag.Usage()
		log.Fatal(func(int) { os.Exit(1) }, "show must be set")
	}

	if isShow() {
		log.Infof("%q show only.", *showName)
	} else {
		log.Infof("%q : Season %d", *showName, *seasonNumber)
	}

	cfg, err := tmdb.GetConfig(tmdb.HttpGetter(http.Get))
	if err != nil {
		flag.Usage()
		log.Fatal(func(int) { os.Exit(2) }, err)
	}

	log.Infof("image base_url = %q", cfg.Images.SecureBaseUrl)

	showQuery, err := tmdb.QueryShows(*showName, tmdb.HttpGetter(http.Get))
	if err != nil {
		flag.Usage()
		log.Fatal(func(int) { os.Exit(3) }, err)
	}

	var bestResult = 0
	var bestJaroDist float64
	if len(showQuery.Results) <= 0 {
		log.Fatalf(func(int) { os.Exit(3) }, "no show results matching %q", *showName)
	} else {
		normShowName := strings.TrimPrefix(strings.ToLower(strings.TrimSpace(*showName)), "the ")
		for i, res := range showQuery.Results {
			jaroDist := smetrics.JaroWinkler(normShowName,
				strings.TrimPrefix(strings.ToLower(strings.TrimSpace(res.Name)), "the "), 0.7, 4)
			if jaroDist > bestJaroDist {
				bestJaroDist = jaroDist
				bestResult = i
			}
		}
	}

	if len(*out) == 0 {
		if isShow() {
			*out = fmt.Sprintf("%s-SD.jpg", showQuery.Results[bestResult].Name)
		} else {
			*out = fmt.Sprintf("Season %d-SD.jpg", *seasonNumber)
		}
	}

	var url string
	if isShow() {
		if len(showQuery.Results[bestResult].PosterPath) == 0 {
			log.Fatal(func(int) { os.Exit(5) },
				fmt.Sprintf("%q (%d) does not have a poster", showQuery.Results[bestResult].Name, bestResult))
		}
		url = cfg.Images.SecureBaseUrl + "original" + showQuery.Results[bestResult].PosterPath
	} else {
		showID := showQuery.Results[bestResult].ID

		log.Info("show id = ", showID)

		season, err := tmdb.GetSeason(showID, *seasonNumber, tmdb.HttpGetter(http.Get))
		if err != nil {
			flag.Usage()
			log.Fatal(func(int) { os.Exit(4) }, err)
		}
		if len(season.PosterPath) == 0 {
			log.Fatal(func(int) { os.Exit(5) },
				fmt.Sprintf("%q season %d does not have a poster", *showName, *seasonNumber))
		}
		url = cfg.Images.SecureBaseUrl + "original" + season.PosterPath
	}
	log.Info("poster path = ", url)

	data, err := download(url)
	if err != nil {
		log.Fatal(func(int) { os.Exit(6) }, err)
	}

	outf, err := os.Create(*out)
	if err != nil {
		log.Fatal(func(int) { os.Exit(7) }, err)
	}
	defer outf.Close()

	if *width > 0 && *height > 0 {
		image, _, err := image.Decode(data)
		if nil != err {
			log.Fatal(os.Exit, err)
		}
		newImage := resize.Resize(uint(*width), uint(*height), image, resize.Lanczos3)
		switch strings.ToLower(filepath.Ext(*out)) {
		case ".jpg", ".jpeg":
			err = jpeg.Encode(outf, newImage, nil)
		case ".png":
			err = png.Encode(outf, newImage)
		default:
			err = errors.New("unknown type \"" + strings.ToUpper(filepath.Ext(*out)) + "\"")
		}
		if nil != err {
			log.Fatal(func(int) { os.Exit(8) }, err)
		}
	} else {
		_, err = io.Copy(outf, data)
		if err != nil {
			log.Fatal(func(int) { os.Exit(9) }, err)
		}
	}

	if isShow() {
		log.Infof("successfully downloaded poster for %q to %q",
			*showName, *out)
	} else {
		log.Infof("successfully downloaded poster for %q season %d to %q",
			*showName, *seasonNumber, *out)
	}
}
