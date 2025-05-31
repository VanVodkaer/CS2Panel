import React, { useState } from "react";
import { Card, Input, Button, Upload, message } from "antd";
import { UploadOutlined } from "@ant-design/icons";

const { TextArea } = Input;

function CommandInput({ execCommand }) {
  const [manualCmd, setManualCmd] = useState("");
  const [manualOutput, setManualOutput] = useState("");

  const handleUpload = async (file) => {
    const text = await file.text();
    const commands = text.split(/\r?\n/).filter((l) => l.trim());
    if (!commands.length) return message.warning("文件为空");
    return execCommand(commands.join("\n"), setManualOutput);
  };

  return (
    <Card title="命令输入 + 批量上传" className="command-input-card">
      <TextArea
        rows={3}
        placeholder="请输入RCON命令，支持多行批量执行"
        value={manualCmd}
        onChange={(e) => setManualCmd(e.target.value)}
        className="command-textarea"
      />

      <div className="command-actions">
        <Button type="primary" onClick={() => execCommand(manualCmd, setManualOutput)}>
          执行命令
        </Button>
        <Upload showUploadList={false} beforeUpload={handleUpload}>
          <Button icon={<UploadOutlined />} className="upload-button">
            上传命令文件
          </Button>
        </Upload>
      </div>

      <TextArea rows={5} readOnly value={manualOutput} className="command-result" placeholder="执行结果" />
    </Card>
  );
}

export default CommandInput;
