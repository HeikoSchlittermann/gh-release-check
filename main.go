package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	//URL = "https://api.github.com/repos/RocketChat/Rocket.Chat/releases/latest"
	//URL = "https://api.github.com/repos/RocketChat/Rocket.Chat.Electron/releases/latest"
	URL = "https://api.github.com/repos/%s/releases/latest"
)

var (
	o = struct {
		debug   *bool
		format  *string
		project *string
	}{
		debug:   flag.Bool("debug", false, "switch on debugging (print the JSON and exit cleanly)"),
		format:  flag.String("format", "+default", "the `format` of the output"),
		project: flag.String("project", "HeikoSchlittermann/github-release-check", "the `organization/package` to check the version"),
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
	  .Deb_url
	Predefined formats:
	- +default %s
	- +assets  %s
`, format["+default"], format["+assets"])
	}
	flag.Parse()

	// check, we have one of the pre-defined formats
	if v, ok := format[*o.format]; ok {
		*o.format = v
	}

	var tt = template.New("output")
	if _, err := tt.Parse(*o.format); err != nil {
		log.Fatal(err)
	}

	response, err := http.Get(fmt.Sprintf(URL, *o.project))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if *o.debug {
		jq := exec.Command("jq", ".")
		jq.Stdout = os.Stdout
		pipe, err := jq.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := jq.Start(); err != nil {
			log.Println("the attempt to start jq returned:", err)
			pipe = os.Stdout
		}
		io.Copy(pipe, response.Body)
		pipe.Close()
		_ = jq.Wait()
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
