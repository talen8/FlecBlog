<script lang="ts" setup>
import { getArticlesForWeb } from '@/composables/api/article';
import { getCategories } from '@/composables/api/category';
import { getTags } from '@/composables/api/tag';
import { parseJSON } from '@/utils/json';

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '关于',
  description: '了解博主的个人信息、经历和故事',
});

const personalityMap: Record<string, { name: string; color: string; image: string; url: string }> =
  {
    INTJ: { name: '建筑师', color: '#885fb8', image: 'intj-architect', url: 'intj' },
    INTP: { name: '逻辑学家', color: '#885fb8', image: 'intp-logician', url: 'intp' },
    ENTJ: { name: '指挥官', color: '#885fb8', image: 'entj-commander', url: 'entj' },
    ENTP: { name: '辩论家', color: '#885fb8', image: 'entp-debater', url: 'entp' },
    INFJ: { name: '提倡者', color: '#56a178', image: 'infj-advocate', url: 'infj' },
    INFP: { name: '调停者', color: '#56a178', image: 'infp-mediator', url: 'infp' },
    ENFJ: { name: '主人公', color: '#56a178', image: 'enfj-protagonist', url: 'enfj' },
    ENFP: { name: '竞选者', color: '#56a178', image: 'enfp-campaigner', url: 'enfp' },
    ISTJ: { name: '物流师', color: '#4298b4', image: 'istj-logistician', url: 'istj' },
    ISFJ: { name: '守卫者', color: '#4298b4', image: 'isfj-defender', url: 'isfj' },
    ESTJ: { name: '总经理', color: '#4298b4', image: 'estj-executive', url: 'estj' },
    ESFJ: { name: '执政官', color: '#4298b4', image: 'esfj-consul', url: 'esfj' },
    ISTP: { name: '鉴赏家', color: '#e4ae3a', image: 'istp-virtuoso', url: 'istp' },
    ISFP: { name: '探险家', color: '#e4ae3a', image: 'isfp-adventurer', url: 'isfp' },
    ESTP: { name: '企业家', color: '#e4ae3a', image: 'estp-entrepreneur', url: 'estp' },
    ESFP: { name: '表演者', color: '#e4ae3a', image: 'esfp-entertainer', url: 'esfp' },
  };

const getPersonalityInfo = (code: string) => {
  const baseType = code.substring(0, 4).toUpperCase();
  const info = personalityMap[baseType]!;

  return {
    name: info.name,
    color: info.color,
    image: `https://www.16personalities.com/static/images/personality-types/avatars/${info.image}.png`,
    url: `https://www.16personalities.com/ch/${info.url}-%E4%BA%BA%E6%A0%BC`,
  };
};

const { siteStats } = useStats();
const { total: articleTotal } = useArticles();
const { total: categoryTotal } = useCategories();
const { total: tagTotal } = useTags();
const { blogConfig, basicConfig } = useSysConfig();

const { data: articlesData } = await useAsyncData('about-articles', async () => {
  const { total: resTotal } = await getArticlesForWeb({ page: 1, page_size: 1 });
  return { total: resTotal };
});

const { data: categoriesData } = await useAsyncData('about-categories', async () => {
  const { total: resTotal } = await getCategories();
  return { total: resTotal };
});

const { data: tagsData } = await useAsyncData('about-tags', async () => {
  const { total: resTotal } = await getTags();
  return { total: resTotal };
});

if (articlesData.value) articleTotal.value = articlesData.value.total || 0;
if (categoriesData.value) categoryTotal.value = categoriesData.value.total || 0;
if (tagsData.value) tagTotal.value = tagsData.value.total || 0;

const info = computed(() => {
  const blog = blogConfig.value;
  const personalityCode = blog.about_personality || 'INFJ-A';
  const personality = getPersonalityInfo(personalityCode);

  return {
    author: basicConfig.value.author || '',
    describe: blog.about_describe || '',
    describeTips: blog.about_describe_tips || '',
    photo: basicConfig.value.author_photo || '',
    exhibitionImg: blog.about_exhibition || '',
    profile: parseJSON<Array<{ label: string; value: string; color: string }>>(
      blog.about_profile,
      []
    ),
    personality: {
      type: personalityCode,
      name: personality.name,
      color: personality.color,
      image: personality.image,
      url: personality.url,
    },
    motto: {
      main: parseJSON<string[]>(blog.about_motto_main, []),
      sub: blog.about_motto_sub || '',
    },
    socialize: parseJSON<Array<{ name: string; url: string }>>(blog.about_socialize, []),
    creation: parseJSON<Array<{ name: string; url: string }>>(blog.about_creation, []),
    versions: parseJSON<Array<{ name: string; version: string }>>(blog.about_versions, []),
    union: parseJSON<Array<{ name: string; url: string }>>(blog.about_unions, []),
    story: blog.about_story || '',
  };
});

const runningDays = computed(() => {
  const established = blogConfig.value.established || '2024-01-01';
  const startDate = new Date(established).getTime();
  const now = Date.now();
  return Math.floor((now - startDate) / 86400000);
});
const runTime = computed(() => `已稳定运行 ${runningDays.value} 天 🚀`);

const formatWords = (words: string) => {
  const n = +words;
  return n >= 1e4 ? (n / 1e4).toFixed(1) + 'w' : n >= 1e3 ? (n / 1e3).toFixed(1) + 'k' : words;
};
</script>

<template>
  <div id="about-page">
    <!-- 个人介绍 -->
    <div class="Personal-Introduction">
      <div class="PI-box-left">
        <h1 class="title">你好！</h1>
        <div v-if="info.author" class="title">我是 {{ info.author }}</div>
        <div v-if="info.describe" class="describe">{{ info.describe }}</div>
        <span v-if="info.describeTips" class="describe-tips">{{ info.describeTips }}</span>
        <div class="PI-button">
          <a href="#one">博主信息</a>
          <a href="#two">本站信息</a>
        </div>
      </div>
      <div class="PI-box-right">
        <NuxtImg :src="info.photo" alt="个人照片" loading="lazy" />
      </div>
    </div>

    <!-- 博主信息 -->
    <div id="one">
      <div class="h1-box">
        <div class="box-top">
          <span>01</span>
          <div class="title-h1">博主信息</div>
        </div>
        <div class="about-layout box-bottom">{{ info.author }}</div>
      </div>
      <div class="information">
        <div class="about-layout Introduction">
          <div v-for="n in Math.ceil(info.profile.length / 3)" :key="n" class="bar-box-row">
            <div
              v-for="item in info.profile.slice((n - 1) * 3, n * 3)"
              :key="item.label"
              class="bar-box"
            >
              <span class="tips">{{ item.label }}</span>
              <div class="title" :style="{ color: item.color }">
                {{ item.value }}
              </div>
            </div>
          </div>
        </div>
        <div v-if="info.exhibitionImg" class="about-layout Exhibition">
          <NuxtImg :src="info.exhibitionImg" alt="展示图片" loading="lazy" />
        </div>
      </div>
    </div>

    <!-- 性格与座右铭 -->
    <div class="Philosophical">
      <div class="about-layout P-box-left">
        <div class="tips">性格</div>
        <div class="title">{{ info.personality.name }}</div>
        <div class="title" :style="{ color: info.personality.color }">
          {{ info.personality.type }}
        </div>
        <NuxtImg class="image" :src="info.personality.image" alt="性格类型" loading="lazy" />
        <div class="tips-bottom">
          在
          <a href="https://www.16personalities.com/ch" target="_blank">16Personalities</a>
          了解关于
          <a :href="info.personality.url" target="_blank">{{ info.personality.name }}</a
          >&ensp;的更多信息
        </div>
      </div>
      <div class="about-layout P-box-right">
        <div class="tips">座右铭</div>
        <span
          v-for="(text, index) in info.motto.main"
          :key="index"
          class="title"
          :style="{
            opacity: index === info.motto.main.length - 1 ? 1 : 0.6,
            marginBottom: index < info.motto.main.length - 1 ? '8px' : '0',
          }"
        >
          {{ text }}
        </span>
        <div v-if="info.motto.sub" class="tips-bottom">{{ info.motto.sub }}</div>
      </div>
    </div>

    <!-- 联系方式与创作平台 -->
    <div class="Platform">
      <div class="about-layout Socialize">
        <div class="tips">账号</div>
        <div class="title">联系方式</div>
        <div class="S-box">
          <a
            v-for="item in info.socialize"
            :key="item.name"
            class="btn-layout"
            :href="item.url"
            target="_blank"
            >{{ item.name }}</a
          >
        </div>
      </div>
      <div class="about-layout Creation">
        <div class="tips">订阅</div>
        <div class="title">创作平台</div>
        <div class="S-box">
          <a
            v-for="item in info.creation"
            :key="item.name"
            class="btn-layout"
            :href="item.url"
            target="_blank"
            >{{ item.name }}</a
          >
        </div>
      </div>
    </div>

    <!-- 本站信息 -->
    <div id="two">
      <div class="h1-box">
        <div class="box-top">
          <span>02</span>
          <div class="title-h1">本站信息</div>
        </div>
        <div class="about-layout box-bottom">{{ runTime }}</div>
      </div>
      <div class="information">
        <div class="about-layout Version">
          <div v-for="v in info.versions" :key="v.name" class="V-box">
            <div class="title">{{ v.name }}</div>
            <div class="tips-v">V{{ v.version }}</div>
          </div>
        </div>
        <div class="about-layout Statistics">
          <span>{{ articleTotal }}篇文章</span>
          <span>{{ categoryTotal }}个分类</span>
          <span>{{ tagTotal }}个标签</span>
          <span v-if="siteStats.total_words">{{ formatWords(siteStats.total_words) }}字</span>
        </div>
      </div>
    </div>

    <!-- 访问统计与站长联盟 -->
    <div class="data">
      <div class="about-layout statistic">
        <div class="tips">浏览</div>
        <div class="title">访问统计</div>
        <div id="statistic">
          <div>
            <span class="tips">今日访客</span><span>{{ siteStats.today_visitors || 0 }}</span>
          </div>
          <div>
            <span class="tips">今日访问</span><span>{{ siteStats.today_pageviews || 0 }}</span>
          </div>
          <div>
            <span class="tips">昨日访客</span><span>{{ siteStats.yesterday_visitors || 0 }}</span>
          </div>
          <div>
            <span class="tips">昨日访问</span><span>{{ siteStats.yesterday_pageviews || 0 }}</span>
          </div>
          <div>
            <span class="tips">本月访问</span><span>{{ siteStats.month_pageviews || 0 }}</span>
          </div>
        </div>
        <a class="T-btn" href="/statistics">更多统计</a>
      </div>
      <div class="about-layout union">
        <div class="tips">共创</div>
        <div class="title">站长联盟</div>
        <div class="U-box">
          <a
            v-for="item in info.union"
            :key="item.name"
            class="btn-layout"
            :href="item.url"
            target="_blank"
            >{{ item.name }}</a
          >
        </div>
      </div>
    </div>

    <!-- 心路历程 -->
    <div v-if="info.story" class="about-layout content">
      <div class="tips">心路历程</div>
      <div class="title">关于本站的介绍</div>
      <p>{{ info.story }}</p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#about-page {
  @extend .cardHover;
  padding: 40px;

  .about-layout {
    @extend .cardHover;
    border-radius: 12px;
    position: relative;
    padding: 1rem 2rem;
    overflow: hidden;

    &:hover {
      border-color: var(--theme-color);
      transform: translateY(-2px);
    }
  }

  .title {
    font-size: 2.25rem;
    font-weight: 700;
    line-height: 1.2;
  }

  .tips {
    opacity: 0.8;
    font-size: 0.75rem;
    line-height: 1.2;
    margin-bottom: 0.75rem;
  }

  .tips-bottom {
    font-size: 0.875rem;
    position: absolute;
    bottom: 1rem;
    left: 2rem;

    a {
      font-weight: 600;
      text-decoration: none;
      color: var(--font-color);

      &:hover {
        color: var(--theme-color);
      }
    }
  }

  .btn-layout {
    @extend .cardHover;
    padding: 6px 18px;
    margin: 0 18px 18px 0;
    color: var(--font-color);
    text-decoration: none;
    display: inline-block;

    &:hover {
      background: var(--theme-color);
      color: #fff;
    }
  }

  .h1-box {
    display: flex;
    flex-direction: column;
    justify-content: flex-end;

    .box-top {
      margin: auto;
      position: relative;
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;

      span {
        font-size: 100px;
        position: absolute;
        top: 0;
        left: -64px;
        opacity: 0.4;
      }

      .title-h1 {
        font-size: 42px;
        font-weight: 700;
        line-height: 1.1;
        margin: 0.5rem 0;
        letter-spacing: 0.2rem;
        color: var(--font-color);
      }
    }

    .box-bottom {
      padding: 1.25rem 2rem;
      display: inline-flex;
      justify-content: center;
      font-size: 18px;
    }
  }

  // 个人介绍
  .Personal-Introduction {
    display: flex;
    justify-content: space-between;
    padding: 2rem 0;

    .PI-box-left {
      margin-top: 1.5rem;
      color: var(--font-color);
      width: 60%;
      z-index: 1;

      .title {
        font-size: 42px;
        margin: 0.5rem 0;
        letter-spacing: 0.2rem;
      }

      .describe {
        font-size: 18px;
        letter-spacing: 0.2rem;
        margin-top: 2.25rem;
        opacity: 0.9;
      }

      .describe-tips {
        font-size: 16px;
        opacity: 0.4;
      }

      .PI-button {
        position: relative;
        top: 50px;
        display: flex;

        a {
          @extend .cardHover;
          padding: 6px 18px;
          margin-right: 16px;
          text-decoration: none;
          color: var(--font-color);

          &:hover {
            background: var(--theme-color);
            color: #fff;
          }
        }
      }
    }

    .PI-box-right {
      height: 550px;
      width: 40%;
      overflow: hidden;
      display: flex;
      justify-content: center;

      img {
        height: 100%;
        object-fit: cover;
        border-radius: 12px;
      }
    }
  }

  // 博主信息
  #one {
    margin-top: 32px;
    display: flex;
    flex-direction: row-reverse;
    scroll-margin-top: 100px;

    .h1-box {
      width: 50%;
      margin-left: 16px;
      aspect-ratio: 1 / 1;
    }

    .information {
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .Introduction {
        display: flex;
        flex-direction: column;
        gap: 16px;
        padding: 1rem;
        flex: 1;

        .bar-box-row {
          display: flex;
          justify-content: space-between;
          gap: 16px;
        }

        .bar-box {
          flex: 1;
          text-align: center;
          padding: 1rem;
        }
      }

      .Exhibition {
        padding: 0;
        height: 76px;
        margin-top: 16px;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.3s;
        }

        &:hover img {
          transform: scale(1.05);
        }
      }
    }
  }

  // 性格与座右铭
  .Philosophical {
    margin-top: 16px;
    min-height: 240px;
    display: flex;
    gap: 16px;

    .P-box-left {
      width: 60%;
      padding: 2.25rem 2rem;

      &:hover .image {
        transform: rotate(-10deg);
      }

      .image {
        position: absolute;
        right: 10px;
        top: 10px;
        width: 200px;
        transition: transform 2s cubic-bezier(0.13, 0.45, 0.21, 1.02);
      }
    }

    .P-box-right {
      width: 40%;
      padding: 2.25rem 2rem;
      display: flex;
      flex-direction: column;
    }
  }

  // 联系方式与创作平台
  .Platform {
    margin-top: 16px;
    display: flex;
    gap: 16px;

    .Socialize,
    .Creation {
      padding: 2.25rem 2rem;

      .S-box {
        margin-top: 2.25rem;
        display: flex;
        flex-wrap: wrap;
      }
    }

    .Socialize {
      width: 40%;
    }

    .Creation {
      width: 60%;
    }
  }

  // 本站信息
  #two {
    margin-top: 32px;
    display: flex;
    scroll-margin-top: 100px;

    .h1-box {
      width: 50%;
      margin-right: 16px;
      aspect-ratio: 1 / 1;
    }

    .information {
      width: 100%;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .Version {
        display: flex;
        flex: 1;
        margin-bottom: 16px;
        align-items: center;
        justify-content: space-around;

        .V-box {
          display: flex;
          flex-direction: column;
          align-items: center;
        }

        .title {
          color: var(--font-color);
          z-index: 1;
        }

        .tips-v {
          font-size: 0.875rem;
          color: var(--theme-meta-color);
          margin-top: 0.5rem;
        }
      }

      .Statistics {
        padding: 1.25rem 2rem;
        display: inline-flex;
        justify-content: space-between;

        span {
          font-size: 18px;
        }
      }
    }
  }

  // 访问统计与站长联盟
  .data {
    margin-top: 16px;
    display: flex;
    justify-content: space-between;

    .statistic {
      width: calc(65% - 8px);
      padding: 2.25rem 2rem;
      display: flex;
      flex-direction: column;
      background: linear-gradient(135deg, #0c1c2c 0%, #1a3a52 100%);
      color: #fff;
      position: relative;

      &::before {
        content: '';
        position: absolute;
        inset: 0;
        background:
          radial-gradient(circle at 30% 50%, rgba(73, 177, 245, 0.15) 0%, transparent 50%),
          radial-gradient(circle at 70% 80%, rgba(120, 194, 244, 0.1) 0%, transparent 50%);
        pointer-events: none;
      }

      & > * {
        z-index: 1;
      }

      #statistic {
        display: flex;
        justify-content: space-between;
        margin: auto 0;

        div {
          margin: 0 16px 16px 0;

          span:last-child {
            font-size: 36px;
            font-weight: 700;
            color: #fff;
            display: block;
            margin-bottom: 0.5rem;
          }
        }
      }

      .T-btn {
        @extend .cardHover;
        position: absolute;
        bottom: 1rem;
        right: 2rem;
        height: 40px;
        width: 160px;
        border-radius: 20px;
        display: flex;
        align-items: center;
        justify-content: center;
        text-decoration: none;
        color: var(--font-color);

        &:hover {
          background: var(--theme-color);
          color: #fff;
        }
      }
    }

    .union {
      width: calc(35% - 8px);
      padding: 2.25rem 2rem;

      .U-box {
        margin-top: 1.25rem;
        display: flex;
        flex-wrap: wrap;
      }
    }
  }

  // 心路历程
  .content {
    margin-top: 16px;
    padding: 2.25rem 2rem;

    p {
      white-space: pre-line;
      line-height: 1.8;
      margin-top: 1rem;
      color: var(--font-color);
      opacity: 0.9;
    }
  }
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #about-page {
    padding: 30px;

    .about-layout {
      padding: 1rem 1.5rem;
    }

    .title {
      font-size: 2rem;
    }

    .h1-box {
      .box-top {
        span {
          font-size: 80px;
          left: -50px;
        }

        .title-h1 {
          font-size: 36px;
        }
      }
    }

    .Personal-Introduction {
      .PI-box-left {
        .title {
          font-size: 36px;
        }

        .describe {
          font-size: 16px;
        }
      }
    }

    #one,
    #two {
      .h1-box {
        margin: 0 12px;
      }
    }

    .Philosophical {
      .P-box-left {
        .image {
          width: 180px;
        }
      }
    }
  }
}

@media screen and (max-width: 768px) {
  #about-page {
    padding: 18px;

    .h1-box .box-top span {
      display: none;
    }

    .Personal-Introduction {
      .PI-box-left {
        width: 100%;

        .describe {
          font-size: 16px;
          letter-spacing: 0;
        }

        .describe-tips {
          font-size: 12px;
        }

        .PI-button {
          position: static;
          margin-top: 10px;
        }
      }

      .PI-box-right {
        display: none;
      }
    }

    #one,
    #two {
      flex-direction: column;

      .h1-box {
        width: 100%;
        height: 220px;
        margin: 0 0 16px;
      }

      .information {
        .Introduction .bar-box-row {
          flex-direction: column;
          gap: 12px;
        }

        .Exhibition {
          margin: 16px 0;
        }

        .Version {
          flex-direction: column;

          .V-box {
            margin: 16px 0;
          }
        }

        .Statistics {
          flex-wrap: wrap;

          span {
            width: 50%;
            text-align: center;
            margin: 0;
          }
        }
      }
    }

    .Philosophical {
      flex-direction: column;
      gap: 0;
      margin-top: 0;

      .P-box-left,
      .P-box-right {
        width: 100%;
        height: 210px;
        margin-bottom: 16px;
      }

      .P-box-left .image {
        width: 120px;
        right: 18px;
        top: 40px;
      }
    }

    .Platform {
      flex-direction: column;
      gap: 0;

      .Socialize,
      .Creation {
        width: 100%;
        margin-bottom: 16px;
      }
    }

    .data {
      flex-direction: column;

      .statistic {
        width: 100%;
        margin-bottom: 16px;

        #statistic {
          flex-wrap: wrap;
          margin: 1.25rem 0;

          div {
            margin: 0 0 16px;
            width: 50%;
          }
        }

        .T-btn {
          display: none;
        }
      }

      .union {
        width: 100%;
      }
    }
  }
}
</style>
