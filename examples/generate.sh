#!/bin/bash

for fmc in $(ls *.fmc); do
  b=$(basename -s .fmc $fmc)
  echo $b
  ../dreitafel "$(cat $b.fmc)" | dot -Tpng > $b.png
done
