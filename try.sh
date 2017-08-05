#!/bin/bash
set -e

if (( $# != 2 )); then
  echo "Usage: $0 FMCBLOCK-SOURCE PNG-FILENAME"
  echo "Need FMC Block diagram source and PNG filename as parameters, for example:"
  echo "$0 '[Engine] <- (Fuel)' engine.png"
  exit 4
fi

if [[ -z "$1" ]]; then
  echo "Need FMC Block diagram source as first parameter, e.g. '[Engine] <- (Fuel)'"
  exit 4
fi

if [[ -z "$2" ]]; then
  echo "Need PNG filename as second parameter, e.g. 'engine.png'"
  exit 4
fi

make dreitafel

echo "$1" | ./dreitafel | dot -Tpng > $2
