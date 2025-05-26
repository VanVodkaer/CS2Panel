import MainLayout from "./layout/MainLayout";
import { Routes, Route } from "react-router-dom";
import Home from "./pages/Dashboard";
import CreateContainer from "./pages/CreateContainer";
import ServerDetailsPage from "./pages/ServerDetailsPage";

const App = () => {
  return (
    <MainLayout>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/container/create" element={<CreateContainer />} />
        <Route path="/container/detail/:name" element={<ServerDetailsPage />} />
        {/* 其他页面路由 */}
      </Routes>
    </MainLayout>
  );
};

export default App;
