# 抖音极简版
第五届字节青训营后端 极简抖音大项目
## 概览
本项目为第五届字节青训营后端，**服务熔断**队作品。
## 构建
项目中使用`ffmpeg`对视频进行处理，需要先确保运行环境的环境变量有`ffmpeg`。

在`macOS`下可以使用`brew install ffmpeg`命令安装`ffmpeg`，其它环境可以使用相应包管理器进行安装或下载可执行文件，自行添加到环境变量。
```shell
$ git clone https://github.com/amcones/douyin
$ cd douyin
$ go build && ./douyin
```
## 使用指南
需要先查看位于`./config/conf.toml`的数据库配置，根据需要修改。项目默认部署在`127.0.0.1:8080`，可以使用`Apifox`、`Postman`等工具测试接口。
## 如何贡献
本项目使用 Github Forking 工作流，具体参考：[抖音极简版项目代码贡献流程](./docs/zh-CN/Contribute.md)