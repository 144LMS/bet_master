import { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";

const Deposit = () => {
  const [amount, setAmount] = useState(0);
  const [error, setError] = useState("");
  const { user, token } = useContext(AuthContext);
  const navigate = useNavigate();

  const handleDeposit = async (e) => {
    e.preventDefault();
    try {
      const res = await fetch(`/wallets/${user.id}/deposit`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ amount }),
      });
      if (!res.ok) throw new Error("Deposit failed");
      navigate("/user/dashboard");
    } catch (err) {
      setError(err.message);
    }
  };

  return (
    <div className="deposit-form">
      <h1>Deposit Money</h1>
      {error && <p className="error">{error}</p>}
      <form onSubmit={handleDeposit}>
        <input
          type="number"
          min="0.01"
          step="0.01"
          value={amount}
          onChange={(e) => setAmount(parseFloat(e.target.value))}
          required
        />
        <button type="submit">Deposit</button>
      </form>
    </div>
  );
};

export default Deposit;