<template>
  <div class="system-info">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <div class="toolbar">
        <h2>系统信息</h2>
      </div>

      <!-- 内容区域 -->
      <div class="info-content">
        <div class="version-block">
          <div class="version-list">
            <div class="version-item">
              <span class="label">博客系统</span>
              <span class="value">FlecBlog</span>
            </div>
            <div class="version-item">
              <span class="label">当前版本</span>
              <span class="value">{{ staticInfo.app_version || 'dev' }}</span>
            </div>
            <div class="version-item">
              <span class="label">最新版本</span>
              <span class="version-value">
                <template v-if="latestVersion">
                  <span v-if="hasNewVersion" class="value version-link" @click="showUpdateDialog">
                    <span class="new-version-dot"></span>{{ latestVersion }}
                  </span>
                  <span v-else class="value">{{ latestVersion }}</span>
                </template>
                <el-button
                  type="primary"
                  size="small"
                  link
                  :loading="checking"
                  @click="handleCheckUpdate"
                >
                  {{ checking ? '检查中...' : '检查更新' }}
                </el-button>
              </span>
            </div>
          </div>
        </div>

        <div class="info-grid">
          <!-- 服务器 -->
          <div class="info-section">
            <div class="section-header">
              <el-icon class="icon-orange">
                <Monitor />
              </el-icon>
              <span>服务器</span>
            </div>
            <div class="section-body">
              <div class="info-item">
                <span class="label">主机名</span>
                <span class="value">{{ staticInfo.hostname }}</span>
              </div>
              <div class="info-item">
                <span class="label">操作系统</span>
                <span class="value">{{ staticInfo.os }}</span>
              </div>
              <div class="info-item">
                <span class="label">IP</span>
                <span class="value">{{ staticInfo.server_ip || 'N/A' }}</span>
              </div>
              <div class="info-item">
                <span class="label">时区</span>
                <span class="value">{{ staticInfo.timezone || 'N/A' }}</span>
              </div>
              <div class="info-item">
                <span class="label">运行时间</span>
                <span class="value">{{ formatDays(dynamicInfo.host_uptime) }}</span>
              </div>
            </div>
          </div>

          <!-- CPU -->
          <div class="info-section">
            <div class="section-header">
              <el-icon class="icon-blue">
                <Cpu />
              </el-icon>
              <span>CPU</span>
            </div>
            <div class="section-body">
              <div class="info-item">
                <span class="label">核心数</span>
                <span class="value">{{ staticInfo.cpu_core }} 核</span>
              </div>
              <div class="info-item">
                <span class="label">使用率</span>
                <el-progress
                  :percentage="Math.round(dynamicInfo.cpu_usage || 0)"
                  :stroke-width="6"
                  :color="getProgressColor(Math.round(dynamicInfo.cpu_usage || 0))"
                  style="width: 120px"
                />
              </div>
              <div class="info-item">
                <span class="label">型号</span>
                <span class="value">{{ staticInfo.cpu_model || 'N/A' }}</span>
              </div>
              <div class="info-item">
                <span class="label">架构</span>
                <span class="value">{{ staticInfo.cpu_arch }}</span>
              </div>
              <div class="info-item">
                <span class="label">系统负载</span>
                <span class="value"
                  >{{ dynamicInfo.load_1?.toFixed(2) || 'N/A' }} /
                  {{ dynamicInfo.load_5?.toFixed(2) || 'N/A' }} /
                  {{ dynamicInfo.load_15?.toFixed(2) || 'N/A' }}</span
                >
              </div>
            </div>
          </div>

          <!-- 内存 -->
          <div class="info-section">
            <div class="section-header">
              <el-icon class="icon-green">
                <Coin />
              </el-icon>
              <span>内存</span>
            </div>
            <div class="section-body">
              <div class="info-item">
                <span class="label">总容量</span>
                <span class="value"
                  >{{ formatBytes(dynamicInfo.memory_used) }} /
                  {{ formatBytes(staticInfo.memory_total) }}</span
                >
              </div>
              <div class="info-item">
                <span class="label">使用率</span>
                <el-progress
                  :percentage="calcPercent(dynamicInfo.memory_used, staticInfo.memory_total)"
                  :stroke-width="6"
                  :color="
                    getProgressColor(calcPercent(dynamicInfo.memory_used, staticInfo.memory_total))
                  "
                  style="width: 120px"
                />
              </div>
              <div class="info-item">
                <span class="label">未使用</span>
                <span class="value">{{ formatBytes(dynamicInfo.memory_available) }}</span>
              </div>
              <div class="info-item">
                <span class="label">Swap</span>
                <span class="value"
                  >{{ formatBytes(dynamicInfo.swap_used) }} /
                  {{ formatBytes(staticInfo.swap_total) }}</span
                >
              </div>
            </div>
          </div>

          <!-- 磁盘 -->
          <div class="info-section">
            <div class="section-header">
              <el-icon class="icon-red">
                <FolderOpened />
              </el-icon>
              <span>磁盘</span>
            </div>
            <div class="section-body">
              <div class="info-item">
                <span class="label">总容量</span>
                <span class="value">{{ formatBytes(staticInfo.disk_total) }}</span>
              </div>
              <div class="info-item">
                <span class="label">使用率</span>
                <el-progress
                  :percentage="calcPercent(dynamicInfo.disk_used, staticInfo.disk_total)"
                  :stroke-width="6"
                  :color="
                    getProgressColor(calcPercent(dynamicInfo.disk_used, staticInfo.disk_total))
                  "
                  style="width: 120px"
                />
              </div>
              <div class="info-item">
                <span class="label">已使用</span>
                <span class="value">{{ formatBytes(dynamicInfo.disk_used) }}</span>
              </div>
              <div class="info-item">
                <span class="label">未使用</span>
                <span class="value">{{ formatBytes(dynamicInfo.disk_free) }}</span>
              </div>
            </div>
          </div>

          <!-- 数据库 -->
          <div class="info-section">
            <div class="section-header">
              <el-icon class="icon-purple">
                <DataLine />
              </el-icon>
              <span>数据库</span>
            </div>
            <div class="section-body">
              <div class="info-item">
                <span class="label">类型</span>
                <span class="value">{{ staticInfo.db_type }}</span>
              </div>
              <div class="info-item">
                <span class="label">状态</span>
                <el-tag
                  :type="dynamicInfo.db_status === '正常' ? 'success' : 'danger'"
                  size="small"
                >
                  {{ dynamicInfo.db_status }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">大小</span>
                <span class="value">{{ formatBytes(dynamicInfo.db_size) }}</span>
              </div>
              <div class="info-item">
                <span class="label">表数量</span>
                <span class="value">{{ staticInfo.db_tables }}</span>
              </div>
              <div class="info-item">
                <span class="label">连接数</span>
                <span class="value">{{ dynamicInfo.db_conn_count }}</span>
              </div>
            </div>
          </div>

          <!-- 外部连通 -->
          <div class="info-section">
            <div class="section-header">
              <el-icon class="icon-cyan">
                <Connection />
              </el-icon>
              <span>外部连通</span>
            </div>
            <div class="section-body">
              <div class="info-item">
                <span class="label">文件存储</span>
                <el-tag
                  :type="staticInfo.storage_status === '正常' ? 'success' : 'danger'"
                  size="small"
                >
                  {{ staticInfo.storage_status }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">邮箱通知</span>
                <el-tag
                  :type="
                    staticInfo.email_status === '正常'
                      ? 'success'
                      : staticInfo.email_status === '未配置'
                        ? 'info'
                        : 'danger'
                  "
                  size="small"
                >
                  {{ staticInfo.email_status }}
                </el-tag>
              </div>
              <div class="info-item">
                <span class="label">飞书交互</span>
                <el-tag
                  :type="
                    staticInfo.feishu_status === '正常'
                      ? 'success'
                      : staticInfo.feishu_status === '未配置'
                        ? 'info'
                        : 'danger'
                  "
                  size="small"
                >
                  {{ staticInfo.feishu_status }}
                </el-tag>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 更新内容弹窗 -->
    <el-dialog
      v-model="updateDialogVisible"
      title="版本更新内容"
      width="600px"
      :close-on-click-modal="false"
    >
      <div v-if="updatingVersions.length > 0">
        <div v-for="version in updatingVersions" :key="version.id" class="version-item-detail">
          <div class="version-header">
            <span class="version-tag">{{ version.version }}</span>
            <span class="version-date">{{ version.date }}</span>
          </div>
          <div class="version-changes" v-html="formatChanges(version.changes)"></div>
        </div>
      </div>
      <div v-else class="no-changes">暂无更新内容</div>
      <template #footer>
        <el-button type="primary" @click="updateDialogVisible = false">关闭</el-button>
        <el-button type="primary" @click="goToGitHub">前往 GitHub</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { ElMessage } from 'element-plus';
import { Monitor, Cpu, Coin, DataLine, FolderOpened, Connection } from '@element-plus/icons-vue';
import { getSystemStatic, getSystemDynamic, checkUpdate } from '@/api/system';
import type { SystemStatic, SystemDynamic } from '@/types/system';
import type { VersionInfo } from '@/api/system';

let refreshTimer: ReturnType<typeof setInterval> | null = null;

const staticInfo = ref<SystemStatic>({
  cpu_core: 0,
  cpu_model: '',
  cpu_arch: '',
  hostname: '',
  os: '',
  server_ip: '',
  timezone: '',
  db_type: '',
  memory_total: 0,
  swap_total: 0,
  disk_total: 0,
  db_tables: 0,
  storage_status: '',
  email_status: '',
  feishu_status: '',
  app_version: '',
});

const dynamicInfo = ref<SystemDynamic>({
  cpu_usage: 0,
  load_1: 0,
  load_5: 0,
  load_15: 0,
  memory_used: 0,
  memory_available: 0,
  swap_used: 0,
  host_uptime: 0,
  disk_used: 0,
  disk_free: 0,
  db_status: '',
  db_size: 0,
  db_conn_count: 0,
  version_latest_version: '',
  version_last_check_error: '',
});

const checking = ref(false);
const updateDialogVisible = ref(false);
const panelVersions = ref<VersionInfo[]>([]);
const updatingVersions = ref<VersionInfo[]>([]);

const latestVersion = computed(() => dynamicInfo.value.version_latest_version);

const currentVersion = computed(() => staticInfo.value.app_version || 'dev');

const hasNewVersion = computed(() => {
  if (!latestVersion.value || currentVersion.value === 'dev') return false;
  return compareVersion(latestVersion.value, currentVersion.value) > 0;
});

const compareVersion = (a: string, b: string): number => {
  const parseVersion = (v: string): number[] => {
    return v
      .replace(/^v/, '')
      .split('.')
      .map(n => parseInt(n, 10) || 0);
  };
  const av = parseVersion(a);
  const bv = parseVersion(b);
  for (let i = 0; i < Math.max(av.length, bv.length); i++) {
    const an = av[i] || 0;
    const bn = bv[i] || 0;
    if (an > bn) return 1;
    if (an < bn) return -1;
  }
  return 0;
};

const showUpdateDialog = async () => {
  updatingVersions.value = panelVersions.value
    .filter(version => {
      const comparedCurrent = compareVersion(version.version, currentVersion.value);
      return comparedCurrent > 0;
    })
    .sort((a, b) => compareVersion(b.version, a.version));

  updateDialogVisible.value = true;
};

const formatChanges = (changes: string): string => {
  // 简单的换行处理，支持换行
  return changes.replace(/\n/g, '<br>');
};

const goToGitHub = () => {
  window.open('https://github.com/talen8/FlecBlog/releases', '_blank');
};

const fetchStaticInfo = async () => {
  try {
    staticInfo.value = await getSystemStatic();
  } catch (_error) {
    ElMessage.error('获取系统静态信息失败');
  }
};

const fetchDynamicInfo = async () => {
  try {
    dynamicInfo.value = await getSystemDynamic();
  } catch (_error) {
    ElMessage.error('获取系统动态信息失败');
  }
};

const handleCheckUpdate = async () => {
  checking.value = true;
  try {
    const result = await checkUpdate();
    panelVersions.value = result.versions;
    await fetchDynamicInfo();
    if (result.has_update) {
      ElMessage.success('发现新版本');
    } else {
      ElMessage.success('已是最新版本');
    }
  } catch (_error) {
    ElMessage.error('检查更新失败');
  } finally {
    checking.value = false;
  }
};

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 B';
  const unit = 1024;
  const units = ['B', 'KB', 'MB', 'GB', 'TB', 'PB'];
  const i = Math.floor(Math.log(bytes) / Math.log(unit));
  return (bytes / Math.pow(unit, i)).toFixed(1) + ' ' + units[i];
};

const calcPercent = (used: number, total: number): number => {
  if (total === 0) return 0;
  return Math.round((used / total) * 100);
};

const formatDays = (seconds: number): string => {
  const days = Math.floor(seconds / 86400);
  return `${days} 天`;
};

const getProgressColor = (percentage: number): string => {
  if (percentage < 50) return '#67c23a';
  if (percentage < 80) return '#e6a23c';
  return '#f56c6c';
};

onMounted(() => {
  fetchStaticInfo();
  fetchDynamicInfo();
  refreshTimer = setInterval(fetchDynamicInfo, 10000);
});

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer);
    refreshTimer = null;
  }
});
</script>

<style lang="scss" scoped>
.system-info {
  height: 100%;

  :deep(.el-card) {
    height: 100%;
    display: flex;
    flex-direction: column;

    .el-card__body {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
  }
}

.toolbar {
  margin-bottom: 12px;

  h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 500;
  }
}

.info-content {
  flex: 1;
  overflow-y: auto;

  &::-webkit-scrollbar {
    display: none;
  }

  -ms-overflow-style: none;
  scrollbar-width: none;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.version-block {
  margin-bottom: 20px;
  padding: 16px 20px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  background: #fafafa;
}

.version-list {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  column-gap: 32px;
  row-gap: 12px;
}

.version-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;

  .label {
    color: #909399;
    font-size: 13px;
    line-height: 1.4;
  }

  .value {
    color: #303133;
    font-size: 14px;
    line-height: 1.6;
    word-break: break-all;
  }
}

.version-value {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.info-section {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;

  .section-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 16px;
    background: #f5f7fa;
    border-bottom: 1px solid #e4e7ed;
    font-weight: 500;
    font-size: 14px;

    .icon-blue {
      color: #409eff;
      font-size: 16px;
    }

    .icon-green {
      color: #67c23a;
      font-size: 16px;
    }

    .icon-orange {
      color: #e6a23c;
      font-size: 16px;
    }

    .icon-purple {
      color: #a855f7;
      font-size: 16px;
    }

    .icon-cyan {
      color: #06b6d4;
      font-size: 16px;
    }

    .icon-red {
      color: #f56c6c;
      font-size: 16px;
    }
  }

  .section-body {
    padding: 12px 16px;
  }
}

.info-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px dashed #ebeef5;

  &:last-child {
    border-bottom: none;
  }

  .label {
    color: #909399;
    font-size: 14px;
    flex-shrink: 0;
  }

  .value {
    color: #303133;
    font-size: 14px;
    text-align: right;
    word-break: break-all;
  }

  :deep(.el-progress__text) {
    min-width: auto;
  }
}

.version-link {
  cursor: pointer;
  color: #409eff;
  text-decoration: underline;

  &:hover {
    color: #79bbff;
  }

  .new-version-dot {
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: #f56c6c;
    margin-right: 6px;
    vertical-align: middle;
  }
}

.version-item-detail {
  border-bottom: 1px solid #ebeef5;
  padding: 16px 0;

  &:last-child {
    border-bottom: none;
  }

  .version-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;

    .version-tag {
      display: inline-flex;
      align-items: center;
      padding: 2px 8px;
      background: #ecf5ff;
      color: #409eff;
      border-radius: 4px;
      font-size: 14px;
      font-weight: 500;
    }

    .version-date {
      color: #909399;
      font-size: 13px;
    }
  }

  .version-changes {
    color: #303133;
    font-size: 14px;
    line-height: 1.8;
  }
}

.no-changes {
  text-align: center;
  color: #909399;
  padding: 40px 0;
}

@media (max-width: 768px) {
  .toolbar {
    h2 {
      font-size: 18px;
    }
  }

  .version-block {
    padding: 14px 16px;
  }

  .version-list {
    grid-template-columns: 1fr;
    row-gap: 12px;
  }
}
</style>
