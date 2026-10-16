import { useEffect, useState } from "react";
import api from "../api";
import ProductCard from "../components/ProductCard";

const CATEGORIES = ["全部", "吧唧", "立牌", "挂件", "色纸"];

export default function Home() {
  const [cat, setCat] = useState("全部");
  const [kw, setKw] = useState("");
  const [list, setList] = useState([]);

  useEffect(() => {
    const params = { size: 50 };
    if (cat !== "全部") params.category = cat;
    if (kw.trim()) params.keyword = kw.trim();
    api.get("/products", { params }).then((r) => setList(r.data.list));
  }, [cat, kw]);

  return (
    <div>
      <div className="tabs">
        {CATEGORIES.map((c) => (
          <span key={c} className={`tab ${cat === c ? "active" : ""}`} onClick={() => setCat(c)}>
            {c}
          </span>
        ))}
        <input
          className="input"
          style={{ width: 220, margin: 0 }}
          placeholder="搜索角色 / 作品…"
          value={kw}
          onChange={(e) => setKw(e.target.value)}
        />
      </div>
      {list.length === 0 ? (
        <div className="empty">暂无商品</div>
      ) : (
        <div className="grid">
          {list.map((p) => (
            <ProductCard key={p.id} p={p} />
          ))}
        </div>
      )}
    </div>
  );
}
