import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import api, { getUser } from "../api";
import StarRating from "../components/StarRating";

export default function ProductDetail() {
  const { id } = useParams();
  const user = getUser();
  const [data, setData] = useState(null);
  const [reviews, setReviews] = useState([]);
  const [reviewTotal, setReviewTotal] = useState(0);
  const [ratingFilter, setRatingFilter] = useState("");
  const [qty, setQty] = useState(1);
  const [msg, setMsg] = useState("");
  // 发布评价表单
  const [myRating, setMyRating] = useState(5);
  const [myContent, setMyContent] = useState("");
  const [myFollow, setMyFollow] = useState("");

  const loadProduct = () =>
    api.get(`/products/${id}`).then((r) => setData(r.data));
  const loadReviews = () =>
    api
      .get(`/products/${id}/reviews`, {
        params: ratingFilter ? { rating: ratingFilter } : {},
      })
      .then((r) => {
        setReviews(r.data.list);
        setReviewTotal(r.data.total);
      });

  useEffect(() => {
    loadProduct();
  }, [id]);
  useEffect(() => {
    loadReviews();
  }, [id, ratingFilter]);

  const buy = async () => {
    try {
      const r = await api.post("/orders", { product_id: Number(id), quantity: qty });
      setMsg(r.data.message + `，花费 ¥${r.data.total}`);
      loadProduct();
      const me = await api.get("/me");
      localStorage.setItem("user", JSON.stringify(me.data));
    } catch (e) {
      setMsg(e.response?.data?.error || "购买失败");
    }
  };

  const submitReview = async () => {
    try {
      await api.post(`/products/${id}/reviews`, {
        rating: myRating,
        content: myContent,
        follow_up: myFollow,
      });
      setMyContent("");
      setMyFollow("");
      setMsg("评价发布成功！");
      loadProduct();
      loadReviews();
    } catch (e) {
      setMsg(e.response?.data?.error || "发布失败");
    }
  };

  if (!data) return <div className="empty">加载中…</div>;
  const p = data.product;
  const dist = data.rating_dist || {};

  return (
    <div>
      <div className="detail">
        <div className="detail-img">{p.image}</div>
        <div className="detail-info">
          <h1>
            <span className={`badge ${p.rarity}`}>{p.rarity}</span> {p.name}
          </h1>
          <div className="meta-row">
            角色：{p.character} · 作品：{p.series} · 类型：{p.category}
          </div>
          <div className="meta-row">
            <StarRating value={Math.round(data.rating_avg)} /> {data.rating_avg.toFixed(1)} 分 ·{" "}
            {data.rating_count} 条评价
          </div>
          <div className="detail-price">
            <small>¥</small>
            {p.price}
          </div>
          <div className="meta-row">
            库存 {p.stock} · 已售 {p.sales}
          </div>
          <div style={{ margin: "20px 0" }}>
            数量：
            <button className="tab" onClick={() => setQty(Math.max(1, qty - 1))}>-</button>{" "}
            {qty}{" "}
            <button className="tab" onClick={() => setQty(qty + 1)}>+</button>
          </div>
          <button className="btn btn-outline">加入购物车</button>
          <button className="btn btn-primary" onClick={buy} disabled={!user}>
            立即购买
          </button>
          {!user && <div className="form-error">请先登录后购买</div>}
          {msg && <div className="msg">{msg}</div>}
        </div>
      </div>

      {/* 评价区（仿淘宝） */}
      <div className="review-section">
        <div className="review-header">
          <h2>商品评价（{data.rating_count}）</h2>
          <div className="avg-score">
            <span className="avg-num">{data.rating_avg.toFixed(1)}</span>
            <div>
              <StarRating value={Math.round(data.rating_avg)} />
              <div style={{ fontSize: 12, color: "#999" }}>综合评分</div>
            </div>
          </div>
          <div style={{ flex: 1 }}>
            {[5, 4, 3, 2, 1].map((s) => (
              <div key={s} className="dist-bar">
                <span style={{ width: 30 }}>{s}星</span>
                <div className="dist-track">
                  <div
                    className="dist-fill"
                    style={{
                      width: `${
                        data.rating_count ? ((dist[s] || 0) / data.rating_count) * 100 : 0
                      }%`,
                    }}
                  />
                </div>
                <span style={{ width: 24 }}>{dist[s] || 0}</span>
              </div>
            ))}
          </div>
        </div>

        {/* 星级筛选 */}
        <div className="tabs" style={{ margin: "14px 0" }}>
          <span className={`tab ${ratingFilter === "" ? "active" : ""}`} onClick={() => setRatingFilter("")}>
            全部（{data.rating_count}）
          </span>
          {[5, 4, 3, 2, 1].map((s) => (
            <span key={s} className={`tab ${ratingFilter === String(s) ? "active" : ""}`} onClick={() => setRatingFilter(String(s))}>
              {s}星（{dist[s] || 0}）
            </span>
          ))}
        </div>

        {/* 发布评价 */}
        {user ? (
          <div style={{ padding: "16px 0", background: "#fafafa", borderRadius: 10, marginBottom: 10 }}>
            <div style={{ marginBottom: 10 }}>
              <span className="star-picker stars">
                {[1, 2, 3, 4, 5].map((i) => (
                  <span key={i} className={i <= myRating ? "" : "off"} onClick={() => setMyRating(i)}>
                    ★
                  </span>
                ))}
              </span>
            </div>
            <textarea
              className="input"
              rows={2}
              placeholder="写写你的吃谷感受吧～"
              value={myContent}
              onChange={(e) => setMyContent(e.target.value)}
            />
            <input
              className="input"
              placeholder="追评（可选）：使用一段时间后的补充评价"
              value={myFollow}
              onChange={(e) => setMyFollow(e.target.value)}
            />
            <button className="btn btn-primary" style={{ padding: "8px 28px", fontSize: 14 }} onClick={submitReview}>
              发布评价
            </button>
          </div>
        ) : (
          <div className="form-error" style={{ padding: "8px 0" }}>登录后可以发布评价</div>
        )}

        {/* 评价列表 */}
        {reviews.length === 0 ? (
          <div className="empty" style={{ padding: 30 }}>暂无评价</div>
        ) : (
          reviews.map((r) => (
            <div key={r.id} className="review-item">
              <div className="review-user">
                <div className="avatar">匿</div>
                <div>
                  <div>{r.username}</div>
                  <StarRating value={r.rating} />
                </div>
              </div>
              <div className="review-content">{r.content}</div>
              {r.has_follow && (
                <div className="followup">
                  <span className="followup-tag">[追评]</span>
                  {r.follow_up}
                </div>
              )}
              <div className="review-time">{new Date(r.created_at).toLocaleString()}</div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
