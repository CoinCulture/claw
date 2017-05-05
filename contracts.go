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
)

var (
	VAR_PREFIX     = ".Var."
	EXHIBIT_PREFIX = ".Exhibit."
	SIGN_PREFIX    = ".Sign."
	// [zr] not sure what to do with suffix ...
	VALUE_SUFFIX = ".value"

	VARIABLE_REGEX = regexp.MustCompile(`{{\s*([\w\._-]+)\s*}}`)
)

// for reading a contract template
type ContractTemplate struct {
	TemplateHash string
	Vars         []string
	Exhibits     []string
	Signing      []string
}

type Var struct {
	Date       string `toml:"Date"`
	Consultant string `toml:"Consultant"`
	Schedule   string `toml:"Schedule"`
	StartDate  string `toml:"StartDate"`
	Email      string `toml:"Email"`
}

type Exhibit struct {
	Services     string `toml:"Services"`
	Compensation string `toml:"Compensation"`
	Expenses     string `toml:"Expenses"`
}

type Sign struct {
	Image         string `toml:"Image"`
	CompanySigner string `toml:"CompanySigner"` // TODO multiple signers
}

// for writing a contract template
// this will be ever growing ...
type WriteContractTemplate struct {
	Var     *Var
	Exhibit *Exhibit
	Sign    *Sign
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
		} else if strings.HasPrefix(match, VAR_PREFIX) {
			varName := strings.TrimPrefix(match, VAR_PREFIX)
			contractTemplate.Vars = appendNew(contractTemplate.Vars, varName)
		}
	}
	return contractTemplate
}

// copy the template into a new dir and instantiate params file with empty values
func newEngagement(engagementPath, tmplPath string) error {

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
{{end}}`

//-----------------------------------------

func generateContract(engagementName, outputType string) error {

	// load the params from toml file
	// [zr] this ought to return a struct that can be
	// passed into the template, or we can fill in the struct afterwards ... ?
	params, err := loadConfig(engagementName)
	if err != nil {
		return fmt.Errorf("wtf1:%v", err)
	}

	var contractValues *WriteContractTemplate

	if err := params.Unmarshal(&contractValues); err != nil {
		return fmt.Errorf("wtf2:%v", err)
		//return err
	}

	fmt.Printf("contract values %v", contractValues)

	// read the contract template
	contractTemplateBytes, err := ioutil.ReadFile(filepath.Join(engagementName, "template.md"))
	if err != nil {
		return fmt.Errorf("wtf3:%v", err)
		//return err
	}

	fmt.Println(string(contractTemplateBytes))

	var contractTemplate *template.Template
	var buffer bytes.Buffer

	contractTemplate, err = template.New("contract").Parse(string(contractTemplateBytes))
	if err != nil {
		return fmt.Errorf("wtf4:%v", err)
	}

	if err := contractTemplate.Execute(&buffer, *contractValues); err != nil {
		return fmt.Errorf("wtf5:%v", err)
	}
	fmt.Printf("string: %v", buffer.String())
	markdownOutput := buffer.Bytes()

	/*
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
		})*/

	// error if params is missing anything
	// [zr] not sure how this will get handled by new template format ...
	//if len(missingParams) > 0 {
	//	return fmt.Errorf("Missing params: %v", missingParams)
	//}

	switch outputType {
	case "md":
		if err := ioutil.WriteFile(filepath.Join(engagementName, "contract.md"), markdownOutput, 0600); err != nil {
			return err
		}
	case "html":
		htmlOutput := markdown2html(markdownOutput)
		if err := ioutil.WriteFile(filepath.Join(engagementName, "contract.html"), htmlOutput, 0600); err != nil {
			return err
		}
	case "pdf":
		// requires the md to be written
		mdPath := filepath.Join(engagementName, "contract.md") // doesn't get removed IIRC
		if err := ioutil.WriteFile(mdPath, markdownOutput, 0600); err != nil {
			return err
		}
		cmd := exec.Command("pandoc", mdPath, "--latex-engine=xelatex", "-o", filepath.Join(engagementName, "contract.pdf"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	default:
		return fmt.Errorf("Unknown output format: %v", outputType)
	}

	return nil
}
