<template>
  <div class="mixin-components-container">
    <el-row>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span v-if="isEdit===false">创建租客信息</span>
          <span v-if="isEdit===true">修改租客信息</span>
        </div>
        <div style="margin-bottom:80px;">
          <el-form ref="form" :model="form" label-width="140px">
            <el-form-item label="app_id" class="is-required">
              <el-input v-model="form.app_id" placeholder="6-128位字母数字下划线" :disabled="isEdit===true"></el-input>
            </el-form-item>
            <el-form-item label="租户名称" class="is-required">
              <el-input v-model="form.name" placeholder="最多256字符,必填"></el-input>
            </el-form-item>
            <el-form-item label="32位密钥">
              <el-input v-model="form.secret" placeholder="32位字符串,非必填,自动生成" :disabled="isEdit===true"></el-input>
            </el-form-item>
            <el-form-item label="QPS限流">
              <el-input v-model="form.qps" placeholder="0表示无限制"></el-input>
            </el-form-item>
            <el-form-item label="日请求量限流">
              <el-input v-model="form.qpd" placeholder="0表示无限制"></el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :disable="submitButtonDisabled" @click="handleSubmit">立即提交</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import { appAdd, appDetail, appUpdate } from '@/api/app'
export default {
  name: 'AppCreate',
  data() {
    return {
      submitButtonDisabled: false,
      isEdit: false,
      form: {
        app_id: '',
        name: '',
        secret: '',
        qps: '',
        qpd: ''
      }
    }
  },
  created() {
    // 捕获id
    const id = this.$route.params && this.$route.params.id
    if (id > 0) {
      this.isEdit = true
      this.fecthData(id)
    }
  },
  methods: {
    fecthData(id) {
      const query = { 'id': id }
      appDetail(query).then((response) => {
        this.form.id = response.data.id
        this.form.app_id = response.data.app_id
        this.form.name = response.data.name
        this.form.secret = response.data.secret
        this.form.qps = response.data.qps
        this.form.qpd = response.data.qpd
      }).catch(() => {
      })
    },
    handleSubmit() {
      this.submitButtonDisabled = true
      const addQuery = Object.assign({}, this.form)
      // string转换位int
      addQuery.qps = Number(addQuery.qps)
      addQuery.qpd = Number(addQuery.qpd)
      if (this.isEdit) {
        appUpdate(addQuery).then((response) => {
          this.submitButtonDisabled = false
          this.$notify({
            title: 'Success',
            message: '修改成功',
            type: 'success',
            duration: 2000
          }).catch(() => {
            this.submitButtonDisabled = false
          })
        })
      } else {
        appAdd(addQuery).then((response) => {
          this.submitButtonDisabled = false
          this.$notify({
            title: 'Success',
            message: '添加成功',
            type: 'success',
            duration: 2000
          }).catch(() => {
            this.submitButtonDisabled = false
          })
        })
      }
    }
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
