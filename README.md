# fanlin

[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

English | [日本語](README.ja.md)

## abstract
fanlin is image proxy server in native Go language.

## Support
### OS
* Linux (x86 and amd64)
* macOS

### Go Versions
* go 1.19.x

### Image Format
* JPEG
* PNG
* GIF
* WebP

## Cross compile for amd64 Linux
### go 1.19
```
$ GOOS=linux GOARCH=amd64 go build github.com/livesense-inc/fanlin/cmd/fanlin
```

## testing
```
$ go test -cover ./...
```

## configure
On Unix, Linux and macOS, fanlin programs read startup options from the following files, in the specified order(top files are read first, and precedence).

```
/etc/fanlin.json
/etc/fanlin.cnf
/etc/fanlin.conf
/usr/local/etc/fanlin.json
/usr/local/etc/fanlin.cnf
/usr/local/etc/fanlin.conf
./fanlin.json
./fanlin.cnf
./fanlin.conf
./conf.json
~/.fanlin.json
~/.fanlin.cnf
```

### example

#### fanlin.json
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

## Controling where logs are output to
You can output each log to a file by specifying the path in `access_log_path`/`error_log_path`/`debug_log_path`.
If you want to output logs to standard output, specify `/dev/stdout`.

## Using WebP and Limitations
You can get WebP image format with GET parameter `webp=true` requeest.

Examples:

- JPG image URL:
  - http://localhost:8080/abc.jpg?h=400&w=400&quality=80
- WebP encoded image URL:
  - http://localhost:8080/abc.jpg?h=400&w=400&quality=80&webp=true

fanlin returns lossless WebP image in following conditions.

- GET parameter `quality=100` AND source image format is PNG / GIF / WebP

### Limitations

- Do not support animations


## Server-Timing Support

Add `"use_server_timing": true` at Global parameters in config file.
You will get [Server-Timing](https://www.w3.org/TR/server-timing/) output.
Be careful, your system architecture or perfomance will be exposed to enduser with Server-Timing output.

fanlin outputs following timings:

- f_load: The time for load source image.
- f_decode: The time for decode and format source image.
- f_encode: The time for encode to final image format.


## LICENSE
Written in Go and licensed under [the MIT License](https://opensource.org/licenses/MIT), it can also be used as a library.
