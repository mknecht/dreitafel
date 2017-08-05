#!/bin/bash

for fmc in $(ls *.fmc); do
  b=$(basename -s .fmc $fmc)
  echo $b
  cat $b.fmc | ../dreitafel | dot -Tpng > $b.png
done
