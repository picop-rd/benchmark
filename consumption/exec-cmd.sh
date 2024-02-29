#!/bin/bash -eux

# 引数のチェック
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <cmd> <dir> <timestamp>"
    exit 1
fi

cmd=$1
dir=$2
timestamp=$3

# ディレクトリを必要に応じて作成
mkdir -p "$dir"

filename="${dir}/${timestamp}.txt"

if [ -e $filename ]; then
	echo "duplicated timestamp: $filename"
	exit 1
fi

echo "\$ ${cmd}" > $filename

# コマンドの実行と出力の保存
eval $cmd | tee -a $filename
