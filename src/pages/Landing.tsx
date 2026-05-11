import { motion } from "motion/react";
import { ArrowRight, CheckCircle2, Factory, HandHeart, Info, Rocket, Sparkles, Building2, Users, Store } from "lucide-react";
import { Link } from "react-router-dom";

export default function Landing() {
  return (
    <div className="min-h-screen bg-[#f8faff] text-[#101828] font-sans selection:bg-indigo-100">
      {/* Navigation */}
      <nav className="fixed top-0 w-full z-50 bg-white/80 backdrop-blur-md border-b border-gray-100">
        <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between">
          <div className="flex items-center gap-2 font-bold text-xl tracking-tight text-indigo-600">
            <Sparkles className="w-6 h-6" />
            <span>Platforma SP</span>
          </div>
          <div className="hidden md:flex items-center gap-8 text-sm font-medium text-gray-600">
            <a href="#how-it-works" className="hover:text-indigo-600 transition-colors">Как это работает</a>
            <a href="#about" className="hover:text-indigo-600 transition-colors">О платформе</a>
            <a href="#participants" className="hover:text-indigo-600 transition-colors">Участники</a>
          </div>
          <div className="flex items-center gap-4">
            <Link to="/app" className="text-sm font-semibold text-gray-700 hover:text-indigo-600">Войти</Link>
            <Link to="/app" className="bg-indigo-600 text-white px-5 py-2 rounded-full text-sm font-semibold hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-200">
              Начать работу
            </Link>
          </div>
        </div>
      </nav>

      {/* Hero */}
      <header className="pt-32 pb-20 px-4">
        <div className="max-w-7xl mx-auto grid lg:grid-cols-2 gap-12 items-center">
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6 }}
          >
            <div className="inline-flex items-center gap-2 px-3 py-1 bg-indigo-50 border border-indigo-100 text-indigo-600 rounded-full text-xs font-bold mb-6">
              <Rocket className="w-3 h-3" />
              <span>AI-DRIVEN SOCIAL PARTNERSHIP</span>
            </div>
            <h1 className="text-5xl md:text-7xl font-bold tracking-tight leading-[1.1] mb-6">
              Соединяем задачи бизнеса с <span className="text-transparent bg-clip-text bg-gradient-to-r from-indigo-600 to-violet-600">социальным эффектом</span>.
            </h1>
            <p className="text-lg md:text-xl text-gray-500 max-w-xl mb-10 leading-relaxed">
              Платформа, где корпорации публикуют социальные задачи, а AI-ассистент помогает предпринимателям предлагать лучшие решения.
            </p>
            <div className="flex flex-wrap gap-4">
              <Link to="/app" className="group flex items-center gap-2 bg-indigo-600 text-white px-8 py-4 rounded-2xl font-bold hover:bg-indigo-700 transition-all shadow-xl shadow-indigo-200">
                Создать задачу <ArrowRight className="w-5 h-5 group-hover:translate-x-1 transition-transform" />
              </Link>
              <Link to="/app" className="flex items-center gap-2 bg-white border border-gray-200 text-gray-900 px-8 py-4 rounded-2xl font-bold hover:border-gray-300 transition-all">
                Стать исполнителем
              </Link>
            </div>
          </motion.div>

          <motion.div 
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="bg-white border border-gray-200 rounded-[2rem] p-8 shadow-2xl shadow-indigo-100 relative overflow-hidden"
          >
            <div className="absolute top-0 right-0 w-64 h-64 bg-indigo-50 rounded-full -mr-32 -mt-32 blur-3xl opacity-50" />
            <h3 className="text-xl font-bold mb-6 flex items-center gap-2">
              <Store className="w-5 h-5 text-indigo-600" />
              Активные задачи платформы
            </h3>
            <div className="space-y-4 relative z-10">
              {[
                { title: "Закупка сувениров для ветеранов", count: 12, type: "Retail" },
                { title: "Проведение инклюзивного фестиваля", count: 8, type: "Events" },
                { title: "Организация социальной столовой", count: 15, type: "CSR" },
              ].map((item, idx) => (
                <div key={idx} className="flex items-center justify-between p-4 bg-gray-50 rounded-2xl border border-gray-100 hover:border-indigo-200 transition-colors cursor-pointer group">
                  <div>
                    <div className="text-xs font-bold text-indigo-600 mb-1 uppercase tracking-wider">{item.type}</div>
                    <div className="font-semibold text-gray-900">{item.title}</div>
                  </div>
                  <div className="bg-indigo-100 text-indigo-700 px-3 py-1 rounded-full text-xs font-bold group-hover:bg-indigo-600 group-hover:text-white transition-colors">
                    {item.count} откликов
                  </div>
                </div>
              ))}
            </div>
          </motion.div>
        </div>
      </header>

      {/* Stats/Logos */}
      <section className="py-12 border-y border-gray-100 bg-white" id="participants">
        <div className="max-w-7xl mx-auto px-4">
          <p className="text-center text-sm font-bold text-gray-400 uppercase tracking-[0.2em] mb-10">Нам доверяют лидеры рынка</p>
          <div className="flex flex-wrap justify-center items-center gap-12 grayscale opacity-60">
            <div className="flex items-center gap-2 font-bold text-2xl">SBER</div>
            <div className="flex items-center gap-2 font-bold text-2xl">YANDEX</div>
            <div className="flex items-center gap-2 font-bold text-2xl">AVITO</div>
            <div className="flex items-center gap-2 font-bold text-2xl">Sibur</div>
            <div className="flex items-center gap-2 font-bold text-2xl">Magnit</div>
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="py-24 px-4" id="how-it-works">
        <div className="max-w-7xl mx-auto text-center mb-16">
          <h2 className="text-4xl font-bold tracking-tight mb-4 text-gray-900">Прозрачный процесс в 3 шага</h2>
          <p className="text-gray-500 max-w-xl mx-auto">AI-подбор сокращает время на поиск исполнителя до нескольких минут.</p>
        </div>
        <div className="max-w-7xl mx-auto grid md:grid-cols-3 gap-8">
          {[
            { 
              icon: Building2, 
              title: "Бизнес создает задачу", 
              desc: "Корпорация описывает цель и бюджет. Наш AI-ассистент помогает правильно структурировать запрос." 
            },
            { 
              icon: HandHeart, 
              title: "AI подбирает героев", 
              desc: "Система анализирует базу социальных предпринимателей и предлагает тех, кто идеально подходит под задачу." 
            },
            { 
              icon: CheckCircle2, 
              title: "Запуск проекта", 
              desc: "Вы выбираете лучшего исполнителя и начинаете сотрудничество. Весь документооборот — на платформе." 
            },
          ].map((feature, idx) => (
            <motion.div 
              key={idx}
              whileHover={{ y: -5 }}
              className="bg-white p-8 rounded-3xl border border-gray-200 shadow-sm hover:shadow-xl transition-all"
            >
              <div className="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center mb-6">
                <feature.icon className="w-6 h-6" />
              </div>
              <h3 className="text-xl font-bold mb-4">{feature.title}</h3>
              <p className="text-gray-500 leading-relaxed">{feature.desc}</p>
            </motion.div>
          ))}
        </div>
      </section>

      {/* CTA section */}
      <section className="py-24 px-4 bg-indigo-600 relative overflow-hidden" id="cta">
        <div className="absolute top-0 left-0 w-full h-full opacity-10 pointer-events-none">
          <div className="absolute top-10 right-10 w-96 h-96 bg-white rounded-full blur-[100px]" />
          <div className="absolute bottom-10 left-10 w-96 h-96 bg-indigo-400 rounded-full blur-[100px]" />
        </div>
        <div className="max-w-4xl mx-auto text-center relative z-10">
          <h2 className="text-4xl md:text-5xl font-bold text-white mb-6">Готовы масштабировать позитивные изменения?</h2>
          <p className="text-indigo-100 text-lg mb-10">Присоединяйтесь к платформе социального партнерства нового поколения.</p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link to="/app" className="bg-white text-indigo-600 px-10 py-5 rounded-2xl font-bold hover:bg-gray-50 transition-all shadow-2xl">
              Зарегистрировать бизнес
            </Link>
            <Link to="/app" className="bg-indigo-500 text-white border border-indigo-400 px-10 py-5 rounded-2xl font-bold hover:bg-indigo-400 transition-all">
              Стать партнером
            </Link>
          </div>
        </div>
      </section>

      <footer className="py-8 border-t border-gray-100 text-center text-gray-400 text-sm">
        <p>© 2026 Platforma SP. Разработано для социальных предпринимателей и бизнеса.</p>
      </footer>
    </div>
  );
}
