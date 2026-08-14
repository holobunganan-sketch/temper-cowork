// TemperWorkPanel 展示当前项目的 Formal Work 列表,并支持创建与状态流转。
// 它通过 Wails 绑定调用真实 Go 后端(ListTemperWorks/CreateTemperWork/
// UpdateTemperWorkStatus),不是 mock。
import { useCallback, useEffect, useState } from "react";
import type { TemperWorkView } from "../lib/types";

interface Props {
  projectID: string;
  projectName?: string;
}

const WORK_STATUSES = [
  "draft", "ready", "running", "waiting_user", "blocked",
  "reviewing", "validating", "completed", "failed", "cancelled",
] as const;

export function TemperWorkPanel({ projectID, projectName }: Props) {
  const [works, setWorks] = useState<TemperWorkView[]>([]);
  const [title, setTitle] = useState("");
  const [goal, setGoal] = useState("");
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    try {
      const items = (await window.go?.main?.App?.ListTemperWorks(projectID)) ?? [];
      setWorks(items);
      setError(null);
    } catch (e) {
      setError(String(e));
    }
  }, [projectID]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const createWork = async () => {
    if (!title.trim()) return;
    setCreating(true);
    setError(null);
    try {
      await window.go?.main?.App?.CreateTemperWork(projectID, title.trim(), goal.trim(), "", "default");
      setTitle("");
      setGoal("");
      await refresh();
    } catch (e) {
      setError(String(e));
    } finally {
      setCreating(false);
    }
  };

  const setStatus = async (workID: string, status: string) => {
    try {
      await window.go?.main?.App?.UpdateTemperWorkStatus(workID, status);
      await refresh();
    } catch (e) {
      setError(String(e));
    }
  };

  return (
    <div className="temper-work" aria-label="Temper Works">
      <header className="temper-work__head">
        <h2 className="temper-work__title">
          Works{projectName ? ` — ${projectName}` : ""}
        </h2>
        <span className="temper-work__count">{works.length}</span>
      </header>

      <div className="temper-work__create">
        <input
          className="temper-work__input"
          placeholder="Work title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          aria-label="Work title"
        />
        <input
          className="temper-work__input"
          placeholder="Goal (optional)"
          value={goal}
          onChange={(e) => setGoal(e.target.value)}
          aria-label="Work goal"
        />
        <button
          className="temper-work__create-btn"
          onClick={() => void createWork()}
          disabled={creating || !title.trim()}
        >
          {creating ? "Creating…" : "Create work"}
        </button>
      </div>

      {error && <p className="temper-work__error">{error}</p>}

      <ul className="temper-work__list">
        {works.length === 0 && (
          <li className="temper-work__empty">No formal works yet. Create one to start.</li>
        )}
        {works.map((w) => (
          <li key={w.id} className="temper-work__item">
            <div className="temper-work__item-head">
              <span className="temper-work__item-title">{w.title}</span>
              <span className={`temper-work__badge temper-work__badge--${w.status}`}>{w.status}</span>
            </div>
            {w.goal && <p className="temper-work__item-goal">{w.goal}</p>}
            <div className="temper-work__actions">
              {WORK_STATUSES.filter((s) => s !== w.status && s !== "draft").map((s) => (
                <button
                  key={s}
                  className="temper-work__action"
                  onClick={() => void setStatus(w.id, s)}
                >
                  {s}
                </button>
              ))}
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
