<template>
  <div class="system-page">
    <el-card title="管理员账号管理">
      <el-descriptions title="当前登录信息" :column="2" border>
        <el-descriptions-item label="用户名">{{ store.username }}</el-descriptions-item>
        <el-descriptions-item label="角色">管理员</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card title="预警规则配置" style="margin-top: 20px">
      <el-form label-width="160px">
        <el-form-item label="未申报预警触发时间">
          <el-select :model-value="'每月最后5天'" disabled>
            <el-option label="每月最后5天" value="last5" />
          </el-select>
        </el-form-item>
        <el-form-item label="异常波动阈值">
          <el-input-number :model-value="30" :min="5" :max="100" disabled />%
          <span style="margin-left: 8px; color: #909399;">当前值与前2月均值比较</span>
        </el-form-item>
        <el-form-item label="预警扫描频率">
          <el-select :model-value="'每天上午9点'" disabled>
            <el-option label="每天上午9点" value="daily" />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { useAdminStore } from '@/store'
const store = useAdminStore()
</script>

<style scoped>
.system-page { display: flex; flex-direction: column; gap: 20px; }
</style>
