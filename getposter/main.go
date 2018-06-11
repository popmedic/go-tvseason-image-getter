package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	img "image"
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
)

const (
	exitCodeShowNotSet                       = 10
	exitCodeUnableToGetConfig                = 20
	exitCodeUnableToQueryShows               = 30
	exitCodeNoMatchingShowResults            = 40
	exitCodeUnableToGetSeason                = 50
	exitCodeNoPoster                         = 60
	exitCodeDownloadError                    = 70
	exitCodeUnableToCreateOutFile            = 80
	exitCodeUnableToDecodePosterImage        = 90
	exitCodeUnableToEncodeResizedPosterImage = 100
	exitCodeUnableToCopyPosterImage          = 110
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
	// Make sure that the showName flag is set
	if len(*showName) == 0 {
		flag.Usage()
		log.Fatal(func(int) { os.Exit(exitCodeShowNotSet) }, "show must be set")
	}

	if isShow() {
		log.Infof("get poster for show: %q", *showName)
	} else {
		log.Infof("get poster for show: %q season %d", *showName, *seasonNumber)
	}

	cfg, err := tmdb.GetConfig(tmdb.HttpGetter(http.Get))
	if err != nil {
		flag.Usage()
		log.Fatal(func(int) { os.Exit(exitCodeUnableToGetConfig) }, err)
	}

	showQuery, err := tmdb.QueryShows(*showName, tmdb.HttpGetter(http.Get))
	if err != nil {
		flag.Usage()
		log.Fatal(func(int) { os.Exit(exitCodeUnableToQueryShows) }, err)
	}

	if len(showQuery.Results) <= 0 {
		log.Fatalf(func(int) { os.Exit(exitCodeNoMatchingShowResults) }, "no show results matching %q", *showName)
	}

	showResult := showQuery.GetClosestResult(*showName)
	// See if out flag is passed in, if not set the default (Season X-SD.jpg or showName-SD.jpg)
	if len(*out) == 0 {
		if isShow() {
			*out = fmt.Sprintf("%s-SD.jpg", showResult.Name)
		} else {
			*out = fmt.Sprintf("Season %d-SD.jpg", *seasonNumber)
		}
	}
	// Get the URL to download poster
	var url string
	if isShow() {
		if len(showResult.PosterPath) == 0 {
			log.Fatal(func(int) { os.Exit(exitCodeNoPoster) }, fmt.Sprintf("%q does not have a poster", showResult.Name))
		}
		url = cfg.Images.SecureBaseUrl + "original" + showResult.PosterPath
	} else {
		season, err := tmdb.GetSeason(showResult.ID, *seasonNumber, tmdb.HttpGetter(http.Get))
		if err != nil {
			flag.Usage()
			log.Fatal(func(int) { os.Exit(exitCodeUnableToGetSeason) }, err)
		}
		if len(season.PosterPath) == 0 {
			log.Fatal(func(int) { os.Exit(exitCodeNoPoster) }, fmt.Sprintf("%q season %d does not have a poster", *showName, *seasonNumber))
		}
		url = cfg.Images.SecureBaseUrl + "original" + season.PosterPath
	}

	log.Infof("downloading poster %q to %q", url, *out)

	data, err := download(url)
	if err != nil {
		log.Fatal(func(int) { os.Exit(exitCodeDownloadError) }, err)
	}
	// Create the out file
	outf, err := os.Create(*out)
	if err != nil {
		log.Fatal(func(int) { os.Exit(exitCodeUnableToCreateOutFile) }, err)
	}
	defer outf.Close()
	// Set the resize width and height if given
	if *width > 0 && *height > 0 {
		image, _, err := img.Decode(data)
		if nil != err {
			log.Fatal(func(int) { os.Exit(exitCodeUnableToDecodePosterImage) }, err)
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
			log.Fatal(func(int) { os.Exit(exitCodeUnableToEncodeResizedPosterImage) }, err)
		}
	} else {
		_, err = io.Copy(outf, data)
		if err != nil {
			log.Fatal(func(int) { os.Exit(exitCodeUnableToCopyPosterImage) }, err)
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
