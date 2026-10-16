import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import api, { getUser, setUser } from "../api";

export default function Login() {
  const nav = useNavigate();
  const [mode, setMode] = useState("login");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    if (getUser()) nav("/");
  }, []);

  const submit = async (e) => {
    e.preventDefault();
    setError("");
    try {
      const r = await api.post(`/${mode}`, { username, password });
      localStorage.setItem("token", r.data.token);
      setUser(r.data.user);
      nav("/");
      window.location.reload();
    } catch (err) {
      setError(err.response?.data?.error || "操作失败");
    }
  };

  return (
    <div className="form-card">
      <h2>{mode === "login" ? "登录 谷子小铺" : "注册新账号"}</h2>
      <form onSubmit={submit}>
        <input
          className="input"
          placeholder="用户名（至少2位）"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />
        <input
          className="input"
          type="password"
          placeholder="密码（至少6位）"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
        />
        {error && <div className="form-error">{error}</div>}
        <button className="btn btn-primary" style={{ width: "100%" }}>
          {mode === "login" ? "登录" : "注册并登录（送 10000 金币）"}
        </button>
      </form>
      <div style={{ textAlign: "center", marginTop: 16, fontSize: 14 }}>
        {mode === "login" ? (
          <>
            没有账号？
            <span className="nav-link" style={{ color: "var(--primary)" }} onClick={() => setMode("register")}>
              立即注册
            </span>
            <div style={{ color: "#bbb", marginTop: 10, fontSize: 13 }}>
              演示账号：demo / demo123
            </div>
          </>
        ) : (
          <>
            已有账号？
            <span className="nav-link" style={{ color: "var(--primary)" }} onClick={() => setMode("login")}>
              去登录
            </span>
          </>
        )}
      </div>
    </div>
  );
}
