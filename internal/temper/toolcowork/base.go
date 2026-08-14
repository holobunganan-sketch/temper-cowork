// Package toolcowork 提供 Temper CoWork 工具(temper_*),注册进 Reasonix
// Tool Registry 供 Agent 调用。每个工具实现 tool.Tool 接口,通过 context
// 注入的 Store 访问器读写 CoWork 业务数据(Evidence/Decision/Artifact/
// Validation/Completion)。
//
// 这些工具不发起模型调用、不造第二套 Runtime——它们只是把 Agent 的执行
// 结果持久化为 Temper 的可追溯 CoWork 记录。
package toolcowork

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"reasonix/internal/temper/domain"
)

// StoreAccessor 是工具访问 CoWork store 的最小接口。desktop 通过
// WithStore 注入;测试注入内存实现。
type StoreAccessor interface {
	CreateWork(projectID, title, goal, modelRef, qualityProfile string) (*domain.Work, error)
	GetWork(id string) (*domain.Work, error)
	UpdateWorkStatus(id string, status domain.WorkStatus) error
	SetWorkSessionRef(id, sessionRef string) error
	SetWorkTaskContract(id, contract string) error
	SetWorkFinalArtifact(id, artifactID string) error
	CreateEvidence(workID string, e domain.Evidence) error
	CreateDecision(workID string, d domain.Decision) error
	RegisterArtifact(workID string, a domain.Artifact) error
	RecordAcceptance(workID string, a domain.AcceptanceResult) error
	RecordQualityRun(workID string, q domain.QualityRun) error
}

type ctxKey struct{}

// WithStore 把 Store 访问器注入 context(desktop/测试使用)。
func WithStore(ctx context.Context, s StoreAccessor) context.Context {
	return context.WithValue(ctx, ctxKey{}, s)
}

// FromContext 从 context 取 Store 访问器;未注入时返回 nil。
func FromContext(ctx context.Context) StoreAccessor {
	if ctx == nil {
		return nil
	}
	s, _ := ctx.Value(ctxKey{}).(StoreAccessor)
	return s
}

// toolBase 提供 Tool 接口的公共字段。
type toolBase struct {
	name        string
	description string
	schema      json.RawMessage
}

func (b toolBase) Name() string        { return b.name }
func (b toolBase) Description() string { return b.description }
func (b toolBase) Schema() json.RawMessage {
	return b.schema
}
func (b toolBase) ReadOnly() bool { return true }

// requireStore 从 ctx 取 store,缺失时报错(工具在无 store 的会话中禁用)。
func requireStore(ctx context.Context) (StoreAccessor, error) {
	s := FromContext(ctx)
	if s == nil {
		return nil, fmt.Errorf("temper cowork store is unavailable in this session")
	}
	return s, nil
}

// nowUTC 返回统一 UTC 时间戳(RFC3339)。
func nowUTC() time.Time { return time.Now().UTC() }
