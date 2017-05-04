package main

import (
	"os"
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

This Consulting Agreement (this "Agreement") is made as of {{date}}, by and between My Corp, Inc., a Delaware corporation (the "Company"), and {{consultant}} ("Consultant").

# Consulting Relationship.  

During the term of this Agreement, Consultant will provide consulting services to the Company as described on {{exhibit.services}} hereto (the "Services").  Consultant represents that Consultant is duly licensed (as applicable) and has the qualifications, the experience and the ability to properly perform the Services.  Consultant shall use Consultant’s best efforts to perform the Services such that the results are satisfactory to the Company.  {{schedule}}, or updated with 14 days prior notice.

# Fees.  

As consideration for the Services to be provided by Consultant and other obligations, the Company shall pay to Consultant the amounts specified in {{exhibit.compensation}} hereto at the times specified therein.

# Expenses.  

Consultant shall not be authorized to incur on behalf of the Company any expenses and will be responsible for all expenses incurred while performing the Services except as expressly specified in {{exhibit.expenses}} hereto unless otherwise agreed to by the Company's CEO, which consent shall be evidenced in writing for any such expenses in excess of $0.00.  As a condition to receipt of reimbursement, Consultant shall be required to submit to the Company reasonable evidence that the amount involved was both reasonable and necessary to the Services provided under this Agreement.

# Term and Termination.  

Consultant shall serve as a consultant to the Company for a period commencing on {{start-date}} and terminating on the earlier of (a) the date Consultant completes the provision of the Services to the Company under this Agreement, or (b) the date Consultant shall have been paid the maximum amount of consulting fees as provided in {{exhibit.compensation}} hereto.

\pagebreak

# Signatures

## THE COMPANY

My Corp Inc.

\ ![Company Signature]({{sign.image}})

---

By: {{sign.company-signer}}


## CONSULTANT

{{consultant}}

---

{{email}}


\pagebreak

# {{exhibit.services}}

## DESCRIPTION OF CONSULTING SERVICES

{{exhibit.services.value}}

\pagebreak

# {{exhibit.compensation}}

## COMPENSATION

{{exhibit.compensation.value}}

\pagebreak

# {{exhibit.expenses}}

## ALLOWABLE EXPENSES

{{exhibit.expenses.value}}
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
