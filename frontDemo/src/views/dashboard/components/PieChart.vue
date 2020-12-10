<template>
  <div :class="className" :style="{height:height,width:width}" />
</template>

<script>
import echarts from 'echarts'
require('echarts/theme/macarons') // echarts theme
import resize from './mixins/resize'

export default {
  mixins: [resize],
  props: {
    className: {
      type: String,
      default: 'chart'
    },
    width: {
      type: String,
      default: '100%'
    },
    height: {
      type: String,
      default: '350px'
    },
    chartData: {
      type: Object,
      required: true,
      default() {
        return {
          'title': '服务统计',
          'legend': ['HTTP', 'GRPC', 'TCP'],
          'series': [
            { value: 100, name: 'HTTP' },
            { value: 200, name: 'GRPC' },
            { value: 300, name: 'TCP' }
          ]
        }
      }
    }
  },
  data() {
    return {
      chart: null
    }
  },
  watch: {
    chartData: {
      deep: true,
      handler(val) {
        // this.setOptions(val)
        this.initChart()
      }
    }
  },
  mounted() {
    this.$nextTick(() => {
      this.initChart()
    })
  },
  beforeDestroy() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    initChart() {
      this.chart = echarts.init(this.$el, 'macarons')

      this.chart.setOption({
        title: {
          text: this.chartData.title,
          textStyle: {
            fontSize: 16
          }
        },
        tooltip: {
          trigger: 'item',
          formatter: '{a} <br/>{b} : {c} ({d}%)'
        },
        legend: {
          left: 'center',
          bottom: '10',
          data: this.chartData.legend
        },
        series: [
          {
            name: '服务占比',
            type: 'pie',
            // roseType: 'radius',
            radius: '50%',
            center: ['50%', '45%'],
            data: this.chartData.series,
            animationEasing: 'cubicInOut',
            animationDuration: 2600
          }
        ]
      })
    }
  }
}
</script>
