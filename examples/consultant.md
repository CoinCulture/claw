# My Corp Inc.
# CONSULTING AGREEMENT

This Consulting Agreement (this "Agreement") is made as of {{ .Var.Date}}, by and between My Corp, Inc., a Delaware corporation (the "Company"), and {{ .Var.Consultant}} ("Consultant").

# Consulting Relationship.  

During the term of this Agreement, Consultant will provide consulting services to the Company as described on {{ .Exhibit.Services}} hereto (the "Services").  Consultant represents that Consultant is duly licensed (as applicable) and has the qualifications, the experience and the ability to properly perform the Services.  Consultant shall use Consultant’s best efforts to perform the Services such that the results are satisfactory to the Company.  {{ .Var.Schedule}}, or updated with 14 days prior notice.

# Fees.  

As consideration for the Services to be provided by Consultant and other obligations, the Company shall pay to Consultant the amounts specified in {{ .Exhibit.Compensation}} hereto at the times specified therein.

# Expenses.  

Consultant shall not be authorized to incur on behalf of the Company any expenses and will be responsible for all expenses incurred while performing the Services except as expressly specified in {{ .Exhibit.Expenses}} hereto unless otherwise agreed to by the Company's CEO, which consent shall be evidenced in writing for any such expenses in excess of $0.00.  As a condition to receipt of reimbursement, Consultant shall be required to submit to the Company reasonable evidence that the amount involved was both reasonable and necessary to the Services provided under this Agreement.

# Term and Termination.  

Consultant shall serve as a consultant to the Company for a period commencing on {{ .Var.StartDate}} and terminating on the earlier of (a) the date Consultant completes the provision of the Services to the Company under this Agreement, or (b) the date Consultant shall have been paid the maximum amount of consulting fees as provided in {{ .Exhibit.Compensation}} hereto.

\pagebreak

# Signatures

## THE COMPANY

My Corp Inc.

\ ![Company Signature]({{ .Sign.Image}})

---

By: {{ .Sign.CompanySigner}}


## CONSULTANT

{{ .Var.Consultant}}

---

{{ .Var.Email}}


\pagebreak

# {{ .Exhibit.Services}}

## DESCRIPTION OF CONSULTING SERVICES

{{ .Exhibit.Services}}

\pagebreak

# {{ .Exhibit.Compensation}}

## COMPENSATION

{{ .Exhibit.Compensation}}

\pagebreak

# {{ .Exhibit.Expenses}}

## ALLOWABLE EXPENSES

{{ .Exhibit.Expenses}}
