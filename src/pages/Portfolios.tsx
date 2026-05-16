import { useState, useEffect } from "react";
import { motion } from "motion/react";
import { User, MapPin, Briefcase, ExternalLink, ShieldCheck, Loader2 } from "lucide-react";
import { apiFetch } from "@/src/lib/api";

interface Entrepreneur {
  id: number;
  representativeName: string;
  companyName: string;
  email: string;
  phone: string;
  role: string;
  portfolioUrl: string;
}

export default function Portfolios() {
  const [entrepreneurs, setEntrepreneurs] = useState<Entrepreneur[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchPortfolios = async () => {
      try {
        const data = await apiFetch("/api/admin/portfolios");
        setEntrepreneurs(data);
      } catch (err) {
        console.error("Failed to fetch portfolios", err);
      } finally {
        setLoading(false);
      }
    };
    fetchPortfolios();
  }, []);

  if (loading) {
    return (
      <div className="h-96 flex flex-col items-center justify-center text-gray-400 gap-3">
        <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
        <p className="text-xs font-bold uppercase tracking-widest">Загружаем профили...</p>
      </div>
    );
  }

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      <div>
        <h1 className="text-3xl font-bold tracking-tight text-gray-900 serif">Портфолио предпринимателей</h1>
        <p className="text-gray-500 mt-1">Список верифицированных социальных предпринимателей платформы.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {entrepreneurs.length === 0 ? (
          <div className="col-span-full p-12 text-center text-gray-400 serif text-xl">Профилей пока нет.</div>
        ) : entrepreneurs.map((e) => (
          <motion.div 
            key={e.id}
            whileHover={{ y: -5 }}
            className="bg-white p-6 rounded-3xl border border-gray-100 shadow-xl shadow-indigo-50/50 relative overflow-hidden"
          >
            <div className="flex items-center gap-4 mb-6">
              <div className="w-16 h-16 bg-indigo-50 rounded-2xl flex items-center justify-center text-indigo-600">
                <User className="w-8 h-8" />
              </div>
              <div>
                <h3 className="font-bold text-lg text-gray-900">{e.companyName}</h3>
                <p className="text-xs text-gray-400 font-medium">{e.representativeName}</p>
                <div className="flex items-center gap-1.5 text-[10px] font-bold text-green-600 uppercase tracking-widest mt-1">
                  <ShieldCheck className="w-3.5 h-3.5" /> Верифицирован
                </div>
              </div>
            </div>

            <div className="space-y-3 mb-6">
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <Briefcase className="w-4 h-4" />
                Социальный предприниматель
              </div>
              <div className="flex items-center gap-2 text-sm text-gray-500">
                <MapPin className="w-4 h-4" />
                Россия, Москва
              </div>
            </div>

            <div className="pt-4 border-t border-gray-50 flex items-center justify-between">
              <span className="text-xs font-bold text-gray-400 uppercase tracking-widest">Телефон: {e.phone}</span>
              {e.portfolioUrl ? (
                <a 
                  href={e.portfolioUrl} 
                  target="_blank" 
                  rel="noreferrer"
                  className="p-2 bg-indigo-50 text-indigo-600 rounded-lg hover:bg-indigo-100 transition-colors"
                >
                  <ExternalLink className="w-4 h-4" />
                </a>
              ) : (
                <span className="text-[10px] bg-gray-50 text-gray-400 px-2 py-1 rounded-md">Нет ссылки</span>
              )}
            </div>
          </motion.div>
        ))}
      </div>
    </div>
  );
}
