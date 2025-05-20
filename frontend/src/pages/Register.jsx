import { useNavigate } from 'react-router-dom'
import { register } from '../api/api'
import AuthForm from '../components/AuthForm'
import '../style/Register.css'

export default function Register() {
  const navigate = useNavigate()
  
  const handleRegister = async (userData) => {
    try {
      const response = await register(
        userData.username,
        userData.email,
        userData.password
      )
      navigate('/')
    } catch (error) {
      throw error // Error will be handled in AuthForm
    }
  }

  return (
    <div className="register-bg-strict">
      <div className="register-card">
        <div className="register-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h1 className="register-title-strict">Регистрация</h1>
        <p className="register-subtitle">Создайте новый аккаунт</p>
        <AuthForm type="register" onSubmit={handleRegister} />
      </div>
    </div>
  )
}