import { NativeModulesProxy, EventEmitter, Subscription } from 'expo-modules-core';

// Import the native module. On web, it will be resolved to SaExpoWebserver.web.ts
// and on native platforms to SaExpoWebserver.ts
import SaExpoWebserverModule from './SaExpoWebserverModule';
import SaExpoWebserverView from './SaExpoWebserverView';
import { ChangeEventPayload, SaExpoWebserverViewProps } from './SaExpoWebserver.types';

// Get the native constant value.
export const PI = SaExpoWebserverModule.PI;

export function hello(): string {
  return SaExpoWebserverModule.hello();
}

export async function setValueAsync(value: string) {
  return await SaExpoWebserverModule.setValueAsync(value);
}

const emitter = new EventEmitter(SaExpoWebserverModule ?? NativeModulesProxy.SaExpoWebserver);

export function addChangeListener(listener: (event: ChangeEventPayload) => void): Subscription {
  return emitter.addListener<ChangeEventPayload>('onChange', listener);
}

export { SaExpoWebserverView, SaExpoWebserverViewProps, ChangeEventPayload };
