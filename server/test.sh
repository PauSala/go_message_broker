#!/bin/bash

for i in {1..1}
do
   nc localhost 3000 < messages/set.bin &&
   nc localhost 3000 < messages/pub.bin &&
   nc localhost 3000 < messages/sub.bin
done
