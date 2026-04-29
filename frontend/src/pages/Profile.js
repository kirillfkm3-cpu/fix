import React, { useEffect, useState } from 'react';

function Profile() {
  const [user, setUser] = useState(null);
  const [posts, setPosts] = useState([]);
  const [followers, setFollowers] = useState([]);
  const [following, setFollowing] = useState([]);
  const [error, setError] = useState('');
  const [isPublic, setIsPublic] = useState(true);

  useEffect(() => {
    // 1. Берем ID из памяти
    const userId = localStorage.getItem('userId');
    
    if (!userId) {
      setError('Вы не авторизованы. Войдите в систему.');
      return;
    }

    // 2. Делаем запрос к API профиля
    fetch(`/api/profile?id=${userId}`, {
      credentials: 'include'
    })
      .then(res => {
        if (!res.ok) throw new Error('Пользователь не найден');
        return res.json();
      })
      .then(data => {
        setUser(data);
      })
      .catch(err => {
        setError(err.message);
      });

    // 3. Делаем запрос к API постов
    fetch(`/api/posts?user_id=${userId}`, {
      credentials: 'include'
    })
      .then(res => {
        if (!res.ok) throw new Error('Ошибка загрузки постов');
        return res.json();
      })
      .then(data => {
        setPosts(data);
      })
      .catch(err => {
        console.error('Ошибка постов:', err);
      });

    fetch(`/api/followers?user_id=${userId}`, {
      credentials: 'include'
    })
      .then(res => res.json())
      .then(data => setFollowers(data))
      .catch(err => console.error('Ошибка followers:', err));

    fetch(`/api/following?user_id=${userId}`, {
      credentials: 'include'
    })
      .then(res => res.json())
      .then(data => setFollowing(data))
      .catch(err => console.error('Ошибка following:', err));
  }, []);

  const handleTogglePublic = async () => {
    try {
      const response = await fetch(`/api/profile/update`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ is_public: !isPublic })
      });

      if (response.ok) {
        setIsPublic(!isPublic);
        alert('Профиль обновлен');
      } else {
        alert('Ошибка обновления');
      }
    } catch (err) {
      alert('Ошибка сети');
    }
  };

  if (error) return <div className="card" style={{color: '#ff4d4d', textAlign: 'center'}}>{error}</div>;
  if (!user) return <div className="card" style={{textAlign: 'center'}}>Загрузка...</div>;

  return (
    <div className="card" style={{ maxWidth: '600px', margin: '20px auto' }}>
      <div style={{ textAlign: 'center' }}>
        <div style={avatarCircle}>
          {user.first_name?.[0]}{user.last_name?.[0]}
        </div>
        <h2 style={{ color: '#fff' }}>{user.first_name} {user.last_name}</h2>
        <p style={{ color: '#8b5cf6' }}>@{user.nickname || 'id' + user.id.slice(0,4)}</p>
        <button onClick={handleTogglePublic} style={{ padding: '8px 16px', background: 'linear-gradient(135deg, #8b5cf6, #6366f1)', color: 'white', border: 'none', borderRadius: '8px', cursor: 'pointer' }}>
          {isPublic ? 'Сделать приватным' : 'Сделать публичным'}
        </button>
      </div>

      <div style={{ marginTop: '20px', borderTop: '1px solid rgba(255,255,255,0.1)', paddingTop: '20px' }}>
        <p style={label}>Email</p>
        <p style={value}>{user.email}</p>
        
        <p style={label}>Дата рождения</p>
        <p style={value}>{user.date_of_birth}</p>

        <p style={label}>О себе</p>
        <p style={value}>{user.about_me || 'Ничего не указано'}</p>

        <p style={label}>Подписчики ({followers.length})</p>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px' }}>
          {followers.slice(0, 5).map(f => (
            <div key={f.id} style={{ padding: '5px 10px', background: 'rgba(255,255,255,0.1)', borderRadius: '8px', fontSize: '12px' }}>
              {f.first_name} {f.last_name}
            </div>
          ))}
        </div>

        <p style={label}>Подписки ({following.length})</p>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '10px' }}>
          {following.slice(0, 5).map(f => (
            <div key={f.id} style={{ padding: '5px 10px', background: 'rgba(255,255,255,0.1)', borderRadius: '8px', fontSize: '12px' }}>
              {f.first_name} {f.last_name}
            </div>
          ))}
        </div>

        <p style={label}>Посты</p>
        {posts.length === 0 ? (
          <p style={value}>Нет постов</p>
        ) : (
          posts.map(post => (
            <div key={post.id} style={{ border: '1px solid rgba(255,255,255,0.1)', padding: '10px', marginBottom: '10px', borderRadius: '8px' }}>
              <p style={{ color: '#e4e4e7' }}>{post.content}</p>
              <p style={{ color: '#71717a', fontSize: '12px' }}>{new Date(post.created_at).toLocaleString()}</p>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

const avatarCircle = {
  width: '80px', height: '80px', borderRadius: '50%', 
  background: 'linear-gradient(135deg, #8b5cf6, #3b82f6)', 
  margin: '0 auto 10px', display: 'flex', alignItems: 'center', 
  justifyContent: 'center', fontSize: '30px', fontWeight: 'bold'
};
const label = { color: '#71717a', fontSize: '12px', marginBottom: '2px' };
const value = { color: '#e4e4e7', marginBottom: '15px' };

export default Profile;