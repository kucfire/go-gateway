<template>
  <div class="dashboard-editor-container">
    <panel-group :data="PanelGroupData" />
    <el-row :gutter="32">
      <el-col :xs="24" :sm="24" :lg="16">
        <div class="chart-wrapper">
          <line-chart :chart-data="FlowStat" />
        </div>
      </el-col>
      <el-col :xs="24" :sm="24" :lg="8">
        <div class="chart-wrapper">
          <pie-chart :chart-data="ServiceStat" />
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import PanelGroup from './components/PanelGroup'
import LineChart from './components/LineChart'
import PieChart from './components/PieChart'
import { panelGroupData, flowStat, serviceStat } from '@/api/dashboard'

export default {
  name: 'DashboardAdmin',
  components: {
    PanelGroup,
    LineChart,
    PieChart
  },
  data() {
    return {
      PanelGroupData: {
        'service_num': 500,
        'app_num': 500,
        'current_qps': 500,
        'today_request_num': 500
      },
      FlowStat: {
        'title': '今日流量统计',
        'today': [],
        'yesterday': []
      },
      ServiceStat: {
        'title': '服务统计',
        'legend': [],
        'series': []
      }
    }
  },
  created() {
    this.fetchPanelGroupData()
    this.fetchFlowStat()
    this.fetchServceStat()
  },
  methods: {
    fetchPanelGroupData() {
      panelGroupData({}).then(response => {
        this.PanelGroupData = response.data
      }).catch(() => {})
    },
    fetchFlowStat() {
      flowStat({}).then(response => {
        this.FlowStat.today = response.data.today
        this.FlowStat.yesterday = response.data.yesterday
      }).catch(() => {})
    },
    fetchServceStat() {
      serviceStat({}).then(response => {
        this.ServiceStat.legend = response.data.legend
        this.ServiceStat.series = response.data.data
      }).catch(() => {})
    },
    handleSetLineChartData(type) {
      // this.lineChartData = lineChartData[type]
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard-editor-container {
  padding: 32px;
  background-color: rgb(240, 242, 245);
  position: relative;

  .github-corner {
    position: absolute;
    top: 0px;
    border: 0;
    right: 0;
  }

  .chart-wrapper {
    background: #fff;
    padding: 16px 16px 0;
    margin-bottom: 32px;
  }
}

@media (max-width:1024px) {
  .chart-wrapper {
    padding: 8px;
  }
}
</style>
