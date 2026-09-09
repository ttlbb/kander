# Kander

**English** | [简体中文](README-CN.md) | [日本語](README-JA.md)

One person schedules multiple AI agents with a kanban board.

![Kander workflow](docs/workflow-en.svg)

## 1. Quick Start

Running requires Git, plus at least one of Codex, Claude, Grok, Cursor, or Kimi.

Kimi differs from the other agents in two ways. First, kimi-code cannot take its prompt on the command line, only in the panel, so Kimi tasks must launch through `tmux`, `tmux-session`, or `herdr`; `console` and `foreground` are refused with an explanation. Second, the first time kimi-code starts in a repository it asks whether the folder is trusted and accepts no input until answered, so kander waits for the timeout and tells you to answer once in the panel; it stops asking for that repository afterwards. Kimi also takes no reasoning effort: it has no command-line switch for one, and the level comes from `[thinking] effort` in `~/.kimi-code/config.toml`.

Download the latest kander binary from [Releases](https://github.com/dualface/kander/releases) and run it directly. On first launch, if not yet installed, an interactive wizard starts.

Once installation finishes, it is ready to use.

Four steps to get going:

1. Start an agent session and discuss the requirement or task there, making the goal and acceptance criteria clear. The agent's Plan mode is recommended.
2. Once the task is confirmed, the agent asks whether to launch it through the kanban flow. Confirm, and the task launches automatically.
3. When you have multiple requirements, repeat steps 1-2 for each one, continuously scheduling and launching tasks.
4. Check task status with the command-line interface:

```sh
kander
```

![Terminal kanban](docs/kanban-screenshot-01.png)

> The board contents above come from my real project [https://quicktui.ai](https://quicktui.ai). QuickTUI is a tool for remotely operating the agents on your computer; it supports iOS/Android/macOS/Linux/Windows and is free to use.

Further reading: the slides [How to Advance Tasks Efficiently](docs/how-to-advance-tasks-efficiently-en.pdf) (PDF).

## 2. License

This project is under the MIT License; see [LICENSE](LICENSE).
