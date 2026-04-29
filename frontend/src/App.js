import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import './App.css'; // Подключаем красивый дизайн

// Импорты страниц (мы создадим их на следующем шаге)
import Home from './pages/Home';
import Login from './pages/Login';
import Register from './pages/Register';
import Profile from './pages/Profile';
import Groups from './pages/Groups';
import Chats from './pages/Chats';

function App() {
  const handleLogout = () => {
    localStorage.removeItem('userId');
    fetch(`/api/logout`, { 
      method: 'POST',
      credentials: 'include'
    })
      .then(() => window.location.reload())
      .catch(() => window.location.reload());
  };

  const isLoggedIn = localStorage.getItem('userId');

  return (
    <Router>
      <div className="app-container">
        {/* Глобальная навигационная панель (Шапка сайта) */}
        <nav className="navbar">
          <div className="navbar-content">
            <Link to="/" className="nav-brand">SocialNet</Link>
            
            <ul className="nav-links">
              <li><Link to="/">Лента</Link></li>
              <li><Link to="/groups">Группы</Link></li>
              <li><Link to="/chats">Чаты</Link></li>
              {isLoggedIn && <li><Link to="/profile/me">Профиль</Link></li>}
            </ul>

            <div className="nav-auth">
              {isLoggedIn ? (
                <button onClick={handleLogout} className="btn-logout">Выход</button>
              ) : (
                <>
                  <Link to="/login" className="btn-login">Вход</Link>
                  <Link to="/register" className="btn-register">Регистрация</Link>
                </>
              )}
            </div>
          </div>
        </nav>

        {/* Область, где будут отображаться разные страницы в зависимости от URL */}
        <main className="main-content">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            {/* :id позволяет открывать профили разных пользователей, например /profile/123 */}
            <Route path="/profile/:id" element={<Profile />} />
            <Route path="/groups" element={<Groups />} />
            <Route path="/chats" element={<Chats />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;