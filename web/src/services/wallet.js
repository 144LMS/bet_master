export const getBalance = async (userId, token) => {
  const res = await fetch(`/wallets/${userId}/balance`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return await res.json();
};

export const getTransactions = async (userId, token) => {
  const res = await fetch(`/wallets/${userId}/transactions`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  return await res.json();
};

export const deposit = async (userId, amount, token) => {
  const res = await fetch(`/wallets/${userId}/deposit`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ amount }),
  });
  return await res.json();
};

export const withdraw = async (userId, amount, token) => {
  const res = await fetch(`/wallets/${userId}/withdraw`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ amount }),
  });
  return await res.json();
};