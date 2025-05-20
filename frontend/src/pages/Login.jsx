import { useNavigate } from 'react-router-dom'
import { login } from '../api/api'
import AuthForm from '../components/AuthForm'
import '../style/Login.css'

export default function Login() {
  const navigate = useNavigate()
  
  const handleLogin = async (credentials) => {
    try {
      const response = await login({
        email: credentials.email,
        password: credentials.password
      });
      navigate('/', { replace: true });
      setTimeout(() => navigate('/user/me', { replace: true }), 10);
    } catch (error) {
      throw error // Error will be handled in AuthForm
    }
  };

  return (
    <div className="login-bg-strict">
      <div className="login-card">
        <div className="login-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h1 className="login-title-strict">Вход в систему</h1>
        <p className="login-subtitle">Авторизация для управления ставками</p>
        <AuthForm type="login" onSubmit={handleLogin} />
      </div>
    </div>
  );
}