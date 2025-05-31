import { Card, Descriptions } from "antd";

function ServerStatusCard({ name, status, getMapName }) {
  return (
    <Card title={`服务器状态：${name}`} className="server-status-card">
      <Descriptions column={2}>
        <Descriptions.Item label="主机名">{status.hostname}</Descriptions.Item>
        <Descriptions.Item label="地图">{getMapName()}</Descriptions.Item>
        <Descriptions.Item label="版本">{status.version}</Descriptions.Item>
        <Descriptions.Item label="操作系统">{status.os}</Descriptions.Item>
        <Descriptions.Item label="客户端状态">{status.client_status}</Descriptions.Item>
        <Descriptions.Item label="当前状态">{status.current_state}</Descriptions.Item>
        <Descriptions.Item label="本地IP">{status.local_ip}</Descriptions.Item>
        <Descriptions.Item label="服务器IP">{status.public_ip}</Descriptions.Item>
        <Descriptions.Item label="人类玩家数">{status.player_summary?.humans || 0}</Descriptions.Item>
        <Descriptions.Item label="机器人数量">{status.player_summary?.bots || 0}</Descriptions.Item>
      </Descriptions>
    </Card>
  );
}

export default ServerStatusCard;
