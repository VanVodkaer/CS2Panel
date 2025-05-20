import MainLayout from "./layout/MainLayout";
import { Routes, Route } from "react-router-dom";
import Home from "./pages/Dashboard";
import CreateContainer from "./pages/CreateContainer";

const App = () => {
  return (
    <MainLayout>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/container/create" element={<CreateContainer />} />
        {/* 其他页面路由 */}
      </Routes>
    </MainLayout>
  );
};

export default App;
