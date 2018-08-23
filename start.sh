#!/usr/bin/env bash

export GOPATH=$(pwd) 

go test cleaner -count=1
go build cleanerapp
python3 creator.py
./cleanerapp

