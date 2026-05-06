<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="title">经营数据上报系统</h1>
      <p class="subtitle">商务局管理端</p>
      <el-form ref="formRef" :model="form" :rules="rules">
        <el-form-item prop="username">
          <el-input v-model="form.username" placeholder="用户名" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="密码" size="large" show-password @keyup.enter="handleLogin" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" style="width:100%" @click="handleLogin">登 录</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/store'
import { adminLogin } from '@/api'
import { ElMessage } from 'element-plus'

const router = useRouter()
const store = useAdminStore()
const loading = ref(false)
const formRef = ref()

const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

async function handleLogin() {
  if (!form.username || !form.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  loading.value = true
  try {
    const res = await adminLogin(form.username, form.password)
    store.setAuth(res.data.token, form.username)
    router.push('/dashboard')
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  height: 100vh; display: flex; align-items: center; justify-content: center;
  background: #0a0f1a;
  background-image:
    radial-gradient(ellipse at 30% 50%, rgba(0, 212, 255, 0.06) 0%, transparent 60%),
    radial-gradient(ellipse at 70% 50%, rgba(124, 58, 237, 0.06) 0%, transparent 60%);
  position: relative;
}
.login-page::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background:
    radial-gradient(1px 1px at 15% 25%, rgba(0,212,255,0.4), transparent),
    radial-gradient(1px 1px at 85% 20%, rgba(124,58,237,0.4), transparent),
    radial-gradient(1px 1px at 35% 75%, rgba(0,212,255,0.3), transparent),
    radial-gradient(1px 1px at 65% 80%, rgba(124,58,237,0.3), transparent),
    radial-gradient(1px 1px at 50% 50%, rgba(0,212,255,0.3), transparent);
  pointer-events: none;
}
.login-card {
  width: 420px;
  background: rgba(15, 23, 42, 0.85);
  border: 1px solid rgba(0, 212, 255, 0.12);
  border-radius: 16px;
  padding: 48px 40px;
  backdrop-filter: blur(16px);
  box-shadow: 0 20px 60px rgba(0,0,0,0.4), 0 0 40px rgba(0,212,255,0.05);
  position: relative;
  z-index: 1;
}
.title {
  font-size: 24px;
  background: linear-gradient(135deg, #00d4ff, #a78bfa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-align: center;
  margin-bottom: 8px;
  letter-spacing: 2px;
}
.subtitle {
  font-size: 13px;
  color: #64748b;
  text-align: center;
  margin-bottom: 32px;
  letter-spacing: 3px;
  text-transform: uppercase;
}
</style>
