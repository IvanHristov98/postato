# 🏴‍☠️ 🍠 ⚔️ Arghhhh....!!! Cpt. Postato says you're standing on the fuzzy plank 🏴‍☠️

![Run Postato! Run!](docs/images/meme.jpg "Run Postato! Run!")

Postato is a tool that identifies the position of your body with fuzzy logic.

## Fuzzy inference system

The fuzzy rule set is built upon crisp data which is accepted under the form of a dataset where each tuple is in the format `(xWrist, yWrist, zWrist, xThigh, yThigh, zThigh, activity)`. The data used is provided by the [UCI data repo](http://archive.ics.uci.edu/ml/datasets/selfBACK). It is first fuzzified with `soft kMeans++` using Newton's gravity formula producing a super cluster of fuzzy clusters. Once they are obtained a fuzzy rule is generated from the fuzzy boundaries of each cluster. A rule consists of a mapping between a body position axis and a fuzzy number. Currently Postato supports triangular and gaussian fuzzy numbers.

But how is this valuable? 🤔

A fuzzy inferer is required to make use of a fuzzy rule set. Postato uses the one of Mamdani.

## How to use?

At the moment only fuzzy rule set plotting is supported. It can be executed with:

```bash
# git clone this repo
# Fetch all dependencies.
go get ./...

# If you have direnv installed run `direnv allow` otherwise
source .envrc

# For gaussian fuzzy numbers
go run cmd/postato/main.go draw -d data/new.csv -t gaussian

# For triangular fuzzy numbers
go run cmd/postato/main.go draw -d data/new.csv -t triangular
```

All images are generated within the `gen/image` directory.

## Accuracy

It is tested using 10-fold-cross validation which can be tested like executed like this:

```bash
# For gaussian fuzzy numbers. Shows a success rate of ~82%.
go run cmd/postato/main.go test -d data/new.csv -t gaussian

# For triangular fuzzy numbers. Shows a success rate of ~12% which is far worse.
go run cmd/postato/main.go test -d data/new.csv -t triangular
```
