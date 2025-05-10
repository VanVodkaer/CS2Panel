# CS2Panel (开发中)

> 轻量级、用户友好的 CS2（Counter-Strike 2）游戏服务器管理工具
>
> 本仓库为后端部分，使用 **Go (Golang)** 编写

## 📦 安装 & 运行

1. 编辑config目录下的`config.yaml`文件

2. 运行 Docker

3. 运行以下命令
```bash
# 克隆仓库
git clone https://github.com/VanVodkaer/CS2Panel
cd CS2Panel

# 安装依赖
go mod tidy

# 运行
go run ./cmd
```

---

## ⚙️ 配置

默认配置文件路径：`./config.yaml`

参考[文档](https://github.com/VanVodkaer/CS2Panel/blob/main/docs/config.md)

---

## 📚 API 文档

详细的 API [文档](https://github.com/VanVodkaer/CS2Panel/blob/main/docs/index.md)请参考：`/docs`

---

## 🧱 目录结构

```bash
cs2panel/
├── cmd/            # 程序入口
├── config/         # 配置加载逻辑
├── docker/         # docker实例配置
├── server/         # 主要路由和服务
├── utils/          # 工具方法
└── go.mod
```

---

## ✅ 贡献指南

欢迎社区贡献！

1. Fork 本仓库
2. 创建特性分支 `git checkout -b feature/xxx`
3. 提交修改 `git commit -m '新增功能 xxx'`
4. 推送分支并创建 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证，详情请见 [LICENSE](https://github.com/VanVodkaer/CS2Panel/blob/main/LICENSE)。

