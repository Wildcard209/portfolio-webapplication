'use client';

import { useAuth } from '../../lib/hooks/useAuth';

const ProtectedComponent = () => {
  const { isAuthenticated, user, logout, loading } = useAuth();

  if (loading) {
    return <div>Loading...</div>;
  }

  if (!isAuthenticated) {
    return (
      <div style={{ textAlign: 'center', padding: '20px' }}>
        <h3>Access Denied</h3>
        <p>You need to be logged in to view this content.</p>
      </div>
    );
  }

  return (
    <div style={{ padding: '20px' }}>
      <h3>Protected Content</h3>
      <p>
        Hello {user?.username}, this is protected content that only authenticated users can see.
      </p>
      <button
        onClick={logout}
        style={{
          padding: '8px 16px',
          backgroundColor: '#dc3545',
          color: 'white',
          border: 'none',
          borderRadius: '4px',
          cursor: 'pointer',
        }}
      >
        Logout
      </button>
    </div>
  );
};

export default ProtectedComponent;
