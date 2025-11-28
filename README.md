# QQStickerDL

QQStickerDL 是一个用于下载 QQ 表情包的命令行工具，基于 Go 语言开发。

## 功能

- 支持下载 QQ 表情包
- 简单易用，可连续下载
- 支持链接 ID或本地文件

## 使用教程

构建运行```QQStickerDL```后输入链接/ID/本地文件即可下载

- ### 手机获取表情包链接
   在手机QQ的「表情详情」页面，右上角点击复制链接

- ### QQNT获取本地文件
    默认情况下在 <br>
```C:\Users\%USERNAME%\Documents\Tencent Files\%这里填写你的QQ%\nt_qq\nt_data\Emoji\marketface\json```以ID开头的jtmp文件(其实是json)

## 构建

1. 克隆本仓库：

```shell
git clone https://github.com/xyBakaQAQ/QQStickerDL.git
```

2. 进入项目目录并构建：

```shell
cd QQStickerDL
go build -o QQStickerDL main.go
```


## 更新日志

25/11/28 V0.1.2
> 支持连续下载和本地文件<br>
> Debug模式下显示下载链接和保存Json

25/11/28 V0.1.1
> 优化代码更简洁<br>

25/11/28 V0.1
> 实现下载300X300的表情包<br>

## 依赖

- Go 1.24 及以上版本

## 许可证

**MIT License**
