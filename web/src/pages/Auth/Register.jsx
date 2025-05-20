import { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "D:/vsCode/bet_master/web/src/context/AuthContext";
import AuthForm from "D:/vsCode/bet_master/web/src/components/AuthForm";

const Register = () => {
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const { register } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await register(username, email, password);
      navigate("/user/dashboard");
    } catch (err) {
      setError("Registration failed");
    }
  };

  return (
    <AuthForm
      title="Register"
      email={email}
      setEmail={setEmail}
      password={password}
      setPassword={setPassword}
      error={error}
      onSubmit={handleSubmit}
      linkText="Already have an account? Login"
      linkPath="/login"
      showUsername
      username={username}
      setUsername={setUsername}
    />
  );
};

export default Register;