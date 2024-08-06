import { requireNativeViewManager } from 'expo-modules-core';
import * as React from 'react';

import { SaExpoWebserverViewProps } from './SaExpoWebserver.types';

const NativeView: React.ComponentType<SaExpoWebserverViewProps> =
  requireNativeViewManager('SaExpoWebserver');

export default function SaExpoWebserverView(props: SaExpoWebserverViewProps) {
  return <NativeView {...props} />;
}
