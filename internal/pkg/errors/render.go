package errors

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mitchellh/copystructure"
	"github.com/texttheater/golang-levenshtein/levenshtein"
)

const (
	reExecError string = `(?m)template: (?P<Template>.*): executing "(?P<Executing>.*)" at <(?P<At>.*)>: (?P<Err>.*)`
)

// ErrTemplateExecErrors provide an error that we can use to add an As method to
// for providing more robust error handling given the ability to scrape the output
// pattern
type ErrRenderError struct {
	err  error
	vars map[string]interface{}
}

func NewErrRenderError(err error, vars map[string]interface{}) *ErrRenderError {
	newVars, _ := copystructure.Copy(vars)
	return &ErrRenderError{
		err:  err,
		vars: newVars.(map[string]interface{}),
	}
}

func (e ErrRenderError) Error() string {
	return e.err.Error()
}

func (e ErrRenderError) As(target interface{}) bool {
	switch target := target.(type) {
	case *ErrTemplateExecError:
		var re = regexp.MustCompile(reExecError)
		match := re.FindStringSubmatch(e.err.Error())

		if len(match) == 0 {
			return false
		}

		target.Err = e.err
		target.Template = match[re.SubexpIndex("Template")]
		target.Executing = match[re.SubexpIndex("Executing")]
		target.At = match[re.SubexpIndex("At")]
		target.ErrText = match[re.SubexpIndex("Err")]
		target.vars = e.vars
		return true

	case *ErrTemplateNilPointerError:
		var re = regexp.MustCompile(reExecError)
		match := re.FindStringSubmatch(e.err.Error())

		if len(match) == 0 || !strings.Contains(match[re.SubexpIndex("Err")], "nil pointer") {
			return false
		}

		target.Err = e.err
		target.Template = match[re.SubexpIndex("Template")]
		target.Executing = match[re.SubexpIndex("Executing")]
		target.At = match[re.SubexpIndex("At")]
		target.ErrText = match[re.SubexpIndex("Err")]
		target.vars = e.vars
		return true

	default:
		return false
	}
}

type ErrTemplateExecError struct {
	Err       error
	Template  string
	Executing string
	At        string
	ErrText   string
	vars      map[string]interface{}
}

type ErrTemplateNilPointerError struct {
	ErrTemplateExecError
}

func (tnp ErrTemplateExecError) Error() string {
	nilVarName := strings.TrimSuffix(tnp.At, fmt.Sprintf(".%s", tnp.Executing))
	return fmt.Sprintf("error rendering %q in %q (%s): %s", nilVarName, tnp.At, tnp.Template, tnp.ErrText)
}

func (tnp ErrTemplateNilPointerError) Error() string {
	nilVarName := strings.TrimSuffix(tnp.At, fmt.Sprintf(".%s", tnp.Executing))
	return fmt.Sprintf("nil pointer evaluating %q in %q.(%s).%s%s", nilVarName, tnp.At, tnp.Template, tnp.didYouMean(tnp.At), tnp.validValues())
}

func (tnp ErrTemplateExecError) validValues() string {
	keys := getAllKeys(tnp.vars)
	for i, k := range keys {
		keys[i] = fmt.Sprintf("     %q", k)
	}
	return fmt.Sprintf("\n\n Valid values are:\n%s", strings.Join(keys, "\n"))
}

func (tnp ErrTemplateExecError) didYouMean(varName string) string {
	dist := 5
	out := ""
	keys := getAllKeys(tnp.vars)
	for _, key := range keys {
		keyDist := levenshtein.DistanceForStrings([]rune(varName), []rune(key), levenshtein.DefaultOptions)
		if keyDist < dist {
			dist = keyDist
			out = key
		}
	}

	if out != "" {
		return fmt.Sprintf("\n\n Did you mean %q?", out)
	}
	return ""
}

func getAllKeys(inMap map[string]interface{}) []string {
	out := make([]string, 0)
	getKeys(inMap, "", &out)
	return out
}

func getKeys(inMap map[string]interface{}, path string, collector *[]string) {
	for k, v := range inMap {
		*collector = append(*collector, path+"."+k)
		switch x := v.(type) {
		case map[string]interface{}:
			getKeys(x, path+"."+k, collector)
		default:
		}
	}
}
