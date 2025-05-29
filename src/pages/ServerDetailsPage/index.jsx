import "./index.less";
import React, { useState, useEffect, useCallback } from "react";
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
  Row,
  Col,
  Divider,
  Select,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";

const { TextArea } = Input;
const { Text } = Typography;
const { Option } = Select;

// 预设游戏模式组合及描述（与 CreateContainer 保持一致）
const modePresets = [
  { key: "casual", label: "休闲", type: 0, mode: 0, desc: "休闲模式：适合放松的游戏模式，死亡后玩家可以快速复活。" },
  {
    key: "competitive",
    label: "竞技",
    type: 0,
    mode: 1,
    desc: "竞技模式：玩家之间有较高的对抗性，死后需要等待复活或游戏结束。",
  },
  { key: "wingman", label: "决斗", type: 0, mode: 2, desc: "决斗模式：2v2小规模对抗。" },
  { key: "armsrace", label: "军备竞赛", type: 1, mode: 0, desc: "军备竞赛模式：不断更换武器，专注击杀。" },
  { key: "demolition", label: "拆除", type: 1, mode: 1, desc: "拆除模式：无固定买枪回合，靠击杀或拆弹获得武器。" },
  { key: "deathmatch", label: "死亡竞赛", type: 1, mode: 2, desc: "死亡竞赛模式：玩家死亡后立即复活，专注击杀。" },
  { key: "training", label: "训练", type: 2, mode: 0, desc: "训练模式：仅供练习，无记分。" },
  { key: "custom", label: "自定义", type: 3, mode: 0, desc: "自定义模式：可通过插件或脚本定义规则。" },
  { key: "cooperative", label: "合作", type: 4, mode: 0, desc: "合作模式：与 AI 合作完成关卡。" },
  { key: "skirmish", label: "小冲突", type: 5, mode: 0, desc: "小冲突模式：小规模对抗，为休闲娱乐。" },
];

function ServerDetailsPage() {
  const { name } = useParams();

  const [status, setStatus] = useState({});
  const [newMap, setNewMap] = useState("");
  const [maxRounds, setMaxRounds] = useState(0);
  const [freezeTime, setFreezeTime] = useState(0);
  const [startMoney, setStartMoney] = useState(0);
  const [matchTime, setMatchTime] = useState(0);
  const [playerId, setPlayerId] = useState("");
  const [playerMsg, setPlayerMsg] = useState("");
  const [customCmd, setCustomCmd] = useState("");
  const [customOutput, setCustomOutput] = useState("");
  const [manualCmd, setManualCmd] = useState("");
  const [manualOutput, setManualOutput] = useState("");
  const [loading, setLoading] = useState(false);

  const [mapOptions, setMapOptions] = useState([]);
  const [isWarmupPaused, setIsWarmupPaused] = useState(false);
  const [restartDelay, setRestartDelay] = useState(5);
  const [warmupPauseTime, setWarmupPauseTime] = useState();
  const [roundTime, setRoundTime] = useState(0);
  const [buyTime, setBuyTime] = useState(0);
  const [buyAnywhere, setBuyAnywhere] = useState(false);
  const [maxMoney, setMaxMoney] = useState(0);
  const [autoTeamBalance, setAutoTeamBalance] = useState(false);
  const [limitTeams, setLimitTeams] = useState(0);
  const [c4Timer, setC4Timer] = useState(0);

  // 游戏模式相关状态
  const [gameMode, setGameMode] = useState(0);
  const [gameType, setGameType] = useState(0);
  const [presetDesc, setPresetDesc] = useState("");
  const [selectedPreset, setSelectedPreset] = useState("");

  // 通用的配置值解析函数
  const parseConfigValue = (response, configName, valueType = "number") => {
    if (!response) return valueType === "boolean" ? false : 0;

    console.log(`${configName} 响应:`, response); // 调试日志

    // 创建更宽松的正则表达式匹配模式
    const patterns = [
      new RegExp(`${configName}\\s*=\\s*([\\d.]+)`, "i"), // 数字值 (含小数)
      new RegExp(`${configName}\\s*=\\s*(true|false)`, "i"), // 布尔值
      new RegExp(`${configName}\\s*=\\s*(\\S+)`, "i"), // 通用值
    ];

    for (const pattern of patterns) {
      const match = response.match(pattern);
      if (match) {
        const value = match[1];
        console.log(`${configName} 匹配到值:`, value);

        switch (valueType) {
          case "boolean":
            if (value === "true" || value === "false") {
              return value === "true";
            }
            return parseInt(value, 10) === 1;
          case "int":
            return parseInt(value, 10);
          case "float":
            return parseFloat(value);
          default:
            return isNaN(parseFloat(value)) ? 0 : parseFloat(value);
        }
      }
    }

    console.warn(`${configName} 未能解析响应:`, response);
    return valueType === "boolean" ? false : 0;
  };

  // 游戏模式处理函数
  const handlePresetChange = useCallback((key) => {
    const p = modePresets.find((item) => item.key === key);
    if (p) {
      setGameType(p.type);
      setGameMode(p.mode);
      setSelectedPreset(p.key);
      setPresetDesc(p.desc);
    }
  }, []);

  const handleModeValuesChange = useCallback(() => {
    const p = modePresets.find((item) => item.type === gameType && item.mode === gameMode);
    if (p) {
      setSelectedPreset(p.key);
      setPresetDesc(p.desc);
    } else {
      setSelectedPreset("");
      setPresetDesc(`自定义：game_type=${gameType}, game_mode=${gameMode}`);
    }
  }, [gameType, gameMode]);

  useEffect(() => {
    handleModeValuesChange();
  }, [gameType, gameMode, handleModeValuesChange]);

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
      setMaxRounds(parseConfigValue(res.data.response, "mp_maxrounds", "int"));
    } catch (error) {
      console.error("获取最大回合数失败:", error);
    }
  };

  const fetchTimeLimit = async () => {
    try {
      const res = await api.post("/rcon/game/config/timelimit", { name });
      setMatchTime(parseConfigValue(res.data.response, "mp_timelimit", "float"));
    } catch (error) {
      console.error("获取比赛时间失败:", error);
    }
  };

  const fetchRoundTime = async () => {
    try {
      const res = await api.post("/rcon/game/config/roundtime", { name });
      setRoundTime(parseConfigValue(res.data.response, "mp_roundtime", "float"));
    } catch (error) {
      console.error("获取回合时间失败:", error);
    }
  };

  const fetchBuyTime = async () => {
    try {
      const res = await api.post("/rcon/game/config/buytime", { name });
      setBuyTime(parseConfigValue(res.data.response, "mp_buytime", "float"));
    } catch (error) {
      console.error("获取购买时间失败:", error);
    }
  };

  const fetchFreezeTime = async () => {
    try {
      const res = await api.post("/rcon/game/config/freezetime", { name });
      setFreezeTime(parseConfigValue(res.data.response, "mp_freezetime", "int"));
    } catch (error) {
      console.error("获取冻结时间失败:", error);
    }
  };

  const fetchStartMoney = async () => {
    try {
      const res = await api.post("/rcon/game/config/startmoney", { name });
      setStartMoney(parseConfigValue(res.data.response, "mp_startmoney", "int"));
    } catch (error) {
      console.error("获取初始金钱失败:", error);
    }
  };

  const fetchMaxMoney = async () => {
    try {
      const res = await api.post("/rcon/game/config/maxmoney", { name });
      setMaxMoney(parseConfigValue(res.data.response, "mp_maxmoney", "int"));
    } catch (error) {
      console.error("获取最大金钱失败:", error);
    }
  };

  const fetchBuyAnywhere = async () => {
    try {
      const res = await api.post("/rcon/game/config/buyanywhere", { name });
      setBuyAnywhere(parseConfigValue(res.data.response, "mp_buy_anywhere", "boolean"));
    } catch (error) {
      console.error("获取购买位置限制失败:", error);
    }
  };

  const fetchAutoTeamBalance = async () => {
    try {
      const res = await api.post("/rcon/game/config/autoteambalance", { name });
      setAutoTeamBalance(parseConfigValue(res.data.response, "mp_autoteambalance", "boolean"));
    } catch (error) {
      console.error("获取自动队伍平衡失败:", error);
    }
  };

  const fetchLimitTeams = async () => {
    try {
      const res = await api.post("/rcon/game/config/limitteams", { name });
      setLimitTeams(parseConfigValue(res.data.response, "mp_limitteams", "int"));
    } catch (error) {
      console.error("获取队伍限制失败:", error);
    }
  };

  const fetchC4Timer = async () => {
    try {
      const res = await api.post("/rcon/game/config/c4timer", { name });
      setC4Timer(parseConfigValue(res.data.response, "mp_c4timer", "int"));
    } catch (error) {
      console.error("获取C4时长失败:", error);
    }
  };

  const fetchWarmupTime = async () => {
    try {
      const res = await api.post("/rcon/game/warm/time", { name });
      setWarmupPauseTime(parseConfigValue(res.data.response, "mp_warmuptime", "int"));
    } catch (error) {
      console.error("获取热身暂停时间失败:", error);
    }
  };

  const fetchIsWarmupPaused = async () => {
    try {
      const res = await api.post("/rcon/game/warm/pause", { name });
      setIsWarmupPaused(parseConfigValue(res.data.response, "mp_warmup_pausetimer", "boolean"));
    } catch (error) {
      console.error("获取热身暂停状态失败:", error);
    }
  };

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
  const changeMap = () =>
    withLoading(() => {
      if (!newMap) return message.error("请输入地图 internal_name");
      return api.post("/rcon/map/change", { name, map: newMap }).then(() => {
        message.success("地图切换成功");
        setNewMap("");
        fetchStatus();
      });
    });

  const scheduleRestart = () =>
    withLoading(() =>
      api
        .post("/rcon/game/restart", { name, value: restartDelay.toString() })
        .then(() => {
          message.success("游戏重启已调度");
        })
        .catch((error) => {
          message.error("调度游戏重启失败: " + error);
        })
    );

  const startWarmup = () =>
    withLoading(() =>
      api
        .post("/rcon/game/warm/start", { name })
        .then(() => {
          message.success("热身已开始");
        })
        .catch((error) => {
          message.error("开始热身失败: " + error);
        })
    );

  const endWarmup = () =>
    withLoading(() =>
      api
        .post("/rcon/game/warm/end", { name })
        .then(() => {
          message.success("热身已结束");
        })
        .catch((error) => {
          message.error("结束热身失败: " + error);
        })
    );

  const toggleWarmupPause = (val) =>
    withLoading(() =>
      api
        .post("/rcon/game/warm/pause", { name, value: val ? "1" : "0" })
        .then(() => {
          setIsWarmupPaused(val);
          message.success(val ? "热身已暂停" : "热身已取消暂停");
        })
        .catch((error) => {
          message.error("切换热身暂停状态失败: " + error);
        })
    );

  const setupWarmupPauseTime = () =>
    withLoading(() =>
      api
        .post("/rcon/game/warm/time", { name, value: warmupPauseTime.toString() })
        .then(() => {
          message.success("暂停时间设置成功");
        })
        .catch((error) => {
          message.error("设置暂停时间失败: " + error);
        })
    );

  // updateGameConfig 支持不同的API路径
  const updateGameConfig = (configType, value) =>
    withLoading(async () => {
      try {
        let response;
        switch (configType) {
          case "maxrounds":
            response = await api.post("/rcon/game/config/maxrounds", { name, value: value.toString() });
            break;
          case "timelimit":
            response = await api.post("/rcon/game/config/timelimit", { name, value: value.toString() });
            break;
          case "roundtime":
            response = await api.post("/rcon/game/config/roundtime", { name, value: value.toString() });
            break;
          case "freezetime":
            response = await api.post("/rcon/game/config/freezetime", { name, value: value.toString() });
            break;
          case "buytime":
            response = await api.post("/rcon/game/config/buytime", { name, value: value.toString() });
            break;
          case "buyanywhere":
            response = await api.post("/rcon/game/config/buyanywhere", { name, buy_anywhere: value ? "1" : "0" });
            break;
          case "startmoney":
            response = await api.post("/rcon/game/config/startmoney", { name, value: value.toString() });
            break;
          case "maxmoney":
            response = await api.post("/rcon/game/config/maxmoney", { name, value: value.toString() });
            break;
          case "autoteambalance":
            response = await api.post("/rcon/game/config/autoteambalance", { name, value: value ? "1" : "0" });
            break;
          case "limitteams":
            response = await api.post("/rcon/game/config/limitteams", { name, value: value.toString() });
            break;
          case "c4timer":
            response = await api.post("/rcon/game/config/c4timer", { name, value: value.toString() });
            break;
          case "gamemode":
            response = await api.post("/rcon/game/mode", {
              name,
              gamemode: gameMode.toString(),
              gametype: gameType.toString(),
            });
            break;
          default:
            throw new Error("未知的配置类型");
        }
        message.success("游戏规则设置成功");
        console.log("配置更新结果:", response.data);
      } catch (error) {
        message.error("设置失败: " + (error.response?.data?.message || error.message));
        console.error("配置更新失败:", error);
      }
    });

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
          {/* 比赛控制区域 */}
          <Form.Item label="比赛控制">
            <Row gutter={16} align="middle">
              <Col span={8}>
                <Button type="primary" onClick={scheduleRestart} block>
                  重启游戏
                </Button>
              </Col>
              <Col span={8}>
                <InputNumber
                  min={1}
                  value={restartDelay}
                  onChange={setRestartDelay}
                  placeholder="延迟时间"
                  style={{ width: "100%" }}
                  addonAfter="秒"
                />
              </Col>
            </Row>
          </Form.Item>

          <Divider />

          {/* 热身控制区域 */}
          <Form.Item label="热身控制">
            <Row gutter={16}>
              <Col span={12}>
                <Button type="primary" onClick={startWarmup} block>
                  开始热身
                </Button>
              </Col>
              <Col span={12}>
                <Button onClick={endWarmup} block>
                  结束热身
                </Button>
              </Col>
            </Row>
          </Form.Item>

          {/* 热身计时控制区域 */}
          <Form.Item label="热身计时控制">
            <Space direction="vertical" style={{ width: "100%" }} size="middle">
              {/* 暂停时间设置 */}
              <Row gutter={16} align="middle">
                <Col span={12}>
                  <InputNumber
                    min={0}
                    value={warmupPauseTime}
                    onChange={setWarmupPauseTime}
                    placeholder="暂停时间"
                    style={{ width: "100%" }}
                    addonAfter="秒"
                  />
                </Col>
                <Col span={12}>
                  <Button onClick={setupWarmupPauseTime} block>
                    设置热身时间
                  </Button>
                </Col>
              </Row>

              {/* 暂停开关 */}
              <Row align="middle">
                <Col span={16}>
                  <span style={{ fontSize: "14px" }}>暂停热身计时</span>
                </Col>
                <Col span={8} style={{ textAlign: "right" }}>
                  <Switch
                    checked={isWarmupPaused}
                    onChange={toggleWarmupPause}
                    checkedChildren="已暂停"
                    unCheckedChildren="计时中"
                  />
                </Col>
              </Row>
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
          {/* 游戏模式设置区域 */}
          <Form.Item label="游戏模式设置">
            <Row gutter={16} align="top">
              <Col span={16}>
                <Space.Compact style={{ width: "100%" }}>
                  <InputNumber
                    min={0}
                    max={5}
                    value={gameType}
                    onChange={setGameType}
                    placeholder="game_type"
                    style={{ width: "50%" }}
                  />
                  <InputNumber
                    min={0}
                    max={2}
                    value={gameMode}
                    onChange={setGameMode}
                    placeholder="game_mode"
                    style={{ width: "50%" }}
                  />
                </Space.Compact>
              </Col>
              <Col span={8}>
                <Select
                  placeholder="选择模式"
                  value={selectedPreset}
                  onChange={handlePresetChange}
                  style={{ width: "100%" }}>
                  {modePresets.map((p) => (
                    <Option key={p.key} value={p.key}>
                      {p.label}
                    </Option>
                  ))}
                </Select>
              </Col>
            </Row>
            <Row style={{ marginTop: 8 }}>
              <Col span={16}>
                <Text type="secondary" style={{ fontSize: "12px" }}>
                  {presetDesc}
                </Text>
              </Col>
              <Col span={8}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("gamemode")} block>
                  设置游戏模式
                </Button>
              </Col>
            </Row>
          </Form.Item>

          <Divider />

          {/* 回合设置区域 */}
          <Form.Item label="回合设置">
            <Row gutter={16} align="middle">
              <Col span={8}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>最大回合数</span>
                  <InputNumber
                    min={1}
                    value={maxRounds}
                    onChange={setMaxRounds}
                    style={{ width: "100%" }}
                    placeholder="30"
                  />
                </Space>
              </Col>
              <Col span={8}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>比赛时间限制(分钟)</span>
                  <InputNumber
                    min={0}
                    value={matchTime}
                    onChange={setMatchTime}
                    style={{ width: "100%" }}
                    placeholder="0表示无限制"
                  />
                </Space>
              </Col>
              <Col span={8}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>回合时间(分钟)</span>
                  <InputNumber
                    min={0}
                    max={60}
                    value={roundTime}
                    onChange={setRoundTime}
                    step={0.1}
                    style={{ width: "100%" }}
                    placeholder="1.92"
                  />
                </Space>
              </Col>
            </Row>
            <Row gutter={16} style={{ marginTop: 12 }}>
              <Col span={8}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("maxrounds", maxRounds)} block>
                  设置回合数
                </Button>
              </Col>
              <Col span={8}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("timelimit", matchTime)} block>
                  设置时间限制
                </Button>
              </Col>
              <Col span={8}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("roundtime", roundTime)} block>
                  设置回合时间
                </Button>
              </Col>
            </Row>
          </Form.Item>

          <Divider />

          {/* 时间设置区域 */}
          <Form.Item label="时间设置">
            <Row gutter={16} align="middle">
              <Col span={12}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>准备时间(秒)</span>
                  <InputNumber
                    min={0}
                    value={freezeTime}
                    onChange={setFreezeTime}
                    style={{ width: "100%" }}
                    placeholder="15"
                  />
                </Space>
              </Col>
              <Col span={12}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>装备购买时间(秒)</span>
                  <InputNumber
                    min={0}
                    value={buyTime}
                    onChange={setBuyTime}
                    style={{ width: "100%" }}
                    placeholder="20"
                  />
                </Space>
              </Col>
            </Row>
            <Row gutter={16} style={{ marginTop: 12 }}>
              <Col span={12}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("freezetime", freezeTime)} block>
                  设置准备时间
                </Button>
              </Col>
              <Col span={12}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("buytime", buyTime)} block>
                  设置购买时间
                </Button>
              </Col>
            </Row>
          </Form.Item>

          <Divider />

          {/* 经济设置区域 */}
          <Form.Item label="经济设置">
            <Row gutter={16} align="middle">
              <Col span={8}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>初始金钱</span>
                  <InputNumber
                    min={0}
                    value={startMoney}
                    onChange={setStartMoney}
                    style={{ width: "100%" }}
                    placeholder="800"
                  />
                </Space>
              </Col>
              <Col span={8}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>最大金钱</span>
                  <InputNumber
                    min={0}
                    value={maxMoney}
                    onChange={setMaxMoney}
                    style={{ width: "100%" }}
                    placeholder="16000"
                  />
                </Space>
              </Col>
              <Col span={8}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>任意位置购买</span>
                  <Switch
                    checked={buyAnywhere}
                    onChange={(val) => {
                      setBuyAnywhere(val);
                      updateGameConfig("buyanywhere", val);
                    }}
                    checkedChildren="开启"
                    unCheckedChildren="关闭"
                    style={{ width: "100%" }}
                  />
                </Space>
              </Col>
            </Row>
            <Row gutter={16} style={{ marginTop: 12 }}>
              <Col span={12}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("startmoney", startMoney)} block>
                  设置初始金钱
                </Button>
              </Col>
              <Col span={12}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("maxmoney", maxMoney)} block>
                  设置最大金钱
                </Button>
              </Col>
            </Row>
          </Form.Item>

          <Divider />

          {/* 队伍设置区域 */}
          <Form.Item label="队伍设置">
            <Row gutter={16} align="middle">
              <Col span={12}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>队伍人数最大差距</span>
                  <InputNumber
                    min={0}
                    value={limitTeams}
                    onChange={setLimitTeams}
                    style={{ width: "100%" }}
                    placeholder="0表示无限制"
                  />
                </Space>
              </Col>
              <Col span={12}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>自动平衡队伍</span>
                  <Switch
                    checked={autoTeamBalance}
                    onChange={(val) => {
                      setAutoTeamBalance(val);
                      updateGameConfig("autoteambalance", val);
                    }}
                    checkedChildren="开启"
                    unCheckedChildren="关闭"
                    style={{ width: "100%" }}
                  />
                </Space>
              </Col>
            </Row>
            <Row gutter={16} style={{ marginTop: 12 }}>
              <Col span={12}>
                <Button type="primary" size="small" onClick={() => updateGameConfig("limitteams", limitTeams)} block>
                  设置队伍限制
                </Button>
              </Col>
            </Row>
          </Form.Item>

          <Divider />

          {/* 炸弹设置区域 */}
          <Form.Item label="炸弹设置">
            <Row gutter={16} align="middle">
              <Col span={12}>
                <Space direction="vertical" style={{ width: "100%" }}>
                  <span style={{ fontSize: "12px", color: "#666" }}>C4爆炸时长(秒)</span>
                  <InputNumber
                    min={0}
                    value={c4Timer}
                    onChange={setC4Timer}
                    style={{ width: "100%" }}
                    placeholder="40"
                  />
                </Space>
              </Col>
              <Col span={12}>
                <Button
                  type="primary"
                  size="small"
                  onClick={() => updateGameConfig("c4timer", c4Timer)}
                  block
                  style={{ marginTop: 24 }}>
                  设置C4时长
                </Button>
              </Col>
            </Row>
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
        <Card title={`服务器状态：${name}`} className="server-card">
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

        <Card className="server-tabs">
          <Tabs
            defaultActiveKey="map"
            items={tabItems}
            onTabClick={(key) => {
              if (key === "map") {
                fetchMapList();
              } else if (key === "control") {
                fetchWarmupTime();
                fetchIsWarmupPaused();
              } else if (key === "rules") {
                fetchMaxRounds();
                fetchTimeLimit();
                fetchRoundTime();
                fetchBuyTime();
                fetchFreezeTime();
                fetchStartMoney();
                fetchMaxMoney();
                fetchBuyAnywhere();
                fetchAutoTeamBalance();
                fetchLimitTeams();
                fetchC4Timer();
              }
            }}
          />
        </Card>

        <Card title="命令输入 + 批量上传">
          <TextArea
            rows={3}
            placeholder="请输入RCON命令，支持多行批量执行"
            value={manualCmd}
            onChange={(e) => setManualCmd(e.target.value)}
            className="command-input"
          />
          <div className="form-item-inline">
            <Button type="primary" onClick={() => execCommand(manualCmd, setManualOutput)}>
              执行命令
            </Button>
            <Upload showUploadList={false} beforeUpload={handleUpload}>
              <Button icon={<UploadOutlined />} className="command-upload-btn">
                上传命令文件
              </Button>
            </Upload>
          </div>
          <TextArea rows={5} readOnly value={manualOutput} className="command-result" placeholder="执行结果" />
        </Card>
      </div>
    </Spin>
  );
}

export default ServerDetailsPage;
