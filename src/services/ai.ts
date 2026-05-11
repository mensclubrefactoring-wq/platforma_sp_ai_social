import { GoogleGenerativeAI } from "@google/generative-ai";

// Ключ берется из переменных окружения
const genAI = new GoogleGenerativeAI(process.env.GEMINI_API_KEY || "");

// Мы используем модель "gemini-1.5-flash" - она работает быстрее и подходит для MVP
export const aiModel = genAI.getGenerativeModel({ 
  model: "gemini-1.5-flash",
  systemInstruction: `Ты — AI-ассистент платформы социального партнерства "Platforma SP". 
  Твоя задача — помогать крупным компаниям (Сбер, Лукойл и т.д.) находить социальных предпринимателей для реализации их ESG-стратегий.
  Ты должен отвечать в деловом, но вдохновляющем стиле. Используй данные о проектах, если они предоставлены.`
});

export async function askAI(prompt: string, history: any[] = []) {
  const chat = aiModel.startChat({
    history: history.map(h => ({
      role: h.role === "user" ? "user" : "model",
      parts: [{ text: h.content }]
    }))
  });

  const result = await chat.sendMessage(prompt);
  return result.response.text();
}
