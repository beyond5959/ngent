# FRONTEND TASKS

这个文件原本记录了 Web UI 启动阶段的 F0-F9 拆分任务。
随着前端已经完成多轮实现和重设计，原任务单里的大量细节已经不再是现行事实，因此不再保留旧版逐项清单。

当前前端的真实来源请以这些文件为准：

- `docs/FRONTEND_SPEC.md`
- `docs/ACCEPTANCE.md`
- `docs/DECISIONS.md`
- `PROGRESS.md`
- `AGENTS.md`

已从旧任务单中移除的过期实现假设包括：

- 可编辑或持久化的浏览器 `Client ID`
- 基于 `EventSource` 的 SSE 客户端
- 提交 `web/dist/` 到仓库
- Web UI 复用 CLI 的 ASCII `NGENT` 品牌块

如果后续需要重新拆分前端任务，应基于当前实现和上述文档重新生成，而不是恢复这份早期启动计划。
