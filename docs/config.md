# CS2Panel 配置文件说明

## config/config.yaml

### 环境设置 (env)
- `mode`: 运行模式，支持 `debug` 和 `release`
- `log_level`: 日志级别，可选 `debug`、`info`、`warn`、`error`

### 服务器配置 (server)
- `port`: API服务端口，默认 8080
- `panel_data_dir`: 面板数据存储目录
- `web_server`: 是否启用Web服务
- `web_server_port`: Web服务端口

### Docker配置 (docker)
- `image_name`: CS2服务器Docker镜像名称
- `tag`: 镜像标签版本
- `volume_name`: Docker数据卷名称
- `prefix`: 容器名称前缀
- `max_retries`: 连接失败时的重试次数
- `retry_delay`: 重试间隔时间（秒）
- `cs_data_dir`: CS2数据目录路径

### 日志配置 (util)
- `log_dir`: 日志文件存储目录
- `log_filename`: 日志文件名
- `log_max_size`: 单个日志文件最大大小（MB）
- `log_max_backups`: 日志文件最大备份数量
- `log_max_age`: 日志文件保存天数
- `log_compress`: 是否压缩旧日志文件

### 游戏配置 (game)
- `srcds_token`: Steam服务器令牌，需在 [Steam开发者页面](https://steamcommunity.com/dev/managegameservers) 申请
- `rcon_password`: RCON远程控制密码
- `address`: 服务器地址

## 使用说明

1. 首次使用需要申请 `srcds_token`
2. 修改 `rcon_password` 为安全密码
3. 根据需要调整端口和目录路径
4. 生产环境建议设置 `mode` 为 `release`

## .env
- `VITE_API_BASE_URL` : API地址
