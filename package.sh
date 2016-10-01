#!/usr/bin/env bash

rm -f fixedhashmap.zip
zip -r fixedhashmap.zip \
    README.md \
    fixedhashmap.go \
    fixedhashmap_test.go
