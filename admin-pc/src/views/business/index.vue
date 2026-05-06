<template>
  <div class="business-page">
    <el-card>
      <div class="search-bar">
        <el-input v-model="keyword" placeholder="搜索企业名称..." clearable style="width: 300px" @keyup.enter="loadData" />
        <el-button type="primary" @click="loadData">搜索</el-button>
      </div>

      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" label="企业名称" min-width="180" />
        <el-table-column prop="industry_type" label="行业类型" width="100" />
        <el-table-column prop="legal_person" label="法人" width="100" />
        <el-table-column prop="contact_phone" label="联系电话" width="130" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="160">
          <template #default="{ row }">{{ row.created_at?.slice(0, 10) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/business/${row.id}`)">台账</el-button>
            <el-button type="primary" link @click="openEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page" :page-size="pageSize"
        :total="total" layout="total, prev, pager, next"
        @current-change="loadData" style="margin-top: 20px; justify-content: center"
      />
    </el-card>

    <el-dialog v-model="editVisible" title="编辑商户信息" width="500px" @close="resetForm">
      <el-form ref="formRef" :model="editForm" :rules="rules" label-width="80px">
        <el-form-item label="企业名称" prop="name">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="法人" prop="legal_person">
          <el-input v-model="editForm.legal_person" />
        </el-form-item>
        <el-form-item label="行业类型" prop="industry_type">
          <el-select v-model="editForm.industry_type" style="width:100%">
            <el-option label="餐饮" value="餐饮" />
            <el-option label="零售" value="零售" />
            <el-option label="住宿" value="住宿" />
            <el-option label="烟酒" value="烟酒" />
            <el-option label="混合" value="混合" />
          </el-select>
        </el-form-item>
        <el-form-item label="营业执照号" prop="license_no">
          <el-input v-model="editForm.license_no" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contact_phone">
          <el-input v-model="editForm.contact_phone" />
        </el-form-item>
        <el-form-item label="经营地址" prop="address">
          <el-input v-model="editForm.address" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch
            v-model="editForm.status"
            :active-value="'active'"
            :inactive-value="'disabled'"
            active-text="正常"
            inactive-text="禁用"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { onMounted } from 'vue'
import { getBusinesses, updateBusiness } from '@/api'
import type { FormInstance } from 'element-plus'
import { ElMessage } from 'element-plus'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const keyword = ref('')
const loading = ref(false)
const saving = ref(false)
const editVisible = ref(false)
const editId = ref(0)
const formRef = ref<FormInstance>()

const editForm = reactive({
  name: '',
  legal_person: '',
  industry_type: '',
  license_no: '',
  contact_phone: '',
  address: '',
  status: 'active'
})

const rules = {
  name: [{ required: true, message: '请输入企业名称', trigger: 'blur' }],
  legal_person: [{ required: true, message: '请输入法人姓名', trigger: 'blur' }],
  industry_type: [{ required: true, message: '请选择行业类型', trigger: 'change' }]
}

async function loadData() {
  loading.value = true
  try {
    const res = await getBusinesses({ page: page.value, page_size: pageSize.value, keyword: keyword.value })
    list.value = res.data.list
    total.value = res.data.total
  } finally {
    loading.value = false
  }
}

function openEdit(row: any) {
  editId.value = row.id
  editForm.name = row.name
  editForm.legal_person = row.legal_person
  editForm.industry_type = row.industry_type
  editForm.license_no = row.license_no || ''
  editForm.contact_phone = row.contact_phone || ''
  editForm.address = row.address || ''
  editForm.status = row.status
  editVisible.value = true
}

function resetForm() {
  formRef.value?.resetFields()
}

async function handleSave() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  saving.value = true
  try {
    await updateBusiness(editId.value, { ...editForm })
    ElMessage.success('更新成功')
    editVisible.value = false
    loadData()
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || '更新失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => loadData())
</script>

<style scoped>
.business-page { }
.search-bar { display: flex; gap: 12px; margin-bottom: 20px; }
</style>
