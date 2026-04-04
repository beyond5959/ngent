# Ngent

[![CI](https://github.com/beyond5959/ngent/actions/workflows/ci.yml/badge.svg)](https://github.com/beyond5959/ngent/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/beyond5959/ngent)](https://go.dev/)
[![License](https://img.shields.io/github/license/beyond5959/ngent)](LICENSE)

[English](README.md) | [简体中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md)

> **Envoltorio de servicio web para agentes compatibles con ACP**
>
> Ngent envuelve agentes de línea de comandos compatibles con el [Agent Client Protocol (ACP)](https://github.com/beyond5959/acp-adapter) en un servicio web para que puedas acceder a ellos mediante API HTTP y una interfaz web.

## ¿Qué es Ngent?

Ngent actúa como puente entre **agentes compatibles con ACP** (como Codex, Kimi CLI, OpenCode, etc.) y **clientes web**:

```
┌─────────────┐     HTTP/WebSocket     ┌─────────┐     JSON-RPC (ACP)     ┌──────────────┐
│ Interfaz web│ ◄────────────────────► │  Ngent  │ ◄────────────────────► │ Agente CLI   │
│   /v1/* API │   streaming SSE        │ Server  │   stdio                │ (basado en   │
└─────────────┘                        └─────────┘                        │ ACP)         │
                                                                           └──────────────┘
```

### Cómo funciona

1. **Protocolo ACP**: agentes como Codex, Kimi CLI y OpenCode exponen sus capacidades mediante ACP, un protocolo JSON-RPC sobre `stdio`.
2. **Puente de Ngent**: Ngent ejecuta esos agentes CLI como procesos hijo y traduce su protocolo ACP a API HTTP/JSON.
3. **Interfaz web**: proporciona una UI integrada y una API REST para crear conversaciones, enviar prompts y gestionar permisos.

### Características

- 🔌 **Soporte multiagente**: funciona con cualquier agente compatible con ACP (Codex, Kimi CLI, OpenCode, etc.)
- 🌐 **API web**: endpoints HTTP/JSON con Server-Sent Events (SSE) para respuestas en streaming
- 🖥️ **UI integrada**: no hace falta desplegar un frontend aparte; la interfaz web va embebida en el binario con soporte para varios idiomas
- 🔒 **Control de permisos**: sistema de aprobación granular para operaciones de archivos y sistema
- 💾 **Estado persistente**: historial de conversaciones respaldado por SQLite entre sesiones
- 📱 **Compatible con móvil**: código QR para acceder fácilmente desde dispositivos de la misma red
- 🌿 **Control de ramas Git**: consulta la información de ramas por hilo y cambia ramas locales directamente desde la interfaz web integrada
- 🧾 **Inspección de diffs Git**: muestra resúmenes de cambios del árbol de trabajo por sesión, lista archivos tracked y untracked y abre vistas previas por archivo

## Agentes soportados

| Agente | Soportado |
|---|---|
| Codex | ✅ |
| Claude Code | ✅ |
| Cursor CLI | ✅ |
| Gemini CLI | ✅ |
| Kimi CLI | ✅ |
| Qwen Code | ✅ |
| OpenCode | ✅ |
| BLACKBOX AI | ✅ |

## Instalación

### Instalación rápida (recomendada para Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | bash

# O instalar en un directorio personalizado:
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | INSTALL_DIR=~/.local/bin bash
```

## Ejecución

Iniciar con la configuración predeterminada (solo local):

```bash
ngent
```

Modo accesible desde LAN (permite conexiones desde otros dispositivos):

```bash
ngent --allow-public=true
```

Puerto personalizado:

```bash
ngent --port 8080
```

Con autenticación:

```bash
ngent --auth-token "your-token"
```

Directorio de datos personalizado:

```bash
ngent --data-path /path/to/ngent-data
```

Mostrar todas las opciones:

```bash
ngent --help
```

**Rutas predeterminadas:**
- Directorio de datos: `$HOME/.ngent/`
- Base de datos: `$HOME/.ngent/ngent.db`
- Adjuntos: `$HOME/.ngent/attachments/<category>/`

Notas:

- Las solicitudes a `/v1/*` deben incluir `X-Client-ID`.

## Verificación rápida

```bash
curl -s http://127.0.0.1:8686/healthz
curl -s -H "X-Client-ID: demo" http://127.0.0.1:8686/v1/agents
```

## Interfaz web

Abre la URL mostrada al arrancar (por ejemplo, `http://127.0.0.1:8686/`).
