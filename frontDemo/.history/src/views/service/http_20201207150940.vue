<template>
  <div class="mixin-components-container">
    <el-row>
      <el-card class="box-card">
        <div slot="header" class="clearfix">
          <span>创建HTTP服务</span>
        </div>
        <div style="margin-bottom:50px;">
          <el-form ref="form" :model="form" label-width="14service_name0px">
            <el-form-item label="服务名称">
              <el-input v-model="form.service_name" placeholder="请输入由6-128位字母、数字或下划线组成的服务名称"></el-input>
            </el-form-item>
            <el-form-item label="服务描述">
              <el-input v-model="form.service_desc" placeholder="最多255位字符"></el-input>
            </el-form-item>
            <el-form-item label="接入类型">
              <el-input v-model="form.load_type" placeholder="路径格式：/user/,域名格式：www.test.com" class="input-with-select">
                <el-select slot="prepend" 
                v-model="form.rule_type" placeholder="请选择">
                  <el-option label="路径" value="0"></el-option>
                  <el-option label="域名" value="1"></el-option>
                </el-select>
              </el-input>
            </el-form-item>
            <el-form-item label="支持HTTPS">
              <el-switch
                v-model="form.need_https"
                active-color="#13ce66">
              </el-switch>
            </el-form-item>
            <el-form-item label="支持WEBSocket">
              <el-switch
                v-model="form.need_websocket"
                active-color="#13ce66">
              </el-switch>
            </el-form-item>
            <el-form-item label="支持StripURL">
              <el-switch
                v-model="form.need_strip_uri"
                active-color="#13ce66">
              </el-switch>
            </el-form-item>
            <el-form-item label="URL重写">
              <el-input v-model="form.url_rewrite" type="textarea" autosize placeholder="格式:^/gateway/test_service(.*) $1 多条换行"></el-input>
            </el-form-item>
            <el-form-item label="Header转换">
              <el-input v-model="form.header_transfor" type="textarea" autosize placeholder="header转换支持add(增加)/del(删除)/edit(修改) 格式：add header value"></el-input>
            </el-form-item>
            <el-form-item label="开启验证">
              <el-switch
                v-model="form.open_auth"
                active-color="#13ce66">
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
              <el-radio :label="3">random</el-radio>
              <el-radio :label="6">round-robin</el-radio>
              <el-radio :label="9">weight_round-robin</el-radio>
              <el-radio :label="9">ip_hash</el-radio>
            </el-radio-group>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="onSubmit">立即创建</el-button>
              <el-button>取消</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'Dashboard',
  data() {
    return {
      form: {
        name: '',
        region: '',
        date1: '',
        date2: '',
        delivery: false,
        type: [],
        resource: '',
        desc: ''
      }
    }
  },
  methods: {
    onSubmit() {
      console.log('submit!')
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
