// Package contract 把 Temper Work Form 编译为 Reasonix Task Contract。
//
// Work Form 字段:Project / Goal / Materials / Deliverable / Audience /
// Constraints / Acceptance Criteria / Source Policy / Pause Policy /
// Quality / Model。编译结果是纯文本 Task Contract(写入 works.task_contract),
// 由 Reasonix Agent 作为任务上下文执行,Host Completion Gate 据其验收。
//
// 本包不发起模型调用,不造第二套 Prompt Runtime——它只是把 Work Form
// 的结构化字段格式化为 Reasonix Task Contract 的稳定文本协议。
package contract

import (
	"fmt"
	"strings"
)

// WorkForm 是创建 Formal Work 时的用户输入(Master PHASE F)。
type WorkForm struct {
	Project            string   `json:"project"`
	Goal               string   `json:"goal"`
	Materials          []string `json:"materials"`
	Deliverable        string   `json:"deliverable"`
	Audience           string   `json:"audience"`
	Constraints        []string `json:"constraints"`
	AcceptanceCriteria []string `json:"acceptanceCriteria"`
	SourcePolicy       string   `json:"sourcePolicy"`
	PausePolicy        string   `json:"pausePolicy"`
	Quality            string   `json:"quality"`
	Model              string   `json:"model"`
}

// Compile 把 Work Form 编译为 Reasonix Task Contract 文本。
func Compile(f WorkForm) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Task Contract\n\n")
	fmt.Fprintf(&b, "## Context\n")
	fmt.Fprintf(&b, "Project: %s\n", orDefault(f.Project, "(unspecified)"))
	if len(f.Materials) > 0 {
		fmt.Fprintf(&b, "Materials:\n")
		for _, m := range f.Materials {
			if strings.TrimSpace(m) != "" {
				fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(m))
			}
		}
	}
	if f.Audience != "" {
		fmt.Fprintf(&b, "Audience: %s\n", f.Audience)
	}

	fmt.Fprintf(&b, "\n## Request\n%s\n", orDefault(strings.TrimSpace(f.Goal), "(no goal provided)"))

	fmt.Fprintf(&b, "\n## Output format\n")
	fmt.Fprintf(&b, "Deliverable: %s\n", orDefault(f.Deliverable, "a written deliverable in the workspace"))

	fmt.Fprintf(&b, "\n## Constraints\n")
	if len(f.Constraints) > 0 {
		for _, c := range f.Constraints {
			if strings.TrimSpace(c) != "" {
				fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(c))
			}
		}
	} else {
		fmt.Fprintf(&b, "- none specified\n")
	}
	if f.SourcePolicy != "" {
		fmt.Fprintf(&b, "Source policy: %s\n", f.SourcePolicy)
	}
	if f.Quality != "" {
		fmt.Fprintf(&b, "Quality profile: %s\n", f.Quality)
	}
	if f.Model != "" {
		fmt.Fprintf(&b, "Model: %s\n", f.Model)
	}

	fmt.Fprintf(&b, "\n## Acceptance criteria\n")
	if len(f.AcceptanceCriteria) > 0 {
		for i, ac := range f.AcceptanceCriteria {
			if strings.TrimSpace(ac) != "" {
				fmt.Fprintf(&b, "- [ ] AC%d: %s\n", i+1, strings.TrimSpace(ac))
			}
		}
	} else {
		fmt.Fprintf(&b, "- [ ] (no explicit criteria; evaluate against the deliverable)\n")
	}

	fmt.Fprintf(&b, "\n## Pause policy\n%s\n", orDefault(f.PausePolicy, "pause on any request for user decision"))
	return b.String()
}

func orDefault(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return strings.TrimSpace(v)
}
