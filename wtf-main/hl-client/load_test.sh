#!/bin/zsh

N_CLIENTS=$1
LIFETIME=$2
START_DELAY=$3
HOSTNAME=$4

go build .

for i in {1..${N_CLIENTS}}; do
  ./hl-client --protocol tcp --lifetime_sec ${LIFETIME} --host ${HOSTNAME} &
  sleep $3
done
