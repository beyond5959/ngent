# FRONTEND SPEC

## 1. 目标

Ngent 的 Web UI 是随 Go 二进制分发的内嵌 no-framework Vite + TypeScript SPA。
当前目标不是做“聊天 SaaS 外壳”，而是提供一个冷静、专业、桌面工作台式的本地 Agent 操作界面，同时保持协议、SSE、权限、session/history 语义不变。

## 2. 当前界面骨架

- 桌面端采用两段式工作区：
  - 左侧为可折叠的分组 `Threads` 导航 rail；每个 thread header 下面直接内联该 thread 的 session 行，负责线程/会话切换、`New session`、session 刷新与线程操作。
  - 右侧为主聊天工作区，承担头部元信息、消息流、输入区、底部浮动计划卡和流式状态。
- 不再保留独立的 `Session History` 中间栏；thread 和 session 在同一条左侧 rail 中分组呈现。
- 无活跃线程时，主区展示一个锚定式空态面板，只保留一个主操作 `New Agent`。
- `New Agent` 弹窗采用 working-directory-first 流程：先选绝对路径，再选 Agent，`Advanced options` 折叠显示。
- `Settings` 为浏览器级抽屉，只负责 `Server URL`、`Bearer Token`、`Theme`。
- 移动端会折叠这条合并后的左侧 rail，但保持同一套信息层级，不靠“单纯隐藏桌面元素”维持布局。

## 3. 行为约束

- 不允许改后端协议或运行时行为；前端重设计必须是 presentation-only。
- 必须保留 `activeStreamMsgId` 哨兵逻辑，避免订阅回流把流式气泡冲掉。
- 必须保留 history load guards，避免线程或 session 切换后把旧请求写回当前视图。
- 必须继续使用 POST SSE：
  - `fetch` + `ReadableStream`
  - 不能退回 `EventSource`
- 只有 finalized agent 文本才走 Markdown 渲染；流式中的消息气泡只能写入 `textContent`。
- 任何通过 `innerHTML` 注入的 Markdown 结果都必须继续调用 `bindMarkdownControls(...)`。
- 如果调整 DOM 结构，必须同步迁移现有查询选择器与事件绑定，不能破坏线程切换、分组 session 浏览、发送、取消、权限审批、附件、设置面板等交互。

## 4. 技术架构

| 层次 | 当前实现 |
|---|---|
| 构建 | Vite 6 + TypeScript 5.6 |
| 渲染 | 原生 DOM，单入口 `main.ts` 驱动，禁止引入 React/Vue/Svelte |
| 样式 | 单文件 `style.css` + CSS 自定义属性；light/dark 独立调校 |
| 状态 | 轻量 `AppStore` 发布订阅；仅持久化浏览器设置 |
| API | `api.ts` 封装 HTTP/JSON，请求头固定带兼容性 `X-Client-ID: ngent-web-ui` |
| 流式 | `sse.ts` 中 `TurnStream` 通过 `fetch` + `ReadableStream` 解析 `event:/data:` 块 |
| Markdown | `marked` + `highlight.js`；代码块复制/展开控件通过 `bindMarkdownControls(...)` 绑定 |
| 嵌入 | `internal/webui/webui.go` 用 `//go:embed web/dist` 提供 SPA fallback handler |

## 5. 关键源码映射

| 文件 | 角色 |
|---|---|
| `internal/webui/web/src/main.ts` | 应用入口、Shell 渲染、分组 thread/session rail、消息区/输入区 DOM 绑定、流式生命周期协调 |
| `internal/webui/web/src/style.css` | 视觉 token、布局、主题、组件样式、Markdown/代码块排版 |
| `internal/webui/web/src/api.ts` | `/v1/*` API 调用、错误归一化、兼容性请求头 |
| `internal/webui/web/src/sse.ts` | POST SSE 解析、turn 事件分发、断开/终止处理 |
| `internal/webui/web/src/store.ts` | 浏览器内状态容器；只持久化 `authToken`、`serverUrl`、`theme` |
| `internal/webui/web/src/markdown.ts` | Markdown 安全渲染、高亮、代码块复制/展开按钮绑定 |
| `internal/webui/web/src/components/new-thread-modal.ts` | `New Agent` 弹窗、绝对路径校验、Agent 选择、Advanced options 折叠 |
| `internal/webui/web/src/components/settings-panel.ts` | 浏览器级设置抽屉 |
| `internal/webui/web/src/components/permission-card.ts` | 权限审批卡片与倒计时 |

## 6. 状态与持久化

- `localStorage` 只保存：
  - `ngent:authToken`
  - `ngent:serverUrl`
  - `ngent:theme`
- 历史消息、线程列表、分组 rail 所需的 session 列表、流式状态、输入草稿、附件草稿都不持久化到 `localStorage`。
- 启动时会清理旧的 `ngent:clientId` 遗留值；浏览器不再暴露可编辑的 Client ID 设置。

## 7. API 映射

| 前端操作 | API |
|---|---|
| 初始化 | `GET /v1/agents`、`GET /v1/threads` |
| 新建线程 | `POST /v1/threads` |
| 更新线程标题 / agent options | `PATCH /v1/threads/{threadId}` |
| 读取线程历史 | `GET /v1/threads/{threadId}/history?includeEvents=1` |
| 读取 session 列表 | `GET /v1/threads/{threadId}/sessions` |
| 读取 provider transcript replay | `GET /v1/threads/{threadId}/session-history?sessionId=...` |
| 读取 slash commands | `GET /v1/threads/{threadId}/slash-commands` |
| 读取/设置 config options | `GET/POST /v1/threads/{threadId}/config-options` |
| 发起 turn | `POST /v1/threads/{threadId}/turns`（POST SSE） |
| 取消 turn | `POST /v1/turns/{turnId}/cancel` |
| 权限审批 | `POST /v1/permissions/{permissionId}` |

## 8. 视觉与交互原则

- Web UI 使用简洁的 icon + wordmark，不复用 CLI 的 ASCII `NGENT` 品牌块。
- 去除背景网格、装饰性 orb、重玻璃感、普遍悬浮 hover、过度胶囊化 badge。
- 主聊天区是视觉中心；合并后的 thread/session rail 退后为辅助信息层。
- 线程与 session 在同一条分组 rail 中使用工具型列表呈现，而不是拆成多个主容器或整屏卡片堆叠。
- 输入区更接近编辑器工作面板，而不是大圆角聊天气泡。
- reasoning / plan / tool call / permission / markdown 采用统一的文档式分区语言，用排版和边界区分，不靠一堆不一致的卡片外观。
