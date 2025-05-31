import React, { useState } from "react";
import { Space, Input, Button, message } from "antd";

function PlayerManagement({ execCommand, fetchStatus }) {
  const [playerId, setPlayerId] = useState("");
  const [playerMsg, setPlayerMsg] = useState("");

  const kickPlayer = () =>
    execCommand(`kick ${playerId}`, () => {
      setPlayerId("");
      fetchStatus();
    });

  const banPlayer = () => execCommand(`banid ${playerId}`, () => setPlayerId(""));

  const sendMessage = () => {
    if (!playerId || !playerMsg) return message.error("请输入玩家 ID 和消息内容");
    execCommand(`say ${playerId} ${playerMsg}`, () => setPlayerMsg(""));
  };

  return (
    <div className="player-management">
      <div className="player-actions">
        <Space>
          <Input placeholder="玩家 ID" value={playerId} onChange={(e) => setPlayerId(e.target.value)} />
          <Button type="primary" onClick={kickPlayer}>
            踢出
          </Button>
          <Button type="primary" onClick={banPlayer}>
            禁言
          </Button>
        </Space>
      </div>

      <div className="message-section">
        <Space>
          <Input placeholder="消息内容" value={playerMsg} onChange={(e) => setPlayerMsg(e.target.value)} />
          <Button type="primary" onClick={sendMessage}>
            发送消息
          </Button>
        </Space>
      </div>
    </div>
  );
}

export default PlayerManagement;
