import { requireNativeModule } from "expo-modules-core";

import { type SaExpoWebserverModuleType } from "./SaExpoWebserver.types";

// It loads the native module object from the JSI or falls back to
// the bridge module (from NativeModulesProxy) if the remote debugger is on.
export default requireNativeModule(
  "SaExpoWebserver"
) as SaExpoWebserverModuleType;
