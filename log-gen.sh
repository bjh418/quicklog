#!/bin/bash
while true; do
  echo `date -u +%Y-%m-%dT%H:%M:%SZ` 'DEBUG This is a debug line' |nc -U ./quicklogger.sock
  sleep 1
done
