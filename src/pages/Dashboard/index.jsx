import { useEffect, useState } from "react";
import { Row, Col, Typography, Space, Button, Table, Popconfirm, message } from "antd";
import { useNavigate } from "react-router-dom";
import api from "../../config/axiosConfig";
import {
  // getTvConnectCommand,
  getGameConnectCommand,
  getGameConnectUrl,
  // getTvConnectUrl,
} from "../../util/connectServer";
import "./index.less";

const Home = () => {
  // 列表和加载状态
  const [containers, setContainers] = useState([]);
  const [loading, setLoading] = useState(false);

  // 单个操作状态
  const [deletingName, setDeletingName] = useState(null);
  const [stoppingName, setStoppingName] = useState(null);

  // 选中行（存储 trimmedName）
  const [selectedRowKeys, setSelectedRowKeys] = useState([]);

  // 批量操作加载状态
  const [batchStartLoading, setBatchStartLoading] = useState(false);
  const [batchStopLoading, setBatchStopLoading] = useState(false);
  const [batchRemoveLoading, setBatchRemoveLoading] = useState(false);

  const navigate = useNavigate();

  // 从完整名称提取 trimmedName
  const getTrimmedName = (fullName) =>
    fullName.includes("-") ? fullName.substring(fullName.indexOf("-") + 1) : fullName;

  // 拉取容器列表
  const fetchContainers = () => {
    setLoading(true);
    api
      .get("/docker/container/list")
      .then((res) => {
        setContainers(res.data.containers || []);
      })
      .catch((err) => {
        console.error(err);
        message.error("获取容器列表失败");
      })
      .finally(() => {
        setLoading(false);
      });
  };

  useEffect(() => {
    fetchContainers();
  }, []);

  // 单个启动
  const handleStart = (name) => {
    api
      .post("/docker/container/start", { name }, { timeout: 60000 })
      .then(() => {
        message.success("容器已启动");
        fetchContainers();
      })
      .catch((err) => {
        console.error(err);
        message.error("启动失败");
      });
  };

  // 单个停止
  const handleStop = (name) => {
    setStoppingName(name);
    api
      .post("/docker/container/stop", { name }, { timeout: 60000 })
      .then(() => {
        message.success("容器已停止");
        // 停止后稍作延迟再刷新列表
        setTimeout(fetchContainers, 500);
      })
      .catch((err) => {
        console.error(err);
        message.error("停止失败");
      })
      .finally(() => {
        setStoppingName(null);
      });
  };

  // 单个删除
  const handleRemove = (name) => {
    setDeletingName(name);
    api
      .post("/docker/container/remove", { name }, { timeout: 60000 })
      .then(() => {
        message.success("容器已删除");
        setTimeout(fetchContainers, 500);
      })
      .catch((err) => {
        console.error(err);
        message.error("删除失败");
      })
      .finally(() => {
        setDeletingName(null);
      });
  };

  // 批量启动
  const handleBatchStart = () => {
    if (selectedRowKeys.length === 0) return;
    setBatchStartLoading(true);
    api
      .post("/docker/container/start", { names: selectedRowKeys }, { timeout: 60000 })
      .then(() => {
        message.success("批量启动成功");
        setSelectedRowKeys([]);
        fetchContainers();
      })
      .catch((err) => {
        console.error(err);
        message.error("批量启动失败");
      })
      .finally(() => {
        setBatchStartLoading(false);
      });
  };

  // 批量停止
  const handleBatchStop = () => {
    if (selectedRowKeys.length === 0) return;
    setBatchStopLoading(true);
    api
      .post("/docker/container/stop", { names: selectedRowKeys }, { timeout: 60000 })
      .then(() => {
        message.success("批量停止成功");
        setSelectedRowKeys([]);
        setTimeout(fetchContainers, 500);
      })
      .catch((err) => {
        console.error(err);
        message.error("批量停止失败");
      })
      .finally(() => {
        setBatchStopLoading(false);
      });
  };

  // 批量删除
  const handleBatchRemove = () => {
    if (selectedRowKeys.length === 0) return;
    setBatchRemoveLoading(true);
    api
      .post("/docker/container/remove", { names: selectedRowKeys }, { timeout: 60000 })
      .then(() => {
        message.success("批量删除成功");
        setSelectedRowKeys([]);
        setTimeout(fetchContainers, 500);
      })
      .catch((err) => {
        console.error(err);
        message.error("批量删除失败");
      })
      .finally(() => {
        setBatchRemoveLoading(false);
      });
  };

  // 连接/复制指令
  const handleConnection = (name) => {
    getGameConnectUrl(name).then((url) => {
      window.open(url, "_blank");
    });
  };
  const handleCopyConnection = (name) => {
    getGameConnectCommand(name).then((cmd) => {
      navigator.clipboard.writeText(cmd).then(() => {
        message.success("复制连接指令成功");
      });
    });
  };
  // const handleSpectate = (name) => {
  //   getTvConnectUrl(name).then((url) => {
  //     window.open(url, "_blank");
  //   });
  // };
  // const handleCopySpectate = (name) => {
  //   getTvConnectCommand(name).then((cmd) => {
  //     navigator.clipboard.writeText(cmd).then(() => {
  //       message.success("复制观战指令成功");
  //     });
  //   });
  // };

  // 表格列定义
  const columns = [
    {
      title: "容器名称",
      dataIndex: "Names",
      key: "name",
      render: (names) => getTrimmedName(names[0]),
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
        const trimmed = getTrimmedName(record.Names[0]);
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
              disabled={record.State !== "running" || stoppingName === trimmed}
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
        const trimmed = getTrimmedName(record.Names[0]);
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
            {/* <Space size="small">
              <Button size="small" onClick={() => handleSpectate(trimmed)}>
                观战服务器
              </Button>
              <Button size="small" onClick={() => handleCopySpectate(trimmed)}>
                复制观战指令
              </Button>
            </Space> */}
          </Space>
        );
      },
    },
  ];

  // 头部操作按钮
  const headerButtons = (
    <Row justify="space-between" align="middle" className="dashboard-header">
      <Col>
        <Typography.Title level={2} style={{ margin: 0 }}>
          容器列表
        </Typography.Title>
      </Col>
      <Col>
        <Space size="middle">
          <Button type="primary" onClick={() => navigate("/container/create")}>
            创建新容器
          </Button>
          {selectedRowKeys.length > 0 && (
            <>
              <Button onClick={handleBatchStart} loading={batchStartLoading}>
                批量启动
              </Button>
              <Button onClick={handleBatchStop} loading={batchStopLoading}>
                批量停止
              </Button>
              <Button danger onClick={handleBatchRemove} loading={batchRemoveLoading}>
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
        // rowKey 使用 trimmedName
        rowKey={(record) => getTrimmedName(record.Names[0])}
        columns={columns}
        dataSource={containers}
        loading={loading}
        pagination={false}
        rowSelection={{
          selectedRowKeys,
          onChange: setSelectedRowKeys,
          columnWidth: 50,
        }}
      />
    </div>
  );
};

export default Home;
