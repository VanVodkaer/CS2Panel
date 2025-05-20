import { Button, Checkbox, Form, Input, InputNumber, message, Select, Typography } from "antd";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import api from "../../config/axiosConfig";

const { Option } = Select;
const { Text } = Typography;

const CreateContainer = () => {
  const navigate = useNavigate();
  const [form] = Form.useForm(); // 创建表单实例

  const [mapOptions, setMapOptions] = useState([]);
  const [availableModes, setAvailableModes] = useState([]);
  const [gameModeDescription, setGameModeDescription] = useState("");
  const [gameTypeDescription, setGameTypeDescription] = useState("");

  // 获取地图列表
  const fetchMapList = async () => {
    try {
      const mapList = await api.get("/info/map/list");
      setMapOptions(mapList.data.maps);
    } catch (error) {
      console.error("获取地图列表时出错:", error);
    }
  };

  useEffect(() => {
    fetchMapList();
  }, []);

  useEffect(() => {
    handleGameModeChange("0");
    handleGameTypeChange("0");
    form.setFieldsValue({
      cs2_gamemode: "0",
      cs2_gametype: "0",
    }); // 主动设置初始值
  }, [form]);

  const onFinish = (values) => {
    // 将数字类型字段转换为字符串
    values.cs2_port = values.cs2_port?.toString();
    values.tv_port = values.tv_port?.toString();
    values.cs2_rcon_port = values.cs2_rcon_port?.toString();
    values.cs2_maxplayers = values.cs2_maxplayers?.toString();
    values.cs2_bot_quota = values.cs2_bot_quota?.toString();
    values.cs2_tv_delay = values.cs2_tv_delay?.toString();

    // 布尔值转为 "0" 或 "1"
    values.cs2_lan = values.cs2_lan ? "1" : "0";
    values.cs2_cheats = values.cs2_cheats ? "1" : "0";
    values.cs2_tv_enable = values.cs2_tv_enable ? "1" : "0";
    values.cs2_tv_autorecord = values.cs2_tv_autorecord ? "1" : "0";
    values.cs2_competitive_mode = values.cs2_competitive_mode ? "1" : "0";
    values.cs2_logging_enabled = values.cs2_logging_enabled ? "1" : "0";

    console.log("提交的表单数据:", values);

    // 发送请求
    api
      .post("/docker/container/create", values)
      .then((response) => {
        message.success(`容器创建成功: ${response.data.container_id}`);
        setTimeout(() => {
          navigate("/");
        }, 1000);
      })
      .catch((error) => {
        const errorMsg = error.response?.data?.message || error.message;
        message.error(`创建容器失败: ${errorMsg}`);
      });
  };

  const onFinishFailed = (error) => {
    console.error(error);

    const errorMessage = error.errorFields
      .map((field) => `${field.name.join(".")}: ${field.errors.join(" ")}`)
      .join("\n");

    message.error(`创建容器失败:\n${errorMessage}`);
  };

  // 地图变更处理
  const handleMapChange = (selectedMapInternalName) => {
    const selectedMap = mapOptions.find((map) => map.internal_name === selectedMapInternalName);
    if (selectedMap) {
      setAvailableModes(selectedMap.playable_modes);
    }
  };

  // 游戏模式说明更新
  const handleGameModeChange = (selectedMode) => {
    form.setFieldsValue({ cs2_gamemode: selectedMode }); // 显式更新表单字段值

    if (selectedMode === "0") {
      setGameModeDescription(
        "休闲模式：适合放松的游戏模式，死亡后玩家可以快速复活。你可以轻松地享受游戏，不用担心复活时间太长。"
      );
    } else if (selectedMode === "1") {
      setGameModeDescription(
        "竞技模式：玩家之间有较高的对抗性，死后需要等待复活或游戏结束。通常是更为严肃的游戏体验，适合高手玩家。"
      );
    } else if (selectedMode === "2") {
      setGameModeDescription(
        "死亡竞赛模式：玩家每次死亡后会立即复活，专注于击杀和持续战斗。适合喜欢快速复活并继续战斗的玩家。"
      );
    }
  };

  // 游戏类型说明更新
  const handleGameTypeChange = (selectedType) => {
    form.setFieldsValue({ cs2_gametype: selectedType }); // 显式更新表单字段值

    if (selectedType === "0") {
      setGameTypeDescription("普通类型：标准的游戏体验，没有特殊规则。适合大多数玩家，带有常规的游戏玩法和目标。");
    } else if (selectedType === "1") {
      setGameTypeDescription(
        "死亡竞赛类型：专注于快速复活和持续战斗，适合喜欢快速节奏的玩家。玩家每次死亡后立即复活，注重击杀和存活。"
      );
    }
  };

  const randomPort = () => Math.floor(Math.random() * (49151 - 1024 + 1)) + 1024;

  return (
    <Form
      className="create-container-form"
      name="basic"
      labelCol={{ span: 8 }}
      wrapperCol={{ span: 16 }}
      style={{ maxWidth: 600 }}
      form={form} // 绑定表单实例
      initialValues={{
        remember: true,
        cs2_servername: "",
        cs2_port: randomPort(),
        cs2_rcon_port: randomPort(),
        tv_port: randomPort(),
        cs2_maxplayers: 10,
        cs2_lan: false,
        cs2_bot_difficulty: "1",
        cs2_gamemode: "0",
        cs2_gametype: "0",
        cs2_startmap: "",
        cs2_mapgroup: "",
        cs2_pw: "",
        cs2_rconpw: "",
        cs2_cheats: false,
        cs2_tv_enable: false,
        cs2_tv_pw: "",
        cs2_tv_delay: 0,
        cs2_tv_autorecord: false,
        cs2_bot_quota: 0,
        cs2_competitive_mode: false,
        cs2_logging_enabled: true,
      }}
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      autoComplete="off">
      <Form.Item label="容器名称" name="name" rules={[{ required: true, message: "请输入容器名称!" }]}>
        <Input />
      </Form.Item>

      <Form.Item label="服务器名称" name="cs2_servername" rules={[{ required: false }]}>
        <Input />
      </Form.Item>

      <Form.Item label="服务器端口" name="cs2_port" rules={[{ required: false }]}>
        <InputNumber min={1024} max={49151} />
      </Form.Item>

      <Form.Item label="服务器密码" name="cs2_pw" rules={[{ required: false }]}>
        <Input />
      </Form.Item>

      <Form.Item label="RCON端口" name="cs2_rcon_port" rules={[{ required: false }]}>
        <InputNumber min={1024} max={49151} />
      </Form.Item>

      <Form.Item label="RCON密码" name="cs2_rconpw" rules={[{ required: false }]}>
        <Input />
      </Form.Item>

      <Form.Item label="TV端口" name="tv_port" rules={[{ required: false }]}>
        <InputNumber min={1024} max={49151} />
      </Form.Item>

      <Form.Item label="局域网模式" name="cs2_lan" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="最大玩家数" name="cs2_maxplayers" rules={[{ required: false }]}>
        <InputNumber min={1} max={64} />
      </Form.Item>

      <Form.Item label="开始地图" name="cs2_startmap" rules={[{ required: false }]}>
        <Select onChange={handleMapChange}>
          {mapOptions.map((map) => (
            <Option key={map.internal_name} value={map.internal_name}>
              {map.name}
            </Option>
          ))}
        </Select>
      </Form.Item>

      {availableModes.length > 0 && (
        <Form.Item label="可用模式" name="cs2_playable_modes">
          <Text>{availableModes.join(" / ")}</Text>
        </Form.Item>
      )}

      <Form.Item label="游戏模式" name="cs2_gamemode" rules={[{ required: false }]}>
        <div>
          <Select onChange={handleGameModeChange} value={form.getFieldValue("cs2_gamemode")}>
            <Option value="0">休闲</Option>
            <Option value="1">竞技</Option>
            <Option value="2">死亡竞赛</Option>
          </Select>
          <Text type="secondary">{gameModeDescription}</Text>
        </div>
      </Form.Item>

      <Form.Item label="游戏类型" name="cs2_gametype" rules={[{ required: false }]}>
        <div>
          <Select onChange={handleGameTypeChange} value={form.getFieldValue("cs2_gametype")}>
            <Option value="0">普通</Option>
            <Option value="1">死亡竞赛</Option>
          </Select>
          <Text type="secondary">{gameTypeDescription}</Text>
        </div>
      </Form.Item>

      <Form.Item label="允许作弊" name="cs2_cheats" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="启用 SourceTV" name="cs2_tv_enable" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="SourceTV密码" name="cs2_tv_pw" rules={[{ required: false }]}>
        <Input.Password />
      </Form.Item>

      <Form.Item label="SourceTV延迟 (秒)" name="cs2_tv_delay" rules={[{ required: false }]}>
        <InputNumber min={0} />
      </Form.Item>

      <Form.Item label="自动录制 SourceTV" name="cs2_tv_autorecord" valuePropName="checked">
        <Checkbox />
      </Form.Item>

      <Form.Item label="机器人配额" name="cs2_bot_quota" rules={[{ required: false }]}>
        <InputNumber min={0} />
      </Form.Item>

      <Form.Item label="机器人难度" name="cs2_bot_difficulty" rules={[{ required: false }]}>
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

      <Form.Item wrapperCol={{ offset: 8, span: 16 }}>
        <Button type="primary" htmlType="submit">
          创建容器
        </Button>
      </Form.Item>
    </Form>
  );
};

export default CreateContainer;
