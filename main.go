package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	// e.g: %s: RocketChat/Rocket.Chat.Electron
	URL = "https://api.github.com/repos/%s/releases/latest"
)

var (
	o = struct {
		compare *string
		debug   *bool
		format  *string
	}{
		compare: flag.String("compare", "", "`filename` to compare the output with"),
		debug:   flag.Bool("debug", false, "switch on debugging (print the JSON and exit cleanly)"),
		format:  flag.String("format", "+default", "the `format` of the output"),
	}
	// pre-defined formats
	format = map[string]string{
		"+default": `{{.Tag_name}} {{.Name}} {{.Tarball_url}}{{"\n"}}`,
		"+assets":  `{{range .Assets}}{{.Browser_download_url}}{{"\n"}}{{end}}`,
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
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] ORG/PROJECT\n", filepath.Base(os.Args[0]))
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
        .Deb_url
      Predefined formats:
      - +default %s
      - +assets  %s

    Projects for testing: HeikoSchlittermann/github-release-check
                          go-gitea/gitea
                          Exim/exim
`, format["+default"], format["+assets"])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	var project = flag.Arg(0)

	// check, we have one of the pre-defined formats
	if v, ok := format[*o.format]; ok {
		*o.format = v
	}

	var tt = template.New("output")
	if _, err := tt.Parse(*o.format); err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(fmt.Sprintf(URL, project))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if *o.debug {
		jq(response.Body)
		os.Exit(0)
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

	out := new(strings.Builder)
	tt.Execute(out, latest)

	// Make sure we have a "\n" at the end
	if !strings.HasSuffix(out.String(), "\n") {
		out.WriteByte('\n')
	}

	if *o.compare == "" {
		fmt.Print(out)
	} else {
		var line string
		if fh, err := os.Open(*o.compare); err != nil {
			log.Fatal(err)
		} else {
			defer fh.Close()
			line, err = bufio.NewReader(fh).ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
		}
		if out.String() == line {
			log.Print("[✓] ", out)
			os.Exit(0)
		} else {
			log.Print("[ ] ", line)
			log.Print("[✗] ", out)
			os.Exit(1)
		}
	}
}
