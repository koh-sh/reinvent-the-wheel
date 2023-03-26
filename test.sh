#!/usr/bin/env bash

INPUT=$(pbpaste)
go run ./stdin_calc.go < <(echo $INPUT)
