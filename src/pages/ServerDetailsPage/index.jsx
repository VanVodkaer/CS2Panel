import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import api from "../../config/axiosConfig";
import {
  Tabs,
  Card,
  Descriptions,
  Form,
  Input,
  InputNumber,
  Button,
  message,
  Typography,
  Upload,
  Spin,
  AutoComplete,
  Space,
  Switch,
  Tooltip,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";

const { TextArea } = Input;

function ServerDetailsPage() {
  const { name } = useParams();

  const [status, setStatus] = useState({});
  const [newMap, setNewMap] = useState("");
  const [maxRounds, setMaxRounds] = useState(0);
  const [freezeTime, setFreezeTime] = useState(0);
  const [startMoney, setStartMoney] = useState(0);
  const [playerId, setPlayerId] = useState("");
  const [playerMsg, setPlayerMsg] = useState("");
  const [customCmd, setCustomCmd] = useState("");
  const [customOutput, setCustomOutput] = useState("");
  const [manualCmd, setManualCmd] = useState("");
  const [manualOutput, setManualOutput] = useState("");
  const [loading, setLoading] = useState(false);

  const [mapOptions, setMapOptions] = useState([]);
  const [isWarmupPaused, setIsWarmupPaused] = useState(false);
  const [isOfflineWarmupOn, setIsOfflineWarmupOn] = useState(false);

  const [restartDelay, setRestartDelay] = useState(0);
  const [warmupPauseTime, setWarmupPauseTime] = useState(0);

  // --- 数据获取函数 ---
  const fetchMapList = async () => {
    try {
      const mapList = await api.get("/info/map/list");
      setMapOptions(mapList.data.maps);
    } catch (error) {
      console.error("获取地图列表时出错:", error);
    }
  };

  const fetchMaxRounds = async () => {
    try {
      const res = await api.post("/rcon/game/config/maxrounds", { name });
      const match = res.data.response.match(/mp_maxrounds\s*=\s*(\d+)/);
      setMaxRounds(match ? parseInt(match[1], 10) : 0);
    } catch (error) {
      console.error("获取最大回合数失败:", error);
    }
  };

  const fetchIsWarmupPaused = async () => {
    try {
      const res = await api.post("/rcon/game/warm/pause", { name });
      const match = res.data.response.match(/mp_warmup_pausetimer\s*=\s*(\d+)/);
      setIsWarmupPaused(match ? parseInt(match[1], 10) === 1 : false);
    } catch (error) {
      console.error("获取热身暂停状态失败:", error);
    }
  };

  const fetchIsOfflineWarmupOn = async () => {
    try {
      const res = await api.post("/rcon/game/warm/offline", { name });
      const match = res.data.response.match(/mp_warmup_offline\s*=\s*(true|false)/i);

      console.log(match);
      setIsOfflineWarmupOn(match ? match[1].toLowerCase() === "true" : false);
    } catch (error) {
      console.error("获取离线热身状态失败:", error);
    }
  };

  useEffect(() => {
    fetchMapList();
    fetchMaxRounds();
    fetchIsWarmupPaused();
    fetchIsOfflineWarmupOn();
  }, [name]);

  useEffect(() => {
    fetchStatus();
    const timer = setInterval(fetchStatus, 30000);
    return () => clearInterval(timer);
  }, [name]);

  // --- 通用函数 ---
  const fetchStatus = async () => {
    try {
      const res = await api.get("/rcon/game/status", { params: { name } });
      setStatus(res.data.status || {});
    } catch (error) {
      message.error("获取服务器状态失败: " + error);
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

  // --- 操作函数 ---
  const scheduleRestart = () =>
    withLoading(() =>
      new Promise((resolve) => setTimeout(resolve, restartDelay * 1000))
        .then(() => api.post("/rcon/game/restart", { name }))
        .then(() => {
          message.success("游戏重启成功");
          fetchStatus();
        })
    );

  const startWarmup = () =>
    withLoading(() =>
      api.post("/rcon/game/warm/start", { name }).then(() => {
        message.success("热身已开始");
      })
    );

  const endWarmup = () =>
    withLoading(() =>
      api.post("/rcon/game/warm/end", { name }).then(() => {
        message.success("热身已结束");
      })
    );

  const toggleWarmupPause = (val) =>
    withLoading(() =>
      api.post("/rcon/game/warm/pause", { name, value: val ? 1 : 0 }).then(() => {
        setIsWarmupPaused(val);
        message.success(val ? "热身已暂停" : "热身已取消暂停");
      })
    );

  const setWarmupPauseTimer = () =>
    withLoading(() =>
      api.post("/rcon/config/warmuppausetimer", { name, value: warmupPauseTime }).then(() => {
        message.success("暂停时间设置成功");
      })
    );

  const toggleOfflineWarmupOn = (val) =>
    withLoading(() =>
      api.post("/rcon/game/warm/offline", { name, value: val ? 1 : 0 }).then(() => {
        setIsOfflineWarmupOn(val);
        message.success(val ? "离线热身已关闭" : "离线热身已开启");
      })
    );

  const changeMap = () =>
    withLoading(() => {
      if (!newMap) return message.error("请输入地图 internal_name");
      return api.post("/rcon/map/change", { name, map: newMap }).then(() => {
        message.success("地图切换成功");
        setNewMap("");
        fetchStatus();
      });
    });

  const updateGameConfig = (key, value) =>
    withLoading(() =>
      api.post(`/rcon/config/${key}`, { name, value }).then(() => {
        message.success("游戏规则设置成功");
      })
    );

  const execCommand = (command, outputSetter) => {
    if (!command) return message.error("请输入命令");
    return withLoading(() =>
      api.post("/rcon/exec", { name, command }).then((res) => {
        outputSetter(res.data || "无返回值");
      })
    );
  };

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

  const handleUpload = async (file) => {
    const text = await file.text();
    const commands = text.split(/\r?\n/).filter((l) => l.trim());
    if (!commands.length) return message.warning("文件为空");
    return execCommand(commands.join("\n"), setManualOutput);
  };

  const getMapName = () => status.spawngroups?.[0]?.path;

  const tabItems = [
    {
      key: "map",
      label: "地图管理",
      children: (
        <>
          <Typography.Text strong>当前地图：</Typography.Text> {getMapName()}
          <Form layout="inline" style={{ marginTop: 8 }}>
            <Form.Item label="开始地图">
              <AutoComplete
                style={{ width: 200 }}
                options={mapOptions.map((m) => ({ value: m.internal_name, label: m.name }))}
                placeholder="输入地图 internal_name"
                value={newMap}
                onChange={setNewMap}
                filterOption={(input, option) => option.value.includes(input) || option.label.includes(input)}
              />
            </Form.Item>
            <Form.Item>
              <Button type="primary" onClick={changeMap}>
                切换地图
              </Button>
            </Form.Item>
            {newMap && (
              <Form.Item label="可玩模式">
                <Typography.Text>
                  {(mapOptions.find((m) => m.internal_name === newMap)?.playable_modes || []).join("、") || "无"}
                </Typography.Text>
              </Form.Item>
            )}
          </Form>
        </>
      ),
    },
    {
      key: "control",
      label: "游戏控制",
      children: (
        <Form layout="vertical">
          <Form.Item label="比赛控制">
            <Space>
              <Button type="primary" onClick={scheduleRestart}>
                重启游戏
              </Button>
              <InputNumber min={0} value={restartDelay} onChange={setRestartDelay} placeholder="延迟(秒)" />
            </Space>
          </Form.Item>
          <Form.Item label="热身控制">
            <Space align="center">
              <Button type="primary" onClick={startWarmup}>
                开始热身
              </Button>
              <Button type="primary" onClick={endWarmup}>
                结束热身
              </Button>
              <InputNumber min={0} value={warmupPauseTime} onChange={setWarmupPauseTime} placeholder="暂停时间(秒)" />
              <Button size="small" onClick={setWarmupPauseTimer}>
                设置
              </Button>
              <Tooltip title="设置私人游戏 / 离线使用机器人时是否开启热身">
                <Switch checked={isOfflineWarmupOn} onChange={toggleOfflineWarmupOn} /> 离线热身
              </Tooltip>
              <Switch checked={isWarmupPaused} onChange={toggleWarmupPause} /> 暂停热身
            </Space>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: "rules",
      label: "游戏规则设置",
      children: (
        <Form layout="vertical">
          <Form.Item label="最大回合数">
            <InputNumber min={1} value={maxRounds} onChange={setMaxRounds} style={{ width: 140, marginRight: 8 }} />
            <Button type="primary" size="small" onClick={() => updateGameConfig("maxrounds", maxRounds)}>
              设置
            </Button>
          </Form.Item>
          <Form.Item label="冻结时间（秒）">
            <InputNumber min={0} value={freezeTime} onChange={setFreezeTime} style={{ width: 140, marginRight: 8 }} />
            <Button type="primary" size="small" onClick={() => updateGameConfig("freezetime", freezeTime)}>
              设置
            </Button>
          </Form.Item>
          <Form.Item label="初始金钱">
            <InputNumber min={0} value={startMoney} onChange={setStartMoney} style={{ width: 140, marginRight: 8 }} />
            <Button type="primary" size="small" onClick={() => updateGameConfig("startmoney", startMoney)}>
              设置
            </Button>
          </Form.Item>
        </Form>
      ),
    },
    {
      key: "players",
      label: "玩家管理",
      children: (
        <div>
          <Space>
            <Input placeholder="玩家 ID" value={playerId} onChange={(e) => setPlayerId(e.target.value)} />
            <Button type="primary" onClick={kickPlayer}>
              踢出
            </Button>
            <Button type="primary" onClick={banPlayer}>
              禁言
            </Button>
          </Space>
          <div style={{ marginTop: 16 }}>
            <Space>
              <Input placeholder="消息内容" value={playerMsg} onChange={(e) => setPlayerMsg(e.target.value)} />
              <Button type="primary" onClick={sendMessage}>
                发送消息
              </Button>
            </Space>
          </div>
        </div>
      ),
    },
    {
      key: "custom",
      label: "自定义命令",
      children: (
        <div>
          <TextArea
            rows={4}
            placeholder="输入命令，每行一条"
            value={customCmd}
            onChange={(e) => setCustomCmd(e.target.value)}
          />
          <Button type="primary" style={{ marginTop: 8 }} onClick={() => execCommand(customCmd, setCustomOutput)}>
            执行
          </Button>
          <TextArea rows={4} readOnly value={customOutput} style={{ marginTop: 8 }} placeholder="执行结果" />
        </div>
      ),
    },
  ];

  return (
    <Spin spinning={loading}>
      <div>
        <Card title={`服务器状态：${name}`} style={{ marginBottom: 16 }}>
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

        <Card style={{ marginBottom: 16 }}>
          <Tabs defaultActiveKey="map" items={tabItems} />
        </Card>

        <Card title="命令输入 + 批量上传">
          <TextArea
            rows={3}
            placeholder="请输入RCON命令，支持多行批量执行"
            value={manualCmd}
            onChange={(e) => setManualCmd(e.target.value)}
          />
          <div style={{ marginTop: 8 }}>
            <Button type="primary" onClick={() => execCommand(manualCmd, setManualOutput)}>
              执行命令
            </Button>
            <Upload showUploadList={false} beforeUpload={handleUpload}>
              <Button icon={<UploadOutlined />} style={{ marginLeft: 8 }}>
                上传命令文件
              </Button>
            </Upload>
          </div>
          <TextArea rows={5} readOnly value={manualOutput} style={{ marginTop: 8 }} placeholder="执行结果" />
        </Card>
      </div>
    </Spin>
  );
}

export default ServerDetailsPage;
