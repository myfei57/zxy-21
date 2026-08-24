# LIMS 实验室检测样本与报告平台

zxy-21 是一个本地文件持久化的检验实验室流程平台：样本登记、检测任务分配、
仪器连接与派发、原始结果与质控判定、报告签发与归档、仪器并发配额和全流程
审计流水，全部通过内置控制台页面与 JSON API 操作。

## 构建

依赖已 vendor，无需联网：

    go build -mod=vendor ./...
    go vet -mod=vendor ./...
    go test -mod=vendor ./...

## 运行

    LIMS_ADDR=:8080 LIMS_DATA_DIR=./data go run -mod=vendor ./cmd/limsd

启动后访问：

- http://localhost:8080/samples 样本登记台
- http://localhost:8080/tasks 检测任务台
- http://localhost:8080/reports 报告签发台
- http://localhost:8080/audit 审计流水

数据以 JSON 文件保存在 `LIMS_DATA_DIR` 下，重启后自动恢复。
