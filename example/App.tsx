import { StyleSheet, Text, View } from 'react-native';

import * as SaExpoWebserver from '@sa/expo-webserver';

export default function App() {
  return (
    <View style={styles.container}>
      <Text>{SaExpoWebserver.hello()}</Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
