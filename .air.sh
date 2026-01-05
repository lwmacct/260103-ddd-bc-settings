#!/usr/bin/env bash

__main() {
  pkill -f '(go-build|tmp/main)'

  _air_num_procs=$(pgrep air | wc -l)
  echo "_air_num_procs: ${_air_num_procs}"

  if (($(echo "$_air_num_procs > 1" | bc -l))); then
    pgrep air | sort -n | head -n1 | xargs kill
  fi

  # air 2>&1 | tee tmp/run.log
}

__main
