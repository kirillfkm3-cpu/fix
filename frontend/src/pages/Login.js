import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await fetch(`/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ email, password })
      });
      
      // Читаем JSON ответ от сервера
      const data = await response.json();

      if (response.ok) {
        // КРИТИЧЕСКИ ВАЖНО: Сохраняем ID пользователя в localStorage
        // Теперь вкладка "Профиль" сможет его прочитать
        if (data.user_id) {
          localStorage.setItem('userId', data.user_id);
          console.log("Успешный вход, ID сохранен:", data.user_id);
          navigate('/'); // Переходим на главную
        } else {
          setError('Сервер не прислал ID пользователя');
        }
      } else {
        setError(data.message || 'Неверный email или пароль');
      }
    } catch (err) {
      setError('Ошибка сети: проверьте соединение с сервером');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card" style={{ maxWidth: '400px', margin: '100px auto', border: '1px solid rgba(139, 92, 246, 0.2)' }}>
      <h2 style={{ textAlign: 'center', color: '#8b5cf6', marginBottom: '10px' }}>Вход</h2>
      <p style={{ textAlign: 'center', color: '#a1a1aa', fontSize: '14px', marginBottom: '20px' }}>С возвращением!</p>
      
      {error && (
        <div style={{ 
          backgroundColor: 'rgba(255, 77, 77, 0.1)', 
          color: '#ff4d4d', 
          padding: '10px', 
          borderRadius: '8px', 
          fontSize: '13px', 
          marginBottom: '15px',
          textAlign: 'center',
          border: '1px solid rgba(255, 77, 77, 0.2)'
        }}>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
          <label style={{ fontSize: '12px', color: '#71717a', marginLeft: '5px' }}>Email</label>
          <input 
            type="email" 
            placeholder="example@mail.com" 
            value={email}
            onChange={e => setEmail(e.target.value)} 
            required 
            style={inputStyle} 
          />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
          <label style={{ fontSize: '12px', color: '#71717a', marginLeft: '5px' }}>Пароль</label>
          <input 
            type="password" 
            placeholder="••••••••" 
            value={password}
            onChange={e => setPassword(e.target.value)} 
            required 
            style={inputStyle} 
          />
        </div>

        <button 
          type="submit" 
          disabled={loading}
          style={{ 
            ...buttonStyle,
            opacity: loading ? 0.7 : 1,
            cursor: loading ? 'not-allowed' : 'pointer'
          }}
        >
          {loading ? 'Вход...' : 'Войти'}
        </button>

        <p style={{ textAlign: 'center', fontSize: '14px', marginTop: '10px', color: '#a1a1aa' }}>
          Нет аккаунта? <span 
            onClick={() => navigate('/register')} 
            style={{ color: '#8b5cf6', cursor: 'pointer', fontWeight: 'bold' }}
          >
            Создать
          </span>
        </p>
      </form>
    </div>
  );
}

const inputStyle = { 
  padding: '12px 16px', 
  borderRadius: '12px', 
  border: '1px solid rgba(255,255,255,0.1)', 
  backgroundColor: 'rgba(0,0,0,0.3)', 
  color: 'white',
  outline: 'none',
  fontSize: '15px'
};

const buttonStyle = { 
  padding: '12px', 
  background: 'linear-gradient(135deg, #8b5cf6, #6366f1)', 
  color: 'white', 
  border: 'none', 
  borderRadius: '12px', 
  fontWeight: 'bold', 
  fontSize: '16px',
  marginTop: '10px',
  boxShadow: '0 4px 15px rgba(139, 92, 246, 0.3)'
};

export default Login;