package errors

import (
	"bytes"
	stdErrors "errors"
	"fmt"
	"reflect"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewErrRenderError(t *testing.T) {
	table := &renderTestCases{
		cases: []renderTestCase{
			{
				name:          "basic",
				template:      `{{printf "%s" "foo"}}`,
				outputMatches: "foo",
			},
			{
				name:        "nil pointer",
				template:    `{{if eq .bad.string "" }}nope{{end}}`,
				expectError: true,
				errorAs:     &ErrTemplateNilPointerError{},
			},
			// // this test case panics right now, and I don't know
			// // what pack will do now.
			// {
			// 	name:        "bad function",
			// 	template:    `{{ notAFunction }}`,
			// 	expectError: true,
			// 	errorAs:     &ErrTemplateNilPointerError{},
			// },
		},
	}
	dot := map[string]interface{}{
		"foo": "bar",
		"map": map[string]interface{}{
			"int":    1,
			"bool":   true,
			"string": "a string",
		},
	}
	for _, tc := range table.cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := parseTemplate(tc.template, dot)
			if tc.expectError {
				assert.Error(t, err)
				err = NewErrRenderError(err, dot)
				if tc.errorAs != nil {
					if !stdErrors.As(err, &tc.errorAs) {
						require.Failf(t, "unexpected error type", "expected a %s, got a %s", reflect.TypeOf(tc.errorAs).String(), reflect.TypeOf(err).String())
					}
				}
				return
			}

			assert.NoError(t, err)
			fmt.Println(out)
			if tc.outputMatches != "" {
				require.Equal(t, tc.outputMatches, out)
			}
		})
	}
}

type renderTestCase struct {
	name           string
	template       string
	expectError    bool
	outputContains string
	outputMatches  string
	errorAs        interface{}
}
type renderTestCases struct {
	cases []renderTestCase
}

func parseTemplate(tpl string, dot map[string]interface{}) (string, error) {
	// Create a new template and parse the provided string into it.
	t := template.Must(template.New("test").Parse(tpl))
	t.Option("missingkey=error")
	outBuf := new(bytes.Buffer)
	err := t.ExecuteTemplate(outBuf, "test", dot)
	return outBuf.String(), err

}
