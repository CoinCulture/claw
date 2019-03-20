# claw - Command Line Law

## Install

```
go get github.com/CoinCulture/claw
dep ensure
```

You will also need `pandoc` and `xelatex` to output pdfs.
On Ubuntu:

```
sudo apt-get install pandoc

# this one's big
sudo apt-get install texlive-xetex
```

## Run

```
# initiate a new contract
claw new john examples/templates/consultant.md

# change to the newly created directory
cd john

# edit the params
vim params.toml

# save your revisions to the hash log
claw revise

# compile the markdown and output a final pdf using pandoc
claw compile --output pdf
```
