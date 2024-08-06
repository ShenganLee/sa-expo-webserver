import {
  ConfigPlugin,
  createRunOncePlugin,
  withProjectBuildGradle,
} from "expo/config-plugins";

const pkg = require("../../package.json");

const gradleMaven = [
  "",
  `// ${pkg.name}`,
  `allprojects {`,
  `   repositories {`,
  `       flatDir { dirs project(':sa-expo-webserver').file('libs')  }`,
  `   }`,
  `}`,
  `// ${pkg.name}`,
  "",
].join("\n");

const withGoWebServer: ConfigPlugin = (config) => {
  return withProjectBuildGradle(config, (config) => {
    if (config.modResults.language === "groovy") {
      // console.log(config.modResults.contents)
      config.modResults.contents += gradleMaven;
    }

    return config;
  });
};

export default createRunOncePlugin(withGoWebServer, pkg.name, pkg.version);
