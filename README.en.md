🌐 [中文](./README.md) | [English](./README.en.md)

# CS2Panel

> Lightweight and user-friendly management tool for CS2 (Counter-Strike 2) game servers.
>
> This repository contains:
>
> - Backend written in **Go (Golang)**
> - Frontend built with **React**

If you like this project, please give it a ⭐ Star!

## 📦 Installation & Running

1. Clone the repository

```bash
# Clone the repository
git clone https://github.com/VanVodkaer/CS2Panel
cd CS2Panel
```

2. Rename the `config/config.yaml.example` file to `config.yaml` and edit it  
3. Rename the `.env.example` file in the root directory to `.env` and edit it  
4. Run Docker  
5. Install dependencies and build the frontend

```bash
# Install dependencies
npm install
go mod tidy

# Build frontend
npm run build
```

6. Run the application

```bash
# Run
go run ./cmd
```

---

## ⚙️ Configuration

Default configuration file paths: `config/config.yaml` and `.env`

See the [documentation](./docs/config.md) for more details.

---

## 📄 License

This project is licensed under the MIT License. See [LICENSE](./LICENSE) for more information.
