import { AppState, NativeEventSubscription } from "react-native";

import type { SaExpoWebserverProxyType } from "./SaExpoWebserver.types";
import SaExpoWebserverModule from "./SaExpoWebserverModule";

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export const sync = async <T>(
  promise: Promise<T>,
): Promise<[Error | null, T | null]> => {
  try {
    const value = await promise;
    return [null, value];
  } catch (error: unknown) {
    if (!error) return [new Error(), null];

    if (typeof error === "string") return [new Error(error), null];

    // @ts-ignore
    return [error, null];
  }
};

export class WebServer {
  fileDir?: string;
  proxys?: SaExpoWebserverProxyType[];
  addr?: string;
  subscription?: NativeEventSubscription;

  constructor() {
    this.guard();
  }

  guard = (): void => {
    this.subscription = AppState.addEventListener(
      "change",
      async (nextAppState) => {
        if (!this.addr || nextAppState !== "active") return;

        this.restart();
      }
    );
  };

  destory = () => {
    this.stop();
    this.subscription?.remove();
    this.subscription = void 0;
  };

  start = async (
    fileDir: string = "",
    proxys: SaExpoWebserverProxyType[] = []
  ): Promise<string> => {
    this.fileDir = fileDir;
    this.proxys = proxys;

    try {
      this.addr = await SaExpoWebserverModule.start(
        this.fileDir,
        JSON.stringify(this.proxys)
      );

      return this.addr;
    } catch (error) {
      this.addr = void 0;

      this.fileDir = void 0;
      this.proxys = void 0;

      throw error;
    }
  };

  stop = async (): Promise<void> => {
    this.fileDir = void 0;
    this.proxys = void 0;
    if (this.addr) {
      await SaExpoWebserverModule.stop(this.addr);
      this.addr = void 0;
    }
  };

  restart = async (force = false): Promise<void> => {
    if (this.addr) {
      const [err, res] = await sync(
        fetch(this.addr, {
          headers: {
            "Cache-Control": "no-cache",
            Pragma: "no-cache",
          },
        })
      );

      if (!err) {
        if (res?.ok && !force) return;

        if (res?.ok) {
          return SaExpoWebserverModule.restart(this.addr);
        }
      }

      const addr = this.addr;
      const fileDir = this.fileDir;
      const proxys = this.proxys;

      this.addr = void 0;
      let i = 0;
      while (addr !== this.addr && i < 10) {
        if (i !== 0) await SaExpoWebserverModule.stop(addr);
        await this.stop();
        await delay(20);
        await this.start(fileDir, proxys);
        i++;
      }

      if (addr !== this.addr) {
        await this.stop();
        this.addr = addr;
        this.fileDir = fileDir;
        this.proxys = proxys;
        throw new Error("restart failed");
      }
    }
  };
}
