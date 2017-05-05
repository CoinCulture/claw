package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// NOTE: see bottom of file for the constants used in these tests

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// the first step [claw new] requires a sampleContract.md
// that is properly formatted for templating
// and generates a [params.toml] file consumed by the second step
func TestWriteParamsFileFromContractTemplate(t *testing.T) {
	// is called engagementPath in func newEngagement() but that's confusing IMO
	const engagementName = "alice-test"
	const sampleContract = "sampleContractTemplate.md"

	defer os.RemoveAll(engagementName) // remove the dir
	defer os.Remove(sampleContract)    // remove the temp file
	// write the contract template
	if err := ioutil.WriteFile(sampleContract, []byte(sampleContractTemplate), 0600); err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	// run newEngagement [claw new]
	if err := newEngagement(engagementName, sampleContract); err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	// read in the written params.toml
	paramsFileBytes, err := ioutil.ReadFile(filepath.Join(engagementName, "params.toml"))
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	// check that params.toml is what it should be
	if !bytes.Equal(paramsFileBytes, []byte(sampleParamsOutput)) {
		t.Fatalf("Bad params.toml:\nGot: %s\nExpected: %s\n", string(paramsFileBytes), sampleParamsOutput)
	}

	// check that template.md matches contract template
	// this functionality we should revisit ...
	templateDotMDBytes, err := ioutil.ReadFile(filepath.Join(engagementName, "template.md"))
	if err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	if !bytes.Equal(templateDotMDBytes, []byte(sampleContractTemplate)) {
		t.Fatalf("Bad template.md:\nGot: %s\nExpected: %s\n", string(templateDotMDBytes), sampleContractTemplate)
	}
}

// after generating a params.toml from the sampleContract.md,
// the second step is [claw compile] which generates one of:
// a pdf, a markdown file, or an html file
func TestOutputPDFFromCompile(t *testing.T) {
}

func TestOutputMarkdownFromCompile(t *testing.T) {
}

func TestOutputHTMLFromCompile(t *testing.T) {
}

// -------------- test constants -------------------------
const sampleContractTemplate = `# My Corp Inc.
# CONSULTING AGREEMENT

This Consulting Agreement (this "Agreement") is made as of {{ .Var.date}}, by and between My Corp, Inc., a Delaware corporation (the "Company"), and {{ .Var.consultant}} ("Consultant").

# Consulting Relationship.  

During the term of this Agreement, Consultant will provide consulting services to the Company as described on {{ .Exhibit.services}} hereto (the "Services").  Consultant represents that Consultant is duly licensed (as applicable) and has the qualifications, the experience and the ability to properly perform the Services.  Consultant shall use Consultant’s best efforts to perform the Services such that the results are satisfactory to the Company.  {{ .Var.schedule}}, or updated with 14 days prior notice.

# Fees.  

As consideration for the Services to be provided by Consultant and other obligations, the Company shall pay to Consultant the amounts specified in {{ .Exhibit.compensation}} hereto at the times specified therein.

# Expenses.  

Consultant shall not be authorized to incur on behalf of the Company any expenses and will be responsible for all expenses incurred while performing the Services except as expressly specified in {{ .Exhibit.expenses}} hereto unless otherwise agreed to by the Company's CEO, which consent shall be evidenced in writing for any such expenses in excess of $0.00.  As a condition to receipt of reimbursement, Consultant shall be required to submit to the Company reasonable evidence that the amount involved was both reasonable and necessary to the Services provided under this Agreement.

# Term and Termination.  

Consultant shall serve as a consultant to the Company for a period commencing on {{ .Var.start-date}} and terminating on the earlier of (a) the date Consultant completes the provision of the Services to the Company under this Agreement, or (b) the date Consultant shall have been paid the maximum amount of consulting fees as provided in {{ .Exhibit.compensation}} hereto.

\pagebreak

# Signatures

## THE COMPANY

My Corp Inc.

\ ![Company Signature]({{ .Sign.image}})

---

By: {{ .Sign.company-signer}}


## CONSULTANT

{{ .Var.consultant}}

---

{{ .Var.email}}


\pagebreak

# {{exhibit.services}}

## DESCRIPTION OF CONSULTING SERVICES

{{ .Exhibit.services.value}}

\pagebreak

# {{ .Exhibit.compensation}}

## COMPENSATION

{{ .Exhibit.compensation.value}}

\pagebreak

# {{ .Exhibit.expenses}}

## ALLOWABLE EXPENSES

{{ .Exhibit.expenses.value}}
`

// sampleContractTemplate should
// generate exactly this file!
const sampleParamsOutput = `# This is a TOML file containing parameters for this contract

[meta]
# This must match the hash of the local template.md file. DO NOT CHANGE IT
template = "5E43F71626F4D7F2F79B3650E14295760EA190DB51DF6EC9EFF4EF145E2B255E"

[var]
date = ""
consultant = ""
schedule = ""
start-date = ""
email = ""


[exhibit]
services = ""
compensation = ""
expenses = ""


[sign]
image = ""
company-signer = ""
`

// TODO ^ change company-signer from "" to [] (string slice)
