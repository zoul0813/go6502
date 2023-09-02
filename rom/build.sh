#!/bin/sh

rm -rf rom.o rom.bin rom.dbg rom.*.txt

ca65 rom.s -g -o rom.o
ld65 -o rom.bin -C rom.cfg rom.o -m rom.map.txt -Ln rom.labels.txt --dbgfile rom.dbg

echo ""
echo ""
echo "Hex Dump"
echo "------------------------------------------------------------------------------"
hexdump -C rom.bin
echo "------------------------------------------------------------------------------"
echo ""