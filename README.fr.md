# Ngent

[![CI](https://github.com/beyond5959/ngent/actions/workflows/ci.yml/badge.svg)](https://github.com/beyond5959/ngent/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/beyond5959/ngent)](https://go.dev/)
[![License](https://img.shields.io/github/license/beyond5959/ngent)](LICENSE)

[English](README.md) | [简体中文](README.zh.md) | [Español](README.es.md) | [Français](README.fr.md)

> **Enveloppe de service web pour agents compatibles ACP**
>
> Ngent encapsule des agents en ligne de commande compatibles avec l'[Agent Client Protocol (ACP)](https://github.com/beyond5959/acp-adapter) dans un service web afin de les exposer via une API HTTP et une interface web.

## Qu'est-ce que Ngent ?

Ngent sert de passerelle entre des **agents compatibles ACP** (comme Codex, Kimi CLI, OpenCode, etc.) et des **clients web** :

```
┌─────────────┐     HTTP/WebSocket     ┌─────────┐     JSON-RPC (ACP)     ┌──────────────┐
│ Interface   │ ◄────────────────────► │  Ngent  │ ◄────────────────────► │ Agent CLI    │
│ web /v1/*   │   streaming SSE        │ Server  │   stdio                │ (basé sur    │
└─────────────┘                        └─────────┘                        │ ACP)         │
                                                                           └──────────────┘
```

### Fonctionnement

1. **Protocole ACP** : des agents comme Codex, Kimi CLI et OpenCode exposent leurs capacités via ACP, un protocole JSON-RPC au-dessus de `stdio`.
2. **Passerelle Ngent** : Ngent lance ces agents CLI comme processus enfants et traduit leur protocole ACP en API HTTP/JSON.
3. **Interface web** : fournit une UI intégrée et une API REST pour créer des conversations, envoyer des prompts et gérer les permissions.

### Fonctionnalités

- 🔌 **Support multi-agents** : fonctionne avec n'importe quel agent compatible ACP (Codex, Kimi CLI, OpenCode, etc.)
- 🌐 **API web** : endpoints HTTP/JSON avec Server-Sent Events (SSE) pour les réponses en streaming
- 🖥️ **UI intégrée** : aucun déploiement frontend séparé n'est nécessaire, l'interface web est embarquée dans le binaire avec support pour English, 简体中文, Español et Français
- 🔒 **Contrôle des permissions** : système d'approbation fin pour les opérations fichiers/système
- 💾 **État persistant** : historique des conversations basé sur SQLite entre les sessions
- 📱 **Compatible mobile** : QR code pour un accès facile depuis un appareil du même réseau
- 🌿 **Intégration Git** : afficher les infos de branche, changer de branche et afficher les diffs du répertoire de travail par fil

## Agents pris en charge

| Agent | Pris en charge |
|---|---|
| Codex | ✅ |
| Claude Code | ✅ |
| Cursor CLI | ✅ |
| Gemini CLI | ✅ |
| Kimi CLI | ✅ |
| Qwen Code | ✅ |
| OpenCode | ✅ |
| BLACKBOX AI | ✅ |

## Installation

### Installation rapide (recommandée pour Linux/macOS)

```bash
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | bash

# Ou installer dans un répertoire personnalisé :
curl -sSL https://raw.githubusercontent.com/beyond5959/ngent/master/install.sh | INSTALL_DIR=~/.local/bin bash
```

## Exécution

Démarrer avec la configuration par défaut (local uniquement) :

```bash
ngent
```

Mode accessible sur le LAN (autorise les connexions depuis d'autres appareils) :

```bash
ngent --allow-public=true
```

Port personnalisé :

```bash
ngent --port 8080
```

Avec authentification :

```bash
ngent --auth-token "your-token"
```

Répertoire de données personnalisé :

```bash
ngent --data-path /path/to/ngent-data
```

Afficher toutes les options :

```bash
ngent --help
```

**Chemins par défaut :**
- Répertoire de données : `$HOME/.ngent/`
- Base de données : `$HOME/.ngent/ngent.db`
- Pièces jointes : `$HOME/.ngent/attachments/<category>/`

Remarques :

- Les requêtes vers `/v1/*` doivent inclure `X-Client-ID`.

## Vérification rapide

```bash
curl -s http://127.0.0.1:8686/healthz
curl -s -H "X-Client-ID: demo" http://127.0.0.1:8686/v1/agents
```

## Interface web

Ouvrez l'URL affichée au démarrage, par exemple `http://127.0.0.1:8686/`.
