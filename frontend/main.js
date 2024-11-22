import { createSSRApp } from 'vue';
import App from './App.vue';
import { createI18n } from 'vue-i18n'; // Vue 3 中的 i18n 引入方式
import en from './locales/en.json'; // 英文
import zhHans from './locales/zh-Hans.json'; // 中文
import { createPinia } from 'pinia'; // 引入 Pinia

console.log("in Vue3");

export function createApp() {
    const app = createSSRApp(App);

    // 创建 i18n 实例
    const i18n = createI18n({
        legacy: true, // 不使用 Composition API 模式
        locale: 'zh-Hans', // 默认语言
        messages: {
            en,
            'zh-Hans': zhHans,
        },
    });

    // 注册 i18n 和 Pinia 到应用
    app.use(i18n);
    app.use(createPinia());

    return { app };
}
