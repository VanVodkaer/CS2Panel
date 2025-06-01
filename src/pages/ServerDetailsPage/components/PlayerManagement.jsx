import React, { useState, useEffect } from "react";
import { Card, Table, Button, Space, InputNumber, message, Tag, Modal, Typography } from "antd";
import { UserOutlined, RobotOutlined, ExclamationCircleOutlined } from "@ant-design/icons";
import api from "../../../config/axiosConfig";

const { Text } = Typography;
const { confirm } = Modal;

function PlayerManagement({ fetchStatusJson, name, statusjson }) {
  const [players, setPlayers] = useState([]);
  const [loading] = useState(false);
  const [banTime, setBanTime] = useState(30); // 默认封禁30分钟

  useEffect(() => {
    // 从 statusjson.server.clients 获取玩家信息
    if (statusjson?.server?.clients) {
      // 过滤掉无效的记录：
      // 1. steamid为"0"或0的记录
      // 2. steamid为"[I:0:0]"格式的记录（这些也是bot的bug）
      // 3. steamid64为"0"或0的记录
      const validPlayers = statusjson.server.clients.filter((player) => {
        const steamid = player.steamid || "";
        const steamid64 = player.steamid64 || "";

        // 过滤掉这些无效情况
        if (steamid === "0" || steamid === 0) return false;
        if (steamid === "[I:0:0]") return false;
        if (steamid64 === "0" || steamid64 === 0) return false;

        return true;
      });

      const playersWithUniqueId = validPlayers.map((player, index) => ({
        ...player,
        uniqueId: `${player.steamid64 || "unknown"}_${index}`, // 使用steamid64和索引组合创建唯一ID
      }));
      setPlayers(playersWithUniqueId);
    }
  }, [statusjson]);

  const kickPlayer = (playerName) => {
    confirm({
      title: "确认踢出玩家",
      icon: <ExclamationCircleOutlined />,
      content: `确定要踢出玩家 "${playerName}" 吗？`,
      onOk() {
        return api
          .post("/rcon/game/user/kick", {
            name,
            user: playerName,
          })
          .then(() => {
            message.success(`玩家 ${playerName} 已被踢出`);
            fetchStatusJson();
          })
          .catch((error) => {
            message.error("踢出玩家失败: " + (error.response?.data?.message || error.message));
          });
      },
    });
  };

  const banPlayer = (playerName, steamid64) => {
    confirm({
      title: "确认封禁玩家",
      icon: <ExclamationCircleOutlined />,
      content: (
        <div>
          <p>确定要封禁玩家 "{playerName}" 吗？</p>
          <p>
            <Text type="secondary">SteamID: {steamid64}</Text>
          </p>
          <p>
            封禁时长：
            <InputNumber
              min={1}
              max={43200} // 最大30天
              value={banTime}
              onChange={setBanTime}
              style={{ margin: "0 8px" }}
              addonAfter="分钟"
            />
          </p>
        </div>
      ),
      onOk() {
        return api
          .post("/rcon/game/user/banid", {
            name,
            time: banTime.toString(),
            id: steamid64,
          })
          .then(() => {
            message.success(`玩家 ${playerName} 已被封禁 ${banTime} 分钟`);
            fetchStatusJson();
          })
          .catch((error) => {
            message.error("封禁玩家失败: " + (error.response?.data?.message || error.message));
          });
      },
    });
  };

  const columns = [
    {
      title: "类型",
      dataIndex: "bot",
      key: "type",
      width: 60,
      render: (isBot) => (
        <Tag color={isBot ? "orange" : "blue"} icon={isBot ? <RobotOutlined /> : <UserOutlined />}>
          {isBot ? "BOT" : "玩家"}
        </Tag>
      ),
    },
    {
      title: "名称",
      dataIndex: "name",
      key: "name",
      ellipsis: true,
    },
    {
      title: "SteamID",
      dataIndex: "steamid",
      key: "steamid",
      width: 140,
      render: (steamid, record) => {
        // 不显示机器人的steamid
        if (record.bot || !steamid) {
          return <Text type="secondary">-</Text>;
        }
        return (
          <Text code style={{ fontSize: "12px" }}>
            {steamid}
          </Text>
        );
      },
    },
    {
      title: "SteamID64",
      dataIndex: "steamid64",
      key: "steamid64",
      width: 160,
      render: (steamid64, record) => {
        // 不显示机器人的steamid64
        if (record.bot || !steamid64) {
          return <Text type="secondary">-</Text>;
        }
        return (
          <Text code style={{ fontSize: "12px" }}>
            {steamid64}
          </Text>
        );
      },
    },
    {
      title: "操作",
      key: "actions",
      width: 160,
      render: (_, record) => {
        // 不对机器人显示操作按钮
        if (record.bot) {
          return <Text type="secondary">无操作</Text>;
        }

        return (
          <Space size="small">
            <Button type="primary" size="small" danger onClick={() => kickPlayer(record.name)}>
              踢出
            </Button>
            <Button
              type="default"
              size="small"
              danger
              disabled={!record.steamid64} // 如果没有steamid64则禁用封禁按钮
              onClick={() => banPlayer(record.name, record.steamid64)}>
              封禁
            </Button>
          </Space>
        );
      },
    },
  ];

  const humanPlayers = players.filter((p) => !p.bot);

  return (
    <div className="player-management">
      <Space direction="vertical" style={{ width: "100%" }} size="middle">
        {/* 玩家统计信息 */}
        <Card size="small">
          <Space size="large">
            <div>
              <UserOutlined style={{ color: "#1890ff", marginRight: 8 }} />
              <Text strong>人类玩家：{humanPlayers.length}</Text>
            </div>
            <div>
              <Text strong>总计：{players.length}</Text>
            </div>
            <Button type="primary" onClick={fetchStatusJson} loading={loading}>
              刷新玩家列表
            </Button>
          </Space>
        </Card>

        {/* 封禁时长设置 */}
        <Card title="封禁设置" size="small">
          <Space align="center">
            <Text>默认封禁时长：</Text>
            <InputNumber
              min={1}
              max={43200}
              value={banTime}
              onChange={setBanTime}
              addonAfter="分钟"
              style={{ width: 120 }}
            />
            <Text type="secondary">（1分钟 - 30天）</Text>
          </Space>
        </Card>

        {/* 玩家列表 */}
        <Card title="在线玩家列表" size="small">
          <Table
            columns={columns}
            dataSource={players}
            rowKey="uniqueId" // 使用唯一ID作为rowKey
            size="small"
            loading={loading}
            pagination={{
              pageSize: 10,
              showSizeChanger: true,
              showQuickJumper: true,
              showTotal: (total, range) => `第 ${range[0]}-${range[1]} 条，共 ${total} 条`,
            }}
            scroll={{ x: 800 }}
            locale={{
              emptyText: "暂无在线玩家",
            }}
          />
        </Card>

        {/* 快捷操作说明 */}
        <Card title="操作说明" size="small">
          <ul style={{ margin: 0, paddingLeft: 20 }}>
            <li>
              <Text strong>踢出</Text>：立即将玩家从服务器移除，玩家可以重新加入
            </li>
            <li>
              <Text strong>封禁</Text>：根据SteamID64封禁玩家指定时长，期间无法加入服务器
            </li>
            <li>
              <Text type="secondary">注意：机器人无法进行踢出和封禁操作</Text>
            </li>
          </ul>
        </Card>
      </Space>
    </div>
  );
}

export default PlayerManagement;
