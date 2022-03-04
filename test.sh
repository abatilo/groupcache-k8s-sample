#!/bin/bash

i=0
for value in foo bar bing baz foo foo foo foo bar bar bar bing baz
do
  curl "http://localhost:8000/${value}"
  i=$((i+1))
done
