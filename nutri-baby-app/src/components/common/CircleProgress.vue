<template>
  <view class="circle-progress" :style="{ width: width + 'rpx', height: width + 'rpx' }">
    <canvas
      class="progress-canvas"
      :canvas-id="canvasId"
      :style="{ width: width + 'rpx', height: width + 'rpx' }"
    />
    <view class="progress-text" v-if="showPercent">
      <text class="percent-value" :style="{ fontSize: percentTextSize + 'rpx', color: percentTextColor }">
        {{ Math.round(percent) }}
      </text>
      <text class="percent-symbol" :style="{ fontSize: (percentTextSize * 0.5) + 'rpx', color: percentTextColor }">
        %
      </text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'

interface Props {
  percent: number
  strokeWidth?: number
  strokeColor?: string
  trailColor?: string
  width?: number
  radius?: number
  showPercent?: boolean
  percentTextSize?: number
  percentTextColor?: string
  canvasId?: string
}

const props = withDefaults(defineProps<Props>(), {
  strokeWidth: 12,
  strokeColor: '#52c41a',
  trailColor: '#f0f0f0',
  width: 200,
  radius: 90,
  showPercent: true,
  percentTextSize: 48,
  percentTextColor: '#333333',
  canvasId: 'circle-progress-' + Math.random().toString(36).substr(2, 9)
})

// 绘制圆环
const drawCircle = () => {
  const ctx = uni.createCanvasContext(props.canvasId)
  if (!ctx) return

  const centerX = props.width / 2
  const centerY = props.width / 2
  const radius = props.radius

  // 清空画布
  ctx.clearRect(0, 0, props.width, props.width)

  // 绘制底层圆环（轨道）
  ctx.beginPath()
  ctx.arc(centerX, centerY, radius, 0, 2 * Math.PI)
  ctx.setStrokeStyle(props.trailColor)
  ctx.setLineWidth(props.strokeWidth)
  ctx.setLineCap('round')
  ctx.stroke()

  // 绘制进度圆环
  if (props.percent > 0) {
    ctx.beginPath()
    const startAngle = -Math.PI / 2 // 从顶部开始
    const endAngle = startAngle + (2 * Math.PI * props.percent) / 100
    ctx.arc(centerX, centerY, radius, startAngle, endAngle)
    ctx.setStrokeStyle(props.strokeColor)
    ctx.setLineWidth(props.strokeWidth)
    ctx.setLineCap('round')
    ctx.stroke()
  }

  ctx.draw()
}

onMounted(() => {
  setTimeout(() => {
    drawCircle()
  }, 100)
})

watch(() => [props.percent, props.strokeColor], () => {
  drawCircle()
}, { deep: true })
</script>

<style lang="scss" scoped>
.circle-progress {
  position: relative;
  display: inline-block;

  .progress-canvas {
    display: block;
  }

  .progress-text {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    display: flex;
    align-items: baseline;
    justify-content: center;

    .percent-value {
      font-weight: 600;
      line-height: 1;
    }

    .percent-symbol {
      margin-left: 4rpx;
      opacity: 0.8;
    }
  }
}
</style>
