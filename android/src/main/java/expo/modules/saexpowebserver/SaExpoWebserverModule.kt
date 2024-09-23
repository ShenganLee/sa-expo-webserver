package expo.modules.saexpowebserver

import java.util.Timer
import java.util.TimerTask
import expo.modules.kotlin.modules.Module
import expo.modules.kotlin.modules.ModuleDefinition
import expo.modules.kotlin.exception.CodedException

import gowebserver.Gowebserver;

class SaExpoWebserverModule : Module() {

  private val timer = Timer()

  private fun keepWebserverAlive() {
    if (Gowebserver.isRunning() && !Gowebserver.healthy()) {
      Gowebserver.restart()
    }
  }

  // Each module class must implement the definition function. The definition consists of components
  // that describes the module's functionality and behavior.
  // See https://docs.expo.dev/modules/module-api for more details about available components.
  override fun definition() = ModuleDefinition {
    // Sets the name of the module that JavaScript code will use to refer to the module. Takes a string as an argument.
    // Can be inferred from module's class name, but it's recommended to set it explicitly for clarity.
    // The module will be accessible from `requireNativeModule('SaExpoWebserver')` in JavaScript.
    Name("SaExpoWebserver")

    OnCreate {
      val task = object : TimerTask() {
          override fun run() {
              // 这里可以放置需要定时执行的代码
              this@JowoiotExpoWebserverModule.keepWebserverAlive()
          }
      }

      // 计划任务在5秒后执行，并且每隔2秒执行一次
      timer.schedule(task, 5000, 2000)
    }

    OnDestroy {
      Gowebserver.stop()
      timer.cancel()
    }

    OnActivityEntersForeground {
      keepWebserverAlive()
    }

    // Defines a JavaScript function that always returns a Promise and whose native code
    // is by default dispatched on the different thread than the JavaScript runtime runs on.
    Function("start") { serverConfigStr: String ->
       return@Function Gowebserver.start(serverConfigStr)
    }

    Function("stop") {
      Gowebserver.stop()
    }

    Function("restart") {
      Gowebserver.restart()
    }

    Function("isRunning") {
       return@Function Gowebserver.isRunning()
    }

    Function("healthy") {
       return@Function Gowebserver.healthy()
    }

    Function("serverUrl") {
       return@Function Gowebserver.serverUrl()
    }

    Function("setLogFile") { logFile: String ->
       return@Function Gowebserver.setLogFile(logFile)
    }

    Function("logFileClose") {
       return@Function Gowebserver.logFileClose()
    }
  }
}
