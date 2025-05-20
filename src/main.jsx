import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import App from "./App";

import "normalize.css"; // 引入 CSS 重置样式
import "antd/dist/reset.css"; // 引入 Ant Design 样式
import "@ant-design/v5-patch-for-react-19"; // 引入 Ant Design 5 的补丁

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <BrowserRouter>
    <App />
  </BrowserRouter>
);
