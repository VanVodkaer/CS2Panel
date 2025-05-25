import { useEffect, useState } from "react";
import { Row, Col, Typography, Space, Button, Table, Popconfirm, message, Checkbox } from "antd";
import { useNavigate } from "react-router-dom";
import api from "../../config/axiosConfig";
import {
  getTvConnectCommand,
  getGameConnectCommand,
  getGameConnectUrl,
  getTvConnectUrl,
} from "../../util/connectServer";
import "./index.less";

const Home = () => {
  const [containers, setContainers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [deletingName, setDeletingName] = useState(null); // 跟踪正在删除的容器名
  const [stoppingName, setStoppingName] = useState(null); // 跟踪正在停止的容器名
  const [selectedRowKeys, setSelectedRowKeys] = useState([]); // 跟踪选中的容器
  const navigate = useNavigate();

  // 拉取列表
  const fetchContainers = async () => {
    setLoading(true);
    try {
      const res = await api.get("/docker/container/list");
      setContainers(res.data.containers || []);
    } catch (err) {
      console.error(err);
      message.error("获取容器列表失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchContainers();
  }, []);

  // 启动 / 停止 / 删除
  const handleStart = async (name) => {
    try {
      await api.post("/docker/container/start", { name });
      message.success("容器已启动");
      fetchContainers();
    } catch {
      message.error("启动失败");
    }
  };
  const handleStop = async (name) => {
    setStoppingName(name); // 设置当前正在停止的容器名
    try {
      // 设置自定义timeout为30000ms，覆盖默认的10000ms
      await api.post("/docker/container/stop", { name }, { timeout: 30000 });
      message.success("容器已停止");
      setTimeout(() => {
        fetchContainers();
      }, 500);
    } catch {
      message.error("停止失败");
    } finally {
      setStoppingName(null); // 清除正在停止状态
    }
  };
  const handleRemove = async (name) => {
    setDeletingName(name); // 设置当前正在删除的容器名
    try {
      await api.delete("/docker/container/remove", { data: { name } }, { timeout: 40000 });
      message.success("容器已删除");
      setTimeout(() => {
        fetchContainers();
      }, 500);
    } catch {
      message.error("删除失败");
    } finally {
      setDeletingName(null); // 清除正在删除状态
    }
  };

  // 批量操作相关函数
  const handleSelectChange = (selectedKeys) => {
    setSelectedRowKeys(selectedKeys);
  };

  const handleBatchStart = () => {
    // TODO: 实现批量启动逻辑
    console.log("批量启动容器:", selectedRowKeys);
  };

  const handleBatchStop = () => {
    // TODO: 实现批量停止逻辑
    console.log("批量停止容器:", selectedRowKeys);
  };

  const handleBatchRemove = () => {
    // TODO: 实现批量删除逻辑
    console.log("批量删除容器:", selectedRowKeys);
  };

  const handleConnection = (name) => {
    getGameConnectUrl(name).then((cmd) => {
      window.open(cmd, "_blank");
    });
  };

  const handleCopyConnection = (name) => {
    getGameConnectCommand(name).then((cmd) => {
      navigator.clipboard.writeText(cmd).then(() => {
        message.success("复制成功");
      });
    });
  };

  const handleSpectate = (name) => {
    getTvConnectUrl(name).then((cmd) => {
      window.open(cmd, "_blank");
    });
  };

  const handleCopySpectate = (name) => {
    getTvConnectCommand(name).then((cmd) => {
      navigator.clipboard.writeText(cmd).then(() => {
        message.success("复制成功");
      });
    });
  };

  const columns = [
    {
      title: "容器名称",
      dataIndex: "Names",
      key: "name",
      render: (names) => {
        const full = names[0];
        return full.includes("-") ? full.substring(full.indexOf("-") + 1) : full;
      },
    },
    {
      title: "状态",
      dataIndex: "Status",
      key: "status",
    },
    {
      title: "操作",
      key: "action",
      render: (_, record) => {
        const full = record.Names[0];
        const trimmed = full.includes("-") ? full.substring(full.indexOf("-") + 1) : full;
        return (
          <Space size="small">
            <Button
              type="primary"
              size="small"
              disabled={record.State === "running"}
              onClick={() => handleStart(trimmed)}>
              启动
            </Button>
            <Button
              size="small"
              disabled={record.State !== "running" || stoppingName === trimmed} // 如果正在停止，则禁用
              onClick={() => handleStop(trimmed)}>
              {stoppingName === trimmed ? "停止中..." : "停止"}
            </Button>
            <Popconfirm title="确认删除此容器？" onConfirm={() => handleRemove(trimmed)}>
              <Button danger size="small" disabled={deletingName === trimmed}>
                {deletingName === trimmed ? "删除中..." : "删除"}
              </Button>
            </Popconfirm>
            <Button type="link" size="small" onClick={() => navigate(`/container/detail/${trimmed}`)}>
              详情
            </Button>
          </Space>
        );
      },
    },
    {
      title: "复制信息",
      key: "copy",
      render: (_, record) => {
        const full = record.Names[0];
        const trimmed = full.includes("-") ? full.substring(full.indexOf("-") + 1) : full;
        return (
          <Space direction="vertical" size="small">
            <Space size="small">
              <Button size="small" onClick={() => handleConnection(trimmed)}>
                连接服务器
              </Button>
              <Button size="small" onClick={() => handleCopyConnection(trimmed)}>
                复制连接指令
              </Button>
            </Space>
            <Space size="small">
              <Button size="small" onClick={() => handleSpectate(trimmed)}>
                观战服务器
              </Button>
              <Button size="small" onClick={() => handleCopySpectate(trimmed)}>
                复制观战指令
              </Button>
            </Space>
          </Space>
        );
      },
    },
  ];

  const headerButtons = (
    <Row justify="space-between" align="middle" className="dashboard-header">
      <Col>
        <Typography.Title level={2} style={{ margin: 0 }}>
          容器列表
        </Typography.Title>
      </Col>
      <Col>
        <Space size="middle">
          <Button
            type="primary"
            onClick={() => {
              // 先更新地图信息
              api.post("/info/map/update").then(() => {
                navigate("/container/create");
              });
            }}>
            创建新容器
          </Button>
          {selectedRowKeys.length > 0 && (
            <>
              <Button onClick={handleBatchStart}>批量启动</Button>
              <Button onClick={handleBatchStop}>批量停止</Button>
              <Button danger onClick={handleBatchRemove}>
                批量删除
              </Button>
            </>
          )}
        </Space>
      </Col>
    </Row>
  );

  return (
    <div className="dashboard-container">
      {headerButtons}

      <Table
        rowKey="Id"
        columns={columns}
        dataSource={containers}
        loading={loading}
        pagination={false}
        rowSelection={{
          selectedRowKeys: selectedRowKeys,
          onChange: handleSelectChange,
          columnWidth: 50,
        }}
      />
    </div>
  );
};

export default Home;
