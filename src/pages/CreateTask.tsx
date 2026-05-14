import React, { useState } from "react";
import { motion } from "motion/react";
import { useNavigate } from "react-router-dom";
import { Sparkles, Send, Loader2, ArrowLeft, Target, Calendar, CreditCard, MapPin } from "lucide-react";
import { cn } from "@/src/lib/utils";
import { apiFetch } from "@/src/lib/api";
import { askAI } from "@/src/services/ai";

export default function CreateTask() {
  const [title, setTitle] = useState("");
  const [description, setDescription] = useState("");
  const [budget, setBudget] = useState("");
  const [deadline, setDeadline] = useState("");
  const [location, setLocation] = useState("");
  const [category, setCategory] = useState("Экология");
  const [isAIOptimizing, setIsAIOptimizing] = useState(false);
  const [isClassifying, setIsClassifying] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const navigate = useNavigate();

  const handleOptimizeDescription = async () => {
    if (!description.trim()) return;
    setIsAIOptimizing(true);
    try {
      const response = await askAI(`Оптимизируй описание социальной задачи для бизнеса так, чтобы оно было привлекательным для социальных предпринимателей. Сохрани суть, но сделай его структурированным (Цель, Ожидаемый эффект, Требования). \nТекст: ${description}`);
      setDescription(response || description);
    } catch (err) {
      console.error(err);
    } finally {
      setIsAIOptimizing(false);
    }
  };

  const handleAutoClassify = async () => {
    if (!description.trim()) return;
    setIsClassifying(true);
    try {
      const res = await apiFetch("/api/ai/classify", {
        method: "POST",
        body: JSON.stringify({ description })
      });
      if (res.category) setCategory(res.category);
    } catch (err) {
      console.error(err);
    } finally {
      setIsClassifying(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    setIsSubmitting(true);
    try {
      await apiFetch("/api/tasks", {
        method: "POST",
        body: JSON.stringify({
          title,
          description,
          budget,
          deadline,
          location,
          category
        })
      });
      navigate("/app");
    } catch (err) {
      console.error(err);
      alert("Ошибка при сохранении задачи");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="max-w-3xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="flex items-center gap-4">
        <button onClick={() => navigate("/app")} className="p-2 hover:bg-gray-100 rounded-full transition-colors">
          <ArrowLeft className="w-5 h-5 text-gray-400" />
        </button>
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-gray-900 serif">Новая задача</h1>
          <p className="text-gray-500">Опишите проект, чтобы найти идеального социального партнера.</p>
        </div>
      </div>

      <form onSubmit={handleSubmit} className="bg-white rounded-3xl border border-gray-100 shadow-xl shadow-indigo-50 overflow-hidden">
        <div className="p-8 space-y-6">
          {/* Title */}
          <div>
            <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-2 block">Название задачи</label>
            <input 
              required
              type="text" 
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Например: Закупка сувениров из переработанного пластика" 
              className="w-full bg-gray-50 border-none rounded-2xl py-4 px-6 text-lg font-semibold focus:ring-2 focus:ring-indigo-100 transition-all shadow-inner"
            />
          </div>

          {/* Details Grid */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-2 block flex items-center gap-2">
                <CreditCard className="w-3 h-3" /> Бюджет
              </label>
              <input 
                type="text" 
                value={budget}
                onChange={(e) => setBudget(e.target.value)}
                placeholder="Например: 50,000 ₽" 
                className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
              />
            </div>
            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-2 block flex items-center gap-2">
                <Calendar className="w-3 h-3" /> Срок подачи
              </label>
              <input 
                type="text" 
                value={deadline}
                onChange={(e) => setDeadline(e.target.value)}
                placeholder="20 Июня, 2026" 
                className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
              />
            </div>
            <div>
              <label className="text-xs font-bold text-gray-400 uppercase tracking-widest pl-1 mb-2 block flex items-center gap-2">
                <MapPin className="w-3 h-3" /> География
              </label>
              <input 
                type="text" 
                value={location}
                onChange={(e) => setLocation(e.target.value)}
                placeholder="Москва / Дистанционно" 
                className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100"
              />
            </div>
            <div>
              <div className="flex items-center justify-between mb-2 px-1">
                <label className="text-xs font-bold text-gray-400 uppercase tracking-widest block flex items-center gap-2">
                  <Target className="w-3 h-3" /> Категория
                </label>
                <button 
                  type="button"
                  onClick={handleAutoClassify}
                  disabled={isClassifying || !description.trim()}
                  className="text-[9px] font-bold text-indigo-500 hover:text-indigo-700 transition-all uppercase tracking-widest"
                >
                  {isClassifying ? "Классификация..." : "Определить по тексту"}
                </button>
              </div>
              <select 
                value={category}
                onChange={(e) => setCategory(e.target.value)}
                className="w-full bg-gray-50 border-none rounded-xl py-3 px-4 text-sm focus:ring-2 focus:ring-indigo-100 appearance-none"
              >
                {["Экология", "Образование", "Социальное жилье", "Обучение ИТ", "Помощь пожилым", "Инклюзивность"].map(c => (
                  <option key={c} value={c}>{c}</option>
                ))}
              </select>
            </div>
          </div>

          {/* Description with AI */}
          <div>
            <div className="flex items-center justify-between mb-2 px-1">
              <label className="text-xs font-bold text-gray-400 uppercase tracking-widest block">Подробное описание</label>
              <button 
                type="button"
                onClick={handleOptimizeDescription}
                disabled={isAIOptimizing || !description.trim()}
                className="flex items-center gap-1.5 text-[10px] font-bold text-indigo-600 bg-indigo-50 px-3 py-1 rounded-full hover:bg-indigo-100 disabled:opacity-50 transition-all uppercase tracking-wider"
              >
                {isAIOptimizing ? <Loader2 className="w-3 h-3 animate-spin" /> : <Sparkles className="w-3 h-3" />}
                Оптимизировать через AI
              </button>
            </div>
            <textarea 
              required
              rows={8}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Опишите цель проекта, целевую аудиторию и требования к исполнителю..." 
              className="w-full bg-gray-50 border-none rounded-2xl py-4 px-6 text-sm focus:ring-2 focus:ring-indigo-100 transition-all shadow-inner resize-none"
            />
          </div>
        </div>

        <div className="p-8 bg-gray-50 border-t border-gray-100 flex items-center justify-between">
          <p className="text-xs text-gray-400 font-medium max-w-[200px]">Ваша задача будет опубликована в открытом доступе для всех партнеров.</p>
          <button 
            disabled={isSubmitting}
            type="submit"
            className="flex items-center gap-2 bg-indigo-600 text-white font-bold px-8 py-4 rounded-2xl shadow-xl shadow-indigo-100 hover:bg-indigo-700 disabled:opacity-50 transition-all"
          >
            {isSubmitting ? <Loader2 className="w-5 h-5 animate-spin" /> : <><Send className="w-4 h-4" /> Опубликовать задачу</>}
          </button>
        </div>
      </form>
    </div>
  );
}
