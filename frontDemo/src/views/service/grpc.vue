<template>
  <div class="mixin-components-container">
    <el-row>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span v-if="isEdit===false">创建GRPC服务</span>
          <span v-if="isEdit===true">修改GRPC服务</span>
        </div>
        <div style="margin-bottom:80px;">
          <el-form ref="form" :model="form" label-width="140px">
            <el-form-item label="服务名称" class="is-required">
              <el-input v-model="form.service_name" placeholder="请输入由6-128位字母、数字或下划线组成的服务名称" :disabled="isEdit===true"></el-input>
            </el-form-item>
            <el-form-item label="服务描述" class="is-required">
              <el-input v-model="form.service_desc" placeholder="最多255位字符"></el-input>
            </el-form-item>
            <el-form-item label="端口" class="is-required">
              <el-input v-model="form.port" placeholder="需要设置8000~8999范围内的数字"></el-input>
            </el-form-item>
            <el-form-item label="METADATA转换">
              <el-input v-model="form.header_transfor" type="textarea" autosize placeholder="metadata转换支持add(增加)/del(删除)/edit(修改) 格式：add header value 多条换行"></el-input>
            </el-form-item>
            <el-form-item label="开启验证">
              <el-switch v-model="form.open_auth" :active-value="1" :inactive-value="0">
              </el-switch>
            </el-form-item>
            <el-form-item label="IP白名单">
              <el-input v-model="form.white_list" type="textarea" autosize placeholder="格式:127.0.0.1 多条换行，白名单优先级高于黑名单"></el-input>
            </el-form-item>
            <el-form-item label="IP黑名单">
              <el-input v-model="form.black_list" type="textarea" autosize placeholder="格式:127.0.0.1 多条换行"></el-input>
            </el-form-item>
            <el-form-item label="客户端限流">
              <el-input v-model="form.clientip_flow_limit" placeholder="0表示无限制"></el-input>
            </el-form-item>
            <el-form-item label="服务端限流">
              <el-input v-model="form.service_flow_limit" placeholder="0表示无限制"></el-input>
            </el-form-item>
            <el-form-item label="轮询方式">
              <el-radio-group v-model="form.round_type">
                <el-radio :label="0">random</el-radio>
                <el-radio :label="1">round-robin</el-radio>
                <el-radio :label="2">weight_round-robin</el-radio>
                <el-radio :label="3">ip_hash</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item label="IP列表" class="is-required">
              <el-input v-model="form.ip_list" type="textarea" autosize placeholder="格式:127.0.0.1:80 多条换行"></el-input>
            </el-form-item>
            <el-form-item label="权重列表" class="is-required">
              <el-input v-model="form.weight_list" type="textarea" autosize placeholder="格式:50 多条换行"></el-input>
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
import { serviceAddGRPC, serviceDetail, serviceUpdateGRPC } from '@/api/service'
export default {
  name: 'ServiceCreateTCP',
  data() {
    return {
      submitButtonDisabled: false,
      isEdit: false,
      form: {
        service_name: '',
        service_desc: '',
        port: '',
        header_transfor: '',
        open_auth: 0,
        white_list: '',
        black_list: '',
        clientip_flow_limit: '',
        service_flow_limit: '',
        round_type: 0,
        ip_list: '',
        weight_list: ''
      }
    }
  },
  created() {
    const id = this.$route.params && this.$route.params.id
    if (id > 0) {
      this.isEdit = true
      this.fecthData(id)
    }
  },
  methods: {
    fecthData(id) {
      const query = { 'id': id }
      serviceDetail(query).then((response) => {
        this.form.id = response.data.info.id
        this.form.service_name = response.data.info.service_name
        this.form.service_desc = response.data.info.service_desc
        // tcpp rule
        this.form.port = response.data.grpc.port
        this.form.header_transfor = response.data.grpc.header_transfor.replace(/,/g, '\n')
        // access control
        this.form.open_auth = response.data.access_control.open_auth
        this.form.white_list = response.data.access_control.white_list.replace(/,/g, '\n')
        this.form.black_list = response.data.access_control.black_list.replace(/,/g, '\n')
        this.form.clientip_flow_limit = response.data.access_control.clientip_flow_limit
        this.form.service_flow_limit = response.data.access_control.service_flow_limit
        // load balance
        this.form.round_type = response.data.load_balance.round_type
        this.form.ip_list = response.data.load_balance.ip_list.replace(/,/g, '\n')
        this.form.weight_list = response.data.load_balance.weight_list.replace(/,/g, '\n')
      }).catch(() => {
      })
    },
    handleSubmit() {
      this.submitButtonDisabled = true
      const addQuery = Object.assign({}, this.form)
      addQuery.white_list = addQuery.white_list.replace(/\n/g, ',')
      addQuery.black_list = addQuery.black_list.replace(/\n/g, ',')
      addQuery.ip_list = addQuery.ip_list.replace(/\n/g, ',')
      addQuery.weight_list = addQuery.weight_list.replace(/\n/g, ',')
      addQuery.header_transfor = addQuery.header_transfor.replace(/\n/g, ',')
      // string转换位int
      addQuery.port = Number(addQuery.port)
      addQuery.clientip_flow_limit = Number(addQuery.clientip_flow_limit)
      addQuery.service_flow_limit = Number(addQuery.service_flow_limit)
      if (this.isEdit) {
        serviceUpdateGRPC(addQuery).then((response) => {
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
        serviceAddGRPC(addQuery).then((response) => {
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
