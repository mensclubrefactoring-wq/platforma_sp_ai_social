import { ReactNode, useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { LayoutDashboard, Users, BookOpen, Settings, Bell, Search, Sparkles, LogOut, Menu, Plus } from "lucide-react";
import { cn } from "@/src/lib/utils";

interface ShellProps {
  children: ReactNode;
}

export default function Shell({ children }: ShellProps) {
  const location = useLocation();
  const [user, setUser] = useState<any>(null);

  useEffect(() => {
    const data = localStorage.getItem("user_data");
    if (data) setUser(JSON.parse(data));
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("auth_token");
    localStorage.removeItem("user_data");
    window.dispatchEvent(new Event("auth-change"));
  };

  const NAV_ITEMS = [
    { label: "Дашборд", path: "/app", icon: LayoutDashboard },
    { label: "AI Помощник", path: "/app/assistant", icon: Sparkles },
    { label: "Исполнители", path: "/app/executors", icon: Users },
    { label: "База знаний", path: "/app/knowledge", icon: BookOpen },
    { label: "Настройки", path: "/app/settings", icon: Settings },
  ];

  return (
    <div className="min-h-screen bg-[#fcfcfd] flex">
      {/* Sidebar */}
      <aside className="w-64 border-r border-gray-100 bg-white hidden lg:flex flex-col fixed h-full z-20">
        <div className="p-8 flex items-center gap-2 font-bold text-xl tracking-tight text-indigo-600 mb-8">
          <Sparkles className="w-6 h-6" />
          <span>Platforma SP</span>
        </div>
        
        <nav className="flex-1 px-4 space-y-1">
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={cn(
                "flex items-center gap-3 px-4 py-3 rounded-xl transition-all font-semibold text-sm",
                location.pathname === item.path 
                  ? "bg-indigo-50 text-indigo-600" 
                  : "text-gray-400 hover:text-gray-900 hover:bg-gray-50"
              )}
            >
              <item.icon className="w-4 h-4" />
              {item.label}
            </Link>
          ))}
          
          <div className="pt-8 px-4">
             <Link to="/app/create-task" className="flex items-center justify-center gap-2 w-full bg-indigo-600 text-white font-bold py-3 rounded-xl shadow-lg shadow-indigo-100 hover:bg-indigo-700 transition-all text-xs">
              <Plus className="w-3 h-3" /> Новая задача
            </Link>
          </div>
        </nav>

        <div className="p-4 border-t border-gray-50">
          <div className="flex items-center gap-3 p-3 bg-gray-50 rounded-2xl mb-4">
            <div className="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center text-indigo-600 font-bold uppercase">
              {user?.displayName?.[0] || user?.email?.[0] || "U"}
            </div>
            <div className="flex-1 min-w-0">
              <p className="text-xs font-bold text-gray-900 truncate">{user?.displayName || user?.email?.split('@')[0]}</p>
              <p className="text-[10px] text-gray-400 truncate">{user?.email}</p>
            </div>
          </div>
          <button 
            onClick={handleLogout}
            className="w-full flex items-center gap-3 px-4 py-3 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-xl transition-all font-semibold text-sm"
          >
            <LogOut className="w-4 h-4" />
            Выйти
          </button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 lg:ml-64 min-h-screen">
        {/* Top Navbar */}
        <header className="h-20 bg-white/10 backdrop-blur-sm border-b border-gray-100/10 flex items-center justify-between px-8 sticky top-0 z-10">
          <div className="flex lg:hidden items-center gap-4">
            <button className="p-2 bg-white border border-gray-100 rounded-lg"><Menu className="w-5 h-5 text-gray-500" /></button>
            <span className="font-bold text-indigo-600">Platforma SP</span>
          </div>
          
          <div className="hidden md:flex items-center gap-4 bg-white border border-gray-100 px-4 py-2 rounded-xl w-96 shadow-sm">
            <Search className="w-4 h-4 text-gray-300" />
            <input type="text" placeholder="Поиск задач или исполнителей..." className="bg-transparent border-none text-sm outline-none w-full" />
          </div>

          <div className="flex items-center gap-4">
            <div className="h-10 w-10 flex items-center justify-center bg-white border border-gray-100 rounded-xl relative hover:border-indigo-100 transition-colors cursor-pointer">
              <Bell className="w-5 h-5 text-gray-400" />
              <span className="absolute top-2 right-2 w-2 h-2 bg-red-500 rounded-full border-2 border-white" />
            </div>
          </div>
        </header>

        <div className="p-8 max-w-7xl mx-auto">
          {children}
        </div>
      </main>
    </div>
  );
}
