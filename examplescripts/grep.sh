#!/bin/sh

dir=$1
keywords=column21,test21
for f in $(ls $dir) 
do
  file=$dir/$f
  check=$(cat $file | opactl grep -P keywords=$keywords)
  if [ $check ]; then
    echo $file includes all keywords
  fi
done