import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

function Register() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
    first_name: '',
    last_name: '',
    date_of_birth: '',
    nickname: '',
    about_me: '',
  });

  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // Универсальный обработчик для всех полей
  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({
      ...formData,
      [name]: value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await fetch(`/api/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(formData),
      });

      const data = await response.json();

      if (response.ok) {
        // Если всё успешно, отправляем на логин
        navigate('/login');
      } else {
        // Если сервер вернул ошибку (например, email занят)
        setError(data.message || 'Ошибка регистрации. Проверьте данные.');
      }
    } catch (err) {
      setError('Ошибка сети: не удалось связаться с сервером');
      console.error('Register error:', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card" style={{ maxWidth: '500px', margin: '40px auto', border: '1px solid rgba(139, 92, 246, 0.2)' }}>
      <h2 style={{ textAlign: 'center', color: '#8b5cf6', marginBottom: '20px' }}>Регистрация</h2>
      
      {error && (
        <div style={{ 
          backgroundColor: 'rgba(255, 77, 77, 0.1)', 
          color: '#ff4d4d', 
          padding: '10px', 
          borderRadius: '8px', 
          marginBottom: '15px',
          textAlign: 'center',
          fontSize: '14px',
          border: '1px solid rgba(255, 77, 77, 0.2)' 
        }}>
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        <input 
          type="email" 
          name="email" 
          placeholder="Email *" 
          value={formData.email}
          onChange={handleChange} 
          required 
          style={inputStyle} 
        />
        
        <input 
          type="password" 
          name="password" 
          placeholder="Пароль *" 
          value={formData.password}
          onChange={handleChange} 
          required 
          style={inputStyle} 
        />

        <div style={{ display: 'flex', gap: '10px' }}>
          <input 
            type="text" 
            name="first_name" 
            placeholder="Имя *" 
            value={formData.first_name}
            onChange={handleChange} 
            required 
            style={inputStyle} 
          />
          <input 
            type="text" 
            name="last_name" 
            placeholder="Фамилия *" 
            value={formData.last_name}
            onChange={handleChange} 
            required 
            style={inputStyle} 
          />
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '5px' }}>
          <label style={{ fontSize: '12px', color: '#71717a', marginLeft: '5px' }}>Дата рождения *</label>
          <input 
            type="date" 
            name="date_of_birth" 
            value={formData.date_of_birth}
            onChange={handleChange} 
            required 
            style={inputStyle} 
          />
        </div>

        <input 
          type="text" 
          name="nickname" 
          placeholder="Никнейм (необязательно)" 
          value={formData.nickname}
          onChange={handleChange} 
          style={inputStyle} 
        />

        <textarea 
          name="about_me" 
          placeholder="Немного о себе..." 
          value={formData.about_me}
          onChange={handleChange} 
          style={{ ...inputStyle, minHeight: '80px', resize: 'none' }} 
        />

        <button 
          type="submit" 
          disabled={loading}
          style={{ 
            padding: '12px', 
            background: 'linear-gradient(135deg, #8b5cf6, #6366f1)', 
            color: 'white', 
            border: 'none', 
            borderRadius: '10px', 
            fontWeight: 'bold', 
            cursor: loading ? 'not-allowed' : 'pointer',
            opacity: loading ? 0.7 : 1,
            marginTop: '10px'
          }}
        >
          {loading ? 'Создание аккаунта...' : 'Создать аккаунт'}
        </button>

        <p style={{ textAlign: 'center', fontSize: '14px', color: '#a1a1aa' }}>
          Уже есть аккаунт? <span 
            onClick={() => navigate('/login')} 
            style={{ color: '#8b5cf6', cursor: 'pointer', fontWeight: 'bold' }}
          >
            Войти
          </span>
        </p>
      </form>
    </div>
  );
}

const inputStyle = { 
  padding: '12px 16px', 
  borderRadius: '10px', 
  border: '1px solid rgba(255,255,255,0.1)', 
  backgroundColor: 'rgba(0,0,0,0.3)', 
  color: 'white',
  outline: 'none',
  fontSize: '15px'
};

export default Register;