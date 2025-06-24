package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

var (
	Version string
	Commit  string
	Date    string
)

type Build struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"go_version"`
	GoOS      string `json:"go_os"`
	GoArch    string `json:"go_arch"`
}

func (b *Build) String() string {
	return fmt.Sprintf("%+v", *b)
}

func (b *Build) JSON() string {
	m, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling build info: %v", err)
	}

	return string(m)
}

func (b *Build) PrettyPrint() {
	fmt.Printf("Version:\t\t%s\n", b.Version)
	fmt.Printf("Commit:\t\t\t%s\n", b.Commit)
	fmt.Printf("Date:\t\t\t%s\n", b.Date)
	fmt.Println()
	fmt.Printf("Go Version:\t\t%s\n", b.GoVersion)
	fmt.Printf("Go OS:\t\t\t%s\n", b.GoOS)
	fmt.Printf("Go Arch:\t\t%s\n", b.GoArch)
}

func GetBuild() *Build {
	return &Build{
		Version:   Version,
		Commit:    Commit,
		Date:      Date,
		GoVersion: runtime.Version(),
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
	}
}

func init() {
	if Version == "" {
		Version = "dev"
	}

	if Commit == "" {
		Commit = "unknown"
	}

	if Date == "" {
		Date = "unknown"
	}
}
