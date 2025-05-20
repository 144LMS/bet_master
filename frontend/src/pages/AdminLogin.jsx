import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { adminLogin } from "../api/api";
import "../style/AdminLogin.css";

export default function AdminLogin() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      const resp = await adminLogin({ email, password });
      if (resp.data.user && resp.data.user.role === "admin") {
        localStorage.setItem("token", resp.data.token);
        navigate("/admin");
      } else {
        setError("У вас нет прав администратора");
      }
    } catch (err) {
      setError(
        err?.response?.data?.error || "Ошибка авторизации. Проверьте данные."
      );
    }
  };

  return (
    <div className="admin-login-bg-strict">
      <div className="admin-login-card">
        <div className="admin-login-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h2 className="admin-login-title-strict">Вход для администратора</h2>
        <form className="AdminLoginForm" onSubmit={handleSubmit}>
          <input
            type="email"
            placeholder="Email администратора"
            value={email}
            onChange={e => setEmail(e.target.value)}
            required
          />
          <input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={e => setPassword(e.target.value)}
            required
          />
          <button type="submit">Войти</button>
        </form>
        {error && <div className="admin-login-error">{error}</div>}
      </div>
    </div>
  );
}