import { useState, useEffect } from "react";
import { Typography, Form, AutoComplete, Button, message } from "antd";
import api from "../../../config/axiosConfig";

function MapManagement({ name, status, fetchStatus, withLoading }) {
  const [newMap, setNewMap] = useState("");
  const [mapOptions, setMapOptions] = useState([]);

  useEffect(() => {
    fetchMapList();
  }, []);

  const fetchMapList = async () => {
    try {
      const mapList = await api.get("/info/map/list");
      setMapOptions(mapList.data.maps);
    } catch (error) {
      console.error("获取地图列表时出错:", error);
    }
  };

  const changeMap = () =>
    withLoading(() => {
      if (!newMap) return message.error("请输入地图 internal_name");
      return api.post("/rcon/map/change", { name, map: newMap }).then(() => {
        message.success("地图切换成功");
        setNewMap("");
        fetchStatus();
      });
    });

  const getMapName = () => status.spawngroups?.[0]?.path;

  return (
    <div className="map-management">
      <div className="current-map">
        <Typography.Text strong>当前地图：</Typography.Text> {getMapName()}
      </div>

      <Form layout="inline" className="map-form">
        <Form.Item label="开始地图">
          <AutoComplete
            className="map-input"
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
    </div>
  );
}

export default MapManagement;
