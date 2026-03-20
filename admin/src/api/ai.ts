import request from "@/utils/request";
import type {
    AISummaryRequest,
    AISummaryResponse,
    AIAISummaryRequest,
    AIAISummaryResponse,
    AITitleRequest,
    AITitleResponse
} from "@/types/ai";

/**
 * 生成文章摘要（50-100字，创作者角度）
 * @param data 文章内容
 * @returns Promise<AISummaryResponse>
 */
export function generateSummary(data: AISummaryRequest): Promise<AISummaryResponse> {
    return request.post("/admin/ai/summary", data);
}

/**
 * 生成AI摘要（150-200字，旁观者角度）
 * @param data 文章内容
 * @returns Promise<AIAISummaryResponse>
 */
export function generateAISummary(data: AIAISummaryRequest): Promise<AIAISummaryResponse> {
    return request.post("/admin/ai/ai-summary", data);
}

/**
 * 生成标题建议
 * @param data 文章内容
 * @returns Promise<AITitleResponse>
 */
export function generateTitle(data: AITitleRequest): Promise<AITitleResponse> {
    return request.post("/admin/ai/title", data);
}

/**
 * 测试AI配置是否可用
 */
export function testAIConfig(data: { base_url: string; api_key: string; model: string }): Promise<void> {
    return request.post("/admin/ai/test", data);
}
