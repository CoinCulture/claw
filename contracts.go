package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	EXHIBIT_PREFIX = "exhibit."
	VALUE_SUFFIX   = ".value"

	VARIABLE_REGEX = regexp.MustCompile(`{{\s*([\w\._-]+)\s*}}`)
)

// return list of all variables and exhibits found in contract template
func parseTemplate(b []byte) (vars []string, exhibits []string) {
	matches := VARIABLE_REGEX.FindAllStringSubmatch(string(b), -1)
	for _, m := range matches {
		match := m[1]

		if strings.HasPrefix(match, EXHIBIT_PREFIX) {
			exhibitName := strings.TrimPrefix(match, EXHIBIT_PREFIX)
			// only append if it doesn't have .value
			if !strings.HasSuffix(exhibitName, VALUE_SUFFIX) {
				exhibits = appendNew(exhibits, exhibitName)
			}
		} else {
			vars = appendNew(vars, match)
		}
	}

	return
}

// copy the template into a new dir and instantiate params file with empty values
func newContract(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("new expects two args: directory for new engagement and template path")
	}

	engagementPath, tmplPath := args[0], args[1]

	b, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return err
	}

	// compute hash of template
	h := sha256.New()
	h.Write(b)
	templateHash := h.Sum(nil)

	// Make the directory and copy over the template
	if err := os.MkdirAll(engagementPath, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(engagementPath, "template.md"), b, 0600); err != nil {
		return err
	}

	// get list of variables for a blank config file
	vars, exhibits := parseTemplate(b)

	// write the header comment
	paramsFile := "# This is a TOML file containing parameters for this contract\n\n"
	paramsFile += "\n\n"

	// write the template hash
	paramsFile += "[meta]\n"
	paramsFile += "# This must match the hash of the local template.md file. DO NOT CHANGE IT\n"
	paramsFile += fmt.Sprintf(`template = "%X"`, templateHash)
	paramsFile += "\n\n"

	// write the vars
	paramsFile += "[var]\n"
	for _, v := range vars {
		paramsFile += fmt.Sprintf(`%s = ""`, v)
		paramsFile += "\n"
	}
	paramsFile += "\n\n"

	paramsFile += "[exhibit]\n"
	// write the exhibits
	for _, e := range exhibits {
		paramsFile += fmt.Sprintf(`%s = ""`, e)
		paramsFile += "\n"
	}

	// write the params file
	return ioutil.WriteFile(path.Join(engagementPath, "params.toml"), []byte(paramsFile), 0600)
}

//-----------------------------------------

func compileContract(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("compile expects one arg: name")
	}

	name := args[0]

	// load the params from toml file
	params, err := loadConfig(name)
	if err != nil {
		return err
	}

	// read the contract template
	b, err := ioutil.ReadFile(path.Join(name, "template.md"))
	if err != nil {
		return err
	}

	_, exhibits := parseTemplate(b)

	// substitute params into template holes
	var missingParams []string
	markdownOutput := VARIABLE_REGEX.ReplaceAllFunc(b, func(in []byte) []byte {
		paramName := strings.TrimSuffix(strings.TrimPrefix(string(in), "{{"), "}}")

		// if its an exhibit, we replace it with the exhibit number.
		// if its a var, we replace it with its value
		if strings.HasPrefix(paramName, EXHIBIT_PREFIX) {
			exhibitName := strings.TrimPrefix(paramName, EXHIBIT_PREFIX)

			if strings.HasSuffix(exhibitName, VALUE_SUFFIX) {
				exhibitName = strings.TrimSuffix(exhibitName, VALUE_SUFFIX)
				for _, e := range exhibits {
					if exhibitName == e {
						exhibitValue := params.GetString("exhibit." + exhibitName)
						if exhibitValue != "" {
							return []byte(exhibitValue)
						}
					}
				}
			} else {
				for i, e := range exhibits {
					if exhibitName == e {
						return []byte(fmt.Sprintf("Exhibit %d", i+1))
					}
				}
			}
		} else {
			paramVal := params.GetString("var." + paramName)
			if paramVal != "" {
				return []byte(paramVal)
			}
		}

		missingParams = append(missingParams, paramName)
		return []byte("----")
	})

	// error if params is missing anything
	if len(missingParams) > 0 {
		return fmt.Errorf("Missing params: %v", missingParams)
	}

	// TODO: option to write html
	// htmlOutput := markdown2html(markdownOutput)

	if err := ioutil.WriteFile(path.Join(name, "contract.md"), markdownOutput, 0600); err != nil {
		return err
	}
	return nil
}
