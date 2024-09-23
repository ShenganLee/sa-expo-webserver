import type { ServerConfigType } from "./SaExpoWebserver.types";
import SaExpoWebserverModule from "./SaExpoWebserverModule";

export type { ServerConfigType } from "./SaExpoWebserver.types";

export const start = (serverConfig: ServerConfigType): string => {
  const serverConfigStr = JSON.stringify(serverConfig);
  const uri = SaExpoWebserverModule.start(serverConfigStr);
  if (!uri) throw new Error("Failed to start webserver");

  return uri;
};
export const stop = (): void => SaExpoWebserverModule.stop();
export const restart = (): void => SaExpoWebserverModule.restart();
export const isRunning = (): boolean => SaExpoWebserverModule.isRunning();
export const healthy = (): boolean => SaExpoWebserverModule.healthy();
export const serverUrl = (): string => SaExpoWebserverModule.serverUrl();
export const setLogFile = (logFile: string): void => SaExpoWebserverModule.setLogFile(logFile);
export const logFileClose = (): void => SaExpoWebserverModule.logFileClose();
