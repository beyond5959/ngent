# Ngent

[![CI](https://github.com/beyond5959/ngent/actions/workflows/ci.yml/badge.svg)](https://github.com/beyond5959/ngent/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/beyond5959/ngent)](https://go.dev/)
[![License](https://img.shields.io/github/license/beyond5959/ngent)](LICENSE)

> **Web Service Wrapper for ACP-compatible Agents**
>
> Ngent wraps command-line agents that speak the [Agent Client Protocol (ACP)](https://github.com/beyond5959/acp-adapter) into a web service, making them accessible via HTTP API and Web UI.

## What is Ngent?

Ngent acts as a bridge between **ACP-compatible agents** (like Claude Code, Codex, Gemini CLI) and **web clients**:

```
┌─────────────┐     HTTP/WebSocket     ┌─────────┐     JSON-RPC (ACP)     ┌──────────────┐
│  Web UI     │ ◄────────────────────► │  Ngent  │ ◄────────────────────► │  CLI Agent   │
│  /v1/* API  │   SSE streaming        │  Server │   stdio                │  (ACP-based) │
└─────────────┘                        └─────────┘                        └──────────────┘
```

### How it Works

1. **ACP Protocol**: Agents like Claude Code and Codex expose their capabilities through the Agent Client Protocol (ACP) — a JSON-RPC protocol over stdio
2. **Ngent Bridge**: Ngent spawns these CLI agents as child processes and translates their ACP protocol into HTTP/JSON APIs
3. **Web Interface**: Provides a built-in Web UI and REST API for creating conversations, sending prompts, and managing permissions

### Features

- 🔌 **Multi-Agent Support**: Works with any ACP-compatible agent (Codex, Claude Code, Gemini, Qwen, OpenCode)
- 🌐 **Web API**: HTTP/JSON endpoints with Server-Sent Events (SSE) for streaming responses
- 🖥️ **Built-in UI**: No separate frontend deployment needed — the web UI is embedded in the binary
- 🔒 **Permission Control**: Fine-grained approval system for agent file/system operations
- 💾 **Persistent State**: SQLite-backed conversation history across sessions
- 📱 **Mobile Friendly**: QR code for easy access from mobile devices on the same network


## Supported Agents

| Agent | Supported |
|---|---|
| Codex | ✅ |
| Claude Code | ✅ |
| Gemini CLI | ✅ |
| Qwen Code | ✅ |
| OpenCode | ✅ |



## Installation

### Quick Install (recommended for Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | bash

# Or install to a custom directory:
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | INSTALL_DIR=~/.local/bin bash
```

## Run

Start with default settings (LAN-accessible):

```bash
ngent
```

Local-only mode (recommended for security):

```bash
ngent --listen 127.0.0.1:8686 --allow-public=false
```

With authentication:

```bash
ngent --auth-token "your-token"
```

Show all options:

```bash
ngent --help
```

**Default paths:**
- Database: `$HOME/.go-agent-server/agent-hub.db`

Notes:

- `/v1/*` requests must include `X-Client-ID`.

## Quick Check

```bash
curl -s http://127.0.0.1:8686/healthz
curl -s -H "X-Client-ID: demo" http://127.0.0.1:8686/v1/agents
```

## Web UI

Once started, open the URL shown in the startup output:

```
Agent Hub Server started
  [QR Code]
Port: 8686
URL:  http://192.168.1.10:8686/
```

Scan the QR code or open the URL in your browser.

**Features:**

- Create threads with any supported agent
- Chat with streaming responses
- Approve/deny permission requests inline
- Browse conversation history
- Light / dark / system themes
- Works on desktop and mobile

The Web UI is embedded in the binary — no separate installation needed.
