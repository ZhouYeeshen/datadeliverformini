<template>
  <div class="pos-page">
    <el-card>
      <div class="header-bar">
        <span>收银系统接口凭证管理</span>
        <el-button type="primary" @click="showCreateDialog">新增凭证</el-button>
      </div>

      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="business_id" label="商户ID" width="80" />
        <el-table-column prop="api_key" label="API Key" min-width="280" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_called_at" label="最近调用" width="160">
          <template #default="{ row }">{{ row.last_called_at?.slice(0, 16) || '从未调用' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="toggleStatus(row)">
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" @click="showLogs(row)">日志</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="createVisible" title="新增API凭证" width="450px">
      <el-form>
        <el-form-item label="商户ID">
          <el-input-number v-model="newCredBusinessId" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="createCredential" :loading="creating">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="logVisible" title="接口调用日志" width="800px">
      <el-table :data="logs" stripe max-height="400">
        <el-table-column prop="method" label="方法" width="60" />
        <el-table-column prop="path" label="路径" min-width="180" />
        <el-table-column prop="response_code" label="状态码" width="80" />
        <el-table-column prop="duration_ms" label="耗时(ms)" width="80" />
        <el-table-column prop="ip_address" label="IP" width="140" />
        <el-table-column prop="created_at" label="时间" width="160">
          <template #default="{ row }">{{ row.created_at?.slice(0, 16) }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPosCredentials, createPosCredential, updatePosCredential, getPosLogs } from '@/api'

const list = ref<any[]>([])
const loading = ref(false)
const createVisible = ref(false)
const creating = ref(false)
const newCredBusinessId = ref(0)
const logVisible = ref(false)
const logs = ref<any[]>([])

async function loadData() {
  loading.value = true
  try {
    const res = await getPosCredentials({ page: 1, page_size: 100 })
    list.value = res.data.list
  } finally { loading.value = false }
}

function showCreateDialog() { createVisible.value = true }

async function createCredential() {
  creating.value = true
  try {
    await createPosCredential(newCredBusinessId.value)
    createVisible.value = false
    newCredBusinessId.value = 0
    loadData()
    ElMessage.success('创建成功')
  } catch { ElMessage.error('创建失败') } finally { creating.value = false }
}

async function toggleStatus(row: any) {
  const newStatus = row.status === 'active' ? 'disabled' : 'active'
  try {
    await updatePosCredential(row.id, newStatus)
    loadData()
    ElMessage.success('已' + (newStatus === 'active' ? '启用' : '禁用'))
  } catch {}
}

async function showLogs(row: any) {
  try {
    const res = await getPosLogs({ credential_id: row.id, page: 1, page_size: 100 })
    logs.value = res.data.list
    logVisible.value = true
  } catch {}
}

onMounted(() => loadData())
</script>

<style scoped>
.pos-page { }
.header-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
</style>
