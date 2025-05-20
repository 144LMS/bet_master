import { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";

const Withdraw = () => {
  const [amount, setAmount] = useState(0);
  const [error, setError] = useState("");
  const { user, token } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleWithdraw = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch(`/wallets/${user.id}/withdraw`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ amount }),
      });
      if (!res.ok) throw new Error("Withdrawal failed");
      navigate("/user/dashboard");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="withdraw-form">
      <h1>Withdraw Money</h1>
      {error && <p className="error">{error}</p>}
      <form onSubmit={handleWithdraw}>
        <input
          type="number"
          min="0.01"
          step="0.01"
          value={amount}
          onChange={(e) => setAmount(parseFloat(e.target.value))}
          required
        />
        <button type="submit">Withdraw</button>
      </form>
    </div>
  );
};

export default Withdraw;