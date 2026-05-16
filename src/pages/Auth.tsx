import React, { useState } from "react";
import { motion } from "motion/react";
import { useNavigate } from "react-router-dom";
import { Building2, HandHeart, LogIn, UserPlus, Sparkles, AlertCircle } from "lucide-react";
import { cn } from "@/src/lib/utils";
import { apiFetch } from "@/src/lib/api";

export default function Auth() {
  const [isLogin, setIsLogin] = useState(true);
  const [role, setRole] = useState<"business" | "entrepreneur">("business");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [representativeName, setRepresentativeName] = useState("");
  const [companyName, setCompanyName] = useState("");
  const [phone, setPhone] = useState("");
  const [consent, setConsent] = useState(false);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  
  const navigate = useNavigate();

  const handleAuth = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const endpoint = isLogin ? "/api/auth/login" : "/api/auth/register";
      const payload = isLogin ? { email, password } : { 
        email, 
        password, 
        representativeName, 
        companyName, 
        role, 
        phone, 
        consentGiven: consent 
      };
      
      const data = await apiFetch(endpoint, {
        method: "POST",
        body: JSON.stringify(payload),
      });

      localStorage.setItem("auth_token", data.token);
      localStorage.setItem("user_data", JSON.stringify(data.user));
      
      // Вызываем событие, чтобы App.tsx узнал об изменениях (простой способ для MVP)
      window.dispatchEvent(new Event("auth-change"));
      navigate("/app");
    } catch (err: any) {
      setError(err.message || "Ошибка авторизации");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-[#f8faff] flex items-center justify-center p-4">
      <motion.div 
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        className="max-w-md w-full bg-white rounded-3xl shadow-2xl shadow-indigo-100 border border-gray-100 overflow-hidden"
      >
        <div className="p-8">
          <div className="flex justify-center mb-6 text-indigo-600">
            <Sparkles className="w-10 h-10" />
          </div>
          <h2 className="text-2xl font-bold text-center mb-2 serif">
            {isLogin ? "С возвращением" : "Создать аккаунт"}
          </h2>
          <p className="text-gray-500 text-center text-sm mb-8">
            {isLogin ? "Локальный вход Platforma SP" : "Регистрация в локальной сети социального партнерства"}
          </p>

          {!isLogin && (
            <div className="grid grid-cols-2 gap-3 mb-6">
              <button
                type="button"
                onClick={() => setRole("business")}
                className={cn(
                  "flex flex-col items-center gap-2 p-4 rounded-2xl border transition-all text-xs font-bold",
                  role === "business" ? "border-indigo-600 bg-indigo-50 text-indigo-600" : "border-gray-100 text-gray-400"
                )}
              >
                <Building2 className="w-5 h-5" />
                Крупный бизнес
              </button>
              <button
                type="button"
                onClick={() => setRole("entrepreneur")}
                className={cn(
                  "flex flex-col items-center gap-2 p-4 rounded-2xl border transition-all text-xs font-bold",
                  role === "entrepreneur" ? "border-indigo-600 bg-indigo-50 text-indigo-600" : "border-gray-100 text-gray-400"
                )}
              >
                <HandHeart className="w-5 h-5" />
                Предприниматель
              </button>
            </div>
          )}

          <form onSubmit={handleAuth} className="space-y-4">
            {!isLogin && (
              <>
                <div>
                  <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-1 block">ФИО представителя</label>
                  <input 
                    required
                    type="text" 
                    value={representativeName}
                    onChange={(e) => setRepresentativeName(e.target.value)}
                    placeholder="Иванов Иван Иванович" 
                    className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
                  />
                </div>
                <div>
                  <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-1 block">Название компании</label>
                  <input 
                    required
                    type="text" 
                    value={companyName}
                    onChange={(e) => setCompanyName(e.target.value)}
                    placeholder="ООО «Социальный проект»" 
                    className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
                  />
                </div>
              </>
            )}
            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-1 block">Email</label>
              <input 
                required
                type="email" 
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@company.com" 
                className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
              />
            </div>
            
            {!isLogin && (
              <div>
                <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-1 block">Телефон</label>
                <input 
                  required
                  type="text" 
                  value={phone}
                  onChange={(e) => setPhone(e.target.value)}
                  placeholder="+7 (999) 000-00-00" 
                  className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
                />
              </div>
            )}

            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-1 block">Пароль</label>
              <input 
                required
                type="password" 
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••" 
                className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
              />
            </div>

            {!isLogin && (
              <label className="flex items-start gap-3 cursor-pointer group mt-4">
                <input 
                  type="checkbox" 
                  checked={consent}
                  onChange={(e) => setConsent(e.target.checked)}
                  required
                  className="mt-1 accent-indigo-600"
                />
                <span className="text-[11px] text-gray-500 font-medium leading-relaxed">
                  Я подтверждаю, что являюсь реальным человеком и даю согласие на обработку персональных данных.
                </span>
              </label>
            )}

            {error && (
              <div className="flex items-center gap-2 text-red-600 bg-red-50 p-3 rounded-xl text-xs font-bold">
                <AlertCircle className="w-4 h-4" />
                {error}
              </div>
            )}

            <button 
              disabled={loading}
              type="submit"
              className="w-full bg-indigo-600 text-white font-bold py-3 rounded-xl shadow-lg shadow-indigo-100 hover:bg-indigo-700 transition-all flex items-center justify-center gap-2"
            >
              {loading ? "Загрузка..." : isLogin ? <><LogIn className="w-4 h-4" /> Войти</> : <><UserPlus className="w-4 h-4" /> Регистрация</>}
            </button>
          </form>

          <p className="mt-8 text-center text-xs font-bold text-gray-400 uppercase tracking-widest">
            {isLogin ? "Нет аккаунта?" : "Уже зарегистрированы?"}
            <button 
              onClick={() => setIsLogin(!isLogin)}
              className="text-indigo-600 ml-2 hover:underline"
            >
              {isLogin ? "Создать здесь" : "Войти здесь"}
            </button>
          </p>
        </div>
      </motion.div>
    </div>
  );
}
