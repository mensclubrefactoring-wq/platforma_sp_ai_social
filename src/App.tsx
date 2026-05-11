/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import Landing from "./pages/Landing";
import Dashboard from "./pages/Dashboard";
import Assistant from "./pages/Assistant";
import Auth from "./pages/Auth";
import CreateTask from "./pages/CreateTask";
import Shell from "./components/layout/Shell";
import { useEffect, useState } from "react";
import { apiFetch } from "./lib/api";

export default function App() {
  const [user, setUser] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  const checkAuth = async () => {
    const token = localStorage.getItem("auth_token");
    if (!token) {
      setUser(null);
      setLoading(false);
      return;
    }

    try {
      const userData = await apiFetch("/api/auth/me");
      setUser(userData);
    } catch (err) {
      localStorage.removeItem("auth_token");
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkAuth();
    
    // Слушаем событие изменения авторизации
    window.addEventListener("auth-change", checkAuth);
    return () => window.removeEventListener("auth-change", checkAuth);
  }, []);

  if (loading) return <div className="h-screen w-screen flex items-center justify-center bg-[#f8faff]"><div className="w-8 h-8 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin" /></div>;

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/auth" element={!user ? <Auth /> : <Navigate to="/app" />} />
        
        <Route path="/app" element={user ? <Shell children={<Dashboard />} /> : <Navigate to="/auth" />} />
        <Route path="/app/assistant" element={user ? <Shell children={<Assistant />} /> : <Navigate to="/auth" />} />
        <Route path="/app/create-task" element={user ? <Shell children={<CreateTask />} /> : <Navigate to="/auth" />} />
        
        {/* Placeholder routes */}
        <Route path="/app/executors" element={user ? <Shell children={<div className="p-8 text-center text-gray-500 italic serif text-xl">Раздел в разработке</div>} /> : <Navigate to="/auth" />} />
        <Route path="/app/knowledge" element={user ? <Shell children={<div className="p-8 text-center text-gray-500 italic serif text-xl">Раздел в разработке</div>} /> : <Navigate to="/auth" />} />
        <Route path="/app/settings" element={user ? <Shell children={<div className="p-8 text-center text-gray-500 italic serif text-xl">Раздел в разработке</div>} /> : <Navigate to="/auth" />} />
        
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
