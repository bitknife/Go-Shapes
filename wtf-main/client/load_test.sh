#!/bin/zsh

N_CLIENTS=$1
LIFETIME=$2

for i in {1..${N_CLIENTS}}; do
  ./client --headless --lifetime_sec ${LIFETIME} &
  sleep 0.05
done
