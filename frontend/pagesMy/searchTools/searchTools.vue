<template>
  <view class="search-container">
	<image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view class="header">
      <input class="search-box" v-model="searchText" @input="onSearchInput" :placeholder="placeholdertext" />
      <button class="search-button" @click="onSearch"> {{$t('text_search')}} </button>
    </view>
    <view v-if="suggestions.length > 0" class="suggestions">
      <view v-for="(item, index) in suggestions" :key="index" class="suggestion-item" @click="onSelectSuggestion(item)">
        <span v-html="highlightMatch(item)"></span>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const searchText = ref('');
const suggestions = ref([]);
const placeholderText = ref(t('text_search_label')); // 将翻译后的文本赋给变量

// 模拟与后端通信的函数
const fetchSuggestions = async (query) => {
  // 模拟调用后端API
  // 此处可以替换为真实的后端接口调用
  return new Promise((resolve) => {
    setTimeout(() => {
      const mockData = [
        '推荐的搜索结果1',
        '推荐的搜索结果2',
        '推荐的搜索结果3'
      ];
      resolve(mockData);
    }, 500);
  });
};

const onSearchInput = async () => {
  if (searchText.value.trim() !== '') {
    suggestions.value = await fetchSuggestions(searchText.value);
  } else {
    suggestions.value = [];
  }
};

const onSearch = () => {
  // 实现搜索逻辑
  console.log('搜索：', searchText.value);
};

const onSelectSuggestion = (item) => {
  console.log('用户选择了推荐结果:', item);
};

const highlightMatch = (item) => {
  const query = searchText.value.trim();
  if (!query) return item;
  const chars = query.split('').map(char => char.replace(/[.*+?^${}()|[\]\\]/g, '\$&')); // 转义特殊字符
  const regex = new RegExp(`(${chars.join('|')})`, 'gi');
  return item.replace(regex, '<span style="color: green;">$1</span>');
};
</script>

<style scoped>
body {
  font-family: 'Arial', sans-serif;
  background: url('/static/images/index/background_img.jpg') no-repeat center center fixed;
  background-size: cover;
  background-color: #f0f4f7;
  margin: 0;
  padding: 0;
}

.background-image {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  z-index: 0;
  opacity: 0.1;
}
	
.search-container {
  padding: 16px;
}
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.search-box {
  flex: 1;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.search-button {
  margin-left: 8px;
  padding: 8px;
  height: 36px;
  line-height: 20px;
  background-color: #4CAF50;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.suggestions {
  margin-top: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
  background-color: #fff;
}
.suggestion-item {
  padding: 8px;
  cursor: pointer;
}
.suggestion-item:hover {
  background-color: #f0f0f0;
}
</style>
