import React from 'react';
import { Spinner } from '../../frontend/src/component/ui/Spinner';

export default {
  title: 'UI/Spinner',
  component: Spinner
};

export const spinner = () => {
  return (
    <div
      style={{
        backgroundColor: '#000'
      }}
    >
      <Spinner />
    </div>
  );
};
