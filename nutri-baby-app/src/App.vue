<script setup lang="ts">
import { onLaunch, onShow, onHide } from "@dcloudio/uni-app";
import { initialize } from "@/store/user";

onLaunch((options) => {
  console.log("App Launch with options:", options);

  // 应用启动时立即初始化用户状态,从本地存储恢复登录信息
  // 这样可以避免代码重新编译时丢失登录状态
  initialize();

  console.log("[App] User state initialized from storage");

  // 处理小程序扫码进入的场景
  handleQRCodeScan(options);
});

onShow((options) => {
  console.log("App Show with options:", options);

  // 从后台切换到前台时,也检查是否是扫码进入
  handleQRCodeScan(options);
});

onHide(() => {
  console.log("App Hide");
});

/**
 * 处理小程序码扫描进入的场景
 * @param options - 启动参数或显示参数
 */
function handleQRCodeScan(options: any) {
  console.log("[App] handleQRCodeScan called with:", options);

  // scene 1047 表示扫描小程序码进入
  // scene 1048 表示长按小程序码进入
  // scene 1049 表示识别小程序码进入（从图片）
  const qrCodeScenes = [1047, 1048, 1049];

  if (options?.scene && qrCodeScenes.includes(options.scene)) {
    console.log("[App] Detected QR code scan, scene:", options.scene);

    // 获取 scene 参数（格式：c=ABC123）
    const sceneParam = options.query?.scene || options.scene;

    console.log("[App] Scene parameter:", sceneParam);
    if (sceneParam && typeof sceneParam === "string") {
      // 解析 scene 参数
      const shortCode = parseSceneParameter(sceneParam);
      console.log("[App] Parsed short code:", shortCode);
      if (shortCode) {
        console.log("[App] Parsed short code:", shortCode);

        // 跳转到加入页面
        uni.reLaunch({
          url: `/pages/baby/join/join?code=${shortCode}`,
        });
      } else {
        console.error(
          "[App] Failed to parse short code from scene:",
          sceneParam
        );
        uni.showToast({
          title: "二维码格式错误",
          icon: "none",
        });
      }
    } else {
      console.warn("[App] No scene parameter found in options");
    }
  }
}

/**
 * 解析 scene 参数，提取短码
 * @param scene - scene 字符串，格式如 "c=ABC123"
 * @returns 短码，如 "ABC123"，解析失败返回 null
 */
function parseSceneParameter(scene: string): string | null {
  try {
    // 先尝试 URL 解码
    const params = decodeURIComponent(scene);
    if (params.startsWith("c=")) {
      return params.substring(2);
    }

    // scene 格式: c=ABC123
    if (scene.startsWith("c=")) {
      return scene.substring(2); // 提取 "=" 后面的部分
    }
    if (params.startsWith("c=")) {
      return params.substring(2);
    }

    console.warn("[App] Scene format not recognized:", scene);
    return null;
  } catch (error) {
    console.error("[App] Error parsing scene parameter:", error);
    return null;
  }
}
</script>
<style lang="scss">
/* 引入 Wot Design Uni 主题配置 */
@import "@/theme.scss";

/* 组件样式通过按需引入插件自动加载,无需手动导入 */
</style>
