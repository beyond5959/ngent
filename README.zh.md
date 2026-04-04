# Ngent

[![CI](https://github.com/beyond5959/ngent/actions/workflows/ci.yml/badge.svg)](https://github.com/beyond5959/ngent/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/beyond5959/ngent)](https://go.dev/)
[![License](https://img.shields.io/github/license/beyond5959/ngent)](LICENSE)

[English](README.md) | [简体中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md)

> **兼容 ACP 协议 Agent 的 Web 服务封装**
>
> Ngent 将支持[智能体客户端协议（ACP）](https://github.com/beyond5959/acp-adapter)的命令行智能体封装为 Web 服务，让你可以通过 HTTP API 与 Web 界面轻松调用。

## Ngent 是什么？

Ngent 是 **ACP 兼容智能体**（如 Codex、Kimi CLI、OpenCode 等）与 **Web 客户端** 之间的桥梁：

```
┌─────────────┐     HTTP/WebSocket     ┌─────────┐     JSON-RPC (ACP)     ┌──────────────┐
│  Web 界面   │ ◄────────────────────► │  Ngent  │ ◄────────────────────► │  命令行智能体 │
│  /v1/* API  │   SSE 流式传输         │  服务端 │   标准输入输出         │  (基于 ACP)  │
└─────────────┘                        └─────────┘                        └──────────────┘
```

### 工作原理

1. **ACP 协议**：Codex、Kimi CLI、OpenCode 等智能体通过智能体客户端协议（ACP）暴露自身能力——这是一种基于标准输入输出（stdio）的 JSON-RPC 协议。
2. **Ngent 桥接**：Ngent 将这些命令行智能体作为子进程启动，并将其 ACP 协议转译为 HTTP/JSON API。
3. **Web 界面**：提供内置的 Web 界面与 REST API，用于创建对话、发送提示词及管理权限。

### 核心特性

- 🔌 **多智能体支持**：兼容任意 ACP 智能体（Codex、Kimi CLI、OpenCode 等）
- 🌐 **Web API**：HTTP/JSON 端点，支持 SSE 流式响应
- 🖥️ **内置界面**：Web 界面直接嵌入二进制文件，无需单独部署前端，支持多种语言
- 🔒 **权限管控**：对智能体的文件/系统操作提供细粒度的审批机制
- 💾 **持久化状态**：基于 SQLite 的会话历史，跨会话保留
- 📱 **移动端友好**：生成二维码，方便同网络环境下的移动设备快速访问
- 🌿 **Git 分支能力**：支持在内置 Web UI 中查看线程工作目录的分支信息，并直接切换本地分支
- 🧾 **Git Diff 检查能力**：支持按会话展示工作区变更摘要、列出已跟踪与未跟踪文件，并打开单文件预览查看

## 支持的智能体

| 智能体 | 支持状态 |
|---|---|
| Codex | ✅ |
| Claude Code | ✅ |
| Cursor CLI | ✅ |
| Gemini CLI | ✅ |
| Kimi CLI | ✅ |
| Qwen Code | ✅ |
| OpenCode | ✅ |
| BLACKBOX AI | ✅ |

## 安装

### 快速安装（推荐 Linux/macOS 用户）

```bash
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | bash

# 或安装到自定义目录：
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | INSTALL_DIR=~/.local/bin bash
```

## 运行

默认启动（仅本地访问）：

```bash
ngent
```

局域网可访问模式（允许其他设备连接）：

```bash
ngent --allow-public=true
```

自定义端口：

```bash
ngent --port 8080
```

启用身份验证：

```bash
ngent --auth-token "your-token"
```

自定义数据目录：

```bash
ngent --data-path /path/to/ngent-data
```

查看全部选项：

```bash
ngent --help
```

**默认路径：**
- 数据目录：`$HOME/.ngent/`
- 数据库：`$HOME/.ngent/ngent.db`
- 附件：`$HOME/.ngent/attachments/<category>/`

注意事项：

- 访问 `/v1/*` 接口时，请求头中必须携带 `X-Client-ID`。

## 快速验证

```bash
curl -s http://127.0.0.1:8686/healthz
curl -s -H "X-Client-ID: demo" http://127.0.0.1:8686/v1/agents
```

## Web 界面

打开启动日志中显示的地址（例如 `http://127.0.0.1:8686/`）。
