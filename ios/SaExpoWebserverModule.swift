import ExpoModulesCore

import Foundation
import Gowebserver

public class SaExpoWebserverModule: Module {

  private var timer = WebServerkeepAliveTimer()

  private func keepWebserverAlive() {
    // print("keepWebserverAlive GowebserverIsRunning \(GowebserverIsRunning())")
    // print("keepWebserverAlive GowebserverHealthy \(GowebserverHealthy())")
    if GowebserverIsRunning() && !GowebserverHealthy() {
      GowebserverRestart()
    }
  }

  // Each module class must implement the definition function. The definition consists of components
  // that describes the module's functionality and behavior.
  // See https://docs.expo.dev/modules/module-api for more details about available components.
  public func definition() -> ModuleDefinition {
    // Sets the name of the module that JavaScript code will use to refer to the module. Takes a string as an argument.
    // Can be inferred from module's class name, but it's recommended to set it explicitly for clarity.
    // The module will be accessible from `requireNativeModule('SaExpoWebserver')` in JavaScript.
    Name("SaExpoWebserver")

    OnCreate() {
      self.timer.setupTimer(interval: 2.0, repeats: true, action: {  self.keepWebserverAlive() })
    }

    OnDestroy() {
      GowebserverStop()
      self.timer.stopTimer()
    }

    OnAppEntersForeground {
      keepWebserverAlive()
    }

    OnAppBecomesActive {
      keepWebserverAlive()
    }

    Function("start") { (serverConfigStr: String) in
      return GowebserverStart(serverConfigStr)
    }

    Function("stop") { () in
      GowebserverStop()
    }

    Function("restart") { () in
      GowebserverRestart()
    }

    Function("isRunning") {() in
      return GowebserverIsRunning()
    }

    Function("healthy") {() in
      return GowebserverHealthy()
    }

    Function("serverUrl") {() in
      return GowebserverServerUrl()
    }

    Function("setLogFile") {(logFile: String) in
      return GowebserverSetLogFile(logFile)
    }

    Function("logFileClose") {() in
      return GowebserverLogFileClose()
    }
  }
}

private class WebServerkeepAliveTimer {
    func setupTimer(interval: TimeInterval, repeats: Bool, action: @escaping () -> Void) {
        let timer = Timer.scheduledTimer(
          timeInterval: interval,
          target: self,
          selector: #selector(performAction),
          userInfo: ["action": action],
          repeats: repeats
        )
        // 将timer存储起来，以便在需要时可以停止它
        self.timer = timer
    }

    func stopTimer() {
        self.timer?.invalidate()
        self.timer = nil
    }
    
    @objc private func performAction() {
        // 这里放置定时器触发时要执行的代码
        if let userInfo = self.timer?.userInfo as? [String: Any],
           let action = userInfo["action"] as? () -> Void {
            action() // 执行闭包
        }
    }
    
    var timer: Timer?
}
