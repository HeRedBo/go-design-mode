# Go Design Patterns

Go 语言设计模式与并发模式集合，包含 Go 特有模式和经典设计模式的完整实现与测试用例。

## 📁 项目结构

```
design-mode/
├── 01-go-specific/              # Go 语言特有模式
│   ├── functional_options/      # 函数选项模式
│   ├── error_wrapping/          # 错误包装模式
│   ├── defer_cleanup/           # Defer 资源清理
│   ├── embedding/               # 嵌入组合模式
│   └── graceful_shutdown/       # 优雅退出模式
│
├── 02-concurrency/              # Go 并发模式
│   ├── worker_pool/             # 工作池模式
│   ├── pipeline/                # 流水线模式
│   ├── context_cancellation/    # Context 取消模式
│   ├── fan_out_in/              # Fan-Out/Fan-In 模式
│   ├── rate_limiter/            # 限流器模式
│   ├── semaphore/               # 信号量模式
│   ├── pubsub/                  # 发布订阅模式
│   └── backoff_retry/           # 退避重试模式
│
├── 03-creational/               # 创建型设计模式（进行中）
│   └── singleton/               # 单例模式
│
├── .gitignore
├── go.mod
└── README.md
```

## 🎯 已实现的模式

### Go 特有模式 (01-go-specific)

| 模式 | 说明 | 特点 |
|------|------|------|
| **Functional Options** | 函数选项模式 | 优雅处理可选参数，标准库广泛使用 |
| **Error Wrapping** | 错误包装 | 哨兵错误、%w 包装、errors.Is/As |
| **Defer Cleanup** | Defer 资源清理 | 延迟执行、panic 恢复、资源管理 |
| **Embedding** | 嵌入组合 | 方法提升、接口组合、替代继承 |
| **Graceful Shutdown** | 优雅退出 | Context 控制、服务管理、信号监听 |

### 并发模式 (02-concurrency)

| 模式 | 说明 | 特点 |
|------|------|------|
| **Worker Pool** | 工作池 | 并发控制、panic 恢复、批量处理 |
| **Pipeline** | 流水线 | Stage 模式、Map/Filter、错误处理 |
| **Context Cancellation** | Context 取消 | 超时控制、级联取消、传值 |
| **Fan-Out/Fan-In** | 多路分发汇聚 | 并行处理、结果合并 |
| **Rate Limiter** | 限流器 | 令牌桶算法、速率控制 |
| **Semaphore** | 信号量 | 并发限制、资源控制 |
| **Pub/Sub** | 发布订阅 | 主题订阅、消息广播 |
| **Backoff Retry** | 退避重试 | 指数退避、随机化、可取消 |

### 经典设计模式 (03-creational)

| 模式 | 说明 | 特点 |
|------|------|------|
| **Singleton** | 单例模式 | sync.Once 实现、线程安全 |

## 🚀 快速开始

### 运行所有测试

```bash
# 运行全部测试
go test ./... -v

# 运行特定模式测试
go test ./01-go-specific/functional_options -v
go test ./02-concurrency/worker_pool -v
```

### 使用示例

```go
// Functional Options 示例
server := NewServer(
    WithHost("localhost"),
    WithPort(8080),
    WithTimeout(30*time.Second),
    WithDebug(),
)

// Worker Pool 示例
pool := workerpool.New(4, 10)
pool.Start()
pool.Submit(Task{ID: 1, Data: data, Handler: processFn})
pool.Close()
pool.Wait()

// Context Cancellation 示例
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
err := LongRunningTask(ctx)
```

## 📝 每个模式包含

- ✅ 完整的实现代码
- ✅ 详细的注释说明
- ✅ 完整的测试用例
- ✅ 使用示例

## 🛠 技术特点

- **Go 1.23+**
- **并发安全**：所有模式都考虑了并发场景
- **Context 支持**：支持取消和超时控制
- **错误处理**：完善的错误包装和处理
- **测试覆盖**：每个模式都有完整的测试用例

## 📖 学习路径

1. **第一阶段**：Go 特有模式 ✅（已完成）
2. **第二阶段**：并发模式 ✅（已完成）
3. **第三阶段**：经典 23 种设计模式 🚧（进行中）

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
