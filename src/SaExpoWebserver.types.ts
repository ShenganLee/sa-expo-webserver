export type SaExpoWebserverModuleType = {
  start: (ServerConfigStr: string) => string;
  stop: () => void;
  restart: () => void;
  isRunning: () => boolean;
  healthy: () => boolean;
  serverUrl: () => string;
  setLogFile: (logFile: string) => void;
  logFileClose: () => void;
};

export type AddHeadersType = {
  Regex?: string;
  Headers: Record<string, string>;
};

export type HeaderType = {
  RequestHeaders?: AddHeadersType[];
  ResponseHeaders?: AddHeadersType[];
};

export type IndexHeadersType = {
  RequestHeaders?: Record<string, string>;
  ResponseHeaders?: Record<string, string>;
};

export type SkipperType = string[];
export type HeadersType = HeaderType[];

export type ProxyType = {
  Path: string;
  Target: string;
  Skipper?: SkipperType;
  Headers?: HeadersType;
};

export type RouterType = {
  Path: string;
  FilePath: string;
  Headers?: HeadersType;
  IndexHeaders?: IndexHeadersType;
};

export type ProxysType = ProxyType[];
export type RoutersType = RouterType[];

export type ServerConfigType = {
  Port?: number;
  Proxys?: ProxysType;
  Routers?: RoutersType;
};
