import { Card, Tabs } from "antd";
import MapManagement from "./MapManagement";
import GameControl from "./GameControl";
import GameRulesSettings from "./GameRulesSettings";
import PlayerManagement from "./PlayerManagement";

function ServerTabs({ name, status, statusjson, fetchStatus, withLoading, execCommand, fetchStatusJson }) {
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
      children: (
        <PlayerManagement
          execCommand={execCommand}
          fetchStatus={fetchStatus}
          fetchStatusJson={fetchStatusJson}
          name={name}
          statusjson={statusjson}
        />
      ),
    },
  ];

  return (
    <Card className="server-tabs">
      <Tabs defaultActiveKey="map" items={tabItems} />
    </Card>
  );
}

export default ServerTabs;
