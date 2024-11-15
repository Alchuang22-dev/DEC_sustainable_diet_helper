//引入配置文件
import { createI18n } from "vue-i18n"; //引入vue-1i8n
import en from './en.json';   // 英文
import zhHans from './zh-Hans.json';  // 中文
 
const messages = {
	en,
	'zh-Hans': zhHans
};
console.log('Messages object:', messages);
//创建配置
const i18n = createI18n({    
	globalInjection: true,
	locale: uni.getLocale(),
	messages: {
			en,
			'zH-Hans':zhHans
		},
})
 
export default i18n