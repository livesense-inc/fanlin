# fanlin

[![Circle CI](https://circleci.com/gh/jobtalk/fanlin/tree/master.svg?style=shield)](https://circleci.com/gh/jobtalk/fanlin/tree/master)

## 概要

## 環境
* 確認済み
    * go1.4
    * go1.5
    * go1.6

## OS X の環境構築
## master pushの悲劇を防ぐために
[ここを参考に設定すること](http://ganmacs.hatenablog.com/entry/2014/06/18/224132)
## 事前に入れるもの

## Linux用にクロスコンパイルする
go1.5の場合  
1.4はやり方が変わります  
```
$ cd {{repository}}/proxy
$ GOOS=linux GOARCH=amd64 go build -o ../bin/proxy
```

## サーバーに配布するもの
ビルドして作った実行ファイル  
conf.json  

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
例:  
```
http://localhost:8080/path/to/image/?h=400&w=400&rgb=100,100,100
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
    "local_image_path": "ローカルのファイルの設置場所",
    "max_width": 1000,
    "max_height": 1000,
    "404_img_path": "/usr/local/lvimg/404.png",
    "access_log_path": "logの書き出し先",
    "error_log_path": "logの書き出し先",
    "externals": [
        {
            "key" : "proxy先のURL"
        }
    ],
    "include": [
        "include先のパス"
    ],
    "s3_bucket_name": "s3のバケット名",
    "s3_region": "未指定の時は東京になる. 以下のように指定ができる",
    "s3_region": "Tokyo",
    "s3_region": "Tokyo",
    "s3_region": "ap-northeast-1",
    "s3_region": "asia-pacific"
}

```
