import "./index.less";

import { Layout, Menu } from "antd";
import { DashboardOutlined, SettingOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";

const { Header, Sider, Content } = Layout;

const MainLayout = ({ children }) => {
  const navigate = useNavigate();

  const menuItems = [
    { key: "/", icon: <DashboardOutlined />, label: "服务器状态" },
    { key: "/server", icon: <SettingOutlined />, label: "服务器配置" },
  ];

  const handleMenuClick = (e) => {
    navigate(e.key);
  };

  return (
    <Layout className="mainlayout-layout">
      <Sider theme="dark" className="sider">
        <div className="logo">CS2Panel</div>
        <Menu theme="dark" mode="inline" items={menuItems} onClick={handleMenuClick} />
      </Sider>
      <Layout>
        <Header className="layout-header">{/* 顶部导航栏区域 */}</Header>
        <Content className="layout-content">{children}</Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;
