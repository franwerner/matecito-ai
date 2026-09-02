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

// notFoundInPATHDetail is RunVersion's Detail when bin cannot be resolved on
// the running session's PATH. Hoisted into a const so ProbeAt can substitute
// a detail naming the resolved path instead, without RunVersion's own
// behavior changing.
const notFoundInPATHDetail = "no encontrado en PATH"

func RunVersion(name, bin string, args []string, required bool, fixHint string) Result {
	r := Result{Name: name, Required: required}
	if _, err := exec.LookPath(bin); err != nil {
		r.Status = StatusMissing
		r.Detail = notFoundInPATHDetail
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

// ProbeAt reports whether name is installed and healthy at the absolute path
// its own step installs it to. resolve returns that absolute path; when it
// errors the binary is reported Missing — a step that cannot establish where
// its binary would live has no honest way to call it installed. RunVersion
// receives an absolute bin, which Go executes directly without a PATH search,
// so a destination the running session's PATH has not picked up yet still
// counts as installed. Never derives StatusOutdated: neither caller computes a
// minimum version here.
func ProbeAt(name string, resolve func() (string, error), args []string, required bool, fixHint string) Result {
	path, err := resolve()
	if err != nil {
		return Result{
			Name:     name,
			Required: required,
			Status:   StatusMissing,
			Detail:   err.Error(),
			FixHint:  fixHint,
		}
	}
	r := RunVersion(name, path, args, required, fixHint)
	if r.Status == StatusOK {
		// Name the resolved canonical path instead of RunVersion's blank Detail
		// on the OK branch — the point of probing at a location instead of by
		// bare name is that the answer is traceable to where it was found.
		r.Detail = path
	}
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
