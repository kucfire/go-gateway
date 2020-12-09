<template>
  <div class="chart-container">
    <chart height="100%" width="100%" :data="chartData" />
  </div>
</template>

<script>
import Chart from './components/LineStat'
import { appStat, appDetail } from '@/api/app'

export default {
  name: 'AppChart',
  components: { Chart },
  data() {
    return {
      chartData: {
        'title': '',
        'today': [],
        'yesterday': []
      }
    }
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    this.fecthStat(id)
  },
  methods: {
    fecthStat(id) {
      const query = { 'id': id }
      appStat(query).then((responseStat) => {
        appDetail(query).then((responseDetail) => {
          this.chartData = {
            'title': responseDetail.data.name + '服务统计',
            'today': responseStat.data.today,
            'yesterday': responseStat.data.yesterday
          }
        })
      })
    }
  }
}
</script>

<style scoped>
.chart-container{
  position: relative;
  width: 100%;
  height: calc(100vh - 84px);
}
</style>

