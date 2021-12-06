<template>
  <div class="app-container">
    <el-form
      ref="formData"
      :rules="rules"
      :model="formData"
      label-position="left"
      label-width="100px"
      status-icon
    >
      <el-form-item label="邮箱域名：" prop="host" style="width: 300px">
        <el-input v-model="formData.host" placeholder="请输入规则类别" />
      </el-form-item>
      <el-form-item label="端口：" prop="port" style="width: 300px">
        <el-input v-model.number="formData.port" placeholder="请输入端口" />
      </el-form-item>
      <el-form-item label="邮箱用户：" prop="sender" style="width: 300px">
        <el-input v-model="formData.sender" placeholder="请输入规则类别" />
      </el-form-item>
      <el-form-item label="邮箱密码：" prop="password" style="width: 300px">
        <el-input v-model="formData.password" show-password placeholder="请输入规则类别" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: 'Ldap',
  components: { },
  props: {
    params: {
      type: String,
      default: '0'
    },
    name: {
      type: String,
      default: '0'
    }
  },
  data() {
    return {
      formData: {},
      rules: {
        host: [
          { required: true, message: '请输入域名', trigger: 'change' }
        ],
        port: [
          { required: true, message: '请输入端口', trigger: 'change' }
        ],
        sender: [
          { required: true, message: '请输入邮箱用户', trigger: 'change' }
        ],
        password: [
          { required: true, message: '请输入邮箱密码', trigger: 'change' }
        ]
      }
    }
  },
  watch: {
    params(newVal, oldVal) {
      console.log('ldap:', newVal)
      if (newVal === this.name) {
        this.fetchData()
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      console.log('ldap fetch data')
    },
    onSubmit() {
      console.log('submit!')
    }
  }
}
</script>
