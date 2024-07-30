"use client";

import { useState } from "react";

export default function Home() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = async () => {
    const response = await fetch("http://localhost:3000/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
      credentials: "include",
    });
    if (response.ok) {
      const user = await response.json();
      console.log(user);
    }
  };

  const getUser = async () => {
    const response = await fetch("http://localhost:3000/user/5", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
    });
    if (response.ok) {
      const user = await response.json();
      console.log(user);
    }
  };

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
  };

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(e.target.value);
  };

  return (
    <main>
      <header>
        <input
          name="username"
          value={username}
          onChange={handleUsernameChange}
          required
        />
        <input
          name="password"
          value={password}
          onChange={handlePasswordChange}
          required
        />
        <button onClick={handleLogin}>Login</button>
        <button onClick={getUser}>Get user</button>
      </header>
    </main>
  );
}
