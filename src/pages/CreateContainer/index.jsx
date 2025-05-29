import React, { useState, useEffect, useCallback } from "react";
import { Button, Checkbox, Form, Input, InputNumber, message, Select, Typography, Space } from "antd";
import { useNavigate } from "react-router-dom";
import api from "../../config/axiosConfig";

const { Option } = Select;
const { Text } = Typography;

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

const CreateContainer = () => {
  const [form] = Form.useForm();
  const navigate = useNavigate();

  const [mapOptions, setMapOptions] = useState([]);
  const [availableModes, setAvailableModes] = useState([]);
  const [presetDesc, setPresetDesc] = useState("");

  // 预设下拉框变化
  const handlePresetChange = useCallback(
    (key) => {
      const p = modePresets.find((item) => item.key === key);
      if (p) {
        form.setFieldsValue({
          cs2_gametype: p.type,
          cs2_gamemode: p.mode,
          preset: p.key,
        });
        setPresetDesc(p.desc);
      }
    },
    [form]
  );

  // 数值输入变化时同步下拉并更新描述
  const handleValuesChange = useCallback(
    (changed, all) => {
      const type = all.cs2_gametype;
      const mode = all.cs2_gamemode;
      const p = modePresets.find((item) => item.type === type && item.mode === mode);
      if (p) {
        form.setFieldsValue({ preset: p.key });
        setPresetDesc(p.desc);
      } else {
        form.setFieldsValue({ preset: undefined });
        setPresetDesc(`自定义：game_type=${type}, game_mode=${mode}`);
      }
    },
    [form]
  );

  // 获取地图列表
  const fetchMapList = async () => {
    try {
      const resp = await api.get("/info/map/list");
      setMapOptions(resp.data.maps);
    } catch (err) {
      console.error("获取地图列表时出错:", err);
    }
  };

  useEffect(() => {
    fetchMapList();
  }, []);

  // 初始化默认值
  useEffect(() => {
    const init = modePresets[0];
    form.setFieldsValue({
      cs2_gametype: init.type,
      cs2_gamemode: init.mode,
      preset: init.key,
    });
    setPresetDesc(init.desc);
  }, [form]);

  const onFinish = (values) => {
    // 将数值字段转为字符串
    [
      "cs2_port",
      "tv_port",
      "cs2_rcon_port",
      "cs2_maxplayers",
      "cs2_bot_quota",
      "cs2_tv_delay",
      "cs2_gamemode",
      "cs2_gametype",
    ].forEach((k) => {
      if (values[k] != null) values[k] = values[k].toString();
    });
    // 布尔转 "0"/"1"
    [
      "cs2_lan",
      "cs2_cheats",
      "cs2_tv_enable",
      "cs2_tv_autorecord",
      "cs2_competitive_mode",
      "cs2_logging_enabled",
      "steamappvalidate",
    ].forEach((k) => {
      values[k] = values[k] ? "1" : "0";
    });

    api
      .post("/docker/container/create", values)
      .then((res) => {
        message.success(`容器创建成功: ${res.data.container_id}`);
        setTimeout(() => navigate("/"), 1000);
      })
      .catch((err) => {
        const msg = err.response?.data?.message || err.message;
        message.error(`创建容器失败: ${msg}`);
      });
  };

  const onFinishFailed = ({ errorFields }) => {
    const msg = errorFields.map((f) => `${f.name.join(".")}: ${f.errors.join(" ")}`).join("\n");
    message.error(`创建容器失败:\n${msg}`);
  };

  const handleMapChange = (key) => {
    const m = mapOptions.find((m) => m.internal_name === key);
    setAvailableModes(m ? m.playable_modes : []);
  };

  const randomPort = () => Math.floor(Math.random() * (49151 - 1024 + 1)) + 1024;

  return (
    <Form
      form={form}
      name="create-cs2"
      labelCol={{ span: 8 }}
      wrapperCol={{ span: 16 }}
      style={{ maxWidth: 600 }}
      initialValues={{
        cs2_servername: "",
        cs2_port: randomPort(),
        cs2_rcon_port: randomPort(),
        tv_port: randomPort(),
        cs2_maxplayers: 10,
        cs2_lan: false,
        cs2_bot_difficulty: "1",
        cs2_cheats: false,
        cs2_tv_enable: false,
        cs2_tv_pw: "",
        cs2_tv_delay: 0,
        cs2_tv_autorecord: false,
        cs2_bot_quota: 0,
        cs2_competitive_mode: false,
        cs2_logging_enabled: true,
        steamappvalidate: false,
      }}
      onValuesChange={handleValuesChange}
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      autoComplete="off">
      <Form.Item label="容器名称" name="name" rules={[{ required: true, message: "请输入容器名称!" }]}>
        <Input />
      </Form.Item>

      <Form.Item label="服务器名称" name="cs2_servername">
        <Input />
      </Form.Item>

      <Form.Item label="服务器端口" name="cs2_port">
        <InputNumber min={1024} max={49151} />
      </Form.Item>

      <Form.Item label="服务器密码" name="cs2_pw">
        <Input />
      </Form.Item>

      <Form.Item label="RCON 端口" name="cs2_rcon_port">
        <InputNumber min={1024} max={49151} />
      </Form.Item>

      <Form.Item label="RCON 密码" name="cs2_rconpw">
        <Input />
      </Form.Item>

      <Form.Item label="TV 端口" name="tv_port">
        <InputNumber min={1024} max={49151} />
      </Form.Item>

      <Form.Item label="局域网模式" name="cs2_lan" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="最大玩家数" name="cs2_maxplayers">
        <InputNumber min={1} max={64} />
      </Form.Item>

      <Form.Item label="开始地图" name="cs2_startmap">
        <Select onChange={handleMapChange}>
          {mapOptions.map((m) => (
            <Option key={m.internal_name} value={m.internal_name}>
              {m.name}
            </Option>
          ))}
        </Select>
      </Form.Item>

      {availableModes.length > 0 && (
        <Form.Item label="可用模式">
          <Text>{availableModes.join(" / ")}</Text>
        </Form.Item>
      )}

      {/* 游戏模式/类型 统一设置组 */}
      <Form.Item label="游戏模式设置">
        <Space.Compact>
          <Form.Item name="cs2_gametype" noStyle>
            <InputNumber min={0} max={5} placeholder="game_type" />
          </Form.Item>
          <Form.Item name="cs2_gamemode" noStyle>
            <InputNumber min={0} max={2} placeholder="game_mode" />
          </Form.Item>
          <Form.Item noStyle>
            <Select placeholder="选择模式" onChange={handlePresetChange} style={{ width: 120 }}>
              {modePresets.map((p) => (
                <Option key={p.key} value={p.key}>
                  {p.label}
                </Option>
              ))}
            </Select>
          </Form.Item>
        </Space.Compact>
        <Text type="secondary" style={{ display: "block", marginTop: 8 }}>
          {presetDesc}
        </Text>
      </Form.Item>

      <Form.Item label="允许作弊" name="cs2_cheats" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="启用 SourceTV" name="cs2_tv_enable" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="SourceTV 密码" name="cs2_tv_pw">
        <Input.Password />
      </Form.Item>

      <Form.Item label="SourceTV 延迟 (秒)" name="cs2_tv_delay">
        <InputNumber min={0} />
      </Form.Item>

      <Form.Item label="自动录制 SourceTV" name="cs2_tv_autorecord" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="机器人配额" name="cs2_bot_quota">
        <InputNumber min={0} />
      </Form.Item>

      <Form.Item label="机器人难度" name="cs2_bot_difficulty">
        <Select>
          <Option value="0">简单</Option>
          <Option value="1">普通</Option>
          <Option value="2">困难</Option>
          <Option value="3">专家</Option>
        </Select>
      </Form.Item>

      <Form.Item label="比赛模式" name="cs2_competitive_mode" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="启用日志" name="cs2_logging_enabled" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="开启校验" name="steamappvalidate" valuePropName="checked">
        <Space>
          <Checkbox />
          <Typography.Text>开启 Steam 应用校验（验证游戏文件完整性）</Typography.Text>
        </Space>
      </Form.Item>

      {/* 开启 Steam 应用校验（验证游戏文件完整性） */}
      <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
        <Button type="primary" htmlType="submit">
          创建容器
        </Button>
      </Form.Item>
    </Form>
  );
};

export default CreateContainer;
