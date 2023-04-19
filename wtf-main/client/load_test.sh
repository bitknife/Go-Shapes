#!/bin/zsh

N_CLIENTS=$1
LIFETIME=$2
START_DELAY=$3

for i in {1..${N_CLIENTS}}; do
  ./client --headless --lifetime_sec ${LIFETIME} &
  sleep $3
done
