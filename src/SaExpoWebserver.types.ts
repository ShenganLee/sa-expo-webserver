export type SaExpoWebserverModuleType = {
  start: (fileDir: string, proxyStr: string) => Promise<string>;
  stop: (addr: string) => Promise<void>;
  restart: (addr: string) => Promise<void>;
};

export type SaExpoWebserverProxyType = {
  Path: string;
  Target: string;
  Include?: string[];
  Exclude?: string[];
};
