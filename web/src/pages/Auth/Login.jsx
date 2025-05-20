import { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "D:/vsCode/bet_master/web/src/context/AuthContext";
import AuthForm from "D:/vsCode/bet_master/web/src/components/AuthForm";

const Login = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { login } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(email, password);
      navigate("/user/dashboard");
    } catch (err) {
      setError("Invalid credentials");
    }
  };

  return (
    <AuthForm
      title="Login"
      email={email}
      setEmail={setEmail}
      password={password}
      setPassword={setPassword}
      error={error}
      onSubmit={handleSubmit}
      linkText="Don't have an account? Register"
      linkPath="/register"
    />
  );
};

export default Login;