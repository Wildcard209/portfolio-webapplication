"use client";

import { useState, useEffect } from "react";
import { AuthService } from "../../../lib/auth/authService";

interface LoginFlowProps {
  adminToken: string;
}

const LoginFlow = ({ adminToken }: LoginFlowProps) => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  useEffect(() => {
    const checkAuthStatus = async () => {
      const isAuth = await AuthService.checkAuthStatus();
      if (isAuth) {
        setIsLoggedIn(true);
      }
    };
    
    checkAuthStatus();
  }, [adminToken]);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    const result = await AuthService.login(username, password);
    
    if (result.success) {
      setIsLoggedIn(true);
    } else {
      setError(result.error || "Login failed");
    }
    
    setIsLoading(false);
  };

  const handleLogout = async () => {
    setIsLoading(true);
    await AuthService.logout();
    setIsLoggedIn(false);
    setUsername("");
    setPassword("");
    setError("");
    setIsLoading(false);
  };

  if (isLoggedIn) {
    const user = AuthService.getUser();
    return (
      <div style={{ textAlign: "center", padding: "20px" }}>
        <h2>Admin Panel</h2>
        <p>Welcome back, {user?.username}!</p>
        <p>You are now authenticated and can access admin features.</p>
        <button 
          onClick={handleLogout}
          disabled={isLoading}
          style={{
            padding: "10px 20px",
            backgroundColor: "#dc3545",
            color: "white",
            border: "none",
            borderRadius: "4px",
            cursor: isLoading ? "not-allowed" : "pointer",
          }}
        >
          {isLoading ? "Logging out..." : "Logout"}
        </button>
      </div>
    );
  }

  return (
    <div style={{ maxWidth: "400px", margin: "0 auto", padding: "20px" }}>
      <form onSubmit={handleLogin} style={{ textAlign: "center" }}>
        <h2>Admin Login</h2>
        
        {error && (
          <div style={{
            color: "#dc3545",
            backgroundColor: "#f8d7da",
            border: "1px solid #f5c6cb",
            borderRadius: "4px",
            padding: "10px",
            marginBottom: "15px"
          }}>
            {error}
          </div>
        )}

        <div style={{ marginBottom: "15px" }}>
          <label style={{ display: "block", marginBottom: "5px", fontWeight: "bold" }}>
            Username:
          </label>
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
            disabled={isLoading}
            style={{
              width: "100%",
              padding: "8px",
              border: "1px solid #ccc",
              borderRadius: "4px",
              fontSize: "16px"
            }}
          />
        </div>

        <div style={{ marginBottom: "20px" }}>
          <label style={{ display: "block", marginBottom: "5px", fontWeight: "bold" }}>
            Password:
          </label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            disabled={isLoading}
            style={{
              width: "100%",
              padding: "8px",
              border: "1px solid #ccc",
              borderRadius: "4px",
              fontSize: "16px"
            }}
          />
        </div>

        <button 
          type="submit" 
          disabled={isLoading}
          style={{
            width: "100%",
            padding: "12px",
            backgroundColor: isLoading ? "#6c757d" : "#007bff",
            color: "white",
            border: "none",
            borderRadius: "4px",
            fontSize: "16px",
            cursor: isLoading ? "not-allowed" : "pointer",
          }}
        >
          {isLoading ? "Logging in..." : "Login"}
        </button>
      </form>
    </div>
  );
};

export default LoginFlow;
