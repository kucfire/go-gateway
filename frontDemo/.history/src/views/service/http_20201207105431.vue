<template>
  <div class="mixin-components-container">
    <el-row>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span>Buttons</span>
        </div>
        <div style="margin-bottom:50px;">
          <el-col :span="4" class="text-center">
            <router-link class="pan-btn blue-btn" to="/documentation/index">
              Documentation
            </router-link>
          </el-col>
          <el-col :span="4" class="text-center">
            <router-link class="pan-btn light-blue-btn" to="/icon/index">
              Icons
            </router-link>
          </el-col>
          <el-col :span="4" class="text-center">
            <router-link class="pan-btn pink-btn" to="/excel/export-excel">
              Excel
            </router-link>
          </el-col>
          <el-col :span="4" class="text-center">
            <router-link class="pan-btn green-btn" to="/table/complex-table">
              Table
            </router-link>
          </el-col>
          <el-col :span="4" class="text-center">
            <router-link class="pan-btn tiffany-btn" to="/example/create">
              Form
            </router-link>
          </el-col>
          <el-col :span="4" class="text-center">
            <router-link class="pan-btn yellow-btn" to="/theme/index">
              Theme
            </router-link>
          </el-col>
        </div>
      </el-card>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'Dashboard',
  data() {
    return {}
  }
}
</script>

<style scoped>
.mixin-components-container {
  background-color: #f0f2f5;
  padding: 30px;
  min-height: calc(100vh - 84px);
}
.component-item{
  min-height: 100px;
}
</style>
