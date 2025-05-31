import { useState, useEffect } from "react";
import { Form, Button, InputNumber, Row, Col, Switch, Space, Divider, message } from "antd";
import api from "../../../config/axiosConfig";

function GameControl({ name, withLoading }) {
  const [isWarmupPaused, setIsWarmupPaused] = useState(false);
  const [restartDelay, setRestartDelay] = useState(5);
  const [warmupPauseTime, setWarmupPauseTime] = useState();

  useEffect(() => {
    fetchWarmupTime();
    fetchIsWarmupPaused();
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

  return (
    <div className="game-control">
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
    </div>
  );
}

export default GameControl;
