declare module '@qiun/ucharts' {
  interface UChartsConstructorOptions {
    type?: string
    context?: CanvasRenderingContext2D
    canvasId?: string
    width?: number
    height?: number
    background?: string
    categories?: string[]
    series?: Array<{
      name: string
      data: any[]
    }>
    animation?: boolean
    pixelRatio?: number
    color?: string[]
    padding?: [number, number, number, number]
    legend?: any
    xAxis?: any
    yAxis?: any
    extra?: any
    opts?: any
    enableScroll?: boolean
    [key: string]: any
  }

  interface UChartsInstance {
    setOption(options: any): void
    updateData(options: any): void
    show(): void
    hide(): void
    removeOption(index?: number): void
    dispose(): void
    touchLegend?(touch: any): void
    [key: string]: any
  }

  interface ChartOptions {
    [key: string]: any
  }

  interface ChartData {
    [key: string]: any
  }

  class uCharts {
    constructor(options: UChartsConstructorOptions, callback?: (chart: UChartsInstance) => void)
    setOption(options: any): void
    updateData(options: any): void
    show(): void
    hide(): void
    removeOption(index?: number): void
    dispose(): void
    touchLegend?(touch: any): void
    showToolTip?(event: any): void
    touchStart?(event: any): void
    touchMove?(event: any): void
    touchEnd?(event: any): void
    scrollStart?(event: any): void
    scroll?(event: any): void
    scrollEnd?(event: any): void
    [key: string]: any
  }

  export default uCharts
}

