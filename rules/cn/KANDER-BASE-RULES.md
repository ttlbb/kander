# Kander 最小工具协议

- 使用 Kander 前按 `KANDER-AGENTS.md` 读取当前作用域配置. 本文件只约束工具和与用户沟通所用的语言, 不规定交流方式、架构、代码验证、Git 或自动审核流程.
- 与用户沟通, 以及书写卡片、记录和报告, 都用配置中的 `agent_language`; 处理任务卡时用该卡的 `LANGUAGE`, 见 `KANDER-AGENTS.md`「语言」.
- 使用看板命令时读取 `KANDER-KANBAN-RULES.md` 的结构、状态、领取、通知与恢复协议; 无须启用可选模块.

**单次审核**

- `kander review` 是可明确调用的单次审核工具, 参数为 `[agent] [--task <id>]... [--run-id <id>] [--batch-id <id>] [--previous-run-id <id>] [--requirements-file <JSON>] [--advance-file <JSON>] <CWD> <base-commit> <commit> <role> <task-goal|absolute spec path> [review-context] [reviewed-commit]`.
- 目标须为干净 Git worktree, base 须为 commit 祖先.
- 不带 `--task` 时 review 不定位看板. 带任务时, 标志位于 CWD 之前, 重复的任务 ID 去重, 看板从目标 CWD 定位. 卡片须为 working/review 中的目录卡, 语言一致且任务组归属兼容.
- 绑定任务的审核要求明确的批次 ID. 新批次还要求一个 JSON 需求文件, 把四个角色都标为 `required` 或 `N/A: <reason>`. 这些需求由用户/项目规则与环节策略推导; 该文件记录这一决定, 不构成批准.
- 缺少 run ID 时自动生成并打印到 stderr. 只为恢复或完成发布才复用该 ID, 它绝不会再启动一个 Reviewer. 输入变化即冲突. 同一目标上 PM 与 QA 用不同的 run ID. 进程崩溃后的重试记录为中断证据, 绝不记为 PASS.
- 原始输出、日志、输入快照、sidecar 与 manifest 留在每张卡的 `reviews/<run_id>/`. 机器所有的 REVIEWS 区每次运行一行 JSON 索引. 工具执行成功与语义上的 PASS 是两回事. 不得把这些产物当临时报告编辑或删除.
- 发布按卡原子. 部分发布以非零退出, 保留已成功的卡并逐一报告结果. 看板事务中断后, 按其维护要求运行 `kander init`; 然后用同一 run ID 重试完全相同的 review 调用. 未完整发布的运行不能确立完成.
- 命令维持 Reviewer 只读与输出校验.
- 调用不启用完整审核或 Git 流程, 不要求目标分支为 `develop`.

- 开关不改变参数、数据结构、路径校验或进程隔离. 工具边界失败须报告, 禁用普通文件操作或直接控制 Agent 绕过.

**安装与任务文件**

- 安装由二进制自身完成: 首次运行未安装的 `kander` 进入交互向导, 或运行 `kander install` 重跑.
- Windows 不自动修改 `PATH`.
- 含特殊字符的自动化须用进程 API 的 argv 数组直调命令根的 `kander`, 禁拼 PowerShell/cmd 命令字符串.
- 所有平台的执行 Agent 与 Reviewer 均从 UTF-8 临时文件读取完整任务, 启动参数只含 CLI 必需控制项与一句文件路径指令.
- 文件要求 Agent 完成后尝试删除.
- 删除失败或遗留不影响结果.

**权限与清理**

- 审核私有目录与文件仅当前用户可访问: POSIX `0600`/`0700`.
- Windows 创建即使用关闭继承的受保护 DACL, 禁先发布再收紧.
- Windows 审核根句柄不共享 WRITE/DELETE, 持有至敏感文件写入、Reviewer 运行、进程树收集及清理结束, 阻止改名与原地 reparse 切换.
- 清理从固定句柄逐层拒绝 reparse point, 预算有界, 失败即审核失败.
- 配置不检查、迁移或收紧权限: POSIX 保留既有 mode, 新对象遵循 umask.
- Windows 新配置文件与目录继承父 ACL.
- Windows 配置从卷/UNC anchor 逐分量拒绝 reparse point, 读取与原子替换用固定句柄.
- 看板和 Git exclude 同样拒绝符号链接、junction 等 reparse point.
- Git exclude 保留既有 ACL, 在同一固定句柄内去重追加.
- 禁绕过命令直接操作这些边界.
