# fanlin

[![Circle CI](https://circleci.com/gh/jobtalk/fanlin/tree/master.svg?style=shield)](https://circleci.com/gh/jobtalk/fanlin/tree/master)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

English | [日本語](README.ja.md)

## abstract
fanlin is image proxy server in native Go language.

## Support
### OS
* Linux (x86 and amd64)
* OS X

### Go Versions
* go 1.4
* go 1.5
* go 1.6

## Cross compile for amd64 Linux
### go 1.5
```
$ GOOS=linux GOARCH=amd64 go build github.com/jobtalk/fanlin/cmd/fanlin
```

## configure
On Unix, Linux and OS X, fanlin programs read startup options from the following files, in the specified order(top files are read first, and precedence).

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
    "providers": [
        {
            "alias/0" : {
                "type" : "web",
                "src" : "http://aaa.com/bbb"
            }
        },
        {
            "alias/1" : {
                "type" : "web",
                "src" : "https://ccc.com"
            }
        },
        {
            "alias/3" : {
                "type" : "s3",
                "src" : "s3://bucket/path",
                "region" : "ap-northeast-1"
            }
        }
    ]
}
```

## LICENSE
Written in Go and licensed under [the MIT License](https://opensource.org/licenses/MIT), it can also be used as a library.
