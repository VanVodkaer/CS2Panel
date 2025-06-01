import "./index.less";
import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import api from "../../config/axiosConfig";
import { Card, Spin, message } from "antd";
import ServerStatusCard from "./components/ServerStatusCard";
import ServerTabs from "./components/ServerTabs";
import CommandInput from "./components/CommandInput";

function ServerDetailsPage() {
  const { name } = useParams();
  const [status, setStatus] = useState({});
  const [statusjson, setStatusJson] = useState({});
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchStatus();
    fetchStatusJson();
    const timer = setInterval(fetchStatus, 10000);
    return () => clearInterval(timer);
  }, [name]);

  const fetchStatus = async () => {
    try {
      const res = await api.get("/rcon/game/status", { params: { name } });
      setStatus(res.data.status || {});
    } catch (error) {
      message.error("获取服务器status失败,服务器可能正在启动中/未启动，请检查并稍后再试: " + error);
    }
  };

  const fetchStatusJson = async () => {
    try {
      const res = await api.get("/rcon/game/statusjson", { params: { name } });
      setStatusJson(res.data.status || {});
    } catch (error) {
      message.error("获取服务器status_json失败,服务器可能正在启动中/未启动，请检查并稍后再试: " + error);
    }
  };

  const withLoading = async (callback) => {
    setLoading(true);
    try {
      await callback();
    } finally {
      setLoading(false);
    }
  };

  const execCommand = (command, outputSetter) => {
    if (!command) {
      return message.error("请输入命令");
    }
    const commands = command
      .split("\n")
      .map((cmd) => cmd.trim())
      .filter((cmd) => cmd);

    return withLoading(() =>
      api.post("/rcon/exec", { name, cmds: commands }, { timeout: 60000 }).then((res) => {
        console.log("执行命令结果:", res.data);
        const responses = Array.isArray(res.data.responses) ? res.data.responses.join("\n") : res.data.responses;
        outputSetter(responses || "无返回值");
      })
    );
  };

  const getMapName = () => status.spawngroups?.[0]?.path;

  return (
    <div className="server-details-page">
      <Spin spinning={loading}>
        <ServerStatusCard name={name} status={status} getMapName={getMapName} />

        <ServerTabs
          name={name}
          status={status}
          statusjson={statusjson}
          fetchStatus={fetchStatus}
          fetchStatusJson={fetchStatusJson}
          withLoading={withLoading}
          execCommand={execCommand}
        />

        <CommandInput execCommand={execCommand} />
      </Spin>
    </div>
  );
}

export default ServerDetailsPage;
