import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import App from './App.vue'
import { loadPersistedData } from './composables/useStore'
import './assets/css/main.scss'
import 'remixicon/fonts/remixicon.css'

async function bootstrap() {
  await loadPersistedData()
  const app = createApp(App)
  app.use(ElementPlus, { locale: zhCn })
  app.mount('#app')
}

bootstrap()
