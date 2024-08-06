import ExpoModulesCore

import Gowebserver

internal class FailedWebserverStart: Exception {
  override var reason: String {
    "Webserver start failed"
  }
}

public class SaExpoWebserverModule: Module {
  // Each module class must implement the definition function. The definition consists of components
  // that describes the module's functionality and behavior.
  // See https://docs.expo.dev/modules/module-api for more details about available components.
  public func definition() -> ModuleDefinition {
    // Sets the name of the module that JavaScript code will use to refer to the module. Takes a string as an argument.
    // Can be inferred from module's class name, but it's recommended to set it explicitly for clarity.
    // The module will be accessible from `requireNativeModule('SaExpoWebserver')` in JavaScript.
    Name("SaExpoWebserver")

    AsyncFunction("start") { (fileDir: String, proxyStr: String, promise: Promise) in
      let addr = GowebserverStart(fileDir, proxyStr)
      if addr.isEmpty {
        promise.reject(FailedWebserverStart())
      } else {
        promise.resolve(addr)
      }
    }

     AsyncFunction("stop") { (addr: String, promise: Promise) in
      GowebserverStop(addr)
      promise.resolve()
    }

    AsyncFunction("restart") { (addr: String, promise: Promise) in
      GowebserverRestart(addr)
      promise.resolve()
    }
  }
}
