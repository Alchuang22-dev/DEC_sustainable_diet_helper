// #ifndef VUE3
import Vue from 'vue';
import App from './App';
import VueI18n from 'vue-i18n';
import i18n from './locales/index.ts'; // 导入 i18n 实例

console.log("in vue2");

Vue.use(VueI18n);

Vue.config.productionTip = false;

App.mpType = 'app';

// 正确创建 Vue 实例，避免使用展开运算符
const app = new Vue({
  i18n, // 注入 i18n 实例
  render: h => h(App)
});

app.$mount();
// #endif

// #ifdef VUE3
import { createSSRApp } from 'vue';
import App from './App.vue';
import { createI18n } from 'vue-i18n'; // Vue 3 中的 i18n 引入方式
import messages from './locales/index.ts'; // 导入语言包
import en from './locales/en.json';   // 英文
import zhHans from './locales/zh-Hans.json';  // 中文

console.log("in Vue3");
console.log('Messages object:', messages);

export function createApp() {
  const app = createSSRApp(App);

  // 创建 i18n 实例
  const i18n = createI18n({
    legacy: true, // 不使用 Composition API 模式
    locale: 'zh-Hans', // 默认语言
	messages: {
		en,
		'zh-Hans': zhHans
	},
  });

  // 将 i18n 实例注册到应用
  app.use(i18n);
  
  //console.log('app 实例：', app);

  return {
    app
  };
}
// #endif
