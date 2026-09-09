# Kander

[English](README.md) | **简体中文** | [日本語](README-JA.md)

一个人用看板调度多个 AI Agent.

![Kander 工作流](docs/workflow-cn.svg)

## 1. 快速开始

运行需要 Git, 以及 Codex, Claude, Grok, Cursor 或 Kimi 中至少一个.

Kimi 有两点与其它 Agent 不同. 其一, kimi-code 的提示词不能由命令行传入, 只能输入到面板, 因此 Kimi 任务必须用 `tmux`, `tmux-session` 或 `herdr` 启动; `console` 与 `foreground` 会被拒绝并说明原因. 其二, kimi-code 首次在一个仓库里启动时会先问是否信任该文件夹, 此时它还不接受输入, kander 会等待超时并提示你去面板里回答一次; 答过之后该仓库不再询问. 另外 Kimi 不接受推理档位: 它没有对应的命令行开关, 档位由 `~/.kimi-code/config.toml` 的 `[thinking] effort` 决定.

从 [Releases](https://github.com/dualface/kander/releases) 下载 kander 最新二进制后直接运行即可. 首次启动若尚未安装, 会进入交互向导.

安装完成后即可使用.

4 步上手:

1. 新建一个 Agent 会话, 在里面讨论需求或者任务, 说清楚目标和验收条件. 推荐使用 Agent 的 Plan 模式.
2. 任务确认后, Agent 会询问是否用看板流程启动任务. 确认即可自动启动任务.
3. 有多个需求时, 对每个需求重复步骤 1-2, 不断安排并启动任务.
4. 用命令行界面查看任务状态:

```sh
kander
```

![终端看板](docs/kanban-screenshot-01.png)

> 上图看板内容来自我的真实项目 [https://quicktui.ai](https://quicktui.ai). QuickTUI 是一个远程操作电脑上各种 Agent 的工具, 支持 iOS/Android/macOS/Linux/Windows, 免费使用.

进阶阅读: 幻灯片 [如何高效推进任务](docs/how-to-advance-tasks-efficiently-cn.pdf) (PDF).

## 2. 许可

本项目使用 MIT License, 见 [LICENSE](LICENSE).
