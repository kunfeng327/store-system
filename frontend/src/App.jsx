import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import Navbar from "./components/Navbar";
import Home from "./pages/Home";
import ProductDetail from "./pages/ProductDetail";
import BlindBox from "./pages/BlindBox";
import Login from "./pages/Login";
import Orders from "./pages/Orders";

export default function App() {
  return (
    <BrowserRouter>
      <Navbar />
      <div className="container">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/product/:id" element={<ProductDetail />} />
          <Route path="/box" element={<BlindBox />} />
          <Route path="/login" element={<Login />} />
          <Route path="/orders" element={<Orders />} />
          <Route
            path="*"
            element={
              <div className="empty">
                404 未找到 · <Link to="/">返回首页</Link>
              </div>
            }
          />
        </Routes>
      </div>
      <footer className="footer">二次元谷子商店 · 大创项目 Demo</footer>
    </BrowserRouter>
  );
}
