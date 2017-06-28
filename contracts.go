package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
)

var (
	EXHIBIT_PREFIX = "exhibit."
	SIGN_PREFIX    = "sign."
	VALUE_SUFFIX   = ".value"

	VARIABLE_REGEX = regexp.MustCompile(`{{\s*([\w\._-]+)\s*}}`)
)

type ContractTemplate struct {
	TemplateHash string
	Vars         []string
	Exhibits     []string
	Signing      []string
}

// return list of all variables and exhibits found in contract template
func parseTemplate(b []byte) ContractTemplate {
	var contractTemplate ContractTemplate
	matches := VARIABLE_REGEX.FindAllStringSubmatch(string(b), -1)
	for _, m := range matches {
		match := m[1]

		if strings.HasPrefix(match, EXHIBIT_PREFIX) {
			exhibitName := strings.TrimPrefix(match, EXHIBIT_PREFIX)
			// only append if it doesn't have .value
			if !strings.HasSuffix(exhibitName, VALUE_SUFFIX) {
				contractTemplate.Exhibits = appendNew(contractTemplate.Exhibits, exhibitName)
			}
		} else if strings.HasPrefix(match, SIGN_PREFIX) {
			signName := strings.TrimPrefix(match, SIGN_PREFIX)
			contractTemplate.Signing = appendNew(contractTemplate.Signing, signName)
		} else {
			contractTemplate.Vars = appendNew(contractTemplate.Vars, match)
		}
	}
	return contractTemplate
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

	// get list of variables for a blank config file
	tmpl := parseTemplate(b)

	// compute hash of template
	h := sha256.New()
	h.Write(b)
	templateHash := h.Sum(nil)
	tmpl.TemplateHash = fmt.Sprintf("%X", templateHash)

	// make the directory and copy over the template
	if err := os.MkdirAll(engagementPath, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Join(engagementPath, "template.md"), b, 0600); err != nil {
		return err
	}

	paramsFile := generateParamsFile(tmpl)

	// write the params file
	return ioutil.WriteFile(filepath.Join(engagementPath, "params.toml"), paramsFile, 0600)
}

func generateParamsFile(tmpl ContractTemplate) []byte {
	var paramsFileTemplate *template.Template
	var err error
	var buffer bytes.Buffer

	paramsFileTemplate, err = template.New("paramsFile").Parse(paramsFileDefault)
	if err != nil {
		panic(err)
	}

	if err := paramsFileTemplate.Execute(&buffer, tmpl); err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

const paramsFileDefault = `# This is a TOML file containing parameters for this contract

[meta]
# This must match the hash of the local template.md file. DO NOT CHANGE IT
template = "{{ .TemplateHash}}"

[var]
{{range .Vars}}{{.}} = ""
{{end}}

[exhibit]
{{range .Exhibits}}{{.}} = ""
{{end}}

[sign]
{{range .Signing}}{{.}} = ""
{{end}}
`

//-----------------------------------------

func compileContract(cmd *cobra.Command, args []string) error {

	// load the params from toml file
	params, err := loadConfig()
	if err != nil {
		return err
	}

	// read the contract template
	b, err := ioutil.ReadFile("template.md")
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
		if err := ioutil.WriteFile("contract.md", markdownOutput, 0600); err != nil {
			return err
		}
	case "html":
		htmlOutput := markdown2html(markdownOutput)
		if err := ioutil.WriteFile("contract.html", htmlOutput, 0600); err != nil {
			return err
		}
	case "pdf":
		// requires the md to be written
		mdPath := "contract.md"
		if err := ioutil.WriteFile(mdPath, markdownOutput, 0600); err != nil {
			return err
		}
		cmd := exec.Command("pandoc", mdPath, "--latex-engine=xelatex", "-o", "contract.pdf")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	default:
		return fmt.Errorf("Unknown output format: %v", outputType)
	}

	return nil
}

func reviseContract(cmd *cobra.Command, args []string) error {

	// read the params file
	params_b, err := ioutil.ReadFile("params.toml")
	if err != nil {
		return err
	}

	// read the contract template
	template_b, err := ioutil.ReadFile("template.md")
	if err != nil {
		return err
	}

	// open history file (or create history file if it doesn't exist)
	file, err := os.OpenFile("history.toml", os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err = file.WriteString("// Do not modify this file yourself under any circumstances!"); err != nil {
		return err
	}

	// read history file
	history_b, err := ioutil.ReadFile("history.toml")
	if err != nil {
		return err
	}

	// combine params, history, and template data
	all_b := [3][]byte{params_b, template_b, history_b}
	byte_array := make([]byte, 3, 3)

	for _, element := range all_b {
		for _, b := range element {
			byte_array = append(byte_array, b)
		}
	}

	// hash params, history, and template data
	h := sha256.New()
	h.Write(byte_array)
	t := time.Now()
	hashtime := fmt.Sprintf("\n%s: '%X'", t, h.Sum(nil))
	fmt.Println("A hash has been added to your history file; your changes are secure.")

	// write hash to history file
	hfile, err := os.OpenFile("history.toml", os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return err
	}

	defer hfile.Close()

	if _, err = hfile.WriteString(hashtime); err != nil {
		return err
	}

	return nil

}
