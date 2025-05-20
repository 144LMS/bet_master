const TransactionList = ({ transactions }) => {
  return (
    <div className="transaction-list">
      <h2>Recent Transactions</h2>
      <ul>
        {transactions.map((tx) => (
          <li key={tx.id}>
            <span>{tx.type}</span>
            <span>${tx.amount}</span>
            <span>{new Date(tx.createdAt).toLocaleString()}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default TransactionList;