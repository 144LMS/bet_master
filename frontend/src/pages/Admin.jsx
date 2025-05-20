import { useState, useEffect } from "react";
import axios from "axios";
import "../style/Admin.css";

const api = axios.create({
  baseURL: "http://localhost:8080/api/admin",
  withCredentials: true,
});

export default function Admin() {
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [form, setForm] = useState({
    team1: "",
    team2: "",
    start_time: "",
    sport_type: "",
  });
  const [settleId, setSettleId] = useState(null);
  const [winner, setWinner] = useState("");
  const [actionMsg, setActionMsg] = useState("");

  // Загрузка матчей
  const loadMatches = async () => {
    setLoading(true);
    setError("");
    try {
      const resp = await api.get("/matches");
      setMatches(resp.data.matches || []);
    } catch (err) {
      setError(err?.response?.data?.error || "Ошибка загрузки матчей");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadMatches();
  }, []);

  // Создание матча
  const toISOString = (dt) => {
    // "2025-05-22T22:20" -> "2025-05-22T22:20:00Z"
    if (!dt) return "";
    return new Date(dt).toISOString();
  };

  const handleCreate = async (e) => {
    e.preventDefault();
    setActionMsg("");
    try {
      const resp = await api.post("/matches", {
        team1: form.team1,
        team2: form.team2,
        start_time: toISOString(form.start_time),
        sport_type: form.sport_type,
      });
      setActionMsg(resp.data.message || "Матч создан");
      setForm({ team1: "", team2: "", start_time: "", sport_type: "" });
      loadMatches();
    } catch (err) {
      setActionMsg(err?.response?.data?.error || "Ошибка создания матча");
    }
  };

  // Удаление матча
  const handleDelete = async (id) => {
    setActionMsg("");
    if (!window.confirm("Удалить матч?")) return;
    try {
      const resp = await api.delete(`/matches/${id}`);
      setActionMsg(resp.data.message || "Матч удалён");
      loadMatches();
    } catch (err) {
      setActionMsg(err?.response?.data?.error || "Ошибка удаления");
    }
  };

  // Засчитать исход
  const handleSettle = async (e) => {
    e.preventDefault();
    setActionMsg("");
    try {
      const resp = await api.post(`/matches/${settleId}/settle`, {
        winner,
      });
      setActionMsg(resp.data.message || "Матч рассчитан");
      setSettleId(null);
      setWinner("");
      loadMatches();
    } catch (err) {
      setActionMsg(err?.response?.data?.error || "Ошибка расчёта");
    }
  };

  // Сортируем: не-completed сверху, completed снизу
  const sortedMatches = [
    ...matches.filter(m => m.status !== "completed"),
    ...matches.filter(m => m.status === "completed")
  ];

  return (
    <div className="admin-bg-strict">
      <div className="admin-card">
        <div className="admin-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h1 className="admin-title-strict">Панель администратора</h1>
        {error && <div className="form-error">{error}</div>}
        {actionMsg && <div className="form-success">{actionMsg}</div>}

        <h2 className="admin-section-title">Создать матч</h2>
        <form className="admin-match-form" onSubmit={handleCreate}>
          <input
            type="text"
            placeholder="Команда 1"
            value={form.team1}
            onChange={e => setForm(f => ({ ...f, team1: e.target.value }))}
            required
          />
          <input
            type="text"
            placeholder="Команда 2"
            value={form.team2}
            onChange={e => setForm(f => ({ ...f, team2: e.target.value }))}
            required
          />
          <input
            type="datetime-local"
            placeholder="Дата начала"
            value={form.start_time}
            onChange={e => setForm(f => ({ ...f, start_time: e.target.value }))}
            required
          />
          <input
            type="text"
            placeholder="Вид спорта"
            value={form.sport_type}
            onChange={e => setForm(f => ({ ...f, sport_type: e.target.value }))}
            required
          />
          <button type="submit" className="admin-create-btn">Создать</button>
        </form>

        <h2 className="admin-section-title">Список матчей</h2>
        {loading ? (
          <div className="admin-loading">Загрузка...</div>
        ) : (
          <div className="admin-table-wrapper">
            <table className="admin-matches-table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Команда 1</th>
                  <th>Команда 2</th>
                  <th>Старт</th>
                  <th>Статус</th>
                  <th>Вид спорта</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                {sortedMatches.map(match => (
                  <tr key={match.ID}>
                    <td>{match.ID}</td>
                    <td>{match.team1}</td>
                    <td>{match.team2}</td>
                    <td>{new Date(match.start_time).toLocaleString()}</td>
                    <td>{match.status}</td>
                    <td>{match.sport_type}</td>
                    <td>
                      <div className="admin-table-actions">
                        {match.status !== "completed" && (
                            <button className="admin-delete-btn" onClick={() => handleDelete(match.ID)}>
                            Удалить
                            </button>
                        )}
                        {match.status !== "completed" && (
                            <button className="admin-settle-btn" onClick={() => setSettleId(match.ID)}>
                            Рассчитать
                            </button>
                        )}
                        </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}

        {/* Форма расчёта исхода */}
        {settleId && (
          <form className="admin-settle-form" onSubmit={handleSettle}>
            <h3 className="admin-settle-title">Рассчитать матч ID: {settleId}</h3>
            <label>
              Победитель:&nbsp;
              <select
                value={winner}
                onChange={e => setWinner(e.target.value)}
                required
              >
                <option value="">Выбрать</option>
                <option value="team1">Команда 1</option>
                <option value="team2">Команда 2</option>
                <option value="draw">Ничья</option>
              </select>
            </label>
            <div className="admin-settle-actions">
              <button type="submit" className="admin-settle-btn">Рассчитать</button>
              <button type="button" className="admin-cancel-btn" onClick={() => setSettleId(null)}>
                Отмена
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}