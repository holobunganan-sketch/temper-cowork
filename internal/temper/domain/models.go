// Package domain 定义 Temper CoWork 的领域模型。
//
// 这些类型是 Temper-owned 的 CoWork 语义(Project / Work / Evidence /
// Decision / Artifact / Validation / Quality / Completion),与 Reasonix
// Runtime 的 Session/Task 模型正交。持久化由 internal/temper/store 负责。
package domain

import "time"

// Project 是 Temper 工作区注册项。workspaceRoot 必须是真实存在的目录;
// 从 Temper 移除项目绝不删除该目录。
type Project struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	WorkspaceRoot string    `json:"workspaceRoot"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	LastOpenedAt  time.Time `json:"lastOpenedAt,omitempty"`
	Archived      bool      `json:"archived"`
}

// WorkStatus 是 Formal Work 的生命周期状态(Master PHASE F)。
type WorkStatus string

const (
	WorkDraft      WorkStatus = "draft"
	WorkReady      WorkStatus = "ready"
	WorkRunning    WorkStatus = "running"
	WorkWaiting    WorkStatus = "waiting_user"
	WorkBlocked    WorkStatus = "blocked"
	WorkReviewing  WorkStatus = "reviewing"
	WorkValidating WorkStatus = "validating"
	WorkCompleted  WorkStatus = "completed"
	WorkFailed     WorkStatus = "failed"
	WorkCancelled  WorkStatus = "cancelled"
)

// Work 是一份正式工作。goal 是用户意图;reasonixSessionRef 关联 Reasonix
// Session;finalArtifactID 指向完成时的最终交付物。
type Work struct {
	ID                string     `json:"id"`
	ProjectID         string     `json:"projectId"`
	Title             string     `json:"title"`
	Goal              string     `json:"goal"`
	Status            WorkStatus `json:"status"`
	ReasonixSessionRef string    `json:"reasonixSessionRef,omitempty"`
	ModelRef          string     `json:"modelRef,omitempty"`
	QualityProfile    string     `json:"qualityProfile,omitempty"`
	TaskContract      string     `json:"taskContract,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
	StartedAt         *time.Time `json:"startedAt,omitempty"`
	CompletedAt       *time.Time `json:"completedAt,omitempty"`
	FinalArtifactID   string     `json:"finalArtifactId,omitempty"`
}

// WorkEvent 记录 Work 生命周期中的事件(状态变更、证据登记、决策等)。
type WorkEvent struct {
	ID        int64     `json:"id"`
	WorkID    string    `json:"workId"`
	EventType string    `json:"eventType"`
	Detail    string    `json:"detail,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// Evidence 是支持某个结论的可追溯依据。
type Evidence struct {
	ID         string    `json:"id"`
	WorkID     string    `json:"workId"`
	Summary    string    `json:"summary"`
	SourceType string    `json:"sourceType"`
	SourceRef  string    `json:"sourceRef"`
	Supports   string    `json:"supports"`
	Timestamp  time.Time `json:"timestamp"`
}

// Decision 是正式工作中的一个决策。
type Decision struct {
	ID          string    `json:"id"`
	WorkID      string    `json:"workId"`
	Decision    string    `json:"decision"`
	Rationale   string    `json:"rationale"`
	Alternatives string   `json:"alternatives,omitempty"`
	EvidenceIDs []string  `json:"evidenceIds,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// ArtifactKind 是交付物的稳定类型。
type ArtifactKind string

const (
	ArtifactMarkdown ArtifactKind = "md"
	ArtifactText     ArtifactKind = "txt"
	ArtifactJSON     ArtifactKind = "json"
	ArtifactCSV      ArtifactKind = "csv"
	ArtifactHTML     ArtifactKind = "html"
	ArtifactSVG      ArtifactKind = "svg"
	ArtifactPNG      ArtifactKind = "png"
	ArtifactJPEG     ArtifactKind = "jpeg"
	ArtifactSource   ArtifactKind = "source"
)

// Artifact 是真实 workspace 文件 + 元数据。
type Artifact struct {
	ID           string       `json:"id"`
	ProjectID    string       `json:"projectId"`
	WorkID       string       `json:"workId"`
	RelativePath string       `json:"relativePath"`
	Kind         ArtifactKind `json:"kind"`
	Title        string       `json:"title"`
	Description  string       `json:"description,omitempty"`
	SHA256       string       `json:"sha256"`
	Size         int64        `json:"size"`
	Validation   string       `json:"validation,omitempty"`
	IsFinal      bool         `json:"isFinal"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

// AcceptanceStatus 是验收标准的结果。
type AcceptanceStatus string

const (
	AcceptancePending   AcceptanceStatus = "pending"
	AcceptancePass      AcceptanceStatus = "pass"
	AcceptanceFail      AcceptanceStatus = "fail"
	AcceptanceUncertain AcceptanceStatus = "uncertain"
)

// AcceptanceResult 记录单条验收标准的评估结果。
type AcceptanceResult struct {
	ID           int64            `json:"id"`
	WorkID       string           `json:"workId"`
	Criterion    string           `json:"criterion"`
	Status       AcceptanceStatus `json:"status"`
	EvidenceRef  string           `json:"evidenceRef,omitempty"`
	EvaluatedAt  time.Time        `json:"evaluatedAt"`
}

// QualityRun 是一次质量门(Validation/Review)运行。
type QualityRun struct {
	ID        int64     `json:"id"`
	WorkID    string    `json:"workId"`
	GateType  string    `json:"gateType"` // validation | review
	Passed    bool      `json:"passed"`
	Summary   string    `json:"summary,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
