# tinygo-flash

(プログラムの概要説明)

## Description

(プログラムの動作説明)

## Demo

(gifアニメ等でのDemo)

## Requirement

(動作環境)

## Usage

(使用方法)

## Installation

(インストール方法)

## Build

ソースコードは、 `%GOPATH%/src/github.com/sago35/tinygo-flash/main.go` となるように配置してください。  

`go run dist/make_dist.go VERSION` を実行すると、 ./dist/release 以下にファイルが作成されます。  
release用にBuildする際は、Changes.mdに該当Versionの記載が必要です。  

    (例) release用にv1.2.3でbuild
    $ go run dist/make_dist.go 1.2.3

    (例) 開発用にv1.2.3でbuild
    $ go run dist/make_dist.go 1.2.3 --dev

### Environment

* go
* kingpin.v2

## Notice

(注意)

## License

Copyright (c) 2020 sago35  
Released under the MIT license  
https://opensource.org/licenses/mit-license.php  

## Author

sago35 - <sago35@gmail.com>
