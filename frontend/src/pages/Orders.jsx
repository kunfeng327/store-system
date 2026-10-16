import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import api, { getUser } from "../api";

export default function Orders() {
  const user = getUser();
  const [orders, setOrders] = useState([]);

  useEffect(() => {
    if (user) api.get("/orders").then((r) => setOrders(r.data));
  }, []);

  if (!user)
    return (
      <div className="empty">
        请先 <Link to="/login" style={{ color: "var(--primary)" }}>登录</Link> 查看订单
      </div>
    );

  return (
    <div>
      <div className="page-title">我的订单</div>
      {orders.length === 0 ? (
        <div className="empty">还没有订单，去逛逛吧～</div>
      ) : (
        orders.map((o) => (
          <div key={o.id} className="order-card">
            <div className="order-head">
              <span>
                订单号 {o.id} ·{" "}
                <span className={`order-type ${o.type}`}>{o.type === "box" ? "盲盒" : "购买"}</span>
              </span>
              <span>{new Date(o.created_at).toLocaleString()}</span>
            </div>
            {o.items?.map((it) => (
              <div key={it.id} className="order-item-row">
                <span className="order-item-emoji">{it.image}</span>
                <div style={{ flex: 1 }}>
                  <div>
                    <span className={`badge ${it.rarity}`}>{it.rarity}</span> {it.name}
                  </div>
                  <div style={{ fontSize: 13, color: "#999" }}>
                    单价 ¥{it.price} × {it.quantity}
                  </div>
                </div>
                <div className="card-price">¥{it.price * it.quantity}</div>
              </div>
            ))}
            <div style={{ textAlign: "right", marginTop: 8, fontSize: 15 }}>
              合计：<span className="card-price">¥{o.total}</span>
            </div>
          </div>
        ))
      )}
    </div>
  );
}
