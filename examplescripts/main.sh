#!/bin/sh

dir=$1
for f in $(ls $dir | opactl search_spaced -p extension=yaml,word=2) 
do
  file=$dir/$f
  echo [$file]
  cat $file
  echo "\n-------"
done