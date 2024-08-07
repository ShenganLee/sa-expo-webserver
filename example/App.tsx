import { WebServer } from "@sa/expo-webserver";
import { useState, useEffect } from "react";
import { StyleSheet, View } from "react-native";
import { WebView } from "react-native-webview";

export default function App() {
  const [uri, setUri] = useState("");

  useEffect(() => {
    const server = new WebServer();

    server
      .start("", [{ Path: "/", Target: "https://baidu.com" }])
      .then((uri) => {
        setUri(uri);
      });

    return () => {
      server.destory();
    };
  }, []);
  return (
    <View style={styles.container}>{uri && <WebView source={{ uri }} />}</View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#fff",
    alignItems: "center",
    justifyContent: "center",
  },
});
