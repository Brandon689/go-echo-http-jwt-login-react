import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

interface Props {
  onLogin: (token: string) => void;
}

function LoginPage({ onLogin }: Props) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
  
    const response = await fetch('http://localhost:1323/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: new URLSearchParams({ username, password }),
    });
  
    if (!response.ok) {
      console.error('Login failed');
      return;
    }
  
    let data;
    try {
      data = await response.json();
    } catch (err) {
      console.error('Failed to parse JSON', err);
      return;
    }
  
    onLogin(data.token);
    navigate('/restricted');
  };

  return (
    <form onSubmit={handleSubmit}>
      <input value={username} onChange={(e) => setUsername(e.target.value)} />
      <input value={password} onChange={(e) => setPassword(e.target.value)} type="password" />
      <button type="submit">Login</button>
    </form>
  );
}

export default LoginPage;