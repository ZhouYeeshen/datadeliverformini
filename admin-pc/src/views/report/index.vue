<template>
  <div class="report-page">
    <el-card>
      <div class="search-bar">
        <el-input v-model="filters.businessName" placeholder="企业名称" clearable style="width: 200px" />
        <el-input v-model="filters.month" placeholder="月份 (2026-05)" clearable style="width: 180px" />
        <el-select v-model="filters.status" placeholder="状态" clearable style="width: 140px">
          <el-option label="已上报" value="submitted" />
          <el-option label="草稿" value="draft" />
        </el-select>
        <el-button type="primary" @click="loadData">查询</el-button>
      </div>

      <el-table :data="list" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="report_month" label="月份" width="100" />
        <el-table-column label="企业名称" min-width="160">
          <template #default="{ row }">{{ row.business_name || '-' }}</template>
        </el-table-column>
        <el-table-column prop="total_revenue" label="营业总额" width="130">
          <template #default="{ row }">¥{{ row.total_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="source_type" label="来源" width="80" />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'submitted' ? 'success' : 'info'" size="small">
              {{ row.status === 'submitted' ? '已上报' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="附件" width="120">
          <template #default="{ row }">
            <div v-if="row.photo_urls" class="table-photos">
              <img
                v-for="(url, i) in row.photo_urls.split(',').filter(Boolean).slice(0, 3)"
                :key="i"
                :src="getPhotoSrc(url)"
                @click="openImageViewer(row.photo_urls.split(',').filter(Boolean), i)"
                class="photo-thumb"
              />
              <span v-if="row.photo_urls.split(',').filter(Boolean).length > 3" class="photo-more">
                +{{ row.photo_urls.split(',').filter(Boolean).length - 3 }}
              </span>
            </div>
            <span v-else style="color:#ccc">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="submitted_at" label="上报时间" width="160">
          <template #default="{ row }">{{ row.submitted_at?.slice(0, 16) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="primary" link @click="showEdit(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="page" :page-size="pageSize"
        :total="total" layout="total, prev, pager, next"
        @current-change="loadData" style="margin-top: 20px; justify-content: center"
      />
    </el-card>

    <!-- Edit dialog -->
    <el-dialog v-model="editVisible" title="编辑申报记录" width="640px" @close="editFormRef?.resetFields()">
      <el-form v-if="editForm" ref="editFormRef" :model="editForm" label-width="100px">
        <el-form-item label="月份">
          <el-input v-model="editForm.report_month" placeholder="2026-05" style="width: 200px" />
        </el-form-item>
        <el-form-item label="营业总额">
          <el-input-number v-model="editForm.total_revenue" :precision="2" :min="0" :step="100" style="width: 240px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="editForm.status" style="width: 160px">
            <el-option label="已上报" value="submitted" />
            <el-option label="草稿" value="draft" />
          </el-select>
        </el-form-item>

        <!-- Attached photos -->
        <el-form-item label="附件照片" v-if="photoList.length > 0">
          <div class="photo-grid">
            <div v-for="(url, i) in photoList" :key="i" class="photo-item">
              <img :src="getPhotoSrc(url)" @click="openImageViewer(photoList, i)" class="photo-thumb-lg" />
              <el-icon class="photo-delete" @click.stop="removePhoto(i)" v-if="photoList.length > 1"><CircleCloseFilled /></el-icon>
            </div>
          </div>
        </el-form-item>
        <el-form-item label="附件照片" v-else>
          <span style="color:#999">暂无附件</span>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { CircleCloseFilled } from '@element-plus/icons-vue'
import { getAdminReports, getAdminReportDetail, updateReport } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const loading = ref(false)

const editVisible = ref(false)
const editForm = ref<any>(null)
const editFormRef = ref<any>(null)
const editId = ref(0)
const saving = ref(false)

const filters = reactive({ businessName: '', month: '', status: '' })

const photoList = computed(() => {
  if (!editForm.value?.photo_urls) return []
  return editForm.value.photo_urls.split(',').filter(Boolean)
})

function getPhotoSrc(url: string): string {
  if (url.startsWith('http')) return url
  return url.startsWith('/') ? url : '/' + url
}

function removePhoto(i: number) {
  const photos = [...photoList.value]
  photos.splice(i, 1)
  editForm.value.photo_urls = photos.join(',')
}

function openImageViewer(photoUrls: string[], index: number) {
  const urls = photoUrls.map(u => getPhotoSrc(u))
  const win = window.open('', '_blank', 'width=900,height=700,scrollbars=yes,toolbar=no')
  if (!win) return
  const imgs = urls.map((u, i) =>
    `<img src="${u}" style="max-width:100%;display:${i === index ? 'block' : 'none'}" data-idx="${i}" />`
  ).join('\n')
  const nav = urls.length > 1 ? `
    <div class="nav">
      <button id="btnPrev" ${index === 0 ? 'disabled' : ''}>上一张</button>
      <span id="counter">${index + 1} / ${urls.length}</span>
      <button id="btnNext" ${index === urls.length - 1 ? 'disabled' : ''}>下一张</button>
    </div>` : ''
  win.document.write(`
    <!DOCTYPE html>
    <html><head><meta charset="utf-8"><title>查看原图</title>
    <style>
      * { margin:0; padding:0; box-sizing:border-box; }
      body { background:#000; display:flex; flex-direction:column; align-items:center; justify-content:center; min-height:100vh; font-family:sans-serif; }
      .toolbar { position:fixed; top:0; left:0; right:0; display:flex; justify-content:flex-end; padding:8px 16px; background:rgba(0,0,0,0.7); z-index:10; }
      .btn-close { background:#e74c3c; color:#fff; border:none; padding:8px 24px; font-size:16px; border-radius:4px; cursor:pointer; }
      .btn-close:hover { background:#c0392b; }
      .img-wrap { flex:1; display:flex; align-items:center; justify-content:center; padding:60px 20px 20px; width:100%; }
      img { max-width:100%; max-height:80vh; object-fit:contain; }
      .nav { position:fixed; bottom:0; left:0; right:0; display:flex; align-items:center; justify-content:center; gap:16px; padding:12px; background:rgba(0,0,0,0.7); z-index:10; }
      .nav button { padding:6px 20px; border:1px solid #fff; background:transparent; color:#fff; border-radius:4px; cursor:pointer; font-size:14px; }
      .nav button:disabled { opacity:0.3; cursor:default; }
      .nav button:hover:not(:disabled) { background:rgba(255,255,255,0.2); }
      .nav span { color:#ccc; font-size:14px; }
    </style></head><body>
    <div class="toolbar"><button class="btn-close" onclick="window.close()">关闭</button></div>
    <div class="img-wrap">${imgs}</div>
    ${nav}
    <script>
      var imgs = document.querySelectorAll('img[data-idx]');
      var cur = ${index};
      function show(i) {
        imgs.forEach(function(im, j) { im.style.display = j === i ? 'block' : 'none'; });
        cur = i;
        var c = document.getElementById('counter'); if (c) c.textContent = (i + 1) + ' / ${urls.length}';
        var p = document.getElementById('btnPrev'); if (p) p.disabled = i === 0;
        var n = document.getElementById('btnNext'); if (n) n.disabled = i === ${urls.length - 1};
      }
      document.getElementById('btnPrev') && document.getElementById('btnPrev').addEventListener('click', function() { if (cur > 0) show(cur - 1); });
      document.getElementById('btnNext') && document.getElementById('btnNext').addEventListener('click', function() { if (cur < ${urls.length - 1}) show(cur + 1); });
      document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape') window.close();
        if (e.key === 'ArrowLeft' && cur > 0) show(cur - 1);
        if (e.key === 'ArrowRight' && cur < ${urls.length - 1}) show(cur + 1);
      });
    <\\/script></body></html>`);
  win.document.close()
}

async function loadData() {
  loading.value = true
  try {
    const res = await getAdminReports({
      page: page.value, page_size: pageSize.value,
      business_name: filters.businessName, month: filters.month, status: filters.status
    })
    list.value = res.data.list
    total.value = res.data.total
  } finally { loading.value = false }
}

async function showEdit(row: any) {
  try {
    const res = await getAdminReportDetail(row.id)
    editForm.value = { ...res.data }
    editId.value = row.id
    editVisible.value = true
  } catch {}
}

async function saveEdit() {
  saving.value = true
  try {
    await updateReport(editId.value, editForm.value)
    editVisible.value = false
    loadData()
  } catch (err: any) {
    // error handled by interceptor
  } finally {
    saving.value = false
  }
}

onMounted(() => loadData())
</script>

<style scoped>
.report-page { }
.search-bar { display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap; }
.photo-grid { display: flex; gap: 12px; flex-wrap: wrap; }
.photo-item { position: relative; }
.photo-delete { position: absolute; top: -6px; right: -6px; color: #f56c6c; cursor: pointer; font-size: 18px; background: #fff; border-radius: 50%; }
.table-photos { display: flex; align-items: center; }
.photo-thumb { width: 36px; height: 36px; border-radius: 4px; cursor: pointer; margin-right: 4px; object-fit: cover; transition: opacity 0.2s; }
.photo-thumb:hover { opacity: 0.75; }
.photo-thumb-lg { width: 120px; height: 90px; border-radius: 8px; cursor: pointer; object-fit: cover; transition: opacity 0.2s; }
.photo-thumb-lg:hover { opacity: 0.85; }
.photo-more { font-size: 12px; color: #909399; margin-left: 2px; }
</style>
