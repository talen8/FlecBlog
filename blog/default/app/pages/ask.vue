<script lang="ts" setup>
const { basicConfig } = useSysConfig();

const contactEmail = computed(() => basicConfig.value.author_email || '');

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '提问须知',
  description: '学习如何提出高质量的问题，获得更好的技术交流和帮助',
});
</script>

<template>
  <div id="page">
    <div class="page-title">提问须知</div>

    <div class="content">
      <p>本页面旨在帮助你更有效地提问和交流，获得更有价值的回答。</p>

      <p>良好的提问方式不仅能让你更快地解决问题，也能帮助建立更好的技术交流氛围。</p>

      <blockquote>
        <strong>核心理念：</strong>尊重他人的时间，展示你的思考，提供清晰的信息。
      </blockquote>

      <h1>为什么需要这个指南</h1>

      <p>在技术社区中，提问是一门艺术。一个好的问题能：</p>

      <ul>
        <li>让回答者快速理解你的困境</li>
        <li>节省双方的时间和精力</li>
        <li>帮助其他遇到类似问题的人</li>
        <li>促进知识的传播和积累</li>
      </ul>

      <p>相反，一个糟糕的问题可能：</p>

      <ul>
        <li>得不到任何回答</li>
        <li>浪费大家的时间</li>
        <li>降低整体讨论质量</li>
        <li>让你错失学习的机会</li>
      </ul>

      <h1>提问前的准备</h1>

      <p>在提问之前，请先做到以下几点：</p>

      <h2>1. 先尝试自己解决</h2>
      <p>提问前，请先尝试：</p>
      <ul>
        <li>使用搜索引擎（Google、百度、必应）搜索相关问题</li>
        <li>查看相关文章和已有的讨论（如评论区、论坛等）</li>
        <li>阅读官方文档、教程或相关技术资料</li>
        <li>自己动手试验、调试和验证</li>
        <li>检查错误日志和提示信息</li>
      </ul>

      <blockquote>
        表明你已经做过这些努力，会让回答者更愿意帮助你。比如："我已经搜索过XX，查看过XX文档，尝试了XX方法，但还是遇到XX问题..."
      </blockquote>

      <h2>2. 阅读相关内容</h2>
      <ul>
        <li>完整阅读相关文章，不要跳读或只看标题</li>
        <li>查看参考链接和扩展资料</li>
        <li>浏览已有的讨论和回答</li>
      </ul>

      <h2>3. 明确你的问题</h2>
      <ul>
        <li>我到底想知道什么？</li>
        <li>我遇到的具体问题是什么？</li>
        <li>我期望得到什么样的答案？</li>
      </ul>

      <h1>如何提出好问题</h1>

      <h2>使用清晰准确的语言</h2>

      <p><strong>好的提问：</strong></p>
      <ul>
        <li>语句通顺，标点正确</li>
        <li>用词准确，避免模糊表述</li>
        <li>适当分段，方便阅读</li>
      </ul>

      <p><strong>避免：</strong></p>
      <ul>
        <li>错字连篇、语句不通</li>
        <li>全部大写或全部小写</li>
        <li>滥用感叹号、问号</li>
        <li>网络用语、火星文</li>
      </ul>

      <h2>提供必要的背景信息</h2>

      <p>描述问题时，应该包含：</p>

      <ul>
        <li><strong>问题环境</strong>：操作系统、浏览器、软件版本等</li>
        <li><strong>问题现象</strong>：具体发生了什么，有什么错误提示</li>
        <li><strong>预期行为</strong>：你期望发生什么</li>
        <li><strong>已做尝试</strong>：你已经尝试过哪些方法</li>
      </ul>

      <h2>描述问题而非猜测</h2>

      <p><strong>不好的提问：</strong></p>
      <p>"我觉得是XX模块有问题，应该怎么修？"</p>

      <p><strong>好的提问：</strong></p>
      <p>"当我执行XX操作时，系统提示XX错误。我检查了XX部分，发现XX情况。这是什么原因导致的？"</p>

      <blockquote>描述你观察到的事实，而不是你的猜测。让有经验的人来帮你诊断。</blockquote>

      <h2>话不在多而在精</h2>

      <ul>
        <li>简洁地描述问题，不要复制粘贴大量代码或日志</li>
        <li>如果代码较长，只贴出关键部分</li>
        <li>使用代码块格式（如果评论支持）使代码易读</li>
        <li>长篇内容建议使用外部链接（如 GitHub Gist）</li>
      </ul>

      <h1>选择合适的提问方式</h1>

      <ul>
        <li>技术问题和内容讨论，建议在相关文章评论区公开交流</li>
        <li>私密问题或深度交流，可通过邮件联系（{{ contactEmail }}）</li>
        <li>公开提问能让更多人受益，也更容易得到回答</li>
      </ul>

      <h1>交流规范</h1>

      <p><strong>欢迎：</strong></p>
      <ul>
        <li>与主题相关的技术讨论</li>
        <li>有建设性的意见和建议</li>
        <li>分享经验和观点</li>
        <li>礼貌地指出错误</li>
      </ul>

      <p><strong>不欢迎：</strong></p>
      <ul>
        <li>无关内容、广告、垃圾信息</li>
        <li>人身攻击、辱骂性言论</li>
        <li>未经思考的灌水</li>
      </ul>

      <p><strong>基本礼节：</strong></p>
      <ul>
        <li>保持礼貌，尊重他人</li>
        <li>文明讨论，对事不对人</li>
        <li>适度表达，避免过度情绪化</li>
        <li>问题解决后，分享解决方案</li>
      </ul>

      <h1>好问题与坏问题</h1>

      <h2>案例一：技术问题</h2>

      <p><strong>坏问题：</strong></p>
      <p class="bad-example">"急！代码不能运行！怎么办？"</p>
      <p><strong>问题：</strong>没有任何有用信息，无法判断问题所在。</p>

      <p><strong>好问题：</strong></p>
      <p class="good-example">
        "我在使用文章中提到的XX方法时，在第X行出现了'XX错误'。我的环境是Windows 10 + Python
        3.9。我尝试了XX方法但没有解决。请问可能是什么原因？"
      </p>
      <p><strong>优点：</strong>提供了环境、错误信息、已做尝试等关键信息。</p>

      <h2>案例二：概念理解</h2>

      <p><strong>坏问题：</strong></p>
      <p class="bad-example">"XX是什么？"</p>
      <p><strong>问题：</strong>这种问题搜索引擎能更好地回答。</p>

      <p><strong>好问题：</strong></p>
      <p class="good-example">
        "文章中提到XX概念，我查了资料理解为XX意思。但在YY场景下，我不太明白为什么要用XX而不是ZZ，能否解释一下？"
      </p>
      <p><strong>优点：</strong>表明已经做过研究，问的是更深层次的理解。</p>

      <h2>案例三：文章反馈</h2>

      <p><strong>坏问题：</strong></p>
      <p class="bad-example">"写得什么玩意，看不懂。"</p>
      <p><strong>问题：</strong>既不礼貌也没有建设性。</p>

      <p><strong>好问题：</strong></p>
      <p class="good-example">
        "文章的XX部分我理解起来有些困难，能否补充说明一下XX和YY之间的关系？或者推荐一些相关资料？"
      </p>
      <p><strong>优点：</strong>明确指出不理解的地方，并给出建设性建议。</p>

      <h1>如何对待回答</h1>

      <h2>理解回答者</h2>

      <ul>
        <li>回答问题是义务帮助，没有义务回答所有问题</li>
        <li>没有得到回答不代表被忽视，可能暂时无人知道答案</li>
        <li>如果回答看不懂，先尝试自己理解，再礼貌请求说明</li>
        <li>被建议"去看文档"通常是好意，说明答案已存在</li>
        <li>技术讨论往往直接简练，不要误认为不礼貌</li>
      </ul>

      <h2>问题解决后</h2>

      <ul>
        <li>回复说明问题如何解决的</li>
        <li>简单感谢帮助过你的人</li>
        <li>让其他人也能从中受益</li>
      </ul>

      <h1>不适合提问的内容</h1>

      <ul>
        <li><strong>家庭作业</strong>：可以寻求思路提示，但不要要求完整答案</li>
        <li><strong>过于宽泛</strong>："怎么学编程？"这类问题太大，无从回答</li>
        <li><strong>重复问题</strong>：提问前先检查是否已有人问过</li>
      </ul>

      <h1>回复说明</h1>

      <p>为什么有些问题可能得不到回答：</p>

      <ul>
        <li>优先回答质量高、有代表性的问题</li>
        <li>优先回答能帮助更多人的问题</li>
        <li>已经回答过的问题可能不会重复回答</li>
        <li>明显未做准备的提问可能会被建议先自行研究</li>
      </ul>

      <blockquote>这是为了维护讨论质量，鼓励独立思考和主动学习。</blockquote>

      <h1>其他说明</h1>

      <p><strong>发现错误：</strong>欢迎指出文章中的错误，最好能说明正确内容，作者将感激不尽。</p>

      <p>
        <strong>私下交流：</strong>如需私下交流，请通过邮件联系（{{
          contactEmail
        }}）。请在邮件中说明来意，理解可能需要几天才能回复。
      </p>

      <h1>推荐阅读</h1>

      <p>本页面的内容参考和改编自以下资源：</p>

      <ul>
        <li>
          <a
            href="http://www.catb.org/~esr/faqs/smart-questions.html"
            target="_blank"
            rel="noopener noreferrer"
            >《How To Ask Questions The Smart Way》</a
          >
          by Eric S. Raymond & Rick Moen（提问的智慧）
        </li>
        <li>
          <a
            href="https://github.com/ryanhanwu/How-To-Ask-Questions-The-Smart-Way"
            target="_blank"
            rel="noopener noreferrer"
            >《提问的智慧》中文版</a
          >
          - GitHub
        </li>
      </ul>

      <p>强烈建议完整阅读上述文档，它们对在任何技术社区提问都有帮助。</p>

      <h1>总结</h1>

      <p>记住这些原则：</p>

      <ul>
        <li><strong>独立思考</strong>：先尝试自己解决</li>
        <li><strong>清晰表达</strong>：提供足够的信息</li>
        <li><strong>互相尊重</strong>：保持礼貌和谦逊</li>
        <li><strong>质量优先</strong>：一个好问题胜过十个随意提问</li>
        <li><strong>分享知识</strong>：有能力时也帮助他人</li>
      </ul>

      <blockquote>好的提问是双赢的开始 —— 你得到帮助，社区积累知识。</blockquote>

      <p>祝你交流愉快！</p>

      <p class="update-time">最后更新时间：2025年10月</p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;

  .page-title {
    margin: 0 0 10px;
    font-weight: bold;
    font-size: 2rem;
  }

  .content {
    color: var(--font-color);
    line-height: 1.8;
    font-size: 1rem;

    h1 {
      font-size: 1.5rem;
      font-weight: 600;
      margin: 2.5rem 0 1rem;
      color: var(--font-color);
      line-height: 1.3;

      &:first-child {
        margin-top: 1.5rem;
      }
    }

    h2 {
      font-size: 1.25rem;
      font-weight: 600;
      margin: 1.75rem 0 0.875rem;
      color: var(--font-color);
      line-height: 1.3;
    }

    p {
      margin: 1rem 0;
      text-align: justify;
    }

    a {
      color: var(--theme-color);
      text-decoration: none;
      font-weight: 500;
      border-bottom: 1px solid transparent;
      transition: all 0.3s ease;

      &:hover {
        border-bottom-color: var(--theme-color);
      }
    }

    strong {
      font-weight: 600;
      color: var(--font-color);
    }

    ul {
      list-style: disc;
      padding-left: 2rem;
      margin: 1rem 0;

      li {
        margin: 0.5rem 0;
      }
    }

    blockquote {
      margin: 1.5rem 0;
      padding: 0.75em 1.25em;
      background: var(--flec-heavy-bg);
      border-left: 4px solid var(--theme-color);
      border-radius: 0 6px 6px 0;
      font-size: 0.95em;
      color: var(--theme-meta-color);
      line-height: 1.7;
    }

    .bad-example {
      padding: 0.75rem 1rem;
      background: rgba(255, 82, 82, 0.1);
      border-left: 3px solid #ff5252;
      margin: 0.5rem 0;
      font-style: italic;
    }

    .good-example {
      padding: 0.75rem 1rem;
      background: rgba(76, 175, 80, 0.1);
      border-left: 3px solid #4caf50;
      margin: 0.5rem 0;
      font-style: italic;
    }

    .update-time {
      margin-top: 3rem;
      color: var(--theme-meta-color);
      font-size: 0.875rem;
      text-align: right;
      font-weight: 500;
    }
  }
}

// 响应式设�?
@media screen and (max-width: 1024px) {
  #page {
    padding: 30px;

    .page-title {
      font-size: 1.75rem;
    }

    .content {
      h1 {
        font-size: 1.35rem;
        margin: 2rem 0 0.85rem;
      }

      h2 {
        font-size: 1.15rem;
        margin: 1.5rem 0 0.75rem;
      }

      p {
        font-size: 0.95rem;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  #page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
    }

    .content {
      font-size: 0.9rem;

      h1 {
        font-size: 1.25rem;
        margin: 1.75rem 0 0.75rem;
      }

      h2 {
        font-size: 1.1rem;
        margin: 1.35rem 0 0.65rem;
      }

      p {
        font-size: 0.9rem;
      }

      ul {
        padding-left: 1.5rem;

        li {
          margin: 0.4rem 0;
        }
      }

      blockquote {
        padding: 0.65em 1em;
        margin: 1.25rem 0;
        font-size: 0.88rem;
      }

      .bad-example,
      .good-example {
        padding: 0.65rem 0.85rem;
        font-size: 0.88rem;
      }

      .update-time {
        font-size: 0.8rem;
      }
    }
  }
}
</style>
