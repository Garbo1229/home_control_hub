# HomeControlHub

HomeControlHub 是一个基于 Go 语言和 Gin 框架开发的智能家居管理 API。该项目提供了一系列接口用于控制和管理家中的网络设备，如 IP 地址查询与更新、NAS（网络附加存储）设备的唤醒与关闭，以及树莓派设备的重启与关闭。目标是通过简单的 API 实现对多种设备的集中管理。

## 功能

1. IP 地址管理

- GET `/ip`: 查询当前设备的 IP 地址。
- GET `/ip/update`: 更新当前设备的 IP 地址。
- GET `/ip/:ip`: 查询指定 IP 地址。

2. NAS 设备控制

- GET `/nas/open`: 唤醒 NAS 设备。
- GET `/nas/shutdown`: 关闭 NAS 设备。

3. 树莓派控制

- GET `/raspberry/restart`: 重启树莓派设备。
- GET `/raspberry/shutdown`: 关闭树莓派设备。

4. 请求日志记录

所有的请求都会通过 `logrequest.LogRequest` 中间件进行日志记录，方便调试和跟踪请求。

## 安装与运行

1. 克隆项目

```
git clone https://github.com/your-username/HomeControlHub.git
cd HomeControlHub
```

2. 安装依赖

确保你已经安装了 Go 和 Git。

然后在项目根目录运行以下命令来安装项目的依赖：
```
go mod tidy
```

3. 启动服务

执行以下命令启动 API 服务：

```
go run main.go
```
此时，API 服务会在 `:8082` 上运行。

4. 测试 API

你可以使用 Postman 或 curl 来测试以下 API：

- GET `/ip`: 查询当前 IP 地址。
- GET `/ip/update`: 更新当前 IP 地址。
- GET `/ip/:ip`: 查询指定的 IP 地址。
- GET `/nas/open`: 唤醒 NAS 设备。
- GET `/nas/shutdown`: 关闭 NAS 设备。
- GET `/raspberry/restart`: 重启树莓派设备。
- GET `/raspberry/shutdown`: 关闭树莓派设备。

## 配置

为了让项目正确运行，你需要将 `ip2region` 项目放置在与 `HomeControlHub` 项目同级目录下，并确保 `ip2region.db` 文件存在。

项目目录结构示例：
```
Dir/
├── ip2region/
├── HomeControlHub/
```

请复制 config.yaml.example 文件并根据你的设备和文件路径更新配置：

```
# config.yaml.example - 配置文件示例

# NAS 配置
nas:
  url: "http://example-nas.local"
  mac: "00:1A:2B:3C:4D:5E"
  username: "admin"
  password: "password123"

# IP 库配置（ip2region）
ip2region:
  path: "/path/to/ip2region.db"
  maker_path: "/path/to/ip2region-maker"
  db_path: "/path/to/ip2region.db"
```

配置说明：

- nas：NAS 设备的 URL、MAC 地址以及登录凭据。
- ip2region：配置 IP 地址库文件的路径，用于 IP 查询与地理位置获取。

## 感谢

特别感谢 [mkch/qnap-tool](https://github.com/mkch/qnap-tool) 和 [lionsoul2014/ip2region](https://github.com/lionsoul2014/ip2region) 的开源作者，您们的贡献使本项目得以快速完成。