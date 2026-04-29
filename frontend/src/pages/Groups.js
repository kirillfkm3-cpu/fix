import React, { useEffect, useState } from 'react';

function Groups() {
  const [groups, setGroups] = useState([]);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = () => {
    fetch(`/api/groups`, {
      credentials: 'include'
    })
      .then(res => {
        if (!res.ok) throw new Error('Ошибка загрузки групп');
        return res.json();
      })
      .then(data => {
        setGroups(data);
      })
      .catch(err => {
        setError(err.message);
      });
  };

  const handleCreateGroup = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const response = await fetch(`/api/groups/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ name: title, description })
      });

      if (response.ok) {
        setTitle('');
        setDescription('');
        loadGroups();
      } else {
        setError('Ошибка создания группы');
      }
    } catch (err) {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  const handleJoinGroup = async (groupId) => {
    try {
      const response = await fetch(`/api/groups/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ group_id: groupId })
      });

      if (response.ok) {
        alert('Запрос отправлен');
      } else {
        alert('Ошибка вступления');
      }
    } catch (err) {
      alert('Ошибка сети');
    }
  };

  return (
    <div>
      <div className="card">
        <h2>Создать группу</h2>
        {error && <div style={{ color: '#ff4d4d', marginBottom: '10px' }}>{error}</div>}
        <form onSubmit={handleCreateGroup} style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '15px' }}>
          <input 
            type="text" 
            placeholder="Название группы" 
            value={title} 
            onChange={(e) => setTitle(e.target.value)} 
            required 
            style={{ padding: '10px', borderRadius: '5px', border: '1px solid rgba(255,255,255,0.1)', backgroundColor: 'rgba(0,0,0,0.3)', color: 'white' }} 
          />
          <textarea 
            placeholder="Описание группы" 
            value={description} 
            onChange={(e) => setDescription(e.target.value)} 
            required 
            style={{ padding: '10px', borderRadius: '5px', border: '1px solid rgba(255,255,255,0.1)', backgroundColor: 'rgba(0,0,0,0.3)', color: 'white', minHeight: '60px' }} 
          />
          <button 
            type="submit" 
            disabled={loading}
            style={{ padding: '10px', background: 'linear-gradient(135deg, #8b5cf6, #6366f1)', color: 'white', border: 'none', borderRadius: '5px', cursor: loading ? 'not-allowed' : 'pointer', fontWeight: 'bold' }}
          >
            {loading ? 'Создание...' : 'Создать'}
          </button>
        </form>
      </div>

      <div className="card">
        <h2>Каталог групп</h2>
        <div style={{ marginTop: '15px', display: 'flex', flexDirection: 'column', gap: '15px' }}>
          {groups.length === 0 ? (
            <p>Нет групп</p>
          ) : (
            groups.map(group => (
              <div key={group.id} style={{ padding: '15px', border: '1px solid rgba(255,255,255,0.1)', borderRadius: '8px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <div>
                  <h4 style={{ margin: 0, color: '#fff' }}>{group.name}</h4>
                  <p style={{ margin: '5px 0 0 0', color: '#a1a1aa', fontSize: '14px' }}>{group.description}</p>
                </div>
                <button 
                  onClick={() => handleJoinGroup(group.id)}
                  style={{ padding: '8px 16px', background: 'linear-gradient(135deg, #8b5cf6, #6366f1)', color: 'white', border: 'none', borderRadius: '5px', cursor: 'pointer', fontWeight: 'bold' }}
                >
                  Вступить
                </button>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}

export default Groups;