import React, { useState } from "react";
import { Input, Button } from "antd";

const { TextArea } = Input;

function CustomCommand({ execCommand }) {
  const [customCmd, setCustomCmd] = useState("");
  const [customOutput, setCustomOutput] = useState("");

  return (
    <div className="custom-command">
      <TextArea
        rows={4}
        placeholder="输入命令，每行一条"
        value={customCmd}
        onChange={(e) => setCustomCmd(e.target.value)}
        className="command-input"
      />
      <Button type="primary" className="execute-button" onClick={() => execCommand(customCmd, setCustomOutput)}>
        执行
      </Button>
      <TextArea rows={4} readOnly value={customOutput} className="command-output" placeholder="执行结果" />
    </div>
  );
}

export default CustomCommand;
