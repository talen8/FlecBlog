<template>
  <div class="dashboard">
    <!-- 顶部区域 -->
    <el-card class="top-card" shadow="hover">
      <div class="top-content">
        <!-- 左侧：个人信息 -->
        <div class="profile-section">
          <el-avatar :size="64" :src="userAvatar" class="avatar" />
          <div class="profile-info">
            <h2 class="greeting">{{ greeting }}，{{ nickName }}，今天又是充满活力的一天！</h2>
            <p class="weather-info">{{ hitokoto }}</p>
          </div>
        </div>

        <!-- 右侧：统计数字 -->
        <div class="stats-section">
          <div class="stat-item">
            <div class="stat-label">文章</div>
            <div class="stat-value">{{ dashboardData.total_articles }}</div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-label">友链</div>
            <div class="stat-value">{{ dashboardData.total_friends }}</div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-label">动态</div>
            <div class="stat-value">{{ dashboardData.total_moments }}</div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 中间：四个概况卡片 -->
    <el-row :gutter="20" class="overview-cards">
      <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6">
        <el-card class="overview-card" shadow="hover">
          <div class="card-left">
            <div class="card-icon icon-purple">
              <el-icon>
                <View />
              </el-icon>
            </div>
          </div>
          <div class="card-right">
            <div class="card-title">浏览量</div>
            <div class="card-value">
              {{ formatNumber(dashboardData.total_views) }}
            </div>
            <div class="card-stats">
              <span class="today-value">今日: {{ dashboardData.today_views }}</span>
              <span class="growth-rate" :class="getGrowthClass(dashboardData.views_growth)">
                <el-icon v-if="dashboardData.views_growth > 0">
                  <CaretTop />
                </el-icon>
                <el-icon v-else-if="dashboardData.views_growth < 0">
                  <CaretBottom />
                </el-icon>
                {{ Math.abs(dashboardData.views_growth) }}%
              </span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6">
        <el-card class="overview-card" shadow="hover">
          <div class="card-left">
            <div class="card-icon icon-blue">
              <el-icon>
                <User />
              </el-icon>
            </div>
          </div>
          <div class="card-right">
            <div class="card-title">访客量</div>
            <div class="card-value">
              {{ formatNumber(dashboardData.total_visitors) }}
            </div>
            <div class="card-stats">
              <span class="today-value">今日: {{ dashboardData.today_visitors }}</span>
              <span class="growth-rate" :class="getGrowthClass(dashboardData.visitors_growth)">
                <el-icon v-if="dashboardData.visitors_growth > 0">
                  <CaretTop />
                </el-icon>
                <el-icon v-else-if="dashboardData.visitors_growth < 0">
                  <CaretBottom />
                </el-icon>
                {{ Math.abs(dashboardData.visitors_growth) }}%
              </span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6">
        <el-card class="overview-card" shadow="hover">
          <div class="card-left">
            <div class="card-icon icon-green">
              <el-icon>
                <ChatDotRound />
              </el-icon>
            </div>
          </div>
          <div class="card-right">
            <div class="card-title">评论数</div>
            <div class="card-value">
              {{ formatNumber(dashboardData.total_comments) }}
            </div>
            <div class="card-stats">
              <span class="today-value">今日: {{ dashboardData.today_comments }}</span>
              <span class="growth-rate" :class="getGrowthClass(dashboardData.comments_growth)">
                <el-icon v-if="dashboardData.comments_growth > 0">
                  <CaretTop />
                </el-icon>
                <el-icon v-else-if="dashboardData.comments_growth < 0">
                  <CaretBottom />
                </el-icon>
                {{ Math.abs(dashboardData.comments_growth) }}%
              </span>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="12" :md="12" :lg="6" :xl="6">
        <el-card class="overview-card" shadow="hover">
          <div class="card-left">
            <div class="card-icon icon-orange">
              <el-icon>
                <UserFilled />
              </el-icon>
            </div>
          </div>
          <div class="card-right">
            <div class="card-title">用户数</div>
            <div class="card-value">
              {{ formatNumber(dashboardData.total_users) }}
            </div>
            <div class="card-stats">
              <span class="today-value">今日: {{ dashboardData.today_users }}</span>
              <span class="growth-rate" :class="getGrowthClass(dashboardData.users_growth)">
                <el-icon v-if="dashboardData.users_growth > 0">
                  <CaretTop />
                </el-icon>
                <el-icon v-else-if="dashboardData.users_growth < 0">
                  <CaretBottom />
                </el-icon>
                {{ Math.abs(dashboardData.users_growth) }}%
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域1 -->
    <el-row :gutter="20" class="charts-section">
      <el-col :xs="24" :sm="24" :md="16" :lg="15" :xl="15">
        <el-card shadow="hover">
          <template #header>
            <div class="chart-header">
              <span>访问趋势</span>
              <el-radio-group v-model="trendType" size="small" @change="fetchTrendData">
                <el-radio-button value="daily">日</el-radio-button>
                <el-radio-button value="monthly">月</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="trendChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="8" :lg="9" :xl="9">
        <el-card shadow="hover">
          <template #header>
            <div class="chart-header">
              <span>分类统计</span>
            </div>
          </template>
          <div ref="pieChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域2 -->
    <el-row :gutter="20" class="charts-section">
      <el-col :xs="24" :sm="24" :md="8" :lg="9" :xl="9">
        <el-card shadow="hover">
          <template #header>
            <div class="chart-header">
              <span>标签统计</span>
            </div>
          </template>
          <div ref="tagCloudRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="24" :md="16" :lg="15" :xl="15">
        <el-card shadow="hover">
          <template #header>
            <div class="chart-header">
              <span>文章贡献</span>
              <el-select
                v-model="selectedYear"
                size="small"
                style="width: 100px"
                @change="fetchContributionData"
              >
                <el-option v-for="year in availableYears" :key="year" :label="year" :value="year" />
              </el-select>
            </div>
          </template>
          <div ref="calendarChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷访问 -->
    <el-card class="quick-access-card" shadow="hover">
      <template #header>
        <div class="quick-access-header">
          <span>快捷访问</span>
        </div>
      </template>
      <div class="quick-access-content">
        <div class="quick-links">
          <div class="link-item" @click="openLink('https://github.com/talen8')">
            <span class="link-text">GitHub</span>
            <el-icon class="link-icon">
              <Right />
            </el-icon>
          </div>
          <div
            class="link-item"
            @click="openLink('https://ccnlf8xcz6k3.feishu.cn/wiki/space/7618178485001046989')"
          >
            <span class="link-text">文档</span>
            <el-icon class="link-icon">
              <Right />
            </el-icon>
          </div>
          <div class="link-item" @click="openLink('https://qm.qq.com/q/YDynuTZq4o')">
            <span class="link-text">社群</span>
            <el-icon class="link-icon">
              <Right />
            </el-icon>
          </div>
          <div class="link-item" @click="openLink('https://blog.talen.top')">
            <span class="link-text">博客</span>
            <el-icon class="link-icon">
              <Right />
            </el-icon>
          </div>
        </div>
        <div class="quick-illustration">
          <img src="@/assets/img/dashboard.png" alt="dashboard" />
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onUnmounted, nextTick } from 'vue';
import {
  getDashboardStats,
  getTrendData,
  getCategoryStats,
  getTagStats,
  getArticleContribution,
} from '@/api/stats';
import type {
  DashboardStats,
  TrendDataItem,
  CategoryStats,
  TagStats,
  ArticleContribution,
} from '@/types/stats';
import {
  View,
  User,
  ChatDotRound,
  UserFilled,
  CaretTop,
  CaretBottom,
  Right,
} from '@element-plus/icons-vue';
import * as echarts from 'echarts';
import 'echarts-wordcloud';
import { getToday, getDaysAgo, getMonthsAgo, generateDateSeries } from '@/utils/date';
import { getUserInfo } from '@/utils/auth';

// 响应式断点

const DEFAULT_AVATAR = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png';
const nickName = computed(() => getUserInfo()?.nickname || 'Admin');
const userAvatar = computed(() => getUserInfo()?.avatar || DEFAULT_AVATAR);
const hitokoto = ref('加载中...');

const dashboardData = ref<DashboardStats>({
  total_articles: 0,
  total_friends: 0,
  total_moments: 0,
  total_views: 0,
  total_visitors: 0,
  total_comments: 0,
  total_users: 0,
  today_views: 0,
  today_visitors: 0,
  today_comments: 0,
  today_users: 0,
  views_growth: 0,
  visitors_growth: 0,
  comments_growth: 0,
  users_growth: 0,
});

const trendType = ref<'daily' | 'monthly'>('daily');
const trendData = ref<TrendDataItem[]>([]);
const categoryData = ref<CategoryStats[]>([]);
const tagData = ref<TagStats[]>([]);
const contributionData = ref<ArticleContribution[]>([]);
const selectedYear = ref(new Date().getFullYear());
const trendChartRef = ref<HTMLElement>();
const pieChartRef = ref<HTMLElement>();
const tagCloudRef = ref<HTMLElement>();
const calendarChartRef = ref<HTMLElement>();
let trendChart: echarts.ECharts | null = null;
let pieChart: echarts.ECharts | null = null;
let tagCloud: echarts.ECharts | null = null;
let calendarChart: echarts.ECharts | null = null;

const CHART_COLORS = [
  '#5470c6',
  '#91cc75',
  '#fac858',
  '#ee6666',
  '#73c0de',
  '#3ba272',
  '#fc8452',
  '#9a60b4',
  '#ea7ccc',
];

// 生成可选年份列表（从2024年到当前年份）
const availableYears = computed(() => {
  const currentYear = new Date().getFullYear();
  const startYear = 2024;
  const years: number[] = [];

  // 从当前年份倒序到2024年
  for (let year = currentYear; year >= startYear; year--) {
    years.push(year);
  }

  return years;
});

const greeting = computed(() => {
  const hour = new Date().getHours();
  if (hour < 6) return '凌晨好';
  if (hour < 9) return '早安';
  if (hour < 12) return '上午好';
  if (hour < 14) return '中午好';
  if (hour < 17) return '下午好';
  if (hour < 19) return '傍晚好';
  if (hour < 22) return '晚上好';
  return '夜深了';
});

// 格式化数字
const formatNumber = (num: number) => {
  if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
  return num.toString();
};

// 获取增长率样式类
const getGrowthClass = (growth: number) => {
  if (growth > 0) return 'positive';
  if (growth < 0) return 'negative';
  return 'neutral';
};

// 打开链接
const openLink = (url: string) => {
  window.open(url, '_blank');
};

// 获取随机一言
const fetchHitokoto = async () => {
  try {
    const response = await fetch('https://api.pearktrue.cn/api/hitokoto/');
    const text = await response.text();
    hitokoto.value = text || '获取一言失败';
  } catch (error) {
    console.error('获取一言失败:', error);
    hitokoto.value = '获取一言失败';
  }
};

// 获取仪表板数据
const fetchDashboardData = async () => {
  try {
    dashboardData.value = await getDashboardStats();
  } catch (error) {
    console.error('获取仪表板数据失败', error);
  }
};

const fetchTrendData = async () => {
  const endDate = getToday();
  const startDate = trendType.value === 'daily' ? getDaysAgo(6) : getMonthsAgo(6);

  try {
    trendData.value = await getTrendData({
      start_date: startDate,
      end_date: endDate,
      type: trendType.value,
    });
    renderTrendChart();
  } catch (error) {
    console.error('获取趋势数据失败', error);
  }
};

let resizeTimer: number | null = null;

const handleResize = () => {
  if (resizeTimer) clearTimeout(resizeTimer);
  resizeTimer = window.setTimeout(() => {
    trendChart?.resize();
    pieChart?.resize();
    tagCloud?.resize();
    calendarChart?.resize();
  }, 300);
};

const initAllCharts = () => {
  if (trendChartRef.value) trendChart = echarts.init(trendChartRef.value);
  if (pieChartRef.value) pieChart = echarts.init(pieChartRef.value);
  if (tagCloudRef.value) tagCloud = echarts.init(tagCloudRef.value);
  if (calendarChartRef.value) calendarChart = echarts.init(calendarChartRef.value);

  window.addEventListener('resize', handleResize);
};

const renderTrendChart = () => {
  if (!trendChart) return;

  const endDate = getToday();
  const startDate = trendType.value === 'daily' ? getDaysAgo(6) : getMonthsAgo(6);
  const format = trendType.value === 'daily' ? 'YYYY-MM-DD' : 'YYYY-MM';
  const unit = trendType.value === 'daily' ? 'day' : 'month';
  const allDates = generateDateSeries(startDate, endDate, unit, format, 7);

  const dataMap = new Map<string, { pv: number; uv: number }>();
  trendData.value.forEach(item => {
    dataMap.set(item.date, { pv: item.pv_count, uv: item.uv_count });
  });

  const pvData = allDates.map(date => dataMap.get(date)?.pv || 0);
  const uvData = allDates.map(date => dataMap.get(date)?.uv || 0);

  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross',
      },
    },
    legend: {
      data: ['浏览量', '访客量'],
      bottom: 0,
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '60px',
      top: '40px',
      containLabel: true,
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: allDates,
    },
    yAxis: [
      {
        type: 'value',
        name: '浏览量',
        position: 'left',
        axisLabel: {
          formatter: '{value}',
        },
      },
      {
        type: 'value',
        name: '访客量',
        position: 'right',
        axisLabel: {
          formatter: '{value}',
        },
      },
    ],
    series: [
      {
        name: '浏览量',
        type: 'line',
        smooth: true,
        yAxisIndex: 0,
        data: pvData,
        itemStyle: {
          color: '#7c7cff',
        },
        lineStyle: {
          width: 2,
        },
        symbolSize: 6,
      },
      {
        name: '访客量',
        type: 'line',
        smooth: true,
        yAxisIndex: 1,
        data: uvData,
        itemStyle: {
          color: '#409eff',
        },
        lineStyle: {
          width: 2,
        },
        symbolSize: 6,
      },
    ],
  };

  trendChart.setOption(option);
};

const fetchCategoryData = async () => {
  try {
    categoryData.value = await getCategoryStats();
    renderPieChart();
  } catch (error) {
    console.error('获取分类统计数据失败', error);
  }
};

const renderPieChart = () => {
  if (!pieChart) return;

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)',
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center',
      textStyle: {
        fontSize: 12,
      },
      formatter: (name: string) => {
        const item = categoryData.value.find(d => d.name === name);
        return `${name} (${item?.count || 0})`;
      },
    },
    series: [
      {
        name: '文章分类',
        type: 'pie',
        radius: ['40%', '70%'],
        center: ['40%', '50%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2,
        },
        label: {
          show: false,
          position: 'center',
        },
        emphasis: {
          label: {
            show: true,
            fontSize: 18,
            fontWeight: 'bold',
            formatter: '{b}\n{c}篇',
          },
        },
        labelLine: {
          show: false,
        },
        data: categoryData.value.map((item, index) => ({
          value: item.count,
          name: item.name,
          itemStyle: {
            color: CHART_COLORS[index % CHART_COLORS.length],
          },
        })),
      },
    ],
  };

  pieChart.setOption(option);
};

const fetchTagData = async () => {
  try {
    tagData.value = await getTagStats();
    renderTagCloud();
  } catch (error) {
    console.error('获取标签统计数据失败', error);
  }
};

const renderTagCloud = () => {
  if (!tagCloud) return;

  const option = {
    tooltip: {
      formatter: '{b}: {c}篇',
    },
    series: [
      {
        type: 'wordCloud',
        shape: 'circle',
        sizeRange: [14, 35],
        rotationRange: [0, 0],
        gridSize: 5,
        drawOutOfBound: false,
        textStyle: {
          fontFamily: 'sans-serif',
          fontWeight: '500',
          color: function () {
            return CHART_COLORS[Math.floor(Math.random() * CHART_COLORS.length)];
          },
        },
        data: tagData.value.map(item => ({
          name: item.name,
          value: item.count,
        })),
      },
    ],
  };

  tagCloud.setOption(option);
};

const fetchContributionData = async () => {
  try {
    // 调用API时传入year参数，获取指定年份的数据
    const data = await getArticleContribution({ year: selectedYear.value });
    contributionData.value = data || [];
    // 等待响应式数据更新完成后再渲染图表
    await nextTick();
    renderCalendarChart();
  } catch (error) {
    console.error('获取文章贡献数据失败', error);
    contributionData.value = [];
    if (calendarChart) {
      calendarChart.clear();
    }
  }
};

const renderCalendarChart = () => {
  if (!calendarChart) return;

  // 确保数据为数组，即使为空也显示日历网格
  const chartData = contributionData.value || [];

  // 计算最大值，如果没有数据则设为1（避免visualMap显示异常）
  const maxCount = chartData.length > 0 ? Math.max(...chartData.map(item => item.count), 1) : 1;

  const option = {
    tooltip: {
      position: 'top',
      formatter: (params: { data: [string, number] }) => {
        return `${params.data[0]}<br/>文章数: ${params.data[1]}篇`;
      },
    },
    visualMap: {
      min: 0,
      max: maxCount,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      top: 20,
    },
    calendar: {
      range: selectedYear.value,
      left: 55,
      right: 10,
      top: 130,
      bottom: 80,
      cellSize: 12,
      monthLabel: {
        fontSize: 11,
        nameMap: 'cn',
      },
      dayLabel: {
        fontSize: 11,
        nameMap: 'cn',
      },
    },
    series: {
      type: 'heatmap',
      coordinateSystem: 'calendar',
      data: chartData.map(item => [item.date, item.count]),
    },
  };

  // 使用 notMerge: true 确保完全替换配置
  calendarChart.setOption(option, { notMerge: true });
};

onMounted(async () => {
  await fetchDashboardData();
  await nextTick();
  initAllCharts();
  await fetchTrendData();
  await fetchCategoryData();
  await fetchTagData();
  await fetchContributionData();
  await fetchHitokoto();
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
  if (resizeTimer) clearTimeout(resizeTimer);

  trendChart?.dispose();
  pieChart?.dispose();
  tagCloud?.dispose();
  calendarChart?.dispose();
  trendChart = null;
  pieChart = null;
  tagCloud = null;
  calendarChart = null;
});
</script>

<style scoped lang="scss">
.dashboard {
  // 顶部区域
  .top-card {
    margin-bottom: 20px;

    :deep(.el-card__body) {
      padding: 24px;
    }

    .top-content {
      display: flex;
      align-items: center;
      justify-content: space-between;
      gap: 40px;

      .profile-section {
        display: flex;
        align-items: center;
        gap: 20px;
        flex: 1;

        .avatar {
          flex-shrink: 0;
        }

        .profile-info {
          flex: 1;

          .greeting {
            font-size: 18px;
            font-weight: 600;
            color: #303133;
            margin: 0 0 8px 0;
          }

          .weather-info {
            font-size: 14px;
            color: #909399;
            margin: 0;
          }
        }
      }

      .stats-section {
        display: flex;
        align-items: center;
        gap: 40px;
        margin-right: 20px;

        .stat-item {
          text-align: center;

          .stat-label {
            font-size: 14px;
            color: #909399;
            margin-bottom: 8px;
          }

          .stat-value {
            font-size: 28px;
            font-weight: bold;
            color: #303133;
          }
        }

        .stat-divider {
          width: 1px;
          height: 40px;
          background-color: #e4e7ed;
        }
      }
    }
  }

  // 概况卡片
  .overview-cards {
    .el-col {
      margin-bottom: 20px;
    }

    .overview-card {
      transition: all 0.3s;

      :deep(.el-card__body) {
        padding: 16px;
        display: flex;
        align-items: center;
        gap: 16px;
      }

      .card-left {
        .card-icon {
          width: 48px;
          height: 48px;
          border-radius: 8px;
          display: flex;
          align-items: center;
          justify-content: center;
          font-size: 24px;
          flex-shrink: 0;

          &.icon-purple {
            background-color: #f0f0ff;
            color: #7c7cff;
          }

          &.icon-blue {
            background-color: #e8f4ff;
            color: #409eff;
          }

          &.icon-green {
            background-color: #e8f8f0;
            color: #67c23a;
          }

          &.icon-orange {
            background-color: #fff3e0;
            color: #e6a23c;
          }
        }
      }

      .card-right {
        flex: 1;
        min-width: 0;

        .card-title {
          font-size: 13px;
          color: #909399;
          margin-bottom: 4px;
        }

        .card-value {
          font-size: 24px;
          font-weight: bold;
          color: #303133;
          line-height: 1.2;
          margin-bottom: 4px;
        }

        .card-stats {
          display: flex;
          align-items: center;
          justify-content: space-between;
          font-size: 12px;

          .today-value {
            color: #909399;
          }

          .growth-rate {
            display: flex;
            align-items: center;
            gap: 2px;
            font-weight: 600;

            &.positive {
              color: #67c23a;
            }

            &.negative {
              color: #f56c6c;
            }

            &.neutral {
              color: #909399;
            }

            .el-icon {
              font-size: 12px;
            }
          }
        }
      }
    }
  }

  // 图表区域
  .charts-section {
    .el-col {
      margin-bottom: 20px;
    }

    .chart-header {
      height: 24px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: 600;
      font-size: 15px;
    }

    .chart-container {
      width: 100%;
      height: 320px;
    }

    .empty-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      height: 100%;
      color: #909399;

      p {
        margin-top: 16px;
        font-size: 14px;
      }
    }
  }

  // 快捷访问
  .quick-access-card {
    :deep(.el-card__body) {
      padding: 24px;
    }

    .quick-access-header {
      font-weight: 600;
      font-size: 15px;
    }

    .quick-access-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 40px;
      padding: 0 50px;

      .quick-links {
        flex: 1;
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 16px;
        max-width: 500px;

        .link-item {
          padding: 16px 20px;
          background: #f5f7fa;
          border-radius: 8px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          cursor: pointer;
          transition: background 0.2s;

          &:hover {
            background: #e8ecf1;
          }

          .link-text {
            font-size: 15px;
            font-weight: 500;
            color: #303133;
          }

          .link-icon {
            font-size: 16px;
            color: #909399;
          }
        }
      }

      .quick-illustration {
        flex-shrink: 0;
        width: 350px;
        height: 240px;
        display: flex;
        align-items: center;
        justify-content: center;

        img {
          height: 100%;
        }
      }
    }
  }

  // 移动端优化（<992px）
  @media (max-width: 991px) {
    // 禁用移动端的 hover 效果
    :deep(.el-card.is-hover-shadow:hover) {
      box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    }

    .top-card {
      :deep(.el-card__body) {
        padding: 16px;
      }

      .top-content {
        flex-direction: column;
        gap: 16px;

        .profile-section {
          .profile-info {
            .greeting {
              font-size: 16px;
            }

            .weather-info {
              font-size: 13px;
            }
          }
        }

        .stats-section {
          width: 100%;
          justify-content: space-around;
          margin-right: 0;
        }
      }
    }

    .overview-cards {
      .el-col {
        margin-bottom: 12px;
      }
    }

    .charts-section {
      margin-top: 8px;

      .chart-container {
        height: 280px;
      }
    }

    .quick-access-card {
      :deep(.el-card__body) {
        padding: 16px;
      }

      .quick-access-content {
        flex-direction: column;
        padding: 0 20px;
        gap: 24px;

        .quick-links {
          max-width: 100%;
          grid-template-columns: repeat(2, 1fr);
        }

        .quick-illustration {
          display: none;
        }
      }
    }
  }

  // 小屏幕优化（<768px）
  @media (max-width: 767px) {
    .top-content {
      .stats-section {
        .stat-item {
          .stat-label {
            font-size: 12px;
          }

          .stat-value {
            font-size: 20px;
          }
        }

        .stat-divider {
          height: 30px;
        }
      }
    }

    .overview-card {
      .card-left {
        .card-icon {
          width: 50px;
          height: 50px;

          .el-icon {
            font-size: 24px;
          }
        }
      }

      .card-right {
        .card-value {
          font-size: 20px;
        }
      }
    }

    .charts-section {
      margin-top: 8px;

      .chart-container {
        height: 240px;
      }
    }
  }
}
</style>
