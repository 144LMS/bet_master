import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

export default function AuthForm({ type, onSubmit }) {
  const [formData, setFormData] = useState({
    username: type === 'register' ? '' : undefined,
    email: '',
    password: ''
  })
  const [error, setError] = useState('')
  const navigate = useNavigate()

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    })
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      const response = await onSubmit(formData)
      // For register, we get user and wallet data
      if (type === 'register') {
        console.log('Registered user:', response.user)
        console.log('Wallet created:', response.wallet)
      }
      navigate('/')
    } catch (error) {
      setError(error.response?.data?.message || 'An error occurred')
    }
  }

  return (
    <form onSubmit={handleSubmit}>
      {error && <div className="error">{error}</div>}
      
      {type === 'register' && (
        <div className="form-group">
          <label>Username</label>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            required
            minLength={3}
          />
        </div>
      )}
      
      <div className="form-group">
        <label>Email</label>
        <input
          type="email"
          name="email"
          value={formData.email}
          onChange={handleChange}
          required
        />
      </div>
      
      <div className="form-group">
        <label>Password</label>
        <input
          type="password"
          name="password"
          value={formData.password}
          onChange={handleChange}
          required
          minLength={8}
        />
      </div>
      
      <button type="submit">
        {type === 'login' ? 'Login' : 'Register'}
      </button>
    </form>
  )
}