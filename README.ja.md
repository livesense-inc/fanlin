# fanlin

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

[English](README.md) | 日本語

## 概要
fanlinはGo言語で作られた軽量画像プロキシです.
Amazon S3と外部HTTPサーバー上の画像をリアルタイムで加工することができます.

## 環境
### OS
* Linux (x86 and amd64)
* macOS

### Go Versions
* go 1.19.x

## 対応画像フォーマット
* JPEG
* PNG
* GIF
* WebP

## macOS の環境構築
## master pushの悲劇を防ぐために
[ここを参考に設定すること](http://ganmacs.hatenablog.com/entry/2014/06/18/224132)

## Linux用にクロスコンパイルする
### go1.19
```
$ GOOS=linux GOARCH=amd64 go build github.com/livesense-inc/fanlin/cmd/fanlin
```

## サーバーに配布するもの
ビルドして作った実行ファイル  
設定ファイル

## API
getパラメータに値を渡して操作する  
`w`画像の横幅  
`h`画像の縦幅  
色を指定しない時`w`と`h`を指定した時は小さい方に合わせて縮尺を変更する  
この時アスペクト比は維持する  
また一方が`0`の時は`0`ではない値を基準に縮尺を変更する  
`w`及び`h`が`0`であるときは元のサイズで表示する  
あまりにも大きいサイズが指定された時は設定ファイルにかかれている上限の大きさで拡大する  
`rgb`で色を指定した場合`w`と`h`で指定された大きさの画像を生成する  
この時画像のアスペクト比が違うときは隙間を指定した色で塗りつぶす
`quality`で`0`から`100`までの数値を指定した場合その数値にクオリティを設定した画像を生成する
それ以外の数値が指定された場合は`jpeg.DefaultQuality`の値(`75`)が指定される
`crop`で`true`が指定された場合`w`と`h`で指定した比に合わせて、画像中央を基準としてはみ出した部分クロッピングして画像を生成する
例:  
```
http://localhost:8080/path/to/image/?h=400&w=400&rgb=100,100,100&quality=80&crop=true
```

## testing
```
$ go test -cover ./...
```

## 設定項目に関して
だいたいこんな感じでかけます
```
{
    "port": 8080,
    "max_width": 1000,
    "max_height": 1000,
    "404_img_path": "/path/to/404/image",
    "access_log_path": "/path/to/access/log",
    "error_log_path": "/path/to/error/log",
    "use_server_timing": true,
    "providers": [
        {
            "alias/0" : {
                "type" : "web",
                "src" : "http://aaa.com/bbb",
                "priority" : 10
            }
        },
        {
            "alias/1" : {
                "type" : "web",
                "src" : "https://ccc.com",
                "priority" : 20
            }
        },
        {
            "alias/3" : {
                "type" : "s3",
                "src" : "s3://bucket/path",
                "region" : "ap-northeast-1",
                "priority" : 30
            }
        }
    ]
}
```

## ログの出力先を制御する
設定項目の`access_log_path`/`error_log_path`/`debug_log_path`にパスを指定することで、それぞれのログをファイルに出力できます。
標準出力にログを出力したい場合は、`/dev/stdout`を指定してください。

## WebPフォーマットの利用方法と制限事項
fanlinへのリクエストに `webp=true` getパラメータを付与することで、WebPフォーマットの画像を返すことが出来ます.

例:

- JPG画像のURL:
  - http://localhost:8080/abc.jpg?h=400&w=400&quality=80
- WebPに変換した画像のURL:
  - http://localhost:8080/abc.jpg?h=400&w=400&quality=80&webp=true

また、以下の条件の場合は、Lossless WebPフォーマットで変換します.

- getパラメータにて `quality=100` を指定かつ元画像のフォーマットが PNG / GIF / WebP のいずれか

### 制限事項

- アニメーションには対応していません


## Server-Timingのサポートに関して

設定ファイルのグローバル設定値に `"use_server_timing": true` を入れることで[Server-Timing](https://www.w3.org/TR/server-timing/)が出力されます.
Server-Timingの出力によって、システムの内部構成やパフォーマンスがエンドユーザーに見えてしまう可能性があります.利用に際してはご注意ください.

現在の出力項目は以下:

- f_load: ソース画像のロード時間
- f_decode: ソース画像のデコードと加工時間
- f_encode: 最終出力フォーマットへのエンコード時間
