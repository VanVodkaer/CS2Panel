import { Card, Tabs } from "antd";
import MapManagement from "./MapManagement";
import GameControl from "./GameControl";
import GameRulesSettings from "./GameRulesSettings";
import PlayerManagement from "./PlayerManagement";
import CustomCommand from "./CustomCommand";

function ServerTabs({ name, status, fetchStatus, withLoading, execCommand }) {
  const tabItems = [
    {
      key: "map",
      label: "地图管理",
      children: <MapManagement name={name} status={status} fetchStatus={fetchStatus} withLoading={withLoading} />,
    },
    {
      key: "control",
      label: "游戏控制",
      children: <GameControl name={name} withLoading={withLoading} />,
    },
    {
      key: "rules",
      label: "游戏规则设置",
      children: <GameRulesSettings name={name} withLoading={withLoading} />,
    },
    {
      key: "players",
      label: "玩家管理",
      children: <PlayerManagement execCommand={execCommand} fetchStatus={fetchStatus} />,
    },
    {
      key: "custom",
      label: "自定义命令",
      children: <CustomCommand execCommand={execCommand} />,
    },
  ];

  return (
    <Card className="server-tabs">
      <Tabs defaultActiveKey="map" items={tabItems} />
    </Card>
  );
}

export default ServerTabs;
