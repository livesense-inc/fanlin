# fanlin

[![Circle CI](https://circleci.com/gh/jobtalk/fanlin/tree/master.svg?style=shield)](https://circleci.com/gh/jobtalk/fanlin/tree/master)

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
On Unix, Linux and OS X, fanlin programs read startup options from the following files, in the specified order.

```
/etc/fanlin.json
/etc/fanlin.cnf
/etc/fanlin.conf
/usr/local/etc/fanlin.json
/usr/local/etc/fanlin.cnf
/usr/local/etc/fanlin.conf
/usr/local/lvimg/fanlin.json
/usr/local/lvimg/fanlin.cnf
/usr/local/lvimg/fanlin.conf
./fanlin.json
./fanlin.cnf
./fanlin.conf
./conf.json
~/.fanlin.json
~/.fanlin.cnf
```

### example
```
{
    "port": 8080,
    "local_image_path": "{{local_image_path}}",
    "max_width": 1000,
    "max_height": 1000,
    "404_img_path": "/usr/local/lvimg/404.png",
    "access_log_path": "log path",
    "error_log_path": "error log path",
    "externals": [
        {
            "key" : "{{external contents server path}}"
        }
    ],
    "include": [
        "include configure path"
    ],
    "s3_bucket_name": "bucket name",
    "s3_region": "Tokyo",
    "s3_region": "Tokyo",
    "s3_region": "ap-northeast-1",
    "s3_region": "asia-pacific"
}
```
