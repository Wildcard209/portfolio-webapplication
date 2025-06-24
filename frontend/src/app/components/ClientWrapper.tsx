'use client';

interface ClientWrapperProps {
  children: React.ReactNode;
}

if (typeof window !== 'undefined') {
  // Client-side code here if needed
}

const ClientWrapper: React.FC<ClientWrapperProps> = ({ children }) => {
  return <>{children}</>;
};

export default ClientWrapper;
