package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	EXHIBIT_PREFIX = "exhibit."
	SIGN_PREFIX    = "sign."
	VALUE_SUFFIX   = ".value"

	VARIABLE_REGEX = regexp.MustCompile(`{{\s*([\w\._-]+)\s*}}`)
)

type Template struct {
	Vars     []string
	Exhibits []string
	Signing  []string
}

// return list of all variables and exhibits found in contract template
func parseTemplate(b []byte) Template {
	var template Template
	matches := VARIABLE_REGEX.FindAllStringSubmatch(string(b), -1)
	for _, m := range matches {
		match := m[1]

		if strings.HasPrefix(match, EXHIBIT_PREFIX) {
			exhibitName := strings.TrimPrefix(match, EXHIBIT_PREFIX)
			// only append if it doesn't have .value
			if !strings.HasSuffix(exhibitName, VALUE_SUFFIX) {
				template.Exhibits = appendNew(template.Exhibits, exhibitName)
			}
		} else if strings.HasPrefix(match, SIGN_PREFIX) {
			signName := strings.TrimPrefix(match, SIGN_PREFIX)
			template.Signing = appendNew(template.Signing, signName)
		} else {
			template.Vars = appendNew(template.Vars, match)
		}
	}
	return template
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
	tmpl := parseTemplate(b)

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
	for _, v := range tmpl.Vars {
		paramsFile += fmt.Sprintf(`%s = ""`, v)
		paramsFile += "\n"
	}
	paramsFile += "\n\n"

	paramsFile += "[exhibit]\n"
	// write the exhibits
	for _, e := range tmpl.Exhibits {
		paramsFile += fmt.Sprintf(`%s = ""`, e)
		paramsFile += "\n"
	}
	paramsFile += "\n\n"

	paramsFile += "[sign]\n"
	// write the signing values
	for _, e := range tmpl.Signing {
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

	tmpl := parseTemplate(b)
	exhibits := tmpl.Exhibits

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
		} else if strings.HasPrefix(paramName, SIGN_PREFIX) {
			signName := strings.TrimPrefix(paramName, SIGN_PREFIX)
			paramVal := params.GetString("sign." + signName)
			if paramVal != "" {
				return []byte(paramVal)
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

	switch outputType {
	case "md":
		if err := ioutil.WriteFile(path.Join(name, "contract.md"), markdownOutput, 0600); err != nil {
			return err
		}
	case "html":
		htmlOutput := markdown2html(markdownOutput)
		if err := ioutil.WriteFile(path.Join(name, "contract.html"), htmlOutput, 0600); err != nil {
			return err
		}
	case "pdf":
		// requires the md to be written
		mdPath := path.Join(name, "contract.md")
		if err := ioutil.WriteFile(mdPath, markdownOutput, 0600); err != nil {
			return err
		}
		cmd := exec.Command("pandoc", mdPath, "--latex-engine=xelatex", "-o", path.Join(name, "contract.pdf"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	default:
		return fmt.Errorf("Unknown output format: %v", outputType)
	}

	return nil
}
