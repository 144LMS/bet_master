import { useContext } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { AuthContext } from '../context/AuthContext';
import '../styles/Header.css'; // Стили (создайте этот файл)

const Header = () => {
  const { user, token, logout } = useContext(AuthContext);
  const navigate = useNavigate();

  // Выход из системы
  const handleLogout = () => {
    logout();
    navigate('/login'); // Перенаправляем на страницу входа
  };

  return (
    <header className="header">
      <div className="header-container">
        {/* Логотип с ссылкой на главную */}
        <Link to="/" className="logo">
          BetMaster
        </Link>

        {/* Навигация */}
        <nav className="nav">
          {token ? (
            // Если пользователь авторизован
            <>
              <Link to="/user/dashboard" className="nav-link">
                Личный кабинет
              </Link>
              <button onClick={handleLogout} className="logout-button">
                Выйти
              </button>
              {user && (
                <span className="username">Привет, {user.username}!</span>
              )}
            </>
          ) : (
            // Если не авторизован
            <>
              <Link to="/login" className="nav-link">
                Вход
              </Link>
              <Link to="/register" className="nav-link">
                Регистрация
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
};

export default Header;