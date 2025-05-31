import { useState, useEffect, useCallback } from "react";
import { Form, InputNumber, Button, Row, Col, Switch, Space, Divider, Select, Typography, message } from "antd";
import api from "../../../config/axiosConfig";

const { Text } = Typography;
const { Option } = Select;

// 预设游戏模式组合及描述
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

function GameRulesSettings({ name, withLoading }) {
  // 游戏模式相关状态
  const [gameMode, setGameMode] = useState(0);
  const [gameType, setGameType] = useState(0);
  const [presetDesc, setPresetDesc] = useState("");
  const [selectedPreset, setSelectedPreset] = useState("");

  // 游戏配置状态
  const [maxRounds, setMaxRounds] = useState(0);
  const [matchTime, setMatchTime] = useState(0);
  const [roundTime, setRoundTime] = useState(0);
  const [freezeTime, setFreezeTime] = useState(0);
  const [buyTime, setBuyTime] = useState(0);
  const [buyAnywhere, setBuyAnywhere] = useState(false);
  const [startMoney, setStartMoney] = useState(0);
  const [maxMoney, setMaxMoney] = useState(0);
  const [autoTeamBalance, setAutoTeamBalance] = useState(false);
  const [limitTeams, setLimitTeams] = useState(0);
  const [c4Timer, setC4Timer] = useState(0);

  useEffect(() => {
    fetchAllConfigs();
  }, []);

  const parseConfigValue = (response, configName, valueType = "number") => {
    if (!response) return valueType === "boolean" ? false : 0;

    const patterns = [
      new RegExp(`${configName}\\s*=\\s*([\\d.]+)`, "i"),
      new RegExp(`${configName}\\s*=\\s*(true|false)`, "i"),
      new RegExp(`${configName}\\s*=\\s*(\\S+)`, "i"),
    ];

    for (const pattern of patterns) {
      const match = response.match(pattern);
      if (match) {
        const value = match[1];
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
    return valueType === "boolean" ? false : 0;
  };

  const fetchAllConfigs = async () => {
    try {
      const configs = await Promise.allSettled([
        api.post("/rcon/game/config/maxrounds", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/timelimit", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/roundtime", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/freezetime", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/buytime", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/buyanywhere", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/startmoney", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/maxmoney", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/autoteambalance", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/limitteams", { name }, { timeout: 60000 }),
        api.post("/rcon/game/config/c4timer", { name }, { timeout: 60000 }),
      ]);

      if (configs[0].status === "fulfilled")
        setMaxRounds(parseConfigValue(configs[0].value.data.response, "mp_maxrounds", "int"));
      if (configs[1].status === "fulfilled")
        setMatchTime(parseConfigValue(configs[1].value.data.response, "mp_timelimit", "float"));
      if (configs[2].status === "fulfilled")
        setRoundTime(parseConfigValue(configs[2].value.data.response, "mp_roundtime", "float"));
      if (configs[3].status === "fulfilled")
        setFreezeTime(parseConfigValue(configs[3].value.data.response, "mp_freezetime", "int"));
      if (configs[4].status === "fulfilled")
        setBuyTime(parseConfigValue(configs[4].value.data.response, "mp_buytime", "float"));
      if (configs[5].status === "fulfilled")
        setBuyAnywhere(parseConfigValue(configs[5].value.data.response, "mp_buy_anywhere", "boolean"));
      if (configs[6].status === "fulfilled")
        setStartMoney(parseConfigValue(configs[6].value.data.response, "mp_startmoney", "int"));
      if (configs[7].status === "fulfilled")
        setMaxMoney(parseConfigValue(configs[7].value.data.response, "mp_maxmoney", "int"));
      if (configs[8].status === "fulfilled")
        setAutoTeamBalance(parseConfigValue(configs[8].value.data.response, "mp_autoteambalance", "boolean"));
      if (configs[9].status === "fulfilled")
        setLimitTeams(parseConfigValue(configs[9].value.data.response, "mp_limitteams", "int"));
      if (configs[10].status === "fulfilled")
        setC4Timer(parseConfigValue(configs[10].value.data.response, "mp_c4timer", "int"));
    } catch (error) {
      console.error("获取配置失败:", error);
    }
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

  return (
    <div className="game-rules-settings">
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
                <InputNumber min={0} value={buyTime} onChange={setBuyTime} style={{ width: "100%" }} placeholder="20" />
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
                <InputNumber min={0} value={c4Timer} onChange={setC4Timer} style={{ width: "100%" }} placeholder="40" />
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
    </div>
  );
}

export default GameRulesSettings;
