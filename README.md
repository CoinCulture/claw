# Contract Flow

## Install

```
go get github.com/ebuchman/cflow
```

You will also need `pandoc` and `xelatex` to output pdfs.
On Ubuntu:

```
sudo apt-get install pandoc
sudo apt-get install texlive-xetex # this one's big
```

## Run

```
# initaite a new contract
cflow new john template.md

# edit the params
vi john/params.toml

# compile the final markdown
cflow compile john

# compile the pdf
pandoc john/contract.md --latex-engine=xelatex -o john/contract.pdf
```
