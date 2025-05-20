const AuthForm = ({
  title,
  email,
  setEmail,
  password,
  setPassword,
  error,
  onSubmit,
  linkText,
  linkPath,
  showUsername = false,
  username,
  setUsername,
}) => {
  return (
    <div className="auth-container">
      <h1>{title}</h1>
      {error && <p className="error">{error}</p>}
      <form onSubmit={onSubmit}>
        {showUsername && (
          <input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
        )}
        <input
          type="email"
          placeholder="Email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">{title}</button>
      </form>
      <a href={linkPath}>{linkText}</a>
    </div>
  );
};

export default AuthForm;