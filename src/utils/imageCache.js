// 简单的图片缓存模块，基于模块级 Map 缓存 Promise<string>（objectURL）
const cache = new Map();

export async function getCachedImage(url) {
  if (!url) return null;
  // 如果已有缓存（promise或string），返回它
  if (cache.has(url)) {
    const cached = cache.get(url);
    // 如果是字符串（已解析的objectURL），直接返回
    if (typeof cached === 'string') return cached;
    // 否则是 Promise，await 它
    return await cached;
  }

  // 否则创建 fetch Promise 并缓存
  const p = (async () => {
    try {
      const resp = await fetch(url, { cache: 'force-cache' });
      if (!resp.ok) throw new Error('Image fetch failed')
      const blob = await resp.blob();
      const objectUrl = URL.createObjectURL(blob);
      // 将最终的 string 存入 cache（覆盖 promise）
      cache.set(url, objectUrl);
      return objectUrl;
    } catch (e) {
      // 失败时移除缓存项，避免永久失败
      cache.delete(url);
      throw e;
    }
  })();

  // 先缓存 promise，防止并发重复请求
  cache.set(url, p);
  return await p;
}

// 可选：在应用卸载前释放所有 object URLs
export function revokeAllCachedObjectURLs() {
  for (const [key, val] of cache.entries()) {
    if (typeof val === 'string') {
      try { URL.revokeObjectURL(val); } catch (e) {}
    }
  }
  cache.clear();
}
