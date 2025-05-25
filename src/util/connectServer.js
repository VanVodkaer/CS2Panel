import api from "../config/axiosConfig";

async function getGameConnectUrl(name) {
  try {
    const addrRes = await api.get("/info/network/addr");
    const portRes = await api.get("/info/network/gameport", { params: { name } });
    const passwdRes = await api.get("/info/network/gamepasswd", { params: { name } });

    const addr = addrRes.data.addr;
    const port = portRes.data.port;
    const passwd = passwdRes.data.passwd;

    return `steam://connect/${addr}:${port}/${passwd}`;
  } catch (error) {
    console.error("获取连接信息失败:", error);
  }
}

async function getTvConnectUrl(name) {
  try {
    const addrRes = await api.get("/info/network/addr");
    const portRes = await api.get("/info/network/tvport", { params: { name } });
    const passwdRes = await api.get("/info/network/tvpasswd", { params: { name } });

    console.log(passwdRes);

    const addr = addrRes.data.addr;
    const port = portRes.data.port;
    const passwd = passwdRes.data.passwd;

    return `steam://connect/${addr}:${port}/${passwd}`;
  } catch (error) {
    console.error("获取连接信息失败:", error);
  }
}

async function getGameConnectCommand(name) {
  console.log(name);

  try {
    const addrRes = await api.get("/info/network/addr");
    const portRes = await api.get("/info/network/gameport", { params: { name } });
    const passwdRes = await api.get("/info/network/gamepasswd", { params: { name } });

    const addr = addrRes.data.addr;
    const port = portRes.data.port;
    const passwd = passwdRes.data.passwd;

    if (passwd) {
      return `connect ${addr}:${port}; password ${passwd}`;
    } else {
      return `connect ${addr}:${port}`;
    }
  } catch (error) {
    console.error("获取连接信息失败:", error);
    return "";
  }
}

async function getTvConnectCommand(name) {
  try {
    const addrRes = await api.get("/info/network/addr");
    const portRes = await api.get("/info/network/tvport", { params: { name } });
    const passwdRes = await api.get("/info/network/tvpasswd", { params: { name } });

    const addr = addrRes.data.addr;
    const port = portRes.data.port;
    const passwd = passwdRes.data.passwd;

    if (passwd != "") {
      return `connect ${addr}:${port}; password ${passwd}`;
    } else {
      return `connect ${addr}:${port}`;
    }
  } catch (error) {
    console.error("获取连接信息失败:", error);
  }
}

export { getGameConnectCommand, getTvConnectCommand, getGameConnectUrl, getTvConnectUrl };
