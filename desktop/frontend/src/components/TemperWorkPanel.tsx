// TemperWorkPanel 展示当前项目的 Formal Work 列表,并支持创建与状态流转。
// 它通过 Wails 绑定调用真实 Go 后端(ListTemperWorks/CreateTemperWork/
// UpdateTemperWorkStatus),不是 mock。文案走 i18n 字典(en/zh/zh-TW)。
import { useCallback, useEffect, useState } from "react";
import { useI18n } from "../lib/i18n";
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
  const { t } = useI18n();
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
    <div className="temper-work" aria-label={t("temper.works")}>
      <header className="temper-work__head">
        <h2 className="temper-work__title">
          {t("temper.projectLabel", { project: projectName || "" })}
        </h2>
        <span className="temper-work__count">{works.length}</span>
      </header>

      <div className="temper-work__create">
        <input
          className="temper-work__input"
          placeholder={t("temper.workTitlePlaceholder")}
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          aria-label={t("temper.workTitlePlaceholder")}
        />
        <input
          className="temper-work__input"
          placeholder={t("temper.workGoalPlaceholder")}
          value={goal}
          onChange={(e) => setGoal(e.target.value)}
          aria-label={t("temper.workGoalPlaceholder")}
        />
        <button
          className="temper-work__create-btn"
          onClick={() => void createWork()}
          disabled={creating || !title.trim()}
        >
          {creating ? t("temper.creatingWork") : t("temper.createWork")}
        </button>
      </div>

      {error && <p className="temper-work__error">{error}</p>}

      <ul className="temper-work__list">
        {works.length === 0 && (
          <li className="temper-work__empty">{t("temper.noWorksYet")}</li>
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
