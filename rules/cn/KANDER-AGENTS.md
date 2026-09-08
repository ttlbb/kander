# Kander 工作流规则入口

本入口只负责定位与按需加载, 不默认接管开发流程. 安装器复制同一套 Markdown 原文到规则根, 不生成定制规则文件.

## 作用域

本文件所在目录即「规则根」, 据此确定作用域与下列路径:

| 逻辑名   | 全局安装                       | 项目安装                            |
| -------- | ------------------------------ | ----------------------------------- |
| 规则根   | `~/.agents`                    | `<主 worktree>/.kander/rules`       |
| 命令根   | `~/.local/bin`                 | `<主 worktree>/.kander/bin`         |
| 配置文件 | `~/.config/kander/config.json` | `<主 worktree>/.kander/config.json` |
| 资源目录 | `~/.local/share/kander`        | `<主 worktree>/.kander/share`       |

- 全部设置存配置文件.
- 项目安装的载荷只在主树 `.kander/`, 任务 worktree 共享, 不建副本、镜像或符号链接.
- 下文及各分册中的 `kander` 均指当前作用域入口. 全局安装可用命令根绝对路径或已加入 PATH 的 `kander`; 项目安装必须用 `<命令根>/kander` 的绝对路径 (Windows 为 `<命令根>\kander`), 不替换为 PATH 中的全局命令.

## 先读取配置

- 每个新会话先运行当前作用域 `kander config --json`, 读取规范化配置, 再读取同目录的 `KANDER-BASE-RULES.md`. 最小工具协议不受可选模块开关控制.
- 配置读取或校验失败时停止受影响的 Kander 操作并报告, 不猜测开关值; 用户自有流程和无关问答继续.
- 按下表开关与任务需要加载. 关闭模块不读、不执行、不经引用加载; 用户明确要求可例外加载.

## 语言

- 本套规则有英文与中文两份. 安装器按配置中的 `language` 释出对应的一份 (`cn` 为中文, 其余为英文), 该设置变化后 `kander doctor` 会切换已安装的一份. 两份内容一致, 命令、键名、字段名与取值在两份中完全相同.
- 配置中的 `agent_language` 是与用户沟通的语言. 对用户的每次回复、卡片标题与正文、执行记录、完成报告、审核报告, 以及传给 `kander notify` 与 `kander resume` 的消息都用它. 该值缺失或为空时, 用用户书写所用的语言.
- 任务卡的 `LANGUAGE` 字段对该卡的一切优先于配置, 因此一张卡跨会话、接管和配置变更都保持同一语言. `kander new` 在建卡时记录当时的配置值; 没有该字段的卡用配置.
- `kander start`、`resume` 与 `notify` (直投与恢复) 的提示词包含一句明确的英文指令, 要求用该卡语言 (字段缺失时用配置的 `agent_language`) 沟通, 因此 Agent 在读到本套规则之前就已收到它.
- 提交备注、代码注释与标识符遵循项目自身约定, 不受 `agent_language` 影响.
- 配置中的 `language` 决定 `kander` 命令自身的界面语言以及安装本套规则的哪一份; 不影响 Agent 与用户沟通所用的语言.

## 可选模块

| 配置项                | 分册                            | 何时读                          |
| --------------------- | ------------------------------- | ------------------------------- |
| `rules.code`          | `KANDER-CODE-RULES.md`          | 启用时, 修改代码或验证          |
| `rules.collaboration` | `KANDER-COLLABORATION-RULES.md` | 启用时, 开始协作任务            |
| `rules.git`           | `KANDER-GIT-RULES.md`           | 启用时, 操作分支、提交或集成    |
| `rules.task_intake`   | `KANDER-TASK-INTAKE-RULES.md`   | 启用时, 收到 Bug 或功能开发需求 |
| `rules.task_groups`   | `KANDER-TASK-GROUP-RULES.md`    | 启用时, 规划或执行任务组        |
| `rules.review`        | `KANDER-REVIEW-RULES.md`        | 启用时, 判断审核触发及执行流程  |
| `rules.reporting`     | `KANDER-REPORTING-RULES.md`     | 启用时, 任务结束汇报            |

- 使用看板命令时需要读 `KANDER-KANBAN-RULES.md`.
- `task_groups` 依赖 `git`.
- 关 `task_intake` 照样能手动建卡并跑单卡.
- 关 `git` 后单卡沿用用户自己的工作目录、分支、交付流程.
- 关 `review` 后不自动要审核, 仍可明确调 `kander review`.
- 关 `reporting` 不免除真实填卡片结果和必要执行记录.

## 规则优先级

- 当前会话明确的用户指令 > 离目标文件最近的项目级 `AGENTS.md` 或 `CLAUDE.md` > 用户自己的全局规则 > 已启用的 Kander 模块及当前作用域配置 > 模块默认值.
- 同目录 `AGENTS.md` 与 `CLAUDE.md` 冲突且用户指令未消解时, 只停止受影响操作并询问用户.
