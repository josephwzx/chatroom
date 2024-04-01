import React, { useState } from 'react';
import Header from '../Header';

const Register = ({ onRegister }) => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      if (!response.ok) {
        throw new Error('Registration failed');
      }
      onRegister(); 
    } catch (error) {
      setErrorMessage(error.message);
    }
  };

  return (
    <div className='app-container'>
      <Header />
        <div className='login'>
        <form onSubmit={handleSubmit} className="login-form">
          <div className="input-group">
            <label htmlFor="username">Username:</label>
            <input
              id="username"
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div className="input-group">
            <label htmlFor="password">Password:</label>
            <input
              id="password"
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          <button type="submit">Register</button>
          {errorMessage && <p>{errorMessage}</p>}
        </form>
        </div>
    </div>
  );
};

export default Register;
