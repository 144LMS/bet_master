import { useNavigate } from 'react-router-dom';
import '../style/Home.css';

export default function Home() {
  const navigate = useNavigate();

  return (
    <div className="home-bg-strict">
      <div className="home-main-card">
        <div className="home-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h1 className="home-title-strict">Добро пожаловать в BetMaster</h1>
        <p className="home-subtitle">Лучшая платформа для ставок на спорт</p>
        
        <div className="home-buttons">
          <button 
            className="home-btn login-btn"
            onClick={() => navigate('/login')}
          >
            Войти
          </button>
          <button 
            className="home-btn register-btn"
            onClick={() => navigate('/register')}
          >
            Зарегистрироваться
          </button>
          <button 
            className="home-btn admin-btn"
            onClick={() => navigate('/admin-login')}
          >
            Если вы админ
          </button>
        </div>
        
        <div className="home-features">
          <div className="feature">
            <div className="feature-icon">🏆</div>
            <p>Лучшие коэффициенты</p>
          </div>
          <div className="feature">
            <div className="feature-icon">⚡</div>
            <p>Мгновенные выплаты</p>
          </div>
          <div className="feature">
            <div className="feature-icon">🛡️</div>
            <p>Безопасность</p>
          </div>
        </div>
      </div>
      <div className="home-footer-strict">
        <p>© 2024 BetMaster. Все права защищены.</p>
      </div>
    </div>
  );
}