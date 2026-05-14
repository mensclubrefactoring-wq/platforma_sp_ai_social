import { apiFetch } from "@/src/lib/api";

export async function askAI(prompt: string, history: any[] = []) {
  const res = await apiFetch("/api/ai/ask", {
    method: "POST",
    body: JSON.stringify({ prompt, history })
  });
  return res.response || res.proposal || "Нет ответа от AI";
}
