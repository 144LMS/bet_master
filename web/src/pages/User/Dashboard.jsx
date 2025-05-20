import { useState, useEffect, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";
import BalanceCard from "../../components/BalanceCard";
import TransactionList from "../../components/TransactionList";
import "D:/vsCode/bet_master/web/src/pages/UserDashboard.css";

const Dashboard = () => {
  const { user, token } = useContext(AuthContext);
  const [balance, setBalance] = useState(0);
  const [transactions, setTransactions] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) navigate("/login");

    const fetchData = async () => {
      try {
        // Запрос баланса
        const balanceRes = await fetch(`/wallets/${user.id}/balance`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        const balanceData = await balanceRes.json();
        setBalance(balanceData.balance);

        // Запрос транзакций
        const transactionsRes = await fetch(`/wallets/${user.id}/transactions`, {
          headers: { Authorization: `Bearer ${token}` },
        });
        const transactionsData = await transactionsRes.json();
        setTransactions(transactionsData);
      } catch (err) {
        console.error(err);
      }
    };

    fetchData();
  }, [user, token, navigate]);

  return (
    <div className="user-dashboard">
      <h1>Welcome, {user.username}</h1>
      <BalanceCard balance={balance} />
      <div className="actions">
        <button onClick={() => navigate("/user/deposit")}>Deposit</button>
        <button onClick={() => navigate("/user/withdraw")}>Withdraw</button>
      </div>
      <TransactionList transactions={transactions} />
    </div>
  );
};

export default Dashboard;