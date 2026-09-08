# Kander 看板命令协议

## 适用范围

- 实际使用看板命令时读取本文件. 先按 `KANDER-AGENTS.md` 读取配置, 再按需加载已启用模块; 本文不默认要求建卡引导、Git 工作流、审核或固定报告.
- 本文件约束当前作用域 `kander` 命令管理的看板. 用户指令和目标项目规则优先; 卡片只保存任务契约和执行记录, 不覆盖用户决策, 项目规则或安全门禁.
- Agent 操作看板前先完整读取规则根下的本文件, 再读目标卡片.

  下文命令入口按 `KANDER-AGENTS.md`「作用域」选择.

## 存储与定位

- `kanban/` 是不进 Git 的本机共享数据, 唯一实例位于主 worktree 根目录, 只供同主机同文件系统的 Agent 使用.

  任务 worktree 不建副本, 镜像或符号链接.

  远程 Agent 不可见.

  Windows 上符号链接、junction 和其他 reparse point 一律视为不安全入口, `kander` 通过已校验的 Win32 句柄读写和迁移.

  POSIX 继续使用 no-follow 文件操作.

  任一安全校验失败都停止, 禁用文件管理器或普通路径 API 绕过.

- 定位顺序是 `KANBAN_DIR` -> 当前 Git 仓库主 worktree 的 `kanban/` -> 从当前目录向上查找 `kanban/`.

  `KANBAN_DIR` 仅用于测试, 非 Git 项目或明确覆盖.

  正常 Git 项目从任意 worktree 这样定位:

```sh
MAIN_WORKTREE="$(git worktree list --porcelain | sed -n '1s/^worktree //p')"
KANBAN_DIR="$MAIN_WORKTREE/kanban"
```

- `kanban/` 不属于 Kander 安装载荷, 定位不因全局或项目安装而改变.

  看板操作本身不建分支, 不提交, 不 push, 不审核.

  卡片对应的代码任务仍按项目规则执行.

  禁止提交 `kanban/` 或修改项目 `.gitignore` 传播它.

## 命令契约

创建入口、查询和迁移状态只用 `kander`, 禁用 `mv`, `cp` 或文件管理器替代; 正文编辑与同状态形态升级按「入口与文档」「任务规模与分组」执行.

```text
kander config [--json]
kander install
kander doctor
kander init [--maintenance] [project-path]
kander list [--mobile] [backlog|todo|working|review|done|archived|trash]
kander show [--json] <task-id>
kander new [--large] [--language <agent language>] <feature|bug|chore|research> <slug> <title...>
kander move <task-id> <backlog|todo|working|review|done|archived|trash> [--owner <agent>] [--result <result>] [--reason <reason> --decision <reference>] [--duplicate-of <task-id>] [--expect-revision <revision>] [--dispatch-id <id> --execution-epoch <epoch>] [--delivery-commit <full-SHA>] [--disposition <card-relative-path>]
kander pick [task-id]
kander update <task-id> --document <relative-path> --file <UTF8-input> --expect-revision <revision> [--contract-decision-file <UTF8-decision>] [--dispatch-id <id> --execution-epoch <epoch>]
kander start [--agent <configured-agent>] [--launcher auto|tmux|tmux-session|herdr|foreground|console] [task-id]
kander resume [--agent <configured-agent>] [--timeout SECONDS] (--message TEXT | --message-file FILE) [--launcher ...] [--dispatch-id <id>] [--kind fix|sync|wrap-up] [--base <full-SHA>] [--evidence-file <JSON>] <task-id>
kander notify [--pane HERDR-PANE-ID] [--timeout SECONDS] (--message TEXT | --message-file FILE) [--dispatch-id <id>] [--kind fix|sync|wrap-up] [--base <full-SHA>] [--evidence-file <JSON>] <task-id>
kander dispatch prepare <absolute-UTF8-intent.json>
kander dispatch authorize-wrap-up <absolute-UTF8-request.json>
kander dispatch show <task-id> <dispatch-id>
kander dispatch fail|cancel <task-id> <dispatch-id> <dispatch-revision> <reason>
kander dismiss [--timeout SECONDS] <task-id>
kander check [--all] [task-id ...]
kander guard-write <path>
kander subscribe [--refresh SECONDS] [--heartbeat SECONDS] <task-group> <task-id>... [--watch <task-id|task-group-id>...]
kander coordinator show|claim|reconcile ...       # 见「协调者检查点」
kander review ...                                 # 单一 review 入口: KANDER-BASE-RULES.md 与「审核证据完成门禁」
```

`kander show` 在卡片正文前输出当前状态与绝对路径, 供写卡前重新定位; `kander move` 成功输出迁移后的新路径. `kander show --json` 返回已提交的 `text`, `revision`, `operation_id` 与 `entry` 位置. `kander guard-write` 是建议性的写入前检查: 放行退出 0, 拒绝非零, 但该检查与外部写入不是原子的. 它不覆盖任意 shell 命令, 外部工具或内部写入. Agent 必须用 `update` 编辑卡片.

`kander pick [task-id]` 经与 `kander move <task-id> todo` 相同的门禁把 `backlog/` 卡迁入 `todo/`; 两者可互换. 不给任务 ID 时列出 `backlog/` 卡并询问选哪一张; 自动化传入 ID. 状态名之后的 move 选项见「入口与文档」与「持久派发」. `kander review` 之后的 `plan`, `extend-plan`, `assign`, `disposition`, `map-legacy`, `aggregate`, `advance`, `close` 与 `progress` 选择证据子命令; 可选的 Reviewer 参数只接受 `KANDER-REVIEW-RULES.md` 中的四个 Reviewer 名.

新写使用 `@kander_session`、`@kander_project` 与 `# kander-notify:`.

读取、反查与身份匹配时, `@kander_session` 与 `@onevoke_session` 同为会话标记, `@kander_project` 与 `@onevoke_project` 同为项目标记, `# kander-notify:` 与 `# onevoke-notify:` 同为 notify 指令前缀.

任一侧非空且与卡片会话一致即命中, 两侧皆空才算标记缺失.

`resume` 及 `notify` 的恢复/接管通道成功后必须把新容器地址回写到卡片 `WINDOW`.

foreground/console 归一为 launcher 名.

启动或存活校验失败时, 仅在本操作仍持有当前 revision 时恢复调用前原文; 否则保留实际状态与更新的正文并报告冲突.

`notify` 与 `resume` 都不迁移卡片. 在未绑定的兼容流程中: `review -> working` 由被通知或被唤醒的执行 Agent 在处理事项前自行执行 `kander move <task-id> working`, 该状态变化即为旧版「已收到并开工」的回执. 持久派发则要求其原子的已接受回执.

`notify` 探查已记录的 `WINDOW` 时按三类处理. 本分类优先于后文的恢复概述.

- (1) **忙态**: pane 存在且 Agent/会话匹配, 但 herdr 非 `idle`/`done` 或 tmux 在 copy-mode.
  - 在同一 `--timeout` 预算内按固定间隔重试探查.
  - 超时非零返回「目标 Agent 忙, 未投递」.
  - 不启动恢复实例、不创建正文载荷、不改卡片状态或正文.
  - tmux 会话标记缺失属身份不可证, 不算忙或过期, 走恢复通道.
- (2) **地址过期**: pane 不存在/已退出, Agent/会话不匹配, 或 tmux 前台进程不匹配.
  - herdr 按 Agent 与会话从 `pane list` 反查.
  - tmux 用 `list-panes -a -F` 按非空 `@kander_session` 或 `@onevoke_session`、`pane_dead=0`、Agent 可执行名三重过滤.
  - 唯一命中后重跑完整目标校验, 通过才回写兼容地址并直投.
  - `tmux` 回写 session id; `tmux-session` 回写 session name.
  - 空 `WINDOW` 的 tmux 旧卡不反查.
- (3) **反查失败**: 0 个/多个命中、反查不可用或新 pane 复检失败时, 按既有链路恢复.
  - 同时报原地址过期原因与反查原因.

显式 `--pane` 覆盖不做过期地址反查.

**持久派发**

- 本节优先于本文件的旧版通知条款. `review/` 中属于任务组或已带 `DISPATCH_ID` 的卡自动使用持久派发. 显式 `--kind fix|sync|wrap-up`, `--dispatch-id`, `--base <full-SHA>` 或 `--evidence-file` 为 `working/` 卡及非组卡选用它; 都不给时, 发给 `working/` 卡的消息是未绑定的旧版消息. 默认 kind 为 fix. 集成后的完成须显式选择 wrap-up. 不带 base 的新请求使用当前工作目录的 Git HEAD.
- `notify` 与 `resume` 接受这些参数. 省略 ID 时在任何发送前先输出生成的 ID; 记录该 ID 并在重试时复用. 同一 ID 必须保持相同的任务, 消息, kind, 基线与引用. 准备是幂等的; 输入改变即冲突. 重试默认复用原基线与截止时间.
- `dispatch prepare` 的显式 UTF-8 intent 含 `dispatch_id` (可选), `task_id`, `kind`, `message`, `base`, 可选 `references` (任务 ID 加卡片相对 `path`), `created_at` 与 `confirm_by`. 默认接受截止时间为创建后 120 秒; 新的 notify/resume 请求用各自的 timeout. 它不因重启/重试重置, 也不是工作完成截止时间.
- 新 fix intent 须有 `evidence.fix`: 批次 ID, 非空的 run/finding 引用及精确的 `previous_run_id`, 以及对这些 finding 上所有既有作者原件及其显式谱系的引用. 每个作者引用保留其原作者, run/finding ID, 记录 ID, 任务 ID 与卡片相对 artifact 路径. 不要求尚未写出的 disposition. 当前批次目标必须与 base 一致; 分派必须包含该任务; 所有已发布的原件与副本必须通过校验. 已被取代的轮次, 缺失/外来的 finding 与不完整的发布在发送前拒绝. 旧版散文只用审核生产者的显式映射, 决不用推断的 ID.
- 通过 `--evidence-file <JSON>` 向 notify/resume 传入绑定, 或在 dispatch prepare 中作为 `evidence` 传入. 同 ID 重试保留已冻结的绑定; 复用时省略 evidence 文件. 之后的作者记录决不改写更早的派发. 通用 `references` 不能替代语义证据. 历史未绑定 intent 仍可读取以供对账, 但不能作为新的 fix/wrap-up 动作发送.
- 新 wrap-up intent 须有 `evidence.wrap_up.git`, 含绝对 CWD, source_commit, target_commit, target_ref (`refs/heads/develop` 或 `refs/remotes/origin/develop`), author 与 basis. 命令绑定首个已规划批次的 base 与计划中为最后一个批次记录的 target (即该批次加入时的 target; 之后的 fix 推进不反映在那里, 见「审核证据完成门禁」), 校验 source 与 target 到所选 develop 引用的祖先关系, 再重读该 ref. `source_commit` 必须等于派发 base, 因此只要它不同于命令运行目录的 HEAD, 就显式传 `--base <source_commit>`. 默认 source 必须等于计划中为最后一个批次记录的 target; 不包含额外未审核的 commit. 经授权 rebase 后, 显式提供 rebased_base: 命令校验 base 祖先关系并比较记录范围与 rebase 后范围的完整补丁, 只归一化 blob hash 与 hunk 行偏移, 保留空白, 上下文, 模式与二进制改动. 补丁不同即拒绝; 主控随即停止并按 `KANDER-TASK-GROUP-RULES.md`「合回与清理前置」报告. fetch 与集成授权仍是主控的职责. 不得以声称的 PR 合并或等价性替代此证据. 集成证据原子存放在同一派发的卡片相对 `dispatches/<id>/integration.json`; 以该 source_commit 完成. 看板校验结构, Git 感知的命令校验 Git.
- 代做 wrap-up 须用 `dispatch authorize-wrap-up <request.json>`, 含 task_id, dispatch_id, expected_revision (派发 revision), author 与 reason. 不存在 intent 时附带完整的 wrap-up `intent`, 使用相同的显式任务/派发 ID. 命令先创建它, 在投递租约下对账回执, 再要求新的身份有效的已停止观察. 缺 SESSION, 普通非零返回, 投递未知, 执行者活跃, 确认超时与租约过期都不能证明已退出. 可选的 reclaim_decision 引用此前用户对已完成回收的授权; 它既不执行回收, 也不替代已停止观察.
- 授权事务对卡片/派发 revision 做 CAS, 归档并隔离旧 epoch, 再签发仅限 wrap-up 的 epoch, 附独立的 120 秒接受截止时间, 保留原 intent 截止时间. 已有授予时对账它而不再签发. 通过普通的原子 working 回执接受; 只允许此前已授权的清理与仅追加的 wrap-up 记录. 不改代码, 不启动/投递, 不做普通接管升级, 不写运行时身份, 不做原作者结论. 追加 WRAP_UP_RECORDS 前保留整个 spec, 或创建不可变的 `wrap-up/<dispatch-id>-<epoch>.md` 记录. 旧 epoch 的写入仍被拒绝. 这些隔离覆盖受控入口, 不覆盖任意本地文件系统或 Git 操作.
- 状态为 prepared, delivery-unknown, accepted, completed, failed 与 cancelled. 在发送/启动边界前持久化 delivery-unknown. 传输返回值与 marker 回显决不表示已接受. 发送失败不确定时保留 intent 与载荷; 不立即启动另一个执行者.
- 重试前先查询同一 ID 的回执. 已接受/已完成时直接返回, 不发送. 否则只有当前的, 身份有效的已停止观察才允许恢复; 身份未知或缺失不允许. 已知就绪的会话可再次接收同一 ID, 因为接受是幂等的. 忙态等待共用持久化的截止时间. 到期无回执是非零的 pending 结果, 不是业务完成.
- 执行 Agent 先运行 `kander move <task-id> working --dispatch-id <id> --execution-epoch <epoch>`. 其 JSON 含 `dispatch` 与 `replayed`. 只在 `replayed=false` 时开工; 重放或失败时报告回执, 不重复工作. 即使首轮在 notify 返回前已回到 review, 此规则仍然成立.
- 每次绑定的作者更新都携带同一 ID/epoch 与当前期望 revision. 绑定的审核 disposition JSON 还携带 `authorization: {"dispatch_id":"...","epoch":1}`; 旧版记录不变. 以 `move review` 完成 fix/sync, 或以 `move done --result completed` 完成 wrap-up, 携带 `--dispatch-id`, `--execution-epoch`, `--delivery-commit <final-full-SHA>` 及适用的 `--disposition <card-relative-artifact>`. 完成与其回执原子提交. 既有的 done/review 证据门禁仍适用; 记录的 SHA 不是 Git 集成校验.
- 每张卡只有一个活跃的执行授予. 新 intent 只替换已完成, 失败或取消的 intent. 显式 `resume --agent` 保留既有的用户授权要求, 并为未完成的同 ID 派发轮换 epoch, 保持其原载荷与截止时间. 旧 epoch 不能接受, 更新正文/WINDOW 或完成. 没有跨 Agent 会话的卡片锁.
- `dispatch fail|cancel` 以其期望派发 revision 与原因记录显式的 intent 决定; 它不取消或迁移任务, 不授予接管, 不丢弃原件. 已变更/已过期的 intent 须先有显式处置才能创建另一个 ID; 仅凭不确定不构成取消授权.
- 普通未绑定的 working 消息与非组旧版通知保留原有行为, 不产生历史业务回执. 下文的旧版投递/marker 条款只适用于这些未绑定消息. 绑定的卡片写入不能省略执行授予.
- 派发原件与 epoch 回执是生产者所有的 `dispatches/` 附件. 待发布项使用既有的显式 init 恢复与维护规则. 保证只覆盖受控的卡片操作, 决不对任意外部副作用保证 exactly-once.

**已配置的执行 Agent**

- 可选的 `agents` 配置声明执行名, 可执行路径, pane 进程名, CLI 方言或 argv 模板, 以及会话模式. 自定义 Agent 仅用于执行; Reviewer 隔离与 `*_REVIEW_BIN` 保持独立.
- 模板在 argv 元素内替换 `{model}`, `{effort}`, `{session}`; Kander 最后追加 prompt. 不做 shell 插值.
- `generated` 会话使用 UUID; `allocated` 会话从配置的 argv 命令取得 ID. 为 `none` 时拒绝 resume, notify 使用无直投的新进程恢复; 既有的持久已停止观察与回执门禁仍适用. 配置/check 输出对此限制给出告警.
- tmux 前台检查使用配置的进程名, 缺省回落到可执行文件 basename. herdr 仍需其自身的 Agent 识别.

**启动参数与元数据**

- `start` 的 Agent、launcher、模型档位默认取 Kander 配置, 初始化未完成则用默认值.
- `--agent` 与 `--launcher` 只覆盖本次.
- `SIZE: large` 的任务用 `kanban_agents.large`, `SIZE: small` 的任务用 `kanban_agents.small`, 缺省均取 `kanban_agent`.
- 成功输出规模和实际 Agent.
- `start` 默认免确认, 将会话标识写入 `SESSION`:
  - Claude/Grok 为 UUID; Cursor 为 chat id; Codex 与 Kimi 只记 Agent 名.
- 紧邻的 `WINDOW` 字段写投递地址:
  - herdr: `herdr:<tab-id>:<pane-id>`.
  - tmux/tmux-session: `<launcher>:<session-id>:<window-id>:<pane-id>`.
  - foreground/console: launcher 名.
- tmux/tmux-session 先建占位 window/pane, 持久化 `WINDOW`, 再 `respawn-pane` 启动 Agent, 用 `tmux set-option -p -t <pane-id> @kander_session <会话-id>` 写 pane 标记.
- Claude/Grok/Cursor 用卡片 id, Codex 与 Kimi 复用 `notify`/`resume` 的会话存储解析.
- 地址写入、启动或标记写入失败均关闭本次 window 并回滚卡片.
- 旧卡缺两字段时按序插在 `OWNER` 后, 不批量改写未启动旧卡.

**恢复原会话**

- `resume` 按卡片 `SESSION` 唤醒原 Agent, 保留上下文:
  - Claude/Grok 用 `--resume <uuid>`; Cursor 用 `--resume <chat-id>`; Kimi 用 `--session <session-id>`, 首次启动不传任何值, 因为 kimi-code 自行生成 id.
  - Codex 用 `codex resume <session-id>`.
  - Codex session id 在 `CODEX_HOME` (默认 `~/.codex`) 的 rollout 记录中检索.
  - 只匹配以该任务 start/resume prompt 开头的用户消息, 不匹配仅提到任务 ID 的主控会话.
  - 找不到则失败.
- 只接受 `review/` 或 `working/` 中的卡, 必须且只能给非空的 `--message` 或 `--message-file`.
- `--timeout` 必须是大于 60 的有限秒数且默认 120 秒.
- `resume` 不迁移卡片状态: `review/` 卡的 prompt 会要求被唤醒的 Agent 先执行「命令契约」所述的 `review -> working` 迁移 (未绑定消息为普通迁移, 持久派发为 ID/epoch 绑定的迁移) 再处理事项.
- 启动或存活校验失败时, 仅在本操作仍持有当前 revision 时恢复原文档, 否则报告冲突并保留更新的正文; 文档恢复失败时报错写明卡片实际所在目录.
- launcher 与 `start` 相同.
- 拉起后复用 `notify` 恢复分支的同一存活判据, herdr 与 tmux/tmux-session 校验可寻址终端, foreground 与 console 在完整 timeout 观察期内要求进程不退出.
- herdr 校验本次 `tab create` 直接返回的 pane: `agent` 必须匹配且状态只接受 `idle`, `working`, `blocked`.
- pane 上报非空 `agent_session.value` 时还必须与卡片会话精确匹配, 没有上报有效身份时不以缺失本身判失败.
- 这里判断的是已知新 pane 是否存活, 与直投和反查必须靠会话身份确定目标的职责不同.
- 该降级不是同用户安全边界.
- 秒退时清理新实例, 非零退出并附可取得的 Agent 原始输出, 只有校验通过才按 `start` 的输出格式报告 `已唤醒`.
- 没有 `SESSION` 记录的卡 (未经 `start` 启动) 不能 `resume`.
- `start`, `resume` 和需要恢复进程的 `notify` 在所有平台都把完整 prompt 写入 UTF-8 临时任务文件, 内含任务 ID、固定要求和消息正文.
- Agent 命令行只接收一句包含该绝对路径的指令.
- 任务文件内要求 Agent 完成后尝试删除, 删除失败或遗留不影响结果.
- 这类文件不做 POSIX 权限或 Windows ACL 检查与收紧.
- 原生 Windows 优先使用 Agent `.exe`.
- Codex, Claude, Grok, Cursor 或 Kimi 只有 `.cmd`/`.bat` 时, 通过显式 `cmd.exe /d /s /v:off /c` 和 Agent 适配层的参数编码启动.

**接管新会话**

- `resume` 缺省按「恢复原会话」保留原 Agent 与上下文.
- 显式 `--agent <name>` 表示用户授权接管, 即使名称与原 Agent 相同也分配全新会话, 不迁移旧 CLI 上下文: 新 Agent 先从卡片、任务 worktree、Git 状态和实施记录重建进度, 再处理消息.

  接管对照卡片, 任务分支与审核原件核实交付事实. 它不重置验收条件, 不改写契约, 也不重开已被核实证据关闭的 finding, 除非引用具体缺口 (命令, 输出, commit). 撤回前任的完成声明作为独立的 `TAKEOVER_AUDIT` 条目记录, 附使其失效的证据, 卡片保留前任的原始声明作为历史.
- 未经 `start`、没有原 `SESSION` 记录的卡仍不可接管.
- Cursor `create-chat` 等新会话准备失败时卡片不变.
- 启动前以新 Agent、新会话和新 launcher 覆写 `OWNER`/`SESSION`/`WINDOW`, 不改 `STARTED_AT`.
- 启动或存活校验失败时, 仅在本操作仍持有当前 revision 时恢复原文并复原原状态; 否则报告冲突并保留更新的记录.
- tmux/tmux-session 的 Codex 接管会发现本次精确 session id, 写为 `codex <id>` 并设置 pane 标记.
- herdr 与进程型 Codex 缺少发现通道, 保留 `codex` 并在后续恢复时按 rollout mtime 解析.
- 新 Agent 通过存活校验后、输出 `已接管` 前, 命令按 `dismiss` 的身份和单 pane 拓扑门禁尝试优雅退出并关闭原 herdr/tmux 容器.
- 原 Agent 已死时仅在容器拓扑可证时关闭, 原地址为空、foreground/console 或与新地址相同时记为 N/A.
- 校验、退出或关闭失败只输出 `原容器保留` 及原因, 不强杀、不回滚新 Agent, 命令仍成功.
- 清理与存活观察各可使用一次完整 `--timeout`, 最坏耗时为约 2 倍 timeout.
- 成功先输出 `已清理原容器: ...`, 再输出 `已接管: ...`.

**通知与恢复**

- `notify` 是主控向原执行 Agent 派回事项的单一接口.
- 它与 `resume` 一样只接受 `review/` 或 `working/` 卡, 必须且只能给非空的 `--message` 或 `--message-file`, `--timeout` 必须是大于 60 的有限秒数且默认 120 秒.
- 地址优先级为显式 `--pane` 覆盖、卡片 `WINDOW` 快路径、缺窗口时按 Agent 与会话 id 扫描 `herdr pane list`.
- 覆盖与反查均继续用 `pane get` 验证 pane 存在、Agent 和 `agent_session.value` 完全匹配且既有 pane 状态为 `idle` 或 `done`.
- 卡片已记录 id 的 Claude/Grok/Cursor 直接比对; 缺 id 的 Codex 卡, 以及首次发现之前的 Kimi 卡, 复用 `resume` 的会话检索.
- Kander 不设 Agent 白名单, herdr 反查覆盖范围取决于当前版本及各 `source: herdr:<agent>` 集成是否实际报告会话身份.
- 唯一命中后把 `herdr:<tab-id>:<pane-id>` 写回 `WINDOW`.
- 0 个或多个命中不投递.
- 直投与反查负责确定既有目标身份, 因此 pane 未上报有效会话身份时继续拒绝并进入恢复链, 不使用恢复存活校验的降级判据.
- tmux/tmux-session 的旧卡缺窗口时仍不能反查.
- 有地址时同时验证 `pane_dead=0`, `pane_in_mode=0`, `pane_current_command` 与 Agent 可执行名一致, 并要求 pane 的 `@kander_session` 或 `@onevoke_session` 用户选项至少其一非空, 且该非空值与卡片解析出的会话 id 完全一致.
- 两者皆空时因身份不可证按无直投通道处理并回落, 不降级为进程名级放行.
- 选项不一致属于地址过期, 先按本节三类探查结果反查, 仅在唯一命中并复检通过后直投.
- tmux 用户选项和 herdr `agent_session` 都处于同用户权限内, 只用于避免误投, 不构成抵御同用户恶意伪造的安全边界.
- 只有直投地址及探查通过后才创建正文载荷.
- Windows 临时根只取 GetTempPathW 返回的词法路径并由 no-follow 边界逐分量拒绝 reparse point.
- 正文写入仅当前用户可访问、创建时即收紧的 `0700` 临时目录及 `0600` 文件, foreground/console 回落不创建载荷.
- 终端只收到一行以 `# kander-notify:` 开头、含绝对路径和 marker 的指令.
- 匹配该指令行时同时识别 `# onevoke-notify:`.
- herdr 必须用 `agent prompt <pane-id> <instruction>` 投给 pane 内已在运行的 Agent TUI: 它按 pane 实际的 bracketed-paste 模式送正文, 再延时补一次编码后的 Enter, 并在 Agent 已停在审批或提问 UI 时先行拒绝而不是把正文塞进那个对话框.
- 不得改用面向 shell 的 `pane run` (正文与结尾 CR 同批写入, Cursor 等 TUI 不把它当提交, 正文会停在输入栏), 也不得把文本交给只接受按键名的 `pane send-keys`.
- tmux 继续用 `send-keys -l` 后单独发送 `Enter`.
- herdr 用 `pane wait-output --match <marker> --source recent` 做字面子串匹配, 允许 TUI 在 marker 前加渲染前缀.
- tmux 用有界 `capture-pane` 并按字面子串确认 marker, 同样允许渲染前缀.
- 投递动作成功即成功返回, 命令不迁移卡片; `review/` 卡的消息正文会前置「先执行 `review -> working` 迁移 (未绑定消息为普通迁移, 持久派发为 ID/epoch 绑定的迁移) 再处理事项」的要求, 该迁移即为开工回执, 见「命令契约」.
- marker 超时只警告「已投递, 未在超时内确认」, 不恢复第二个进程.
- foreground/console、无直投通道、探查或投递动作失败才由命令内部恢复原会话.
- process 型恢复必须在完整 timeout 观察期内保持存活, foreground 验证后继续占用当前终端直至 Agent 退出.
- herdr 恢复存活按 `resume` 的新 pane 判据执行, 其中状态仍只接受 `idle`, `working`, `blocked`, 拒绝 `done` 与 `unknown`.
- 恢复校验失败后的 tab/window/process 清理失败必须与原错误合并报告, 并提示新 Agent 可能仍存活.
- 主控不另行调用 `resume`.
- 两路都失败时非零退出并同时报告原因, 卡片正文与原状态不变.

**遣散与终端清理**

- `dismiss` 只接受 `done/` 或 `archived/` 卡, 不改卡片正文和状态.
- `--timeout` 必须是大于 60 的有限秒数且默认 120 秒.
- 它按卡片 `WINDOW` 定位 herdr pane 或 tmux/tmux-session pane.
- 缺 `WINDOW` 的旧 herdr 卡, 以及 pane 不存在/已死, Agent, 会话标记或前台进程不匹配的过期地址, 复用 `notify` 的唯一会话反查.
- herdr 以命中 pane 实际所属 tab 为容器, tmux 以 `display-message` 读取命中 pane 实际所属 session/window, 不回写卡片 `WINDOW`.
- 0 个或多个命中都拒绝.
- 忙态不反查且继续拒绝.
- 投递前复用 `notify` 的 Agent 与会话精确匹配: herdr 另要求 `agent_status` 为 `idle` 或 `done`, tmux 另要求 pane 存活、不在 copy-mode 且前台进程匹配.
- 被校验 pane 的当前 tab 或 session/window 必须与定位出的容器精确一致, 且容器只能包含该 pane.
- pane 被移动或容器另有 pane 时必须在投递前拒绝, 等待退出期间再次验证归属并在关闭前复核容器拓扑.
- Claude/Codex/Kimi 送 `/exit`, Grok/Cursor 送 `/quit`.
- herdr 用 `agent prompt`, tmux 用 `send-keys -l` 后单独发送 `Enter`.
- 只有确认 Agent 进程已退出才关 herdr tab 或 tmux window.
- tmux window 已随 Agent 自动消失时视为已关闭.
- 任何身份、状态、容器归属、投递、退出确认或关闭失败, 以及超时时都非零返回并保留当时现场, 不强杀、不降级关容器.
- `foreground`/`console` 没有可关的终端容器, 在不做部分动作的前提下报错.

**初始化看板**

- `init` 幂等创建看板及 7 个状态目录 (`backlog`, `todo`, `working`, `review`, `done`, `archived`, `trash`), 既有看板重跑一次即补建缺失目录, Git 项目只更新本地 `info/exclude`.
- Windows 新目录必须相对固定父句柄以 `CREATE_NEW` 创建并在创建时应用当前用户独占的 protected DACL, 创建竞态失败须拒绝操作.
- 既有目录只迁移叶目录 ACL.
- Git exclude 的父链逐分量拒绝 reparse point, 既有 ACL 不变, 去重读取和追加在同一固定叶句柄及文件锁内完成.

**启动方式**

- 六种 launcher: `auto` 在启动当时解析, 不把结果写回配置.
- auto 选择顺序: 处于 herdr (`HERDR_ENV=1`) 时按 `herdr` 启动, 否则处于 tmux 时按 `tmux` 启动, 同时处于两者时 herdr 优先, 两者都不在则失败且不领取, 不回落到 `tmux-session`, `foreground` 或 `console`.
- `tmux` 在启动者当前 session 里后台建任务 window, 要求 `start` 本身跑在 tmux 内.
- `tmux-session` 按主树路径确定专属 session (`kb-<目录名>-<路径摘要>`):
  - 不存在则创建; 已存在且 `@kander_project` 或 `@onevoke_project` 匹配本项目时复用.
  - 同项目共用一个 session, 每卡一个后台 window.
  - 不要求 `start` 在 tmux 内运行; 启动后不切换客户端.
  - 输出 session 名、window id 和 attach 提示.
- `herdr` 要求 `HERDR_ENV=1`, `HERDR_WORKSPACE_ID` 且 herdr 在 PATH, 在当前 workspace 后台新建 tab (`--no-focus`, 标签复用 `window_name()`) 后先等根 pane 就绪, 再在该 pane 执行与 tmux 相同的 Agent 命令, 不使用 `herdr agent start`.
- `foreground` 在当前终端前台运行并等待 Agent 退出.
- `console` 仅支持原生 Windows, 在独立控制台窗口启动 Agent 后立即返回 PID.
- `console` 没有 session/window 复用、attach 或输出抓取能力, 不是 tmux 或 `tmux-session` 的等价实现.
- POSIX 默认 `auto`, Windows 默认 `console`.
- Windows 拒绝 `tmux` 和 `tmux-session`; herdr 有原生 Windows 版本, `herdr` 在 Windows 可用.
- 配置同样接受 Windows 上的 `auto`, 但它在 Windows 只会落到 herdr.
- 送进终端容器的 Agent 命令要经该容器的 shell 再解析一次: POSIX 按 sh 拼接; Windows 假定 herdr pane 是 PowerShell, argv 一律编码进 `%VAR%` 变量后由 `cmd.exe /d /s /v:off /c` 还原, 不依赖 PowerShell 向原生程序传参.
- Agent 命令含换行或 NUL 时拒绝启动, 不向容器发送半条命令.

**herdr 会话上报**

- herdr 的 `pane run` 成功且会话 reference 非空时, `start` 与 `resume`/`notify` 恢复共用会话上报路径:
  - 通过 `HERDR_SOCKET_PATH` 调用 `pane.report_agent_session`.
  - 传本次 pane id、卡片 Agent、reference、`source=herdr:<agent>` 与单调递增的 `seq`.
  - 在有界预算内用 `pane get` 读回相同 `agent_session.value`.
- 空 reference (包括新启动且尚未发现 id 的 Codex) 不上报.
- socket、响应或读回失败只输出一条告警, 不使启动失败, 不回滚卡片或关闭 tab.
- herdr 自身集成仍可上报同一身份.
- Kander 的路径用于补齐未触发集成钩子的启动方式.

**检查与存活分类**

- `check` 默认检查除 `done/` `archived/` 外的无效入口, 有错非零退出.
- 对 `todo/`, `working/`, `review/` 卡另按与 `todo/` 入口门禁相同的规则检查契约完整性: `GOAL`, `EXPECTED_OUTCOME`, `ACCEPTANCE_CRITERIA` 或 `OUT_OF_SCOPE` 缺失或为空, 这四个章节中任一残留 `<FILL_IN>` 占位符, 或验收条件没有 `- [ ]` 条目, 均计入无效项.
- `--all` 纳入两栏.
- 指定任务 ID 时仅检查目标及跨状态/形态冲突, 无关无效入口不影响结果, 目标在 `done/` 或 `archived/` 也检查.
- 均解析适用卡片的 `PREREQUISITES`, 确认引用存在、依赖无环.
- 定向检查遍历可达依赖, 包括跨组环.
- 默认对 `done/` 或 `archived/` 前置卡只确认存在, `--all` 才检查其自身及完整可达图.
- 依赖未满足不使 `check` 失败.
- 无参和 `--all` 探测全部 `working/` 卡.
- 定向只探测指定的 `working/` 卡, 不探测 `review/`.
- 汇总行前输出四态存活段: `alive` 为 Agent 与可用会话身份匹配.
- `alive` 是对 Agent 存在的观察, 不是可接收输入的就绪状态或任务进展的证据; `notify` 另有自己的就绪检查.
- `stopped` 为 pane 消失或 Agent/进程不匹配且有效反查确认 0 个命中. 反查被禁用或会话 reference 为空时, 按既有的直接 pane 分类.
- `drifted` 为会话唯一反查到新 pane, 附新地址.
- `unknown` 为会话/窗口无效、foreground/console、程序不可用、状态或探测失败. 反查错误, 无效输出, 超时与多个命中均为 `unknown`, 保留地址过期原因及反查阶段与原因; 它们不证明会话已停止.
- herdr 身份缺失但 Agent 和状态有效仍为 `alive`, 注明无法直投.
- Codex 空 reference 不反查.
- 探测错误折算 `unknown`, 不写卡、不影响 `check` 退出码.
- 一张卡的正向探测, 会话反查, 复检与进程清理共用一个截止时间 (默认 10 秒). 耗尽或取消时停止后续查询并报告 `unknown`. 取消会关闭继承的输出管道等待并终止所拥有的进程组/job; 进程创建, 内核 I/O 与回收仍取决于操作系统.
- `check` 存活段在一个 10 秒批次预算内收集全部选中的 working 卡, 最多 4 张卡并发探测. 已完成的观察保留可用; 未完成或未开始的卡保留 `unknown` 原因. 输出含观察时间, 有效性与独立的运行时状态. 这些事实不确立投递就绪或业务进展. 预算不含看板读取与输出写入, 也不是硬实时返回保证. 订阅使用下文规定的同一批次收集器.
- `subscribe` 须显式组 ID 和非空成员 ID, 校验成员归属.
- `--watch` 可重复指定外部卡或组 ID; 保留原始组引用, 每次观察时经共享的看板成员读取口重新展开. 外部任务不必属于被订阅的组.
- 外部目标不存在、展开为空、与成员重复或展开相互重复时, 订阅前失败.

**订阅事件**

- `subscribe` 每行输出一个 JSON 对象. 既有的 `event`, `group_id`, `tasks` 及 state-change 的 `changed` 项 (`task_id`, `from`, `to`) 仍然可用.
- 每行都有 `schema_version: 1`, 新生成的 `subscription_id`, 从 1 起的 `seq`, `observed_at` (UTC 快照时间), 以及随已提交状态与文档一同捕获的 `task_revisions`. `seq` 仅限该次订阅内, 决不是重启游标.
- 初始事件为 `snapshot`. 状态差异产生 `state-change`; 无状态差异的 revision 差异产生 `task-update`. 两者都可附 `updated` 任务 ID. 同一 refresh 内的 review-working-review 往返可由 revision 观察到; 仅凭 revision 不证明业务完成或派发已接受. 重启时以新 snapshot 对比已保存的 revision; 不隐含重放日志. 旧版无版本卡的 revision 为 0; 未受控的文件编辑不必推进它.
- 传 `--watch` 时, 事件保留 `watch_references`, 将外部 ID 展开进 `watched`, 并在 `tasks` 与 `task_revisions` 中包含当前全部被监控任务. `watched` 为空时省略. 组引用还携带完整的 `memberships` 与确定性的 `membership_versions`; 版本为排序后成员 ID 的哈希, 启动后被清空的组为 `empty`.
- 组成员变化产生 `membership-change`. 新增立即加入. 将成员移出其原被监控组的移除与归属变化 (即使它仍在另一被监控组中) 携带 `removed`, 并在订阅余下时间内置 `reconciliation_required: true`. 不得从集合缩小推断依赖已满足: 先重新核对契约, 成员版本与实际交付. 初始为空的组与重叠目标被拒绝.
- `membership_complete: true` 与 `read_status: committed` 描述完整观察, 不是放行依赖的授权. 扫描问题或归属不可读时 fail closed, 以终止性的 `membership-unknown` 事件, `membership_complete: false`, `reconciliation_required: true` 与非零退出结束. 该事件上的任务值是上一次完整快照, 或在任何成功快照前为空映射; 决不将其视为新的完整事实. 显式的仅任务 watch 使用定向读取; 组展开检查所有可能成员而不探测无关 Agent.
- 协调事实读取共用 2 秒的锁争用预算. 锁耗尽时产生终止性的 `read-error` 并置 `read_status: maintenance`; 已知的 prepared 事务立即产生 `read-error` 并置 `read_status: recoverable`. 两者都置不完整/对账标记, 非零退出且不修复. maintenance 状态也可能只是普通的写入者争用. OS 文件操作仍受 OS 约束; 这些不是硬实时保证. 恢复按既有的显式 init 维护协议进行.
- 以任务 ID 为键的 `dispatches` 携带来自与卡片状态/revision 同一共享锁快照的当前派发摘要: `dispatch_id`, `task_id`, `kind`, `epoch`, `state`, 派发 `revision`, 原 intent 的 `created_at`, 当前 epoch 的有效接受截止时间 `confirm_by`, `age_seconds`, 以及持久的 `accepted`/`completed` 回执, 附卡片 revision 与适用的投递/disposition 引用. `created_at` 保持为原 intent 的创建时间, `age_seconds` 从它起算. 普通授权使用 intent 的截止时间; 仅限 wrap-up 的授予使用自己独立的截止时间, 不改写原 intent. 省略消息正文. 未绑定的卡没有捏造的回执. 更早被取代的 ID 仍可用 `dispatch show` 查询; 这不是历史列表.
- 同一 refresh 内的 review-working-review/done 往返仍可由下一事件或重启 snapshot 中的同 ID 回执证明, 包括在 notify 返回前完成的情况. 完成回执不是 Git 集成或审核关闭的证明.
- `confirmation_pending` 仅在 prepared/delivery-unknown 时为 true; `confirmation_overdue` 还表示当前 epoch 的有效接受截止时间已过. 已接受的工作没有推断的完成截止时间. 该截止时间独立于 refresh/heartbeat 唤醒订阅, 且不因同 epoch 重试, 无关事件或订阅重启而重置.
- 到期时 (包括初始 snapshot 已过期), 产生 `dispatch-attention`, 在 `attention` 中列出逾期任务 ID, 附当前派发事实, 可用的 `liveness` 与 `reconciliation_required=true`. 调度一次有界探测; 若批次进行中, 保留一个合并的请求直到它加入. 之后的批次完成再次报告仍逾期的事实; 心跳继续探测. 存活与业务进展缺失分开, 观察时长与派发时长分别报告. attention 决不执行恢复, 重试, 卡片变更或依赖放行. 被监控卡的派发证据缺失/损坏使读取失败, 而不是视为不存在.
- 心跳独立于状态, task-update 与成员事件. 有被监控 `working/` 卡或待派发 `review/` 卡时, heartbeat 附 `liveness` 项, 含 `agent`, `status`, `channel`, `detail`, 复用 check 的四态分类器. 探测失败折算 `unknown`, 不终止订阅; 快速 refresh 不入队探测. 当前派发为 prepared/delivery-unknown 的待处理 review 卡也获得 liveness; 普通 review 卡与已完成的 review 卡不获得. 无符合条件的卡时省略 liveness.
- 扫描不等待探测或输出. snapshot 入队后启动一个探测批次; 后续心跳或派发截止时间仅在上一批次结束后才启动新批次. 每批次总预算 10 秒, 最多 4 个 worker; 批次不累积.
- liveness 含 `revision`, `identity`, 可选 `observed_at`, `age_seconds`, `runtime_state`, `observation_valid`, `collection_state` (`pending`, `complete`, `not-observed`), `collecting`, `stale` 与可选 `new_window`. revision 或身份不匹配时丢弃缓存结果. 待处理/未收集的结果保持 unknown. 早于 heartbeat 加 10 秒的观察视为 stale 且 unknown, 保留其真实时间戳. 这些字段不确立就绪或业务完成.
- 输出最多保留 16 行队列加一次进行中的写入, 单行上限 1 MiB, 从入队到写出的截止时间为 2 秒. 队列溢出, 超长行, 短写, 管道断开或输出超时都以显式错误终止. 消费者必须丢弃不完整的末行, 并在重连后以新 snapshot 对账; 没有事件重放保证.
- Ctrl+C, 受支持平台的终止信号与调用方取消会停止并 join 所有自有的探测/输出工作. 全部卡片完成时不自动退出. 文件系统/内核 I/O 与进程回收仍依赖 OS, 不是硬实时保证.
- `--refresh` 默认 1 秒. `--heartbeat` 默认 900 秒, 在 snapshot 入队后开始, 且仅在每次心跳入队后重新计时, 即使没有任务处于 working.
- 两个间隔都必须有限, 至少 `1e-9` 秒 (1 ns), 且小于 `9223372036.854776` 秒. 纳秒以下的小数截断.

- 命令只做结构和机械校验; 授权, 依赖和终止理由由 Agent 按本文件判断.

## 状态模型

目录是状态唯一真源; 卡片正文不设 `status` 字段.

- `backlog/`: 已记录但尚未承诺执行.
- `todo/`: 用户已确认, 契约完整, 尚未领取.
- `working/`: 已领取, 正在实现, 验证, 审核或集成; 任务组卡在修复轮次和集成后的收尾也回到这里.
- `review/`: 仅任务组卡使用. 开发、验证和任务分支交付记录已完成, 等主控将该交付 ff 到组分支, 再安排适用审核与最终集成. 此状态本身不保证交付已进组分支, 主控须核对后才能放行组内依赖. 执行 Agent 迁入后结束本轮响应并保留交互式 CLI 会话; 修复、同步或集成成功后的收尾由主控 `notify` 派回, 原执行 Agent 收到后先自行运行派发 prompt 中 ID/epoch 绑定的 `kander move <task-id> working --dispatch-id <id> --execution-epoch <epoch>` 再处理.
- `done/`: 已满足完成门禁的近期任务.
- `archived/`: 不占活跃看板的完成, 取消, 重复或不修复记录.
- `trash/`: 用户明确要求删除, 但尚未永久清理的入口; 不是任务状态.

```text
backlog <-> todo -> working -> done -> archived        (单卡流程)
                      |  ^
                      v  |  修复轮次与收尾以 ID/epoch 绑定的迁移迁回 working
                    review -> working -> done         (任务组流程; 主控代做收尾
                                                       在其隔离授予下走同一路径)

todo -> backlog                                       取消承诺, 退回待排期
review -> working (--owner)                           用户授权回收未绑定的卡, 见「领取, 启动与协调」
backlog, todo, working, review -> archived            仅限用户授权的终止
除 trash 外任意状态 -> trash                            仅限用户明确要求
```

- 进 `todo/` 须完成 `GOAL`, `EXPECTED_OUTCOME`, `ACCEPTANCE_CRITERIA` (至少一条顶层 `- [ ]` 且有内容的可判定条目) 和 `OUT_OF_SCOPE`, 且这四个章节不残留 `<FILL_IN>` 占位符, 并附「建卡后自审」的 `SELF_REVIEW:` 记录行 (大任务与任务组成员卡另附 `CARD_REVIEW:` 行); 进 `review/` 须已填写 `TASK_BRANCH`; 进 `done/` 的门禁见「执行与完成」与「审核证据完成门禁」, 其余见「终止与清理」.
- 旧版看板没有 `review/`: 其余 6 个状态目录齐全时, 任一 `kander` 命令首次定位看板即自动补建 `review/`, 不要求用户重跑 `init`.

  其他状态目录缺失时停止普通看板操作, 可用前述初始化命令补建.

  `review` 被文件/符号链接占用时始终失败.

## 入口与文档

### 不变量

- 每张新卡都是含普通文件 `spec.md` 的目录 `YYYYMMDD-short-slug-task/`. 紧跟 TYPE 之后的 `SIZE: small|large` 独立于目录形态决定任务规模. 旧版 `<task-id>.md` 入口在显式迁移前仍可读.

  `short-slug` 只含小写 ASCII 字母, 数字和连字符.

  去掉扩展名的入口名即任务 ID.

- 任务 ID 全看板唯一, 不得跨状态重复或同时存在文件和目录形式. 迁移移动整个入口; 入口名创建后不改, 不复制后删, 不留副本. 卡片目录内只用相对链接, 保证迁移后有效.
- 卡片路径随状态迁移变化. 编辑前先取得 `kander show --json <task-id>`, 在独立的 UTF-8 输入文件中准备替换内容, 再以该 revision 用 `kander update` 发布. 不直接写入缓存的卡片路径. revision 冲突时须先读取并合并当前文档再重试. 新建卡片只能经 `kander new`; 缺失的卡片决不由 update 或回滚重建.
- 卡片不得包含 token, 凭据, 敏感服务地址或不应留在本机的个人数据.

### 受控文档与恢复

- 所有已迁移卡片的正文都用 `update --document spec.md`. 小卡与大卡都接受 `plan.md`, `report.md` 与普通附件, 包括相对子目录. 旧版文件卡只读: 变更命令在产生副作用前失败并引导调用者执行 `kander init`. 不得用旧版二进制或直接编辑文件绕过此要求.
- 路径必须规范, 相对, 且不含遍历, 符号链接与 reparse point. 决不通过 update 编辑受管的 `reviews/`, `dispatches/`, manifest, 索引或控制记录. 这些记录由专属生产者所有.
- 全文替换必须保留 `LANGUAGE`, `OWNER`, `SESSION`, `WINDOW`, 创建/启动/完成时间戳, `RESULT`, 执行元数据与受管索引. `TASK_BRANCH` 与普通实施/总结记录通过 update 保持最新.
- 手工领取用 `kander move <task-id> working --owner <agent>` 随状态迁移一并写入 `OWNER` 与 `STARTED_AT`. `start` 与 `resume` 拥有会话/窗口元数据. 派回到 working 的迁移不替换既有负责人.
- 完成用 `kander move <task-id> done --result completed`; 校验, `RESULT` 与 `FINISHED_AT` 为一个事务. 先通过 update 写好完成总结或报告.
- 终止用 `move <task-id> archived --result cancelled|duplicate|wontfix --reason <reason> --decision <user-decision-reference>`. duplicate 还须提供 `--duplicate-of <replacement-id>`. 归档已完成的卡用 `--result completed` 并附原因与决策引用. 移入 trash 用 `move <task-id> trash --result trashed --reason <reason> --decision <user-decision-reference>`.
- 这些引用记录授权依据; 工具不核实用户意图. 既有的用户授权规则仍适用. 迁移必须匹配此前读取的快照时, move 可用 `--expect-revision`.
- 卡片进入 todo 后, 即使退回 backlog 其契约仍保持冻结. 用户授权的契约变更用 update 附 `--contract-decision-file <UTF8-decision>` 与期望 revision. Kander 在受保护的 `CONTRACT_DECISIONS` 章节记录该决策及冻结字段的新旧值. 此选项决不授权修改受管元数据或语言.
- 发布中被中断的进程留下可恢复的事务. 读取方报告未完成的操作, 不修复它, 也不把其中间文件当作普通缺失卡片. 显式运行 `kander init` 完成已准备的写入与迁移. 恢复冲突时保留证据与两份副本; 不猜测主副本, 不自动删除其一.
- 事务保护遵循本协议的命令与 Agent. 它们不阻止任意本地进程绕过命令直接编辑文件. 将看板切换到本协议前, 先与旧版 Agent 或二进制协调维护.

### 小任务模板 (`<task-id>/spec.md`)

```markdown
# <任务标题>

- TYPE: Feature | Bug | Chore | Research
- SIZE: small
- TASK_GROUP:
- LANGUAGE: <Agent 交流语言, 如 en, zh-CN, ja>
- CREATED_AT: YYYY-MM-DD HH:MM
- OWNER:
- SESSION:
- WINDOW:
- STARTED_AT:
- FINISHED_AT:
- TASK_BRANCH:
- RESULT:

## GOAL

<改什么, 为什么改>

## USER_DECISIONS

<用户已确认的方向和取舍; 没有则写 N/A>

## EXPECTED_OUTCOME

<完成后可观察, 可验证的状态>

## ACCEPTANCE_CRITERIA

- [ ] <条件>

## THREAT_MODEL

<安全任务写资产, 可信主体和攻击者能力; 非安全任务写 N/A>

## OUT_OF_SCOPE

- <按既有问题, 加固, 共享契约与文档, 相邻功能四类逐一写明排除或纳入, 每条附理由>

## DISCUSSION

<关键结论; 任务组卡片还要在开头记录 PREREQUISITES>

## IMPLEMENTATION

<计划, 分支, commit, 验证命令, 结果, 环境缺口和阻塞>

## SUMMARY

<实际成果, 偏差, 未处理问题和验收结论; 完成前留空>
```

### 大任务文档

- `SIZE: large` 选定大任务契约. `spec.md` 必需, 含小任务的元数据及契约章节: `GOAL`, `USER_DECISIONS`, `EXPECTED_OUTCOME`, `ACCEPTANCE_CRITERIA`, `THREAT_MODEL`, `OUT_OF_SCOPE`, `DISCUSSION`.
- `plan.md` 按需创建, 记录实施步骤, 影响模块, 验证, 发布和回滚计划, 不得修改 `spec.md` 契约.
- `report.md` 完成时创建, 记录实际改动, 最终 commit, 验证, 偏差, 未处理问题, 风险和验收结论; 不建空文件.

### 契约与记录

- `LANGUAGE` 是为用户就此卡片所写一切内容的语言: 标题与正文, 记录, 报告, 审核报告, 以及传给 `kander notify` 与 `kander resume` 的消息. `kander new` 从配置的 `agent_language` 填写它, 给出 `--language <value>` 时用该值; 取值遵循 `agent_language` 格式. 它在创建时固定并覆盖配置; 没有该字段的旧卡按 `KANDER-AGENTS.md`「语言」回落到当前配置.
- 手工领取与 `start` 按「受控文档与恢复」与「启动参数与元数据」所述写入 `OWNER`, `STARTED_AT`, `SESSION` 与 `WINDOW`; 手工领取的卡 `SESSION` 与 `WINDOW` 留空. `TASK_BRANCH` 通过受控正文入口更新, 无分支用 `N/A`.

  命令迁入 `done/` 时填写 `FINISHED_AT`.

  专用的 move 选项在进入 `done/`, `archived/` 或 `trash/` 时原子填写结果.

- 卡片进入 `todo/` 后, `GOAL`, `USER_DECISIONS`, `EXPECTED_OUTCOME`, `ACCEPTANCE_CRITERIA`, `OUT_OF_SCOPE`, `SIZE` 以及任务组关系冻结. 修改任何一项都要先取得用户明确决策. `THREAT_MODEL` 属于审核任务上下文, 但工具不冻结它; 通过普通 `update` 细化并在 `DISCUSSION` 记录该变更.
- `OUT_OF_SCOPE` 如实界定任务边界, 不把未确认的扩展目标写入 `ACCEPTANCE_CRITERIA`. 审核模块启用时再按 `KANDER-REVIEW-RULES.md` 的审核契约细化范围.
- 实施期只追加关键决策, 验证, 环境缺口, commit, 阻塞和下一步, 不复制会话流水. 每轮最多追加一条带日期的条目; 更早的轮次被取代后各压缩为一行摘要. 审核报告, finding 列表与 disposition 以 `reviews/<run_id>/` 与 `dispatches/` 引用, 决不粘贴进卡片正文; `SUMMARY` (或 `report.md`) 中的未处理项清单每项一行: finding 写明角色, 档位, 状态及其 disposition 记录的卡片相对路径; 未完成的角色, 缺失的报告章节或验证缺口改为写明该 run 的 sidecar 或错误日志, 档位与状态为 `N/A`; 不存在 run 时 (预检失败, 或审核之外的缺口) 写明实际的命令日志或 `IMPLEMENTATION` 验证记录并注明「未产生 run」, 决不写捏造的路径. 全文留在这些产物与用户报告中. 卡片正文超过约 30 KB 表明历史在被复制而非引用. 稳定的架构, API 和长期规则仍须写入仓库文档或项目规则.

## 审核证据存档

- 任务绑定的 `kander review` 在位置参数前使用可重复的 `--task` 参数. 其稳定的 run/batch 与重试协议在最小工具协议中; 仅为使用此命令不需要加载可选的审核工作流.
- 每张目录卡保留不可变的 `reviews/<run_id>/` 原件, sidecar 与 manifest. REVIEWS 是机器所有的章节, 每次 run 一行 JSON, 含 run/batch/角色/执行状态/base/commit/前驱及相对报告路径 (报告不完整或无效时为原始输出).
- `update` 不能修改此章节或其受管附件. 报告随卡片迁移; 决不重建旧状态路径, 也不用摘要替换原件.
- `check` 报告缺失/冲突的 intent, 不完整的发布, 索引/manifest/哈希不匹配, 语言/成员冲突与断裂的前驱链. 它不检查散文来推断 PASS. 证据完整的失败 run 仍是已记录的执行失败.
- 被中断的看板发布须按其维护规则做显式 init 恢复, 再以同一 run ID 重试. 该重试决不重跑 Reviewer; 已成功的卡片回执保持完整. 不得通过直接编辑文件绕过事务恢复.

## 任务规模与分组

- 新卡一律使用目录形态. `new` 写 `SIZE: small` 并包含 IMPLEMENTATION/SUMMARY; `new --large` 写 `SIZE: large`, 完成时要求非空的 report.md. 小卡即使有 report.md 仍要求填好 SUMMARY. 两者进 todo 前都要求 SELF_REVIEW; 大任务与全部组成员另要求 CARD_REVIEW. 进 todo 后修改 SIZE 须走既有的显式契约决策 update 流程. 决不从目录或 report.md 的存在推断规模.
- 卡片需要 `plan.md` 才能保持可审核时选 `large`: 它涉及多个模块或阶段, 需要发布或回滚计划, 或其验证超出单次定向测试运行. 整个改动与验证能放进一条 `IMPLEMENTATION` 条目的卡是 `small`. `SIZE` 还按「启动参数与元数据」选定执行 Agent 档位.

- 卡片保留可选的任务组字段及依赖记录. 关闭 task_groups 时不自动拆组, 独立单卡仍可使用. 开启时按 KANDER-TASK-GROUP-RULES.md 规划和执行, 必须同时开启 git.
- 建卡引导属于 KANDER-TASK-INTAKE-RULES.md, 仅在 rules.task_intake=true 时读取; 用户主动操作看板不要求开启它.
- kander new 在 backlog 创建模板, 调用者按已确认内容填写 `GOAL`, `EXPECTED_OUTCOME` 和 `ACCEPTANCE_CRITERIA`, 不将建议写成用户决定. 卡片文本以卡片的 `LANGUAGE` 书写.

### 显式迁移与维护

- `init` 把全部七个状态中的旧版文件迁移为 `<task-id>/spec.md`, 为无 SIZE 的文件加 `SIZE: small`, 为无 SIZE 的目录加 `SIZE: large`, 并只调整为保持目标所必需的相对 Markdown 链接目的地. ID, 链接标签/标题, 附件及其余全部正文字节保持不变. 既有有效的 SIZE 保留. 重复 init 报告零迁移, 不触碰未变化的卡片内容与修改时间. 无效或重复的 SIZE 拒绝变更与迁移; 只有精确的 `small`/`large` 值有效.
- 链接重定位同时考虑引用方文档与被引用旧版卡的移动. 它保留 URL query/fragment 语义, 处理行内链接, 图片与引用定义, 包括未使用的定义. 代码 span/块, web URL, 根相对 URL 与纯页面锚点/query 引用保持不变. 只扫描看板内的普通卡片 Markdown 文档; 生产者所有的 reviews/dispatches 子树被排除, 其历史原件保持不可变, 包括 journal 重放期间; 不改写看板外的仓库文件. 不支持的 wiki 链接, HTML href/src/srcset, 无效 URL 与不可移植的反斜杠路径仅在引用方文档移动或可能的目标被映射时才在预检失败. 目标未变的静止历史引用与只含绝对 URL 的 srcset 保持逐字节一致. 失败时指明文档与原因, 保留内容供修正而不猜测改写.
- `list`, `show`, `check` 与 `subscribe` 继续读取旧版文件而不迁移. 缺 SIZE 时旧版文件视为 small, 目录视为 large. `check` 对缺 SIZE 的目录要求 init; 其既有状态范围不变. 读取命令决不执行批量迁移.
- 迁移前暂停所有执行 Agent, 外部编辑器, 通知与归档写入者, 包括不遵循事务协议的旧版二进制. 此维护窗口保持到恢复与迁移完成. Kander 对协作的读写方取得看板独占访问, 但无法核实任意外部进程已停止.
- 存在任何 working 或 review 卡且需要迁移时, init 默认拒绝并列出受影响的 ID. 只有在所有写入者已暂停后, 才用 `init --maintenance` 确认这些前提. 恢复被中断的活跃卡迁移需要同样的确认. 不自动终止任何 Agent.
- 迁移先持久化一条 redo 记录, 含完整的新旧路径映射与全部受影响卡片文档, 把每个文件移入同卷 staging, 通过绑定 journal 的 scratch 与原件备份文件写入 SIZE 与链接替换, 发布目录, 再提交 revision 与 journal. 部分完成的替换只在与记录的 after-image 前缀匹配时才续做; 未知残留作为冲突保留. 这些是带内部中间状态的独立步骤, 不是单次原子 rename. 协作的读取方等待维护锁; 进程中断后它们诊断出受管的待处理事务并要求显式 init 恢复. 目录 SIZE 补写与来自既有目录卡的链接 (含 Markdown 附件) 使用同一 journal. 恢复对照 SIZE 插入与记录的路径映射校验已记录的 after-image; 决不从部分迁移的看板重建该映射. 没有读取方自动修复看板. 有效的待处理操作在普通结构扫描前重放. 无需规划迁移时, 普通的杂散非卡片文件产生 check 告警而不阻塞 init; 缺失 spec, 真正重复的 ID, 未知迁移产物与 reparse point 仍失败并保留证据.
- 恢复冲突, 未知 staging 产物, 重复入口或 reparse point 保留证据并显式失败. 决不猜测主副本, 删除未知产物, 或用点前缀改名来隐藏未完成操作. `guard-write` 同时识别过期状态路径与同状态的旧 `.md` 拼写, 但其检查不使外部写入原子化; 请用 update.

### 建卡后自审

- 创建者填写完整契约后必须自审, 在进入 `todo/` 前完成; 发现问题先修正再复核, 未通过不得推进. 此步骤属于建卡流程, 不依赖 `task_intake` 或 `review` 开关.
- 对照用户目标、已确认计划与项目规则检查:
  - 目标与成果一致, 没有遗漏已确认需求, 没有把建议或假设写成用户决策.
  - 任务边界清楚, 本轮范围与排除项不冲突, 不将达成目标所必需的工作排除在外.
  - 约束准确且可执行, 符合用户决策、项目规则及实际接口和环境, 没有互相矛盾的要求.
  - `ACCEPTANCE_CRITERIA` 覆盖目标与成果, 可执行、可判定, 没有遗漏关键条件或引入范围外要求.
- 在 `DISCUSSION` 中以独立一行 `SELF_REVIEW: <结论>` 记录自审结论与修正项 (ASCII 冒号, 可为列表项); 进入 `todo/` 的门禁校验该行存在. 可依据既有决策修正的内容直接修正; 需要新增或改变用户决策的歧义, 明确列出并等待用户决定, 不自行补成契约.
- 大任务目录卡与任务组成员卡在自审之外还须独立审卡: 由不共享建卡会话上下文的独立 Agent (新会话或子 Agent) 只读卡片与用户原始需求, 按上述四条出具结论; 创建者修正后在 `DISCUSSION` 以 `CARD_REVIEW: <结论>` 行记录结论与审卡者, 门禁同样校验该行. 工具只校验记录行存在, 审卡者的独立性与结论质量仍由创建者如实保证. 创建者只为记录独立 Agent 实际执行过的审卡才写 `CARD_REVIEW:` 行; 没有该审卡而为满足门禁写它是被禁止的.

## 领取, 启动与协调

- 未指定任务且 `todo/` 有多张卡时列候选让用户选; 任务组按已确认依赖排序, 不逐卡询问. 开工条件不足时只报缺口, 不领取或退回 `backlog/`.
- 动代码前必须先取得 `working/` 中的唯一入口. 两种领取方式互斥:

```sh
# 委派给新执行 Agent: start 原子领取并启动
kander start [--agent <configured-agent>] [--launcher auto|tmux|tmux-session|herdr|foreground|console] <task-id>

# 用户明确要求当前 Agent 执行既有任务卡: 原子领取并记录归属
kander move <task-id> working --owner <agent>
```

- `kander move <task-id> working` 仅适用于用户明确要求当前 Agent 执行既有任务卡.

  仅在采用已启用的建卡引导时, 选择任一「确认计划并走看板」选项时, 启动卡片必须用 `start`, 除非用户明确要求当前 Agent 自己执行某张卡, 此时用上述 `kander move <task-id> working --owner <agent>`.

  不得先 `move ... working` 再 `start`.

  `start` 只接受 `todo` 卡.

  同文件系统上的入口迁移就是领取原语, 只有迁移成功者取得任务.

  失败后重查, 不建替代卡, 不另加 lock 服务, 数据库或 ID 分配器.

- 回收: `review/` 卡的执行者已停止且用户明确授权新负责人时, 该 Agent 用 `kander move <task-id> working --owner <agent>` 领取该卡. 这只对没有派发绑定的卡 (没有 `DISPATCH_ID`/`EXECUTION_EPOCH` 元数据) 可行: 已绑定的卡拒绝 `--owner`, 且 `dispatch fail|cancel` 不解除绑定, 因此已绑定的卡只能通过「异常恢复」中的接管路径易手. 回收改写 `OWNER` 与 `STARTED_AT`; 因为执行周期按分钟精度从 `STARTED_AT` 导出, 只在 `review progress` 报告 `requirements-needed` 时才视周期已改变 (见「审核证据完成门禁」). 它就是 wrap-up 授权与计划重绑规则所引用的回收. `working/` 卡决不这样回收: 用 `resume --agent` 接管, 它保留 `STARTED_AT`.

**启动检查与回滚**

- `start` 在启动前检查 Agent, launcher 和 TTY.
- `auto` 按「启动方式」只解析为 `herdr` 或 `tmux`; 启动前检查实际 launcher 的前置条件:
  - `tmux`: 已在 tmux session 内.
  - `tmux-session`: tmux 可用, 启动时选定项目 session 名.
  - `herdr`: `HERDR_ENV=1`, herdr 在 PATH, 且有 `HERDR_WORKSPACE_ID`.
  - `foreground`: 三个标准流均为 TTY.
  - `console`: 原生 Windows.
- 前置检查失败不领取.
- 创建进程, tmux session, tmux window, herdr tab, herdr pane 就绪等待或 `pane run` 失败时恢复文档并迁回 `todo/`.
- herdr 就绪等待或 `pane run` 失败还须关闭本次新建的 tab, tmux/tmux-session 的 Codex 会话发现或 pane 会话标记写入失败还须关闭本次新建的 window.
- 新建 tab 的 shell 接管终端前送入的命令文本会被丢弃, 因此 `pane run` 必须在 pane 渲染出首帧输出之后, 就绪等待有上限, 超时按失败处理.
- tmux/tmux-session 只有在会话发现与 pane 标记写入成功后才算启动成功.
- herdr 的成功条件见本节后面的 best-effort 条款.
- foreground/console 在进程创建成功后即算启动成功, 后续退出不自动回滚.
- `console` 成功时输出 PID 后立即返回.
- 成功输出实际 launcher.
- `auto` 须显示解析结果 (`herdr` 或 `tmux`).

- herdr 在 `pane run` 成功后即算启动成功; 后续会话身份上报和读回是 best-effort, 失败只告警, 不进入 `LaunchFailure` 的 tab 关闭与卡片回滚路径.
- `start` 的临时任务文件只写任务 ID 和固定要求, Agent 命令行只传一句读取该文件的指令.

  执行 Agent 先按入口核对配置, 再读本协议、卡片和项目规则, 按适用流程准备实际工作目录.

  确实使用任务分支时填写 `TASK_BRANCH`.

  任务组的启动任务文件中“退出”表示结束本轮响应, 不要求主动退出交互式 Agent CLI 或关闭终端容器. 非交互式调用自然结束时, 后续派回由 `notify` 选择恢复; 交互式会话保留到用户明确同意遣散.

- 领取后只有执行负责人可修改或迁移 `working/` 入口; 协调和编排 Agent 只读督办. 明确交接后由新负责人接管, 不得并发写.
- 启动后的协调责任按启动方式分:
  - foreground 单卡: 启动者在 Agent 退出后检查结果, 直到任务完成或明确交接.
  - tmux、tmux-session 或 herdr 单卡: 执行 Agent 在独立 window 或 tab 直接向用户汇报, 启动者不巡检.

    启动成功后立即告知用户本会话不跟踪该任务进度, 当前 session 可以结束, 下一个任务另开会话.

    `tmux-session` 还要一并给出 session 名和 attach 命令.

    `herdr` 还要给出 tab id 和 pane id.

    `auto` 解析为 `herdr` 或 `tmux` 后按对应单卡规则协调.

    用户明确要求跟踪时改按 foreground 单卡协调.

  - console 单卡: 执行 Agent 在独立 Windows 控制台直接向用户汇报, 启动者不抓取输出.

    启动成功后告知用户 PID 及本会话不跟踪进度.

    该 PID 只用于只读判断进程是否仍存在, 不能用于 attach 或恢复输出.

    用户明确要求由启动者跟踪时改按 foreground 单卡协调.

  - 任务组: 仅在模块已启用且依赖满足时, 按 `KANDER-TASK-GROUP-RULES.md` 编排, 进行适用的审核、集成和收尾, 启动成功不解除该责任.

## 执行与完成

- 单卡按任务契约、用户和项目规则, 以及当前已启用的 Kander 模块完成交付.

  Git 模块关闭时不要求 develop、worktree、提交、push 或合回.

  审核模块关闭时不自动要求审核.

  用户自己的审核、PR 或验收条件仍须满足.

- `rules.git=true` 时, 授权执行独立单卡也就授权了其合回 `develop` 与清理. 验收条件, 验证与适用的审核门禁满足后, 执行 Agent 自动继续完成集成, 适用的 push, 本地同步, 分支/worktree 清理与下述完成命令. Git 操作按 `KANDER-GIT-RULES.md` 执行; 不另行询问合回确认, 也不仅因代码已提交或审核已通过就停下. 显式暂停, 用户验收, PR 或保留分支的要求仍有约束力. 单卡在此流程中始终留在 `working/`; `review/` 仍专供任务组. Git 模块关闭时本段不增加任何 Git 要求.

- 根据卡片记录确认实际工作目录, 记录实施与验证及未处理问题.

  完成任务契约和所有适用交付步骤后, 通过 update 写入 `SUMMARY` 或 report.md, 满足「审核证据完成门禁」(已封存且每个批次都已关闭的计划, 或无审核适用时的显式 N/A 计划), 再执行 `kander move <task-id> done --result completed` 和 `kander check`. 组内卡按 `KANDER-TASK-GROUP-RULES.md`「执行 Agent 收尾」完成, 并运行定向的 `kander check <task-id>` 而非无目标的 check.

- 失败或暂停时保持实际状态并记录阻塞和解除条件. 不适用的 Git 或审核步骤写 N/A, 不把未执行的验证写为通过.
- 任务组成员仅在启用 task_groups 和 git 时按 KANDER-TASK-GROUP-RULES.md 执行 review、集成和收尾. 不能把这些门禁应用到独立单卡.
- 对用户的固定报告格式只在 rules.reporting=true 时读取 KANDER-REPORTING-RULES.md. 关闭时用用户自己的汇报形式, 卡片结果和必要执行记录不省略.
- 已进入 done 的卡片不因事后发现新问题而退回或复用; 新问题另建卡并指向原卡.

## 终止与清理

- 用户明确取消, 判定重复, 决定不修或接受替代方向后, 才可将 `backlog/`, `todo/` 或 `working/` 卡 (以及 review 状态卡) 直接归档.

  实现困难, 验证失败或暂时阻塞不算授权.

  结果只能是 `cancelled`, `duplicate` 或 `wontfix`, 均写原因, `duplicate` 还须指向替代卡.

  `completed` 只用于 `done -> archived`.

- 卡片迁入 `archived/` 或 `trash/` 后, 执行或操作该卡的 Agent 按用户约定汇报结果 (仅在 `rules.reporting=true` 时使用 `KANDER-REPORTING-RULES.md` 模板), 末行状态写实际去向和结果.
- `done/` 保留近期完成项, 用户确认无需展示后再归档. 只有用户明确要求删除具体卡片时才移入 `trash/`; 用专用的 trash 迁移选项原子记录 `RESULT: trashed`, 原因, 决策引用和时间. 不自动清空或永久删除; 永久删除须逐项授权.

## 异常恢复

- `working/` 卡中断, 无负责人或长期无进展时, 协调 Agent 先通知原执行 Agent. 没有派发绑定的卡用 `kander notify <task-id> --message <现状与要求>`: 不带 `--kind` 的未绑定消息; 之后的任何 fix, sync 或 wrap-up 轮次另行带其 `--kind` 派发. 已绑定的卡 (已记录 `DISPATCH_ID`) 旧版消息无法恢复会话: 先读 `kander dispatch show`. 派发仍为 prepared 或 delivery-unknown 且接受截止时间未过时, 按「持久派发」用同一 ID 与其原消息重试 (`notify --dispatch-id <id> --message-file <original>`). 派发一旦已接受, 同 ID 的 `notify` 只返回回执, 不恢复任何东西; 执行者被证明已停止时, 用户先用 `dispatch fail <task-id> <dispatch-id> <dispatch-revision> <reason>` 处置该轮, 然后才以新 ID 创建同 kind 的新派发: `fix` 重新绑定其 finding 与迄今写出的每个作者原件, `wrap-up` 重新绑定集成证据, `sync` 不带证据并在消息中引用既有的 disposition 记录.

  命令自行选择直投或恢复, 非零退出时由用户决定交接或终止.

  用户决定换 Agent 时只用 `resume --agent <name>` 建立接管新会话: 未绑定的卡用 `kander resume --agent <name> <task-id> --message <现状与要求>`; 派发尚未接受且未过截止时间的已绑定卡用 `kander resume --agent <name> --dispatch-id <id> --message-file <original> <task-id>`, 它轮换 epoch; 执行者已停止的已接受轮次, 则通过 `resume --agent` 附其 `--kind` 及 (`fix` 或 `wrap-up` 的) `--evidence-file` 创建上述新派发. 不手工迁移会话、不再次 `start`.

  其他 Agent 不得自行接管, 迁移或归档.

  进程退出不改变 `working/`, 不退回 `todo/`.

- 出现重复 ID, 跨状态副本, 文件与目录同 ID, 目录卡缺 `spec.md`, 目标冲突, 状态目录缺失或不可写时, 停止受影响操作并保留现场. 不通过删除, 改名或移动来绕过报错.
- 看板无 Git 历史; 误删先查 `trash/` 和本机备份, 不伪造内容.

## 协调者检查点

```text
kander coordinator show <group-id>
kander coordinator claim <claim.json>
kander coordinator reconcile <observations.json>
```

这些单二进制生产者通过既有事务协议存储组检查点. 它们读取任务/派发/审核原件, 决不 notify, 迁移卡片, 集成 Git 或授予任务执行权限. 组工作流策略仅在 task_groups 模块启用时加载. `show` 返回已提交的 schema, revision, 协调者授权与成员事实. `claim` 须有 group_id, expected_revision, expected_epoch, owner, 唯一 token, basis 与 members. `reconcile` 须有 group_id, expected_revision, authority 与以任务 ID 为键的 members 对象; 每个成员写明 revision, 以及在绑定时的 dispatch_id, epoch, base 与其已完成的 delivery_commit. 可选的绝对 cwd 选定一个尚存的仓库做只读 Git 校验; 首次未绑定的交付需要它. authority 含成功 claim 得到的 owner, token 与 epoch. 响应丢失后重试同一 JSON; 输入改变或协调者授权过期即冲突. show 或 reconcile 决不修复待处理事务. 只在既有维护前提下使用显式 init 恢复. 检查点是恢复游标, 决不是接受或集成授权. 确认来自同轮回执, 包括快照中的回执.

## 审核证据完成门禁

执行周期是对一张卡的一次领取, 由其任务 ID 与 `STARTED_AT` (分钟精度) 标识: `start` 与 `move working --owner` 写入 `STARTED_AT`, `resume --agent` 保留它. 活跃执行周期在 `move done` 前要求显式的审核计划, 即使 REVIEWS 为空或没有 Reviewer 运行过. 将每个角色记录为必需或 N/A, 附实际原因与规则依据. 审核被禁用或没有任何触发时, 最小序列为: `review plan` 附一个从审核 base 到最终交付 commit 的已封存批次, 把四个角色都写为 `N/A: <reason and rule basis>`, 对该批次 `review aggregate`, `review close` 绑定其视图哈希, 再 `move done`; 这不加载任何已禁用的审核模块. 计划至少有一个批次, 其成员在创建时固定, 每个周期一张卡最多属于一个计划. 只有在每个成员都处于 `working/` 或 `review/` 时才能创建计划, 因此组计划在最后一个成员启动后创建, 且计划存在前不运行任何审核批次: 计划在首个批次的首次 run 前命名它, 之后的每个批次都在其前驱关闭后, 自身首次 run 前用 `extend-plan` 追加 (`extend-plan` 不能收编已运行过的批次). 关闭批次要求它已规划且 worktree 在该批次的最终目标处干净, `review advance` 也要求已规划的批次. 修复轮次推进批次的运行时目标, 而计划保留批次加入时记录的目标; 历史未被改写时, wrap-up 证据把首个已规划批次的 base 与最后一个批次的该记录目标绑定为其 `source_commit`, 因此主控另按 `KANDER-GIT-RULES.md` 用 `git merge-base --is-ancestor` 校验已关闭最终目标的祖先关系, 并在 wrap-up 通知中写明两个 commit. 历史已被改写且最后一个批次在记录后又推进过时, 补丁比较只覆盖记录范围而会遗漏修复: 不要绑定它; 停止, 保持组状态并报告. 已完成且没有计划的历史卡片仍可作为 legacy-untracked 读取, 决不作为捏造的历史 PASS.

受控的审核证据命令, 全部在既有的单一 review 入口下:

```text
kander review plan <absolute-CWD> <absolute-plan.json>
kander review extend-plan <absolute-CWD> <absolute-extension.json>
kander review assign <absolute-CWD> <absolute-assignment.json>
kander review disposition <absolute-CWD> <absolute-author-record.json> <expected-card-revision>
kander review map-legacy <absolute-CWD> <absolute-map.json>
kander review aggregate <absolute-CWD> <batch-id>
kander review advance <absolute-CWD> <absolute-advance-request.json>
kander review close <absolute-CWD> <absolute-close-request.json>
kander review progress <absolute-CWD> <task-id>
```

计划, 分派, 作者原件, 生成的 disposition 与关闭产物是生产者所有的 reviews 附件; 普通 update 不能替换它们. 作者在 working 期间只提交自己被分派的条目, 附期望 revision. 修订只追加; 主控不能覆盖或冒充作者的 disposition. 生成的完整批次视图可发布给每个成员, 而不必仅为复制记录去通知没有 finding 的成员.

check 与完成使用同一结构校验器. 等待结论的有效当前计划为 pending, 没有计划的旧版活跃卡需要 requirements; 已有证据格式错误, 身份错误, 副本缺失或关闭绑定过期是错误. done 另要求已封存的计划, 覆盖每个周期成员, 所有必需角色的成功结论, 作者覆盖, 以及每个批次在其最终目标处关闭. 显式 N/A 有效; 空索引, 全部失败的 run, 空结论, 任意旧角色 PASS 与部分发布无效.

审核关闭在审核层校验 Git, 并存储绑定到最终目标的证据. 看板重新校验结构与绑定, 不解释 Git, 也不声称已集成. 既有授权, 实际的源分支交付与最终 Git 校验仍是分开的强制工作流职责. 需要时通过 `review progress` 查看机器进度; 继续执行不要求仅为让 check 通过而捏造语义 PASS.

对每个角色都显式 N/A 的非 Git 项目, 计划可对 base 与 target_commit 都使用 `N/A`. 关闭时把 Git 记为不适用, 不声称已校验 commit. 任何必需角色仍需要真实的 commit 目标.

- 重新领取已有计划的卡 (见「领取, 启动与协调」) 可能改变其执行周期. `review progress` 随后报告 `requirements-needed` 与完整的 `rebind_cycles` 映射. 用 `review extend-plan` 附既有计划 ID, 期望 revision, 该映射, author 与 basis, 为整个计划恢复相同的要求. 此操作不能改变批次, 封存, 角色要求, 成员状态或更早的证据. 旧的失败与 finding 仍有约束力; 禁止另建计划来丢弃它们. 后继 OWNER 可对被分派的 finding 追加自己的 disposition, 保留旧作者不可变的原件与记录谱系. advance 与 extend-plan 使用计划的精确 CWD.
