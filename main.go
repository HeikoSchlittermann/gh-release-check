package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	URL = "https://api.github.com/repos/RocketChat/Rocket.Chat.Electron/releases/latest"
)

var (
	o = struct {
		debug  *bool
		format *string
	}{
		debug:  flag.Bool("debug", false, "switch on debugging"),
		format: flag.String("format", "{{.Name}} {{.Tarball_url}}\n", "the `format` of the output"),
	}
)

type latest struct {
	Tarball_url  string
	Published_at time.Time
	Name         string
	Tag_name     string
	Assets       []struct {
		Name                 string
		Browser_download_url string
		Content_type         string
		Updated_at           time.Time
	}
	Deb_url string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), `
    Format string elements:
      .Name             Name of the latest release
      .Tag_name         Name of the release tag
      .Tarball_url      URL of the Tarball for download
      .Published_date   Date of publication
      .Assets []
        .Name
        .Browser_downlaod_url
        .Content_type
        .Updated_at
`)
	}
	flag.Parse()

	var tt = template.New("output")
	if _, err := tt.Parse(*o.format); err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if *o.debug {
		io.Copy(os.Stdout, response.Body)
		return
	}

	var latest latest
	if err := json.NewDecoder(response.Body).Decode(&latest); err != nil {
		log.Fatal(err)
	}

	for _, a := range latest.Assets {
		if strings.HasSuffix(a.Name, ".deb") {
			latest.Deb_url = a.Browser_download_url
			break
		}
	}

	tt.Execute(os.Stdout, latest)

}
