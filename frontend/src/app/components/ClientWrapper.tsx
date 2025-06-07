"use client";

if (typeof window !== "undefined") {
  // Client-side code here if needed
}

const ClientWrapper: React.FC<{
  children: React.ReactNode;
}> = ({ children }) => {
  return <>{children}</>;
};

export default ClientWrapper;
