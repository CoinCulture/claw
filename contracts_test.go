package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// NOTE: see bottom of file for the additional constants used in these tests

// is called engagementPath in func newEngagement() but that's confusing IMO
const engagementName = "alice-test"
const sampleContract = "sampleContractTemplate.md"

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

// the first step [claw new] requires a sampleContract.md
// that is properly formatted for templating
// and generates a [params.toml] file consumed by the second step
func TestWriteParamsFileFromContractTemplate(t *testing.T) {

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

func TestOutputTypeMarkdown(t *testing.T) {
	const outputType = "md"

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

	// a params.toml and a template.md was written by the previous func
	// the former would normally be manually edited prior to compiling;
	// the latter is equivalent to the  markdown file that was passed in
	// on [claw new] when  generating the params.toml. instead, we're
	// going to over-write the params.toml using a mock with filed in fields
	// and test it against the compiled markdown of a contract we'd expect

	if err := os.Remove(filepath.Join(engagementName, "params.toml")); err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	if err := ioutil.WriteFile(filepath.Join(engagementName, "params.toml"), []byte(filledOutParamsToml), 0600); err != nil {
		t.Fatalf("Error: %v\n", err)
	}

	// the function we're actually testing
	if err := generateContract(engagementName, outputType); err != nil {
		t.Fatalf("Error: %v\n", err)
	}

}

func TestOutputTypePDF(t *testing.T) {
}

func TestOutputTypeHTML(t *testing.T) {
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

# {{ .Exhibit.services}}

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
template = "7E1C6EC0F1D68D1E00410C8F4DDBD99913FE23AA9C75FA5D299AE823B05DE65E"

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

const filledOutParamsToml = `# This is a TOML file containing parameters for this contract

[meta]
# This must match the hash of the local template.md file. DO NOT CHANGE IT
template = "7E1C6EC0F1D68D1E00410C8F4DDBD99913FE23AA9C75FA5D299AE823B05DE65E"

[var]
date = "2017-05-04"
consultant = "John Smith"
schedule = "Full Time"
start-date = "2017-06-05"
email = "john@smith.com"


[exhibit]
services = "Software Development"
compensation = "$100/hr"
expenses = "$200/month"


[sign]
image = "examples/franklin.png"
company-signer = "Ben Franklin, President, bf@usa.gov"
`
