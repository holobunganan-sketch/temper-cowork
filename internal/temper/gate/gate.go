// Package gate 实现 Temper 的 Host Completion Gate 与 Quality Gate。
//
// Completion Gate(Master PHASE H)在 Work 完成前检查:
//
//	Task Contract exists
//	required deliverable exists
//	final artifact exists
//	file exists
//	hash current
//	validation no fail
//	no blocking error
//	acceptance evaluated
//
// Quality Gate 检查 Validation / Review 质量门运行记录。
package gate

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"reasonix/internal/temper/domain"
	"reasonix/internal/temper/store"
)

// CheckResult 是 Completion Gate 的逐项检查结果。
type CheckResult struct {
	Item   string `json:"item"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail,omitempty"`
}

// CompletionReport 汇总 Completion Gate 结果。
type CompletionReport struct {
	Passed bool          `json:"passed"`
	Checks []CheckResult `json:"checks"`
}

// QualitySummary 汇总 Quality Gate 结果。
type QualitySummary struct {
	Passed     bool                `json:"passed"`
	Validation []domain.QualityRun `json:"validation"`
	Reviews    []domain.QualityRun `json:"reviews"`
}

// CompletionGate 执行 Host Completion Gate 检查。
// workspaceRoot 是工件文件校验的基准目录;为空时跳过文件级检查。
func CompletionGate(s *store.Store, workID, workspaceRoot string) (*CompletionReport, error) {
	w, err := s.GetWork(workID)
	if err != nil {
		return nil, fmt.Errorf("gate: get work: %w", err)
	}
	report := &CompletionReport{}
	check := func(item string, passed bool, detail string) {
		report.Checks = append(report.Checks, CheckResult{Item: item, Passed: passed, Detail: detail})
		if !passed {
			report.Passed = false
		}
	}
	report.Passed = true

	// 1. Task Contract exists
	check("task_contract", strings.TrimSpace(w.TaskContract) != "", "Task Contract 已编译")

	// 2. Final artifact exists(元数据)
	check("final_artifact", w.FinalArtifactID != "", "最终交付物已登记")

	// 3. required deliverable:final artifact 对应的记录存在
	artifacts, err := s.ListArtifacts(workID)
	if err != nil {
		return nil, fmt.Errorf("gate: list artifacts: %w", err)
	}
	var finalArtifact *domain.Artifact
	for i := range artifacts {
		if artifacts[i].ID == w.FinalArtifactID {
			finalArtifact = &artifacts[i]
			break
		}
	}
	check("artifact_record", finalArtifact != nil, "最终交付物记录存在")

	// 4-5. file exists + hash current
	if finalArtifact != nil && workspaceRoot != "" {
		full := filepath.Join(workspaceRoot, filepath.FromSlash(finalArtifact.RelativePath))
		data, err := os.ReadFile(full)
		if err != nil {
			check("file_exists", false, fmt.Sprintf("无法读取 %s: %v", finalArtifact.RelativePath, err))
		} else {
			check("file_exists", true, finalArtifact.RelativePath)
			sum := sha256.Sum256(data)
			actual := hex.EncodeToString(sum[:])
			check("hash_current", strings.EqualFold(actual, finalArtifact.SHA256),
				fmt.Sprintf("sha256 匹配(%s)", finalArtifact.RelativePath))
		}
	} else if finalArtifact == nil {
		check("file_exists", false, "无最终交付物可校验")
	} else {
		check("file_exists", true, "workspace 未提供,跳过文件级校验")
	}

	// 6. validation no fail:验收结果中不得有 fail(uncertain 不算 fail)
	acceptance, err := s.ListAcceptance(workID)
	if err != nil {
		return nil, fmt.Errorf("gate: list acceptance: %w", err)
	}
	anyFail := false
	for _, a := range acceptance {
		if a.Status == domain.AcceptanceFail {
			anyFail = true
			break
		}
	}
	check("validation_no_fail", !anyFail, "验收无 fail")

	// 7. no blocking error:Work 状态不是 blocked/failed
	switch w.Status {
	case domain.WorkBlocked, domain.WorkFailed:
		check("no_blocking_error", false, "Work 处于阻塞/失败状态")
	default:
		check("no_blocking_error", true, "无阻塞错误")
	}

	// 8. acceptance evaluated:至少一条验收记录
	check("acceptance_evaluated", len(acceptance) > 0, "验收已评估")

	return report, nil
}

// QualityGate 汇总质量门运行记录(Validation/Review)。
func QualityGate(s *store.Store, workID string) (*QualitySummary, error) {
	runs, err := s.ListQualityRuns(workID)
	if err != nil {
		return nil, fmt.Errorf("gate: list quality runs: %w", err)
	}
	sum := &QualitySummary{Passed: true}
	for _, r := range runs {
		switch r.GateType {
		case "validation":
			sum.Validation = append(sum.Validation, r)
		case "review":
			sum.Reviews = append(sum.Reviews, r)
		}
		if !r.Passed {
			sum.Passed = false
		}
	}
	return sum, nil
}
