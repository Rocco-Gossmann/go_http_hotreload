const Hotreload = (function () {
  let socket;

  let forceRefresh = false;
  let connected = false;
  let reconnectTimeout = -1;

  function removeAllListeners() {
    if (!socket) return;
    socket.close();
    socket.removeEventListener("open", onOpen);
    socket.removeEventListener("error", onError);
    socket.removeEventListener("close", onClose);
  }

  function onError() {
    forceRefresh = connected;
    tryReconnect();
  }

  function onOpen() {
    socket.addEventListener("close", onClose);
    socket.send("ready");
    connected = true;
  }

  function onClose() {
    forceRefresh = true;
    tryReconnect();
  }

  function tryReconnect() {
    window.clearTimeout(reconnectTimeout);
    removeAllListeners();

    // Check if the Hotreload module is active on the server
    fetch("/__hotreload.ws", { method: "HEAD" })
      .then((res) => {
        if (res.status == 200) {
          // if the hotreload module is active on the server
          if (forceRefresh) {
            // and the connection was closed before
            removeAllListeners();
            window.location.href = window.location.href;
          } else {
            // if the connection was not closed before connect to the websocket
            socket = new WebSocket("/__hotreload.ws");
            socket.addEventListener("error", onError);
            socket.addEventListener("open", onOpen);
          }
        } else {
          // if the module is not active on the server, treat it as an error
          throw "";
        }
      })
      .catch(() => {
        // In case of any error, retry the connect
        reconnectTimeout = window.setTimeout(tryReconnect, 500);
      });
  }

  tryReconnect();

  return {
    Stop: () => socket.close(),
  };
})();
