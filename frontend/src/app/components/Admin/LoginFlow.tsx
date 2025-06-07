"use client";

import { useState } from "react";
import TwoFactorAuthentication from "./TwoFactorAuthentication";

const LoginFlow = () => {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();

    if (username === "admin" && password === "password123") {
      setIsLoggedIn(true);
    } else {
      alert("Invalid credentials. Please try again!");
    }
  };

  if (isLoggedIn) {
    return <TwoFactorAuthentication />;
  }

  return (
    <form onSubmit={handleLogin} style={{ textAlign: "center" }}>
      <h2>Login</h2>
      <div style={{ marginBottom: "10px" }}>
        <label>Username:</label>
        <input
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
      </div>
      <div style={{ marginBottom: "10px" }}>
        <label>Password:</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>
      <button type="submit">Login</button>
    </form>
  );
};

export default LoginFlow;
