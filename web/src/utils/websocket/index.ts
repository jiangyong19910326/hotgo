import { SocketEnum } from '@/enums/socketEnum';
import { useUserStoreWidthOut } from '@/store/modules/user';
import { isJsonString } from '@/utils/is';
import { registerGlobalMessage } from '@/utils/websocket/registerMessage';

// WebSocket消息格式
export interface WebSocketMessage {
  event: string;
  data: any;
  code: number;
  timestamp: number;
}

type WebSocketMessageHandler = (message: WebSocketMessage) => void;

let socket: WebSocket;
let isActive: boolean;
const messageHandler: Map<string, WebSocketMessageHandler> = new Map();

function normalizeWsAddr(wsAddr = '') {
  if (!wsAddr) return '';
  try {
    const url = new URL(wsAddr);
    const pageHost = window.location.hostname;
    const wsHost = url.hostname;
    const isPageLocal = ['localhost', '127.0.0.1', '::1'].includes(pageHost);
    const isWsLocal = ['localhost', '127.0.0.1', '::1'].includes(wsHost);
    if (isWsLocal && !isPageLocal) {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      return `${protocol}//${window.location.host}${url.pathname}`;
    }
    if (window.location.protocol === 'https:' && url.protocol === 'ws:') {
      url.protocol = 'wss:';
      return url.toString().replace(/\/$/, '');
    }
    return wsAddr;
  } catch {
    return wsAddr;
  }
}

export default () => {
  const heartCheck = {
    timeout: 5000,
    timeoutObj: setTimeout(() => {}),
    serverTimeoutObj: setInterval(() => {}),
    reset: function () {
      clearTimeout(this.timeoutObj);
      clearTimeout(this.serverTimeoutObj);
      return this;
    },
    start: function () {
      // eslint-disable-next-line @typescript-eslint/no-this-alias
      const self = this;
      clearTimeout(this.timeoutObj);
      clearTimeout(this.serverTimeoutObj);
      this.timeoutObj = setTimeout(function () {
        socket.send(
          JSON.stringify({
            event: SocketEnum.EventPing,
          })
        );
        self.serverTimeoutObj = setTimeout(function () {
          console.log('[WebSocket] 关闭服务');
          socket.close();
        }, self.timeout);
      }, this.timeout);
    },
  };

  const useUserStore = useUserStoreWidthOut();
  let lockReconnect = false;
  let timer: ReturnType<typeof setTimeout>;
  const createSocket = () => {
    console.log('[WebSocket] createSocket...');
    if (useUserStore.token === '' || useUserStore.config?.wsAddr == '') {
      console.error('[WebSocket] 用户未登录，稍后重试...');
      resetReconnect();
      return;
    }
    try {
      const wsAddr = normalizeWsAddr(useUserStore.config?.wsAddr);
      socket = new WebSocket(`${wsAddr}?authorization=${useUserStore.token}`);
      init();
      if (lockReconnect) {
        lockReconnect = false;
      }
    } catch (e) {
      console.error(`[WebSocket] createSocket err: ${e}`);
      resetReconnect();
      return;
    }
  };

  const resetReconnect = () => {
    if (lockReconnect) {
      lockReconnect = false;
    }
    reconnect();
  };

  const reconnect = () => {
    console.log('[WebSocket] lockReconnect:' + lockReconnect);
    if (lockReconnect) return;
    lockReconnect = true;
    clearTimeout(timer);
    timer = setTimeout(() => {
      createSocket();
    }, 1000 * 10);
  };

  const init = () => {
    socket.onopen = function (_) {
      console.log('[WebSocket] 已连接');
      heartCheck.reset().start();
      isActive = true;
    };

    socket.onmessage = function (event) {
      isActive = true;
      // console.log('WebSocket:收到一条消息', event.data);

      if (!isJsonString(event.data)) {
        console.log('[WebSocket] message incorrect format:' + JSON.stringify(event));
        return;
      }

      heartCheck.reset().start();

      const message = JSON.parse(event.data) as WebSocketMessage;
      onMessage(message);
    };

    socket.onerror = function (_) {
      console.log('[WebSocket] 发生错误');
      reconnect();
      isActive = false;
    };

    socket.onclose = function (_) {
      console.log('[WebSocket] 已关闭');
      heartCheck.reset();
      reconnect();
      isActive = false;
    };

    window.onbeforeunload = function () {
      socket.close();
      isActive = false;
    };
  };

  createSocket();
  registerGlobalMessage();
};

function onMessage(message: WebSocketMessage) {
  let handled = false;
  messageHandler.forEach((value: WebSocketMessageHandler, key: string) => {
    if (message.event === key || key === '*') {
      handled = true;
      value.call(null, message);
    }
  });

  if (!handled) {
    console.log('[WebSocket] messageHandler not registered. message:' + JSON.stringify(message));
  }
}

// 发送消息
export function sendMsg(event: string, data: any = null, isRetry = true) {
  if (socket === undefined || !isActive) {
    if (!isRetry) {
      console.log('[WebSocket] 连接异常，发送失败！');
      return;
    }
    console.log('[WebSocket] 连接异常，等待重试..');
    setTimeout(() => {
      sendMsg(event, data);
    }, 200);
    return;
  }

  try {
    socket.send(JSON.stringify({ event, data }));
  } catch (err: any) {
    console.log('[WebSocket] 发送消息失败，err：', err.message);
    if (!isRetry) {
      return;
    }

    console.log('[WebSocket] 等待重试..');
    setTimeout(() => {
      sendMsg(event, data);
    }, 100);
  }
}

// 添加消息处理
export function addOnMessage(key: string, value: WebSocketMessageHandler): void {
  messageHandler.set(key, value);
}

// 移除消息处理
export function removeOnMessage(key: string): boolean {
  return messageHandler.delete(key);
}

// 查看所有消息处理
export function getAllOnMessage(): Map<string, WebSocketMessageHandler> {
  return messageHandler;
}
