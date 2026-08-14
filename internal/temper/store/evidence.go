package store

import (
	"strings"
	"time"

	"reasonix/internal/temper/domain"
)

// CreateEvidence 登记一条证据。
func (s *Store) CreateEvidence(workID string, e domain.Evidence) error {
	if e.ID == "" {
		e.ID = newID("ev")
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	_, err := s.db.Exec(`INSERT INTO evidence(id,work_id,summary,source_type,source_ref,supports,timestamp)
		VALUES(?,?,?,?,?,?,?)`,
		e.ID, workID, trim(e.Summary), trim(e.SourceType), trim(e.SourceRef), trim(e.Supports),
		fmtTime(e.Timestamp))
	return err
}

// ListEvidence 列出工作的全部证据。
func (s *Store) ListEvidence(workID string) ([]domain.Evidence, error) {
	rows, err := s.db.Query(`SELECT id,work_id,summary,source_type,source_ref,supports,timestamp
		FROM evidence WHERE work_id=? ORDER BY timestamp`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Evidence
	for rows.Next() {
		var e domain.Evidence
		var ts string
		if err := rows.Scan(&e.ID, &e.WorkID, &e.Summary, &e.SourceType, &e.SourceRef, &e.Supports, &ts); err != nil {
			return nil, err
		}
		e.Timestamp, _ = time.Parse(time.RFC3339, ts)
		out = append(out, e)
	}
	return out, rows.Err()
}

// CreateDecision 记录一条决策。
func (s *Store) CreateDecision(workID string, d domain.Decision) error {
	if d.ID == "" {
		d.ID = newID("dec")
	}
	if d.Timestamp.IsZero() {
		d.Timestamp = time.Now().UTC()
	}
	evidenceIDs := strings.Join(d.EvidenceIDs, ",")
	_, err := s.db.Exec(`INSERT INTO decisions(id,work_id,decision,rationale,alternatives,evidence_ids,timestamp)
		VALUES(?,?,?,?,?,?,?)`,
		d.ID, workID, trim(d.Decision), trim(d.Rationale), trim(d.Alternatives), evidenceIDs,
		fmtTime(d.Timestamp))
	return err
}

// ListDecisions 列出工作的全部决策。
func (s *Store) ListDecisions(workID string) ([]domain.Decision, error) {
	rows, err := s.db.Query(`SELECT id,work_id,decision,rationale,alternatives,evidence_ids,timestamp
		FROM decisions WHERE work_id=? ORDER BY timestamp`, workID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Decision
	for rows.Next() {
		var d domain.Decision
		var evidenceIDs, ts string
		if err := rows.Scan(&d.ID, &d.WorkID, &d.Decision, &d.Rationale, &d.Alternatives, &evidenceIDs, &ts); err != nil {
			return nil, err
		}
		d.Timestamp, _ = time.Parse(time.RFC3339, ts)
		if evidenceIDs != "" {
			d.EvidenceIDs = splitComma(evidenceIDs)
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
