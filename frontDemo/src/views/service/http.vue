<template>
  <div class="mixin-components-container">
    <el-row>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span v-if="isEdit===false">创建HTTP服务</span>
          <span v-if="isEdit===true">修改HTTP服务</span>
        </div>
        <div style="margin-bottom:80px;">
          <el-form ref="form" :model="form" label-width="140px">
            <el-form-item label="服务名称" class="is-required">
              <el-input v-model="form.service_name" placeholder="请输入由6-128位字母、数字或下划线组成的服务名称" :disabled="isEdit===true"></el-input>
            </el-form-item>
            <el-form-item label="服务描述" class="is-required">
              <el-input v-model="form.service_desc" placeholder="最多255位字符"></el-input>
            </el-form-item>
            <el-form-item label="接入类型" class="is-required">
              <el-input v-model="form.rule" :disabled="isEdit===true" placeholder="路径格式：/user/,域名格式：www.test.com" class="input-with-select">
                <el-select slot="prepend" v-model="form.rule_type" placeholder="请选择" style="width:80px" :disabled="isEdit===true">
                  <el-option label="路径" :value="0" />
                  <el-option label="域名" :value="1" />
                </el-select>
              </el-input>
            </el-form-item>
            <el-form-item label="支持HTTPS">
              <el-switch v-model="form.need_https" :active-value="1" :inactive-value="0">
              </el-switch>
              <span style="color:#606266;font-weight: 700; width: 300px; padding-left: 15px; padding-right: 15px; text-align: right;">支持WEBSocket</span>
              <el-switch v-model="form.need_websocket" :active-value="1" :inactive-value="0">
              </el-switch>
              <span style="color:#606266;font-weight: 700; width: 300px; padding-left: 15px; padding-right: 15px; text-align: right;">支持StripURL</span>
              <el-switch v-model="form.need_strip_uri" :active-value="1" :inactive-value="0">
              </el-switch>
            </el-form-item>
            <el-form-item label="URL重写">
              <el-input v-model="form.url_rewrite" type="textarea" autosize placeholder="格式:^/gateway/test_service(.*) $1 多条换行"></el-input>
            </el-form-item>
            <el-form-item label="Header转换">
              <el-input v-model="form.header_transfor" type="textarea" autosize placeholder="header转换支持add(增加)/del(删除)/edit(修改) 格式：add header value"></el-input>
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
            <el-form-item label="建立连接超时">
              <el-input v-model="form.upstream_connect_timeout" placeholder="单位s，0表示无限制"></el-input>
            </el-form-item>
            <el-form-item label="获取header超时">
              <el-input v-model="form.upstream_header_timeout" placeholder="单位s，0表示无限制"></el-input>
            </el-form-item>
            <el-form-item label="链接最大空闲时间">
              <el-input v-model="form.upstream_idle_timeout" placeholder="单位s，0表示无限制"></el-input>
            </el-form-item>
            <el-form-item label="最大空闲链接数">
              <el-input v-model="form.upstream_max_idle" placeholder="0表示无限制"></el-input>
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
import { serviceAddHTTP, serviceDetail, serviceUpdateHTTP } from '@/api/service'
export default {
  name: 'ServiceCreateHTTP',
  data() {
    return {
      submitButtonDisabled: false,
      isEdit: false,
      form: {
        service_name: '',
        service_desc: '',
        rule: '',
        rule_type: 0,
        need_https: 0,
        need_websocket: 0,
        need_strip_uri: 1,
        url_rewrite: '',
        header_transfor: '',
        open_auth: 0,
        white_list: '',
        black_list: '',
        clientip_flow_limit: '',
        service_flow_limit: '',
        round_type: 0,
        ip_list: '',
        weight_list: '',
        upstream_connect_timeout: '',
        upstream_header_timeout: '',
        upstream_idle_timeout: '',
        upstream_max_idle: ''
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
        // hhtp rule
        this.form.rule = response.data.http.rule
        this.form.rule_type = response.data.http.rule_type
        this.form.need_https = response.data.http.need_https
        this.form.need_websocket = response.data.http.need_websocket
        this.form.need_strip_uri = response.data.http.need_strip_uri
        this.form.url_rewrite = response.data.http.url_rewrite.replace(/,/g, '\n')
        this.form.header_transfor = response.data.http.header_transfor.replace(/,/g, '\n')
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
        this.form.upstream_connect_timeout = response.data.load_balance.upstream_connect_timeout
        this.form.upstream_header_timeout = response.data.load_balance.upstream_header_timeout
        this.form.upstream_idle_timeout = response.data.load_balance.upstream_idle_timeout
        this.form.upstream_max_idle = response.data.load_balance.upstream_max_idle
      })
    },
    handleSubmit() {
      this.submitButtonDisabled = true
      const addQuery = Object.assign({}, this.form)
      addQuery.white_list = addQuery.white_list.replace(/\n/g, ',')
      addQuery.black_list = addQuery.black_list.replace(/\n/g, ',')
      addQuery.url_rewrite = addQuery.url_rewrite.replace(/\n/g, ',')
      addQuery.ip_list = addQuery.ip_list.replace(/\n/g, ',')
      addQuery.weight_list = addQuery.weight_list.replace(/\n/g, ',')
      addQuery.header_transfor = addQuery.header_transfor.replace(/\n/g, ',')
      // string转换位int
      addQuery.clientip_flow_limit = Number(addQuery.clientip_flow_limit)
      addQuery.service_flow_limit = Number(addQuery.service_flow_limit)
      addQuery.upstream_connect_timeout = Number(addQuery.upstream_connect_timeout)
      addQuery.upstream_header_timeout = Number(addQuery.upstream_header_timeout)
      addQuery.upstream_idle_timeout = Number(addQuery.upstream_idle_timeout)
      addQuery.upstream_max_idle = Number(addQuery.upstream_max_idle)
      if (this.isEdit) {
        serviceUpdateHTTP(addQuery).then((response) => {
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
        serviceAddHTTP(addQuery).then((response) => {
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
