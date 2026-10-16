import { Link } from "react-router-dom";
import StarRating from "./StarRating";

export default function ProductCard({ p, rating }) {
  return (
    <Link to={`/product/${p.id}`} className="card">
      <div className="card-img">{p.image}</div>
      <div className="card-body">
        <div className="card-name">
          <span className={`badge ${p.rarity}`}>{p.rarity}</span> {p.name}
        </div>
        <div className="card-meta">
          <span className="card-price">¥{p.price}</span>
          <span className="card-sales">已售 {p.sales}</span>
        </div>
        {rating ? <StarRating value={Math.round(rating)} /> : null}
      </div>
    </Link>
  );
}
