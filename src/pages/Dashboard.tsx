import { useState, useEffect } from "react";
import { motion } from "motion/react";
import { Plus, Search, Filter, MoreVertical, Building2, MapPin, Calendar, CreditCard, ChevronRight, MessageSquareCode, Sparkles, Users, Loader2 } from "lucide-react";
import { cn } from "@/src/lib/utils";
import { Link } from "react-router-dom";
import { apiFetch } from "@/src/lib/api";

interface Task {
  id: string;
  title: string;
  status: string;
  budget: string;
  deadline: string;
  location: string;
  category: string;
  creatorId: string;
}

export default function Dashboard() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState<any>(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [categoryFilter, setCategoryFilter] = useState("Все");

  const fetchTasks = async (search = "", category = "Все") => {
    setLoading(true);
    try {
      const query = new URLSearchParams();
      if (search) query.append("search", search);
      if (category !== "Все") query.append("category", category);
      
      const t = await apiFetch(`/api/tasks?${query.toString()}`);
      setTasks(t);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const u = await apiFetch("/api/auth/me");
        setUser(u);
        fetchTasks();
      } catch (err) {
        console.error(err);
      }
    };
    fetchData();
  }, []);

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchTasks(searchTerm, categoryFilter);
  };

  const categories = ["Все", "Экология", "Образование", "Социальное жилье", "Обучение ИТ", "Помощь пожилым", "Инклюзивность"];

  return (
    <div className="space-y-8 animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-gray-900 serif">Рабочий стол {user?.role === 'business' ? 'компании' : 'партнера'}</h1>
          <p className="text-gray-500 mt-1">Управляйте социальными задачами и находите партнеров.</p>
        </div>
        <div className="flex gap-3">
          <Link to="/app/assistant" className="group flex items-center gap-2 bg-indigo-50 border border-indigo-200 text-indigo-700 px-5 py-2.5 rounded-xl font-bold hover:bg-indigo-100 transition-all">
            <Sparkles className="w-4 h-4" />
            AI Помощник
          </Link>
          <Link to="/app/create-task" className="flex items-center gap-2 bg-indigo-600 text-white px-5 py-2.5 rounded-xl font-bold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100">
            <Plus className="w-4 h-4" />
            Создать задачу
          </Link>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {[
          { label: "Всего задач", value: tasks.length.toString(), sub: "На платформе", icon: Building2 },
          { label: "Ваши отклики", value: "0", sub: "В ожидании", icon: Users },
          { label: "Доступный бюджет", value: "≈1.2M ₽", sub: "Суммарно", icon: CreditCard },
        ].map((stat, idx) => (
          <div key={idx} className="bg-white p-6 rounded-3xl border border-gray-100 shadow-sm border-b-4 border-b-indigo-500">
            <div className="flex items-center justify-between mb-4">
              <div className="p-2 bg-indigo-50 rounded-xl text-indigo-600">
                <stat.icon className="w-5 h-5" />
              </div>
              <span className="text-[10px] font-bold text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded-md uppercase tracking-wider">{stat.sub}</span>
            </div>
            <div className="text-sm font-medium text-gray-400 mb-1">{stat.label}</div>
            <div className="text-3xl font-bold text-gray-900">{stat.value}</div>
          </div>
        ))}
      </div>

      {/* Tasks List */}
      <div className="bg-white rounded-3xl border border-gray-100 shadow-xl shadow-indigo-50/50 overflow-hidden">
        <div className="p-8 border-b border-gray-50 flex items-center justify-between bg-white">
          <div className="flex items-center gap-6">
            <h2 className="font-bold text-xl serif">Маркетплейс задач</h2>
            <div className="flex bg-gray-50 p-1.5 rounded-xl overflow-x-auto max-w-[500px] hide-scrollbar">
              {categories.map((cat) => (
                <button 
                  key={cat}
                  onClick={() => {
                    setCategoryFilter(cat);
                    fetchTasks(searchTerm, cat);
                  }}
                  className={cn(
                    "px-4 py-1.5 text-xs font-bold whitespace-nowrap rounded-lg transition-all",
                    categoryFilter === cat ? "bg-white shadow-sm border border-gray-200 text-indigo-600" : "text-gray-400 hover:text-gray-900"
                  )}
                >
                  {cat}
                </button>
              ))}
            </div>
          </div>
          <form onSubmit={handleSearch} className="flex gap-2">
            <div className="relative">
              <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
              <input 
                type="text"
                placeholder="Поиск..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="bg-gray-50 border-none rounded-xl py-2 pl-9 pr-4 text-xs w-48 focus:ring-2 focus:ring-indigo-100"
              />
            </div>
            <button type="submit" className="p-2.5 text-indigo-600 bg-indigo-50 rounded-xl transition-all hover:bg-indigo-100">
              <Filter className="w-5 h-5" />
            </button>
          </form>
        </div>

        <div className="divide-y divide-gray-50">
          {loading ? (
            <div className="p-12 flex flex-col items-center justify-center text-gray-400 gap-3">
              <Loader2 className="w-8 h-8 animate-spin text-indigo-600" />
              <p className="text-xs font-bold uppercase tracking-widest">Загружаем задачи...</p>
            </div>
          ) : tasks.length === 0 ? (
            <div className="p-12 text-center text-gray-400 serif text-lg lowercase">Задач пока нет. Будьте первыми!</div>
          ) : tasks.map((task) => (
            <motion.div 
              key={task.id}
              whileHover={{ backgroundColor: "rgba(249, 250, 251, 0.8)" }}
              className="p-8 transition-all flex items-center justify-between group cursor-pointer border-l-4 border-l-transparent hover:border-l-indigo-500"
            >
              <div className="flex items-start gap-6">
                <div className={cn(
                  "w-14 h-14 rounded-2xl flex items-center justify-center shadow-inner",
                  task.status === "active" ? "bg-green-50 text-green-600" : "bg-gray-100 text-gray-500"
                )}>
                  <MessageSquareCode className="w-7 h-7" />
                </div>
                <div>
                  <div className="flex items-center gap-4 mb-2">
                    <h3 className="font-bold text-xl text-gray-900 group-hover:text-indigo-600 transition-colors">{task.title}</h3>
                    <span className={cn(
                      "text-[10px] font-bold uppercase tracking-widest px-2.5 py-1 rounded-full",
                      task.status === "active" ? "bg-green-100 text-green-700" : "bg-gray-200 text-gray-600"
                    )}>
                      {task.status}
                    </span>
                  </div>
                  <div className="flex flex-wrap items-center gap-y-2 gap-x-6 text-xs font-bold text-gray-400 uppercase tracking-wider">
                    <span className="flex items-center gap-1.5"><MapPin className="w-3.5 h-3.5" /> {task.location}</span>
                    <span className="flex items-center gap-1.5"><Calendar className="w-3.5 h-3.5" /> До {task.deadline}</span>
                    <span className="text-indigo-600 bg-indigo-50 px-3 py-1 rounded-lg">{task.budget}</span>
                  </div>
                </div>
              </div>
              <div className="flex items-center gap-4">
                <button 
                  onClick={(e) => {
                    e.stopPropagation();
                    // Функция генерации
                    alert("Генерируем предложение...");
                    apiFetch("/api/ai/generate-proposal", {
                      method: "POST",
                      body: JSON.stringify(task)
                    }).then(res => alert("Предложение сформировано: \n\n" + res.proposal))
                      .catch(err => alert("Ошибка: " + err.message));
                  }}
                  className="flex items-center gap-2 px-4 py-2 bg-indigo-50 text-indigo-600 rounded-xl font-bold text-xs hover:bg-indigo-100 transition-all"
                >
                  <Sparkles className="w-3 h-3" /> Сформировать предложение
                </button>
                <button className="p-2.5 text-gray-300 hover:text-indigo-600 hover:bg-white rounded-xl transition-all shadow-sm">
                  <MoreVertical className="w-5 h-5" />
                </button>
                <div className="w-10 h-10 rounded-xl border border-gray-100 flex items-center justify-center bg-white shadow-xl shadow-indigo-100/50">
                  <ChevronRight className="w-5 h-5 text-gray-400" />
                </div>
              </div>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
}
