# Contract Flow

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
cflow new john template.md

# edit the params
vim john/params.toml

# compile the markdown and output a final pdf using pandoc
cflow compile --output pdf john
```
