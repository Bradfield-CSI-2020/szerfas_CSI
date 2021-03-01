#!/bin/bash

trap "echo signal received!" SIGINT

echo "The script pid is $$"
sleep 10
