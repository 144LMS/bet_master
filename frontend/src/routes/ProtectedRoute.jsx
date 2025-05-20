import { Navigate } from "react-router-dom";
import { getProfile } from "../api/api";
import { useState, useEffect } from "react";

export default function ProtectedRoute({ children }) {
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
    console.log("ProtectedRoute: checking auth");
      try {
        await getProfile();
        setIsAuthenticated(true);
      } catch (err) {
        setIsAuthenticated(false);
      } finally {
        setLoading(false);
      }
    };
    checkAuth();
  }, []);

  if (loading) return <div>Проверка авторизации...</div>;
  if (!isAuthenticated) return <Navigate to="/login" replace />;

  return children;
}