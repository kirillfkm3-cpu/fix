import React, { useState } from 'react';

function Chats() {
  const [message, setMessage] = useState('');

  const handleSend = (e) => {
    e.preventDefault();
    setMessage('');
  };

  return (
    <div className="card" style={{ display: 'flex', height: '70vh', padding: 0, overflow: 'hidden', backgroundColor: 'rgba(24, 24, 32, 0.8)' }}>
      {/* Левая панель - Список чатов */}
      <div style={{ width: '30%', borderRight: '1px solid rgba(255,255,255,0.05)', display: 'flex', flexDirection: 'column', backgroundColor: 'rgba(15, 15, 20, 0.5)' }}>
        <div style={{ padding: '20px', borderBottom: '1px solid rgba(255,255,255,0.05)', fontWeight: 'bold', color: '#8b5cf6' }}>Чаты</div>
        <div style={{ overflowY: 'auto', flex: 1 }}>
          {/* Активный чат */}
          <div style={{ padding: '15px', borderBottom: '1px solid rgba(255,255,255,0.03)', cursor: 'pointer', backgroundColor: 'rgba(139, 92, 246, 0.15)', borderLeft: '4px solid #8b5cf6' }}>
            <div style={{ fontWeight: 'bold', color: '#fff' }}>Нурбек Нурбеков</div>
            <div style={{ fontSize: '12px', color: '#a1a1aa' }}>Последнее сообщение...</div>
          </div>
          {/* Другой чат */}
          <div style={{ padding: '15px', borderBottom: '1px solid rgba(255,255,255,0.03)', cursor: 'pointer' }}>
            <div style={{ fontWeight: 'bold', color: '#e4e4e7' }}>Группа Golang</div>
            <div style={{ fontSize: '12px', color: '#71717a' }}>Новое сообщение</div>
          </div>
        </div>
      </div>

      {/* Правая панель - Окно сообщений */}
      <div style={{ width: '70%', display: 'flex', flexDirection: 'column', backgroundColor: 'transparent' }}>
        <div style={{ padding: '15px', borderBottom: '1px solid rgba(255,255,255,0.05)', fontWeight: 'bold', color: '#fff', display: 'flex', alignItems: 'center', gap: '10px' }}>
          <div style={{ width: '35px', height: '35px', borderRadius: '50%', background: 'linear-gradient(135deg, #8b5cf6, #3b82f6)' }}></div>
          Нурбек Нурбеков
        </div>
        
        {/* Область сообщений */}
        <div style={{ flex: 1, padding: '20px', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {/* Входящее */}
          <div style={{ textAlign: 'left' }}>
            <span style={{ display: 'inline-block', padding: '12px 16px', borderRadius: '18px', backgroundColor: 'rgba(255,255,255,0.05)', color: '#fff', maxWidth: '70%', border: '1px solid rgba(255,255,255,0.05)' }}>
              Привет, как дела?
            </span>
          </div>
          {/* Исходящее */}
          <div style={{ textAlign: 'right' }}>
            <span style={{ display: 'inline-block', padding: '12px 16px', borderRadius: '18px', background: 'linear-gradient(135deg, #6366f1, #8b5cf6)', color: 'white', maxWidth: '70%', boxShadow: '0 4px 15px rgba(139, 92, 246, 0.2)' }}>
              Привет! Всё отлично, делаю проект.
            </span>
          </div>
        </div>

        {/* Поле ввода */}
        <div style={{ padding: '20px', borderTop: '1px solid rgba(255,255,255,0.05)' }}>
          <form onSubmit={handleSend} style={{ display: 'flex', gap: '12px' }}>
            <button type="button" style={{ fontSize: '20px', background: 'transparent', border: 'none', cursor: 'pointer' }}>😊</button>
            <input 
              type="text" 
              placeholder="Введите сообщение..." 
              value={message} 
              onChange={(e) => setMessage(e.target.value)} 
              style={{ 
                flex: 1, 
                padding: '12px 20px', 
                borderRadius: '25px', 
                border: '1px solid rgba(255,255,255,0.1)', 
                backgroundColor: 'rgba(0,0,0,0.2)', 
                color: '#fff',
                outline: 'none'
              }} 
            />
            <button type="submit" style={{ padding: '10px 25px', backgroundColor: '#8b5cf6', color: 'white', border: 'none', borderRadius: '25px', cursor: 'pointer', fontWeight: 'bold', transition: '0.3s' }}>
              Отправить
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}

export default Chats;