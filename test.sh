#!/bin/bash

i=0
for value in foo bar bing baz foo foo foo foo bar bar bar bing baz
do
  kubectl --context kind-sample run "curl-${i}" --rm --restart=Never --image=curlimages/curl:latest -it -- "http://sample/${value}"
  i=$((i+1))
done
