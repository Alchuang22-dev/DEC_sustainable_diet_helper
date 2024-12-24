<template>
  <image src="/static/images/index/background_index_new.png" class="background-image"></image>
  <view class="recipe-page">
    <image src="https://cloud.tsinghua.edu.cn/thumbnail/2cea5a0d546d4319a1ef/1024/id_1-1.png" alt="Recipe Image" class="recipe-image" mode="widthFix"></image>
    <text class="recipe-name">{{ name }}</text>
    <view class="ingredients-container">
      <view
        v-for="(ingredient, index) in parsedIngredients"
        :key="index"
        :class="['ingredient-tag', `tag-color-${index % tagColors.length}`]"
      >
        {{ ingredient }}
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad as uniOnLoad } from '@dcloudio/uni-app'
import { useI18n } from 'vue-i18n';
const { t, locale } = useI18n();

// 定义响应式变量
const name = ref('');
const ingredients = ref('');
const imageUrl = ref('');
const parsedIngredients = ref([]);

// 定义可用的标签颜色
const tagColors = ['red', 'green', 'blue', 'orange', 'purple', 'cyan'];

// 解析原料信息的函数
const parseIngredients = (ingredientsStr) => {
  try {
    const ingredientsData = JSON.parse(ingredientsStr);
    if (Array.isArray(ingredientsData)) {
      return ingredientsData;
    } else if (typeof ingredientsData === 'object' && ingredientsData !== null) {
      return Object.keys(ingredientsData).map(key => t(key)) // 如果需要翻译，可以在这里集成翻译函数
    } else {
      return ['Ingredients unavailable'];
    }
  } catch (e) {
    console.error('解析原料失败:', e);
    return ['Ingredients unavailable'];
  }
}

// 页面加载时获取参数并初始化数据
uniOnLoad((options) => {
  name.value = decodeURIComponent(options.name || '');
  ingredients.value = decodeURIComponent(options.ingredients || '');
  imageUrl.value = decodeURIComponent(options.image_url || '')

  // 解析原料信息
  parsedIngredients.value = parseIngredients(ingredients.value)
})
</script>

<style scoped>
	
/* 全屏背景图片 */
.background-image {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	object-fit: cover;
	z-index: 0;
	/* 将背景图片置于最底层 */
	opacity: 0.1;
	/* 调整透明度以不干扰内容 */
}

.recipe-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
}

.recipe-image {
  width: 100%;                /* 宽度占满容器 */
  max-width: 400px;           /* 最大宽度限制 */
  border: 2px solid #ccc;     /* 在图片外增加边框 */
  border-radius: 8px;
  margin-bottom: 20px;         
  object-fit: contain;        /* 保持图片完整性，防止裁剪 */
}

.recipe-name {
  font-size: 24px; /* 2.2 将name适当放大 */
  font-weight: bold;
  margin-bottom: 15px;
  text-align: center;
}

.ingredients-container {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.ingredient-tag {
  padding: 5px 10px;
  border-radius: 15px;
  color: #fff;
  font-size: 14px;
  /* 默认背景颜色，如果需要可以覆盖 */
  background-color: #666;
}

/* 2.3 将ingredients用不同颜色的tag展示 */
.tag-color-0 {
  background-color: #f44336; /* 红色 */
}

.tag-color-1 {
  background-color: #4caf50; /* 绿色 */
}

.tag-color-2 {
  background-color: #2196f3; /* 蓝色 */
}

.tag-color-3 {
  background-color: #ff9800; /* 橙色 */
}

.tag-color-4 {
  background-color: #9c27b0; /* 紫色 */
}

.tag-color-5 {
  background-color: #00bcd4; /* 青色 */
}
</style>

