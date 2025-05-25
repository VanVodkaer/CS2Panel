import api from "../config/axiosConfig";
function getGameConnectCommand(name) {
  let addr, port, passwd;
  api.get("/info/network/addr", { name }).then((res) => {
    addr = res.data.addr;
    console.log(res.data.addr);
  });
  api.get("/info/network/gameport", { name }).then((res) => {
    port = res.data.port;
  });
  api.get("/info/network/gamepasswd", { name }).then((res) => {
    passwd = res.data.passwd;
  });

  if (passwd != "") {
    return `connect ${addr}:${port}; password ${passwd}`;
  } else {
    return `connect ${addr}:${port}`;
  }
}

function getTvConnectCommand(name) {
  let addr, port, passwd;
  api.get("/info/network/addr", { name }).then((res) => {
    addr = res.data.addr;
  });
  api.get("/info/network/tvport", { name }).then((res) => {
    port = res.data.port;
  });
  api.get("/info/network/tvpasswd", { name }).then((res) => {
    passwd = res.data.passwd;
  });

  if (passwd != "") {
    return `connect ${addr}:${port}; password ${passwd}`;
  } else {
    return `connect ${addr}:${port}`;
  }
}

export { getGameConnectCommand, getTvConnectCommand };
