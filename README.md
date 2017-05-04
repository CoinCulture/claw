# claw - Command Line Law

## Install

```
go get github.com/CoinCulture/claw
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
# initaite a new contract
claw new john examples/template/consultant.md

# edit the params
vim john/params.toml

# compile the markdown and output a final pdf using pandoc
claw compile --output pdf john
```
