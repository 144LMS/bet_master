import { useState, useEffect } from "react";
import { getWallet, getProfile, deposit } from "../api/api";
import axios from "axios";
import "../style/Transactions.css";

function formatDate(dateStr) {
  if (!dateStr) return "";
  let [date, timeWithZone] = dateStr.split(' ');
  if (!timeWithZone) return dateStr;
  let match = timeWithZone.match(/([+-]\d{2})(\d{2})$/);
  if (match) {
    timeWithZone = timeWithZone.replace(/([+-]\d{2})(\d{2})$/, `$1:$2`);
  }
  let iso = `${date}T${timeWithZone}`;
  let d = new Date(iso);
  return isNaN(d) ? dateStr : d.toLocaleString();
}

// Функция вывода средств через API
const withdraw = async (amount) => {
  // Используй тот же базовый url что и в других API
  return axios.post(
    "http://localhost:8080/api/wallets/withdraw",
    { amount: parseFloat(amount) },
    { withCredentials: true }
  );
};

export default function Transactions() {
  const [wallet, setWallet] = useState(null);
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [depositAmount, setDepositAmount] = useState('');
  const [depositLoading, setDepositLoading] = useState(false);
  const [depositError, setDepositError] = useState('');
  const [depositSuccess, setDepositSuccess] = useState('');
  // Состояния для вывода
  const [withdrawLoading, setWithdrawLoading] = useState(false);
  const [withdrawError, setWithdrawError] = useState('');
  const [withdrawSuccess, setWithdrawSuccess] = useState('');

  // Загрузка данных кошелька и профиля
  const loadData = async () => {
    setLoading(true);
    try {
      const [profileRes, walletRes] = await Promise.all([
        getProfile(),
        getWallet(),
      ]);
      setUser(profileRes.data.user);
      setWallet(walletRes.data);
    } catch (err) {
      console.error("Ошибка загрузки данных: ", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadData();
  }, []);

  // Обработка пополнения 
  const handleDeposit = async (e) => {
    e.preventDefault();
    setDepositError('');
    setDepositSuccess('');
    setWithdrawError('');
    setWithdrawSuccess('');
    setDepositLoading(true);
    let amount = parseFloat(depositAmount);
    if (isNaN(amount) || amount <= 0) {
      setDepositError('Введите корректную сумму');
      setDepositLoading(false);
      return;
    }
    try {
      const resp = await deposit(amount);
      setDepositSuccess(`Успешно пополнено на ${amount}`);
      setDepositAmount('');
      // Обновляем кошелек и транзакции после пополнения
      await loadData();
    } catch (err) {
      setDepositError(err?.response?.data?.error || 'Ошибка пополнения');
    } finally {
      setDepositLoading(false);
    }
  };

  // Обработка вывода
  const handleWithdraw = async (e) => {
    e.preventDefault();
    setWithdrawError('');
    setWithdrawSuccess('');
    setDepositError('');
    setDepositSuccess('');
    setWithdrawLoading(true);
    let amount = parseFloat(depositAmount);
    if (isNaN(amount) || amount <= 0) {
      setWithdrawError('Введите корректную сумму для вывода');
      setWithdrawLoading(false);
      return;
    }
    try {
      const resp = await withdraw(amount);
      setWithdrawSuccess(`Успешно выведено ${amount}`);
      setDepositAmount('');
      // Обновляем кошелек и транзакции после вывода
      await loadData();
    } catch (err) {
      setWithdrawError(
        err?.response?.data?.error ||
        err?.message ||
        'Ошибка вывода средств'
      );
    } finally {
      setWithdrawLoading(false);
    }
  };

  if (loading) return (
    <div className="transactions-bg-strict">
      <div className="transactions-card"><div>Загрузка...</div></div>
    </div>
  );

  return (
    <div className="transactions-bg-strict">
      <div className="transactions-card">
        <div className="transactions-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h1 className="transactions-title-strict">Транзакции и кошелек</h1>
        {user && (
          <div className="transactions-user">
            <span className="user-username">{user.username}</span>
            <span className="user-email">&lt;{user.email}&gt;</span>
          </div>
        )}
        {wallet && (
          <>
          <div className="transactions-balance-block">
            <span className="transactions-balance-label">Баланс:</span>
            <span className="transactions-balance">{wallet.balance}</span>
            <span className="transactions-currency">₽</span>
          </div>
          {/* --- Deposit & Withdraw Forms --- */}
          <form className="deposit-form-strict" onSubmit={handleDeposit} style={{marginBottom: 0}}>
            <input
              type="number"
              step="0.01"
              placeholder="Сумма"
              value={depositAmount}
              onChange={e => setDepositAmount(e.target.value)}
              min="0"
              disabled={depositLoading || withdrawLoading}
              required
            />
            <button
              type="submit"
              disabled={depositLoading || withdrawLoading}
              style={{ minWidth: 110 }}
            >
              {depositLoading ? "Пополнение..." : "Пополнить"}
            </button>
            <button
              type="button"
              onClick={handleWithdraw}
              disabled={withdrawLoading || depositLoading}
              className="withdraw-btn"
              style={{ minWidth: 110 }}
            >
              {withdrawLoading ? "Вывод..." : "Вывести"}
            </button>
          </form>
          {/* Сообщения об ошибках и успехе */}
          {(depositError || withdrawError) && (
            <div className="form-error">{depositError || withdrawError}</div>
          )}
          {(depositSuccess || withdrawSuccess) && (
            <div className="form-success">{depositSuccess || withdrawSuccess}</div>
          )}
          {/* --- Transactions Table --- */}
          <h3 className="transactions-table-title">Транзакции</h3>
          {wallet.transactions && wallet.transactions.length > 0 ? (
            <div className="transactions-table-wrapper">
              <table className="transactions-table-strict">
                <thead>
                  <tr>
                    <th>Дата</th>
                    <th>Тип</th>
                    <th>Сумма</th>
                  </tr>
                </thead>
                <tbody>
                  {wallet.transactions.map((tx) => (
                    <tr key={tx.id}>
                      <td>{formatDate(tx.created_at)}</td>
                      <td>
                        <span className={
                          tx.type === 'deposit'
                            ? "tx-type tx-deposit"
                            : tx.type === 'bet'
                            ? "tx-type tx-bet"
                            : tx.type === 'withdraw'
                            ? "tx-type tx-withdraw"
                            : "tx-type tx-other"
                        }>
                          {tx.type === 'deposit' ? 'Пополнение'
                            : tx.type === 'bet' ? 'Ставка'
                            : tx.type === 'withdraw' ? 'Вывод'
                            : tx.type}
                        </span>
                      </td>
                      <td>
                        <span className={
                          tx.type === 'deposit'
                            ? "tx-amount tx-amount-deposit"
                            : tx.type === 'bet'
                            ? "tx-amount tx-amount-bet"
                            : tx.type === 'withdraw'
                            ? "tx-amount tx-amount-withdraw"
                            : "tx-amount"
                        }>
                          {tx.amount}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <p className="no-transactions">Транзакций пока нет.</p>
          )}
          </>
        )}
      </div>
    </div>
  );
}