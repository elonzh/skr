# skr ![Tag](https://img.shields.io/github/tag/elonzh/skr.svg?style=flat-square) [![GolangCI](https://golangci.com/badges/github.com/elonzh/skr.svg)]() [![Build Status](https://img.shields.io/travis/elonzh/skr.svg?style=flat-square)](https://travis-ci.org/elonzh/skr)

一组脚本。


## Skr

```shell
🏁 skr~ skr~

Usage:
  skr [command]

Available Commands:
  douyin      解析抖音名片数据
  gaoxiaojob  抓取 高校人才网(http://gaoxiaojob.com/) 的最近招聘信息并根据关键词推送至钉钉
  help        Help about any command
  merge_score 合并学生成绩单

Flags:
  -c, --config string      配置文件路径
  -h, --help               help for skr
      --log-level uint32    (default 4)

Use "skr [command] --help" for more information about a command.
```

## douyin

```text
$skr douyin --help
爱抖音小助手, 它能帮你解析抖音名片数据

Usage:
  skr douyin [flags]

Flags:
  -c, --config string   配置文件路径(默认为 config.yaml)
  -h, --help            help for douyin
      --silent          静默模式, 只在出错时输出日志
  -u, --urls strings    抖音分享链接
      --version         version for douyin
```

## gaoxiaojob

```text
$skr gaoxiaojob --help
抓取 高校人才网(http://gaoxiaojob.com/) 的最近招聘信息并根据关键词推送至钉钉

Usage:
  skr gaoxiaojob [flags]

Flags:
  -h, --help                   help for gaoxiaojob
  -k, --keywords stringArray   关键词
  -s, --storage string         历史记录数据路径 (default "storage.boltdb")
  -v, --verbose                调试模式
      --version                version for gaoxiaojob
```

### 定时执行

macOS, Linux 系统可以使用 crontab 进行定时执行, 例如

```text
* * * * * <程序路径> <钉钉机器人地址> >> <日志文件路径> 2>&1
```

Windows 可以使用计划任务进行设置


### 示例

- [今天涨粉了吗?](examples/今天涨粉了吗)

> [查看更多示例](examples)

# License

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Felonzh%2Fskr.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Felonzh%2Fskr?ref=badge_large)
