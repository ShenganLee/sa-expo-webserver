package expo.modules.saexpowebserver

import expo.modules.kotlin.Promise
import expo.modules.kotlin.modules.Module
import expo.modules.kotlin.modules.ModuleDefinition
import expo.modules.kotlin.exception.CodedException

import gowebserver.Gowebserver;

class SaExpoWebserverModule : Module() {
  // Each module class must implement the definition function. The definition consists of components
  // that describes the module's functionality and behavior.
  // See https://docs.expo.dev/modules/module-api for more details about available components.
  override fun definition() = ModuleDefinition {
    // Sets the name of the module that JavaScript code will use to refer to the module. Takes a string as an argument.
    // Can be inferred from module's class name, but it's recommended to set it explicitly for clarity.
    // The module will be accessible from `requireNativeModule('SaExpoWebserver')` in JavaScript.
    Name("SaExpoWebserver")

    // Defines a JavaScript function that always returns a Promise and whose native code
    // is by default dispatched on the different thread than the JavaScript runtime runs on.
    AsyncFunction("start") { fileDir: String, proxyStr: String, promise: Promise ->
      val addr = Gowebserver.start(fileDir, proxyStr)
      if (addr.isBlank()) {
        promise.reject(CodedException("Webserver start failed"))
      } else {
        promise.resolve(addr)
      }
    }

    AsyncFunction("stop") { addr: String, promise: Promise ->
      Gowebserver.stop(addr)
      promise.resolve()
    }

    AsyncFunction("restart") { addr: String, promise: Promise ->
      Gowebserver.restart(addr)
      promise.resolve()
    }
  }
}
