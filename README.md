# CS2Panel

> 轻量级、用户友好的 CS2（Counter-Strike 2）游戏服务器管理工具
>
> 本仓库后端部分，使用 **Go (Golang)** 编写
> 前端部分，使用 **React** 编写

 喜欢这个项目的话，请给个Star吧⭐

## 📦 安装 & 运行

1. 拉取代码

```bash
# 克隆仓库
git clone https://github.com/VanVodkaer/CS2Panel
cd CS2Panel
```

2. 重命名config目录下的 `config.yaml.example` 文件 为 `config.yaml` 并编辑
3. 重命名根目录下的 `.env.example` 文件 为 `.env` 并编辑
4. 运行 Docker
5. 安装依赖并编译前端页面

```bash
# 安装依赖
npm install
go mod tidy

# 编译前端
npm run build
```

6. 运行

```bash
# 运行
go run ./cmd
```

---

## ⚙️ 配置

默认配置文件路径：`config/config.yaml` 和 `.env`

参考[文档](./docs/config.md)

---

## 📄 许可证

本项目采用 MIT 许可证，详情请见 [LICENSE](./LICENSE)。
