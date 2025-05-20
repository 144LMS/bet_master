import { useState, useEffect } from "react";
import { getProfile, getWallet, getMatches, placeBet } from "../api/api";
import { useNavigate } from "react-router-dom";
import "../style/UserProfile.css";

export default function UserProfile() {
  const [user, setUser] = useState(null);
  const [wallet, setWallet] = useState(null);
  const [matches, setMatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [betLoading, setBetLoading] = useState(false);
  const [betError, setBetError] = useState('');
  const [betSuccess, setBetSuccess] = useState('');
  const [bets, setBets] = useState({});

  const navigate = useNavigate();

  useEffect(() => {
    const loadData = async () => {
      try {
        const [profileRes, walletRes, matchesRes] = await Promise.all([
          getProfile(),
          getWallet(),
          getMatches()
        ]);
        setUser(profileRes.data.user);
        setWallet(walletRes.data);
        setMatches((matchesRes.data && matchesRes.data.matches) || []);
      } catch (err) {
        console.error("Ошибка загрузки профиля: ", err);
      } finally {
        setLoading(false);
      }
    };
    loadData();
  }, []);

  const handleBetInput = (matchId, field, value) => {
    setBets(prev => ({
      ...prev,
      [matchId]: {
        ...prev[matchId],
        [field]: value
      }
    }));
  };

  const handlePlaceBet = async (matchId) => {
    setBetError('');
    setBetSuccess('');
    setBetLoading(true);
    const bet = bets[matchId];
    if (!bet || !bet.team || !bet.amount) {
      setBetError("Выберите команду и сумму");
      setBetLoading(false);
      return;
    }
    try {
      const walletId = wallet.id; // или wallet.ID
      const odds = 2; // или получи из матча
      const resp = await placeBet({
        match_id: matchId,
        team: bet.team,
        amount: parseFloat(bet.amount),
        wallet_id: walletId,
        odds: odds
      });
      setBetSuccess("Ставка успешно сделана!");
      setWallet(resp.data.wallet);
    } catch (err) {
      setBetError(err?.response?.data?.error || err.message);
    } finally {
      setBetLoading(false);
    }
  };

  if (loading) return <div className="profile-bg-strict"><div className="profile-card"><div>Загрузка...</div></div></div>;

  return (
    <div className="profile-bg-strict">
      <div className="profile-card">
        <div className="profile-logo">
          <span className="bet-brand">BET<span className="bet-accent">MASTER</span></span>
        </div>
        <h1 className="profile-title-strict">Мой профиль</h1>
        {user && (
          <div className="user-info-strict">
            <p><span>Имя:</span> {user.username}</p>
            <p><span>Email:</span> {user.email}</p>
          </div>
        )}

        {wallet && (
          <div className="wallet-info-strict">
            <p>Баланс: <span className="wallet-balance">{wallet.balance}</span> <span className="wallet-currency">₽</span></p>
            <button className="transactions-btn" onClick={() => navigate('/transactions')}>Транзакции и кошелек</button>
          </div>
        )}

        <h2 className="matches-title-strict">Доступные матчи</h2>
        {matches.filter(m => m.status !== "completed").length === 0 ? (
          <p className="no-matches">Нет доступных матчей</p>
        ) : (
          <div className="matches-table-wrapper">
            <table className="matches-table-strict">
              <thead>
                <tr>
                  <th>Команда 1</th>
                  <th>Команда 2</th>
                  <th>Старт</th>
                  <th>Статус</th>
                  <th>Ставка</th>
                </tr>
              </thead>
              <tbody>
                {matches
                  .filter(match => match.status !== "completed")
                  .map(match => (
                    <tr key={match.ID}>
                      <td>{match.team1}</td>
                      <td>{match.team2}</td>
                      <td>{new Date(match.start_time).toLocaleString()}</td>
                      <td>{match.status}</td>
                      <td>
                        <div className="bet-form-inline">
                          <select
                            value={bets[match.ID]?.team || ""}
                            onChange={e => handleBetInput(match.ID, "team", e.target.value)}
                          >
                            <option value="">Выбрать</option>
                            <option value="team1">{match.team1}</option>
                            <option value="team2">{match.team2}</option>
                          </select>
                          <input
                            type="number"
                            min="1"
                            step="1"
                            value={bets[match.ID]?.amount || ""}
                            onChange={e => handleBetInput(match.ID, "amount", e.target.value)}
                            placeholder="Сумма"
                            style={{width: "70px"}}
                          />
                          <button
                            className="bet-btn"
                            onClick={() => handlePlaceBet(match.ID)}
                            disabled={betLoading}
                          >
                            Поставить
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
        )}
        {betError && <div className="form-error">{betError}</div>}
        {betSuccess && <div className="form-success">{betSuccess}</div>}
      </div>
    </div>
  );
}