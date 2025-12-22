
// 环境变量
const BASE_URL = import.meta.env.VITE_API_BASE_URL || "https://api.example.com";

/**
 * 解析图片地址
 * 如果是相对路径 (如 /uploads/...), 则补全为完整 URL
 * 如果是绝对路径 (http/https), 则直接返回
 * 如果是本地静态资源 (/static/...), 则直接返回
 */
export const resolveImageUrl = (url?: string): string => {
    if (!url) return "";

    // 已经是完整 URL
    if (url.startsWith("http://") || url.startsWith("https://") || url.startsWith("wxfile://")) {
        return url;
    }

    // Base64 Data URI
    if (url.startsWith("data:")) {
        return url;
    }

    // 本地静态资源 (通常在 static 目录下)
    if (url.startsWith("/static/") || url.startsWith("@/static/")) {
        return url;
    }

    // 假设是后端返回的相对路径
    let target = url;
    if (!target.startsWith("/")) {
        target = "/" + target;
    }

    // 移除 BASE_URL 末尾的斜杠 (如果存在)
    const prefix = BASE_URL.endsWith("/") ? BASE_URL.slice(0, -1) : BASE_URL;

    return `${prefix}${target}`;
};
