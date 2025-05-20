const BalanceCard = ({ balance }) => {
  return (
    <div className="balance-card">
      <h2>Your Balance</h2>
      <p>${balance.toFixed(2)}</p>
    </div>
  );
};

export default BalanceCard;