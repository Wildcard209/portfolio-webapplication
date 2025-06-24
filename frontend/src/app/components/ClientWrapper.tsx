'use client';

import { ReactNode } from 'react';

interface ClientWrapperProps {
  children: ReactNode;
}

if (typeof window !== 'undefined') {
  // Client-side code here if needed
}

const ClientWrapper: React.FC<ClientWrapperProps> = ({ children }): React.ReactElement => {
  return <>{children}</>;
};

export default ClientWrapper;
