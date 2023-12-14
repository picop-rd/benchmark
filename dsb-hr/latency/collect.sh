#!/bin/bash -eux

# 引数のチェック
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <cmd> <dir>"
    exit 1
fi

cmd=$1
dir=$2

# ディレクトリを必要に応じて作成
mkdir -p "$dir"

# 現在のUNIXタイムスタンプを取得
timestamp=$(date +%s)

filename="${dir}/${timestamp}.txt"

echo "\$ ${cmd}" > $filename

# コマンドの実行と出力の保存
eval $cmd | tee -a $filename
