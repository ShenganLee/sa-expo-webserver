import * as React from 'react';

import { SaExpoWebserverViewProps } from './SaExpoWebserver.types';

export default function SaExpoWebserverView(props: SaExpoWebserverViewProps) {
  return (
    <div>
      <span>{props.name}</span>
    </div>
  );
}
