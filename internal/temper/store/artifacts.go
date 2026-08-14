package store

import (
	"strings"
	"time"

	"reasonix/internal/temper/domain"
)

func trim(s string) string { return strings.TrimSpace(s) }

// RegisterArtifact 登记一个交付物(真实 workspace 文件 + 元数据)。
func (s *Store) RegisterArtifact(workID string, a domain.Artifact) error {
	if a.ID == "" {
		a.ID = newID("art")
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	if a.UpdatedAt.IsZero() {
		a.UpdatedAt = a.CreatedAt
	}
	_, err := s.db.Exec(`INSERT INTO artifacts(id,project_id,work_id,relative_path,kind,title,description,sha256,size,validation,is_final,created_at,updated_at)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		a.ID, a.ProjectID, workID, trim(a.RelativePath), string(a.Kind), trim(a.Title), trim(a.Description),
		a.SHA256, a.Size, trim(a.Validation), boolInt(a.IsFinal), fmtTime(a.CreatedAt), fmtTime(a.UpdatedAt))
	return err
}

// ListArtifacts 列出工作的全部交付物。
func (s *Store) ListArtifacts(workID string) ([]domain.Artifact, error) {
	rows, err := s.db.Query(`SELECT id,project_id,work_id,relative_path,kind,title,description,sha256,size,validation,is_final,created_at,updated_at
		FROM artifacts WHERE work_id=? ORDER BY created_at`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Artifact
	for rows.Next() {
		var a domain.Artifact
		var kind, createdAt, updatedAt string
		var isFinal int
		if err := rows.Scan(&a.ID, &a.ProjectID, &a.WorkID, &a.RelativePath, &kind, &a.Title, &a.Description,
			&a.SHA256, &a.Size, &a.Validation, &isFinal, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		a.Kind = domain.ArtifactKind(kind)
		a.IsFinal = isFinal != 0
		a.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		a.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		out = append(out, a)
	}
	return out, rows.Err()
}

// RecordAcceptance 记录一条验收结果。
func (s *Store) RecordAcceptance(workID string, a domain.AcceptanceResult) error {
	if a.EvaluatedAt.IsZero() {
		a.EvaluatedAt = time.Now().UTC()
	}
	_, err := s.db.Exec(`INSERT INTO acceptance_results(work_id,criterion,status,evidence_ref,evaluated_at)
		VALUES(?,?,?,?,?)`,
		workID, trim(a.Criterion), string(a.Status), trim(a.EvidenceRef), fmtTime(a.EvaluatedAt))
	return err
}

// ListAcceptance 列出工作的全部验收结果。
func (s *Store) ListAcceptance(workID string) ([]domain.AcceptanceResult, error) {
	rows, err := s.db.Query(`SELECT id,work_id,criterion,status,evidence_ref,evaluated_at
		FROM acceptance_results WHERE work_id=? ORDER BY id`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AcceptanceResult
	for rows.Next() {
		var a domain.AcceptanceResult
		var status, evaluatedAt string
		if err := rows.Scan(&a.ID, &a.WorkID, &a.Criterion, &status, &a.EvidenceRef, &evaluatedAt); err != nil {
			return nil, err
		}
		a.Status = domain.AcceptanceStatus(status)
		a.EvaluatedAt, _ = time.Parse(time.RFC3339, evaluatedAt)
		out = append(out, a)
	}
	return out, rows.Err()
}

// RecordQualityRun 记录一次质量门运行。
func (s *Store) RecordQualityRun(workID string, q domain.QualityRun) error {
	if q.CreatedAt.IsZero() {
		q.CreatedAt = time.Now().UTC()
	}
	_, err := s.db.Exec(`INSERT INTO quality_runs(work_id,gate_type,passed,summary,created_at)
		VALUES(?,?,?,?,?)`,
		workID, trim(q.GateType), boolInt(q.Passed), trim(q.Summary), fmtTime(q.CreatedAt))
	return err
}

// ListQualityRuns 列出工作的全部质量门运行。
func (s *Store) ListQualityRuns(workID string) ([]domain.QualityRun, error) {
	rows, err := s.db.Query(`SELECT id,work_id,gate_type,passed,summary,created_at
		FROM quality_runs WHERE work_id=? ORDER BY id`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.QualityRun
	for rows.Next() {
		var q domain.QualityRun
		var gateType, createdAt string
		var passed int
		if err := rows.Scan(&q.ID, &q.WorkID, &gateType, &passed, &q.Summary, &createdAt); err != nil {
			return nil, err
		}
		q.GateType = gateType
		q.Passed = passed != 0
		q.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		out = append(out, q)
	}
	return out, rows.Err()
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func splitComma(s string) []string {
	var out []string
	for _, part := range splitSeq(s, ",") {
		if t := trim(part); t != "" {
			out = append(out, t)
		}
	}
	return out
}

func splitSeq(s, sep string) []string {
	var out []string
	start := 0
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			out = append(out, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	out = append(out, s[start:])
	return out
}
