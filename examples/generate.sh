#!/bin/bash

for fmc in $(ls *.fmc); do
  b=$(basename -s .fmc $fmc)
  echo $b
  cat $b.fmc | ../dreitafel > $b.dot
  cat $b.dot | dot -Tpng > $b.png
done
