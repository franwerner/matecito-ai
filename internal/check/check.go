package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Status int

const (
	StatusOK Status = iota
	StatusMissing
	StatusOutdated
)

type Result struct {
	Name     string
	Required bool
	Status   Status
	Version  string
	Detail   string
	FixHint  string
}

func RunVersion(name, bin string, args []string, required bool, fixHint string) Result {
	r := Result{Name: name, Required: required}
	if _, err := exec.LookPath(bin); err != nil {
		r.Status = StatusMissing
		r.Detail = "no encontrado en PATH"
		r.FixHint = fixHint
		return r
	}
	out, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil {
		r.Status = StatusMissing
		r.Detail = strings.TrimSpace(string(out))
		r.FixHint = fixHint
		return r
	}
	r.Status = StatusOK
	r.Version = ParseVersion(string(out))
	return r
}

var versionRe = regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`)

func ParseVersion(s string) string {
	s = strings.TrimSpace(s)
	if m := versionRe.FindString(s); m != "" {
		return m
	}
	return strings.SplitN(s, "\n", 2)[0]
}

func ParseMajor(v string) (int, bool) {
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	parts := strings.SplitN(v, ".", 2)
	if len(parts) == 0 || parts[0] == "" {
		return 0, false
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, false
	}
	return n, true
}

func ParseMajorMinor(v string) (int, int, bool) {
	parts := strings.Split(v, ".")
	if len(parts) < 2 {
		return 0, 0, false
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, false
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, false
	}
	return major, minor, true
}

func Outdated(name string, required bool, version, requiredVersion string) Result {
	return Result{
		Name:     name,
		Required: required,
		Status:   StatusOutdated,
		Version:  version,
		Detail:   fmt.Sprintf("%s (se requiere ≥ %s)", version, requiredVersion),
	}
}
