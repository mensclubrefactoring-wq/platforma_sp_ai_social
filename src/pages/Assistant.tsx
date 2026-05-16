import { useState, useRef, useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import { Send, Sparkles, BrainCircuit, Bot, User, Loader2, ArrowLeft, CheckCircle2 } from "lucide-react";
import { Link } from "react-router-dom";
import { cn } from "@/src/lib/utils";
import { askAI } from "@/src/services/ai";
import { apiFetch } from "@/src/lib/api";

interface Message {
  role: "assistant" | "user";
  content: string;
}

export default function Assistant() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchHistory = async () => {
      try {
        const u = await apiFetch("/api/auth/me");
        const history = await apiFetch("/api/ai/history");
        if (history && history.length > 0) {
          setMessages(history.map((m: any) => ({
            role: m.role,
            content: m.content
          })));
        } else {
          const welcome = u.role === "business" 
            ? "Здравствуйте! Я ваш AI-ассистент. Я помогу вам сформулировать социальную задачу для бизнеса так, чтобы она была понятна социальным предпринимателям. Что вы планируете реализовать?"
            : "Здравствуйте! Я ваш AI-ассистент. Я помогу вам подготовить качественное предложение для бизнеса и оптимизировать описание ваших социальных проектов. Какой проект вы сейчас развиваете?";
          
          setMessages([{
            role: "assistant",
            content: welcome
          }]);
        }
      } catch (err) {
        console.error("Failed to fetch AI history", err);
      }
    };
    fetchHistory();
  }, []);

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSend = async () => {
    if (!input.trim() || isLoading) return;

    const userMessage = input.trim();
    setInput("");
    setMessages(prev => [...prev, { role: "user", content: userMessage }]);
    setIsLoading(true);

    try {
      const response = await askAI(userMessage, messages.map(m => ({
        role: m.role === "assistant" ? "assistant" : "user",
        content: m.content
      })));

      setMessages(prev => [...prev, { role: "assistant", content: response }]);
    } catch (err) {
      console.error(err);
      setMessages(prev => [...prev, { role: "assistant", content: "Произошла ошибка при обработке запроса. Попробуйте еще раз." }]);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col h-[calc(100vh-120px)] max-w-4xl mx-auto shadow-2xl shadow-indigo-100 rounded-3xl overflow-hidden border border-gray-100 bg-white">
      {/* Assistant Header */}
      <div className="bg-white border-b border-gray-50 px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link to="/app" className="p-2 hover:bg-gray-50 rounded-full transition-colors mr-1">
            <ArrowLeft className="w-5 h-5 text-gray-500" />
          </Link>
          <div className="w-10 h-10 bg-indigo-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-indigo-200">
            <BrainCircuit className="w-6 h-6" />
          </div>
          <div>
            <div className="font-bold text-gray-900 leading-tight">AI Конструктор Задач</div>
            <div className="text-[10px] uppercase font-bold tracking-widest text-indigo-600">Online • Аналитическая модель</div>
          </div>
        </div>
        <div className="flex items-center gap-2 px-3 py-1 bg-green-50 rounded-full border border-green-100">
          <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
          <span className="text-[10px] font-bold text-green-700 uppercase tracking-wider">Ready</span>
        </div>
      </div>

      {/* Messages */}
      <div ref={scrollRef} className="flex-1 overflow-y-auto p-6 space-y-6 scroll-smooth bg-gray-50/30">
        {messages.map((message, idx) => (
          <motion.div 
            key={idx}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            className={cn(
              "flex items-end gap-3",
              message.role === "user" ? "flex-row-reverse" : "flex-row"
            )}
          >
            <div className={cn(
              "w-8 h-8 rounded-full flex items-center justify-center",
              message.role === "assistant" ? "bg-indigo-600 text-white" : "bg-white border text-gray-400"
            )}>
              {message.role === "assistant" ? <Bot className="w-5 h-5" /> : <User className="w-5 h-5" />}
            </div>
            <div className={cn(
              "max-w-[80%] p-4 rounded-2xl text-sm leading-relaxed shadow-sm",
              message.role === "assistant" 
                ? "bg-white text-gray-800 rounded-bl-none border border-gray-100" 
                : "bg-indigo-600 text-white rounded-br-none"
            )}>
              {message.content.split('\n').map((line, i) => (
                <p key={i} className={line ? "mb-2 last:mb-0" : "h-2"}>{line}</p>
              ))}
            </div>
          </motion.div>
        ))}
        {isLoading && (
          <div className="flex items-end gap-3">
            <div className="w-8 h-8 rounded-full bg-indigo-600 text-white flex items-center justify-center">
              <Bot className="w-5 h-5" />
            </div>
            <div className="bg-white p-4 rounded-2xl rounded-bl-none border border-gray-100 shadow-sm">
              <Loader2 className="w-5 h-5 animate-spin text-indigo-600" />
            </div>
          </div>
        )}
      </div>

      {/* Input Area */}
      <div className="p-6 bg-white border-t border-gray-50">
        <div className="relative group">
          <textarea 
            rows={1}
            value={input}
            onChange={(e) => setInput(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter" && !e.shiftKey) {
                e.preventDefault();
                handleSend();
              }
            }}
            placeholder="Опишите вашу социальную инициативу..."
            className="w-full bg-gray-50 border-none rounded-2xl py-4 pl-4 pr-16 text-sm focus:ring-2 focus:ring-indigo-100 resize-none transition-all"
          />
          <button 
            onClick={handleSend}
            disabled={isLoading || !input.trim()}
            className="absolute right-2 top-2 p-3 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 disabled:opacity-50 disabled:grayscale transition-all shadow-lg shadow-indigo-100"
          >
            <Send className="w-4 h-4" />
          </button>
        </div>
        <div className="mt-3 flex items-center justify-center gap-4 text-[10px] text-gray-400 font-bold uppercase tracking-widest">
          <span className="flex items-center gap-1"><CheckCircle2 className="w-3 h-3" /> Авто-структурирование</span>
          <span className="flex items-center gap-1"><CheckCircle2 className="w-3 h-3" /> NLP Анализ</span>
          <span className="flex items-center gap-1"><CheckCircle2 className="w-3 h-3" /> Контекстная память</span>
        </div>
      </div>
    </div>
  );
}
