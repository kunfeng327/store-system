import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import api, { getUser, setUser } from "../api";

const RARITY_TEXT = {
  N: "普通 N",
  R: "稀有 R",
  SR: "超稀有 SR",
  SSR: "传说 SSR ✨",
};

const PAY_METHODS = {
  wechat: "微信支付 💚",
  alipay: "支付宝 💙",
  union: "银联 💳",
};

export default function BlindBox() {
  const user = getUser();
  const [drawing, setDrawing] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState("");
  const [records, setRecords] = useState([]);
  const [payMethod, setPayMethod] = useState("wechat");

  const loadRecords = () => {
    if (!user) return;
    api.get("/box/records").then((r) => setRecords(r.data));
  };

  useEffect(() => {
    loadRecords();
  }, []);

  const draw = () => {
    if (drawing) return;
    setDrawing(true);
    setResult(null);
    setError("");
    setTimeout(async () => {
      try {
        const r = await api.post("/box/draw", { payment_method: payMethod });
        setResult(r.data);
        const u = { ...user, coins: r.data.coins_left };
        setUser(u);
        loadRecords();
      } catch (e) {
        setError(e.response?.data?.error || "抽取失败");
      } finally {
        setDrawing(false);
      }
    }, 900);
  };

  return (
    <div>
      <div className="box-hero">
        <h1 style={{ marginBottom: 8 }}>🎁 魔法盲盒</h1>
        <p style={{ color: "#888", marginBottom: 24 }}>
          60 金币一抽 · 随机获得一件谷子 · 概率公示：N 60% / R 25% / SR 12% / SSR 3%
        </p>

        {!user ? (
          <div className="empty">
            <Link to="/login" style={{ color: "var(--primary)" }}>登录后开盒</Link>
          </div>
        ) : (
          <>
            <div
              className={`box-gift ${drawing ? "drawing" : ""}`}
              onClick={draw}
              title="点击开盒"
            >
              🎁
            </div>
            <div>
              <button className="btn btn-primary" onClick={draw} disabled={drawing}>
                {drawing ? "开启中…" : "开启盲盒（60 🪙）"}
              </button>
            </div>
            <div className="tabs" style={{ justifyContent: "center", marginTop: 14 }}>
              {Object.entries(PAY_METHODS).map(([key, label]) => (
                <span
                  key={key}
                  className={`tab ${payMethod === key ? "active" : ""}`}
                  onClick={() => setPayMethod(key)}
                >
                  {label}
                </span>
              ))}
            </div>
            {error && <div className="form-error" style={{ marginTop: 12 }}>{error}</div>}
          </>
        )}

        {result && (
          <div style={{ marginTop: 30 }}>
            <div className={`result-card ${result.product.rarity}`}>
              <div className="result-emoji">{result.product.image}</div>
              <div style={{ margin: "10px 0" }}>
                <span className={`badge ${result.product.rarity}`}>
                  {RARITY_TEXT[result.product.rarity]}
                </span>
              </div>
              <div style={{ fontSize: 17, fontWeight: 700 }}>{result.product.name}</div>
              <div style={{ color: "#999", fontSize: 13, margin: "6px 0" }}>
                {result.product.character} · {result.product.series}
              </div>
              <div className="card-price" style={{ display: "inline-block" }}>
                价值 ¥{result.product.price}
              </div>
            </div>
          </div>
        )}
      </div>

      {user && (
        <div className="records">
          <h2 style={{ fontSize: 17, marginBottom: 6 }}>我的开盒记录</h2>
          {records.length === 0 ? (
            <div className="empty" style={{ padding: 24 }}>还没有开过盒哦～</div>
          ) : (
            records.map((r) => (
              <div key={r.id} className="record-row">
                <span className="record-emoji">{r.product?.image}</span>
                <div style={{ flex: 1 }}>
                  <div>
                    <span className={`badge ${r.product?.rarity}`}>{r.product?.rarity}</span>{" "}
                    {r.product?.name}
                  </div>
                  <div style={{ fontSize: 12, color: "#bbb" }}>
                    {new Date(r.created_at).toLocaleString()}
                  </div>
                </div>
                <span style={{ color: "#999", fontSize: 13 }}>¥{r.product?.price}</span>
              </div>
            ))
          )}
        </div>
      )}
    </div>
  );
}
