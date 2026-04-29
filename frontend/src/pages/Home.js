import React, { useEffect, useState } from 'react';

function Home() {
  const [posts, setPosts] = useState([]);
  const [error, setError] = useState('');
  const [newPost, setNewPost] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    loadFeed();
  }, []);

  const loadFeed = () => {
    fetch(`/api/feed`, {
      credentials: 'include'
    })
      .then(res => {
        if (!res.ok) throw new Error('Ошибка загрузки ленты');
        return res.json();
      })
      .then(data => {
        setPosts(data);
      })
      .catch(err => {
        setError(err.message);
      });
  };

  const handleCreatePost = async (e) => {
    e.preventDefault();
    if (!newPost.trim()) return;

    setLoading(true);
    try {
      const response = await fetch(`/api/posts/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ content: newPost, privacy: 'public' })
      });

      if (response.ok) {
        setNewPost('');
        loadFeed(); // Перезагрузить ленту
      } else {
        setError('Ошибка создания поста');
      }
    } catch (err) {
      setError('Ошибка сети');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="card" style={{ maxWidth: '600px', margin: '20px auto' }}>
      <h2>Главная страница (Лента постов)</h2>

      {error && <div style={{ color: '#ff4d4d', marginBottom: '10px' }}>{error}</div>}

      <form onSubmit={handleCreatePost} style={{ marginBottom: '20px' }}>
        <textarea
          value={newPost}
          onChange={e => setNewPost(e.target.value)}
          placeholder="Что у вас нового?"
          style={{ width: '100%', minHeight: '80px', padding: '10px', borderRadius: '8px', border: '1px solid rgba(255,255,255,0.1)', backgroundColor: 'rgba(0,0,0,0.3)', color: 'white' }}
        />
        <button type="submit" disabled={loading} style={{ marginTop: '10px', padding: '10px 20px', background: 'linear-gradient(135deg, #8b5cf6, #6366f1)', color: 'white', border: 'none', borderRadius: '8px' }}>
          {loading ? 'Публикация...' : 'Опубликовать'}
        </button>
      </form>

      {posts.length === 0 ? (
        <p>Нет постов</p>
      ) : (
        posts.map(post => (
          <div key={post.id} style={{ border: '1px solid rgba(255,255,255,0.1)', padding: '15px', marginBottom: '15px', borderRadius: '8px' }}>
            <p style={{ color: '#e4e4e7' }}>{post.content}</p>
            <p style={{ color: '#71717a', fontSize: '12px' }}>{new Date(post.created_at).toLocaleString()}</p>
          </div>
        ))
      )}
    </div>
  );
}

export default Home;