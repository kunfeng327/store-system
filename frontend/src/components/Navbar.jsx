import { Link, useLocation, useNavigate } from "react-router-dom";
import { getUser, logout } from "../api";

export default function Navbar() {
  const user = getUser();
  const loc = useLocation();
  const nav = useNavigate();

  return (
    <div className="navbar">
      <div className="navbar-inner">
        <Link to="/" className="logo">🎁 谷子小铺</Link>
        <Link to="/" className={`nav-link ${loc.pathname === "/" ? "active" : ""}`}>商城</Link>
        <Link to="/box" className={`nav-link ${loc.pathname === "/box" ? "active" : ""}`}>魔法盲盒</Link>
        {user && <Link to="/orders" className={`nav-link ${loc.pathname === "/orders" ? "active" : ""}`}>我的订单</Link>}
        <div className="nav-right">
          {user ? (
            <>
              <span className="coins">🪙 {user.coins}</span>
              <span>{user.username}</span>
              <span className="nav-link" onClick={() => { logout(); nav("/login"); }}>退出</span>
            </>
          ) : (
            <Link to="/login" className="nav-link">登录 / 注册</Link>
          )}
        </div>
      </div>
    </div>
  );
}
