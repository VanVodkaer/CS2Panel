# Docker 控制面板 API 文档

## 一、系统状态接口

### 1. 检查 Docker 服务状态

- **接口** ：`ANY /api/docker/ping`
- **说明** ：检查 Docker 服务是否可连接
- **返回示例** ：

```json
{
  "message": "Docker 服务正在运行",
  "ping": {
    "APIVersion": "1.48",
    "OSType": "linux",
    "Experimental": false,
    "BuilderVersion": "2",
    "SwarmStatus": {
      "NodeState": "inactive",
      "ControlAvailable": false
    }
  }
}
```

---

## 二、容器管理接口

### 2. 拉取镜像

- **接口** ：`POST /api/docker/image/pull`
- **说明** ：从配置中指定地址拉取镜像
- **返回** ：

```json
{
    "message": "已开始拉取镜像"
}
```

---

### 2. 获取拉取镜像状态

- **接口** ：`GET /api/docker/image/pull/status`
- **说明** ：获取拉取镜像状态
- **返回** ：

```json
// 正在拉取镜像
{
    "message": "pulling"
}
```

```json
// 未拉取镜像/镜像拉取完成
{
    "message": "not_started"
}
```

---

### 3. 获取容器列表

- **接口** ：`GET /api/docker/container/list`
- **说明** ：列出所有以配置前缀命名的容器
- **返回** ：

```json
{
  "containers": [
    {
      "Id": "ffa7a35aab00d0f2876e410d07196e6e3859fb18eb2ac16fe8ab06d02b0fbf7c",
      "Names": ["/cs2panel-test2"],
      "Image": "joedwards32/cs2",
      "ImageID": "sha256:a664b0f5c88fade529f0627bc1b3ebaf651cb01a156213b0ff32e0bf0d90acb4",
      "Command": "bash entry.sh",
      "Created": 1746753314,
      "Ports": [],
      "Labels": {
        "desktop.docker.io/binds/0/Source": "C:\\Users\\vodkaer\\Desktop\\CS2Panel\\cs2-data",
        "desktop.docker.io/binds/0/SourceKind": "hostFile",
        "desktop.docker.io/binds/0/Target": "/home/steam/cs2-dedicated",
        "maintainer": "joedwards32@gmail.com",
        "org.opencontainers.image.created": "2025-01-26T05:30:10.203Z",
        "org.opencontainers.image.description": "CS2 Dedicated Server Docker Image",
        "org.opencontainers.image.licenses": "MIT",
        "org.opencontainers.image.revision": "a36e8c215f7dd6316b1ad843e48691a2419f9cd0",
        "org.opencontainers.image.source": "https://github.com/joedwards32/CS2",
        "org.opencontainers.image.title": "CS2",
        "org.opencontainers.image.url": "https://github.com/joedwards32/CS2",
        "org.opencontainers.image.version": "latest"
      },
      "State": "created",
      "Status": "Created",
      "HostConfig": {
        "NetworkMode": "bridge"
      },
      "NetworkSettings": {
        "Networks": {
          "bridge": {
            "IPAMConfig": null,
            "Links": null,
            "Aliases": null,
            "MacAddress": "",
            "DriverOpts": null,
            "NetworkID": "",
            "EndpointID": "",
            "Gateway": "",
            "IPAddress": "",
            "IPPrefixLen": 0,
            "IPv6Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "DNSNames": null
          }
        }
      },
      "Mounts": [
        {
          "Type": "bind",
          "Source": "/run/desktop/mnt/host/c/Users/vodkaer/Desktop/CS2Panel/cs2-data",
          "Destination": "/home/steam/cs2-dedicated",
          "Mode": "",
          "RW": true,
          "Propagation": "rprivate"
        }
      ]
    },
    {
      "Id": "93be7085f696e7d78ca44ffc89fb0ea360ad3808d2f72f6119ec33f0e0e92b13",
      "Names": ["/cs2panel-test"],
      "Image": "joedwards32/cs2",
      "ImageID": "sha256:a664b0f5c88fade529f0627bc1b3ebaf651cb01a156213b0ff32e0bf0d90acb4",
      "Command": "bash entry.sh",
      "Created": 1746721902,
      "Ports": [
        {
          "IP": "0.0.0.0",
          "PrivatePort": 27015,
          "PublicPort": 27015,
          "Type": "tcp"
        },
        {
          "IP": "0.0.0.0",
          "PrivatePort": 27015,
          "PublicPort": 27015,
          "Type": "udp"
        },
        {
          "PrivatePort": 27020,
          "Type": "udp"
        }
      ],
      "Labels": {
        "desktop.docker.io/binds/0/Source": "C:\\Users\\vodkaer\\Desktop\\CS2Panel\\cs2-data",
        "desktop.docker.io/binds/0/SourceKind": "hostFile",
        "desktop.docker.io/binds/0/Target": "/home/steam/cs2-dedicated",
        "maintainer": "joedwards32@gmail.com",
        "org.opencontainers.image.created": "2025-01-26T05:30:10.203Z",
        "org.opencontainers.image.description": "CS2 Dedicated Server Docker Image",
        "org.opencontainers.image.licenses": "MIT",
        "org.opencontainers.image.revision": "a36e8c215f7dd6316b1ad843e48691a2419f9cd0",
        "org.opencontainers.image.source": "https://github.com/joedwards32/CS2",
        "org.opencontainers.image.title": "CS2",
        "org.opencontainers.image.url": "https://github.com/joedwards32/CS2",
        "org.opencontainers.image.version": "latest"
      },
      "State": "running",
      "Status": "Up 10 minutes",
      "HostConfig": {
        "NetworkMode": "bridge"
      },
      "NetworkSettings": {
        "Networks": {
          "bridge": {
            "IPAMConfig": null,
            "Links": null,
            "Aliases": null,
            "MacAddress": "8e:5d:65:b3:34:9a",
            "DriverOpts": null,
            "NetworkID": "b6e17a9a59da852a62a6c760139d7d008d11985e71648c5b30610ae48175cfbb",
            "EndpointID": "e45e208eafc97bef4083d2947c65e4fd3e37c2bcb9a645600c8548962b4ae547",
            "Gateway": "172.17.0.1",
            "IPAddress": "172.17.0.2",
            "IPPrefixLen": 16,
            "IPv6Gateway": "",
            "GlobalIPv6Address": "",
            "GlobalIPv6PrefixLen": 0,
            "DNSNames": null
          }
        }
      },
      "Mounts": [
        {
          "Type": "bind",
          "Source": "/run/desktop/mnt/host/c/Users/vodkaer/Desktop/CS2Panel/cs2-data",
          "Destination": "/home/steam/cs2-dedicated",
          "Mode": "",
          "RW": true,
          "Propagation": "rprivate"
        }
      ]
    }
  ]
}
```

---

### 4. 创建容器

- **接口** ：`POST /api/docker/container/create`
- **请求参数** （JSON）：

| 字段        | 类型   | 必填 | 说明                |
| ----------- | ------ | ---- | ------------------- |
| name        | string | 是   | 容器名称            |
| server_name | string | 否   | 游戏服务名称        |
| game_port   | string | 否   | RCON/游戏服务器端口 |
| watch_port  | string | 否   | 观战端口（可选）    |

- **返回** ：

```json
{
  "container_id": "ffa7a35aab00d0f2876e410d07196e6e3859fb18eb2ac16fe8ab06d02b0fbf7c",
  "message": "容器创建成功"
}
```

---

### 5. 启动容器并执行命令

- **接口** ：`POST /api/docker/container/start`
- **请求参数** （JSON）：

| 字段 | 类型     | 必填 | 说明           |
| ---- | -------- | ---- | -------------- |
| name | string   | 是   | 容器名称       |
| cmds | string[] | 否   | 启动后执行命令 |

- **返回** ：

```json
{
  "message": "容器启动成功",
  "responses": ["执行结果1", "执行结果2"]
}
```

---

### 6. 停止容器

- **接口** ：`POST /api/docker/container/stop`
- **请求参数** （JSON）：

```json
{
  "name": "your_container_name"
}
```

- **返回** ：

```json
{
  "message": "容器停止成功"
}
```

---

### 7. 重启容器

- **接口** ：`POST /api/docker/container/restart`
- **请求参数同上**
- **返回** ：

```json
{
  "message": "容器重启成功"
}
```

---

### 8. 删除容器

- **接口** ：`DELETE /api/docker/container/remove`
- **请求参数同上**
- **返回** ：

```json
{
  "message": "容器删除成功"
}
```

---

### 9. 执行容器命令（RCON）

- **接口** ：`POST /api/docker/container/exec`
- **请求参数** ：

| 字段 | 类型     | 必填 | 说明             |
| ---- | -------- | ---- | ---------------- |
| name | string   | 是   | 容器名称         |
| cmds | string[] | 是   | 要执行的命令列表 |

- **返回** ：

```json
{
  "message": "执行命令成功",
  "responses": ["命令1输出", "命令2输出"]
}
```
