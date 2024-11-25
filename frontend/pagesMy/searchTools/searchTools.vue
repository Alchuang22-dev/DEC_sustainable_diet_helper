<template>
  <view class="search-container">
    <view class="header">
      <input class="search-box" v-model="searchText" @input="onSearchInput" :placeholder="placeholderText" />
      <button class="search-button" @click="onSearch"> {{$t('text_search')}} </button>
    </view>
    <view v-if="suggestions.length > 0" class="suggestions">
      <view v-for="(item, index) in suggestions" :key="index" class="suggestion-item" @click="onSelectSuggestion(item)">
        <span v-html="highlightMatch(item)"></span>
      </view>
    </view>
    <view v-if="historySearches.length > 0" class="history-tags">
      <text class="section-title">历史搜索</text>
      <view class="tags">
        <view v-for="(item, index) in historySearches" :key="index" class="tag" @click="onSelectHistory(item)">
          {{ item }}
        </view>
      </view>
    </view>
    <view v-if="popularSearches.length > 0" class="popular-tags">
      <text class="section-title">热门搜索</text>
      <view class="tags">
        <view v-for="(item, index) in popularSearches" :key="index" class="tag" @click="onSelectPopular(item)">
          {{ item }}
        </view>
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
const placeholderText = ref(t('text_search_label'));

// 历史搜索数据
const historySearches = ref(['历史搜索1', '历史搜索2', '历史搜索3']);
// 热门搜索数据
const popularSearches = ref(['热门搜索1', '热门搜索2', '热门搜索3']);

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
  // 将搜索记录添加到历史搜索中
  if (searchText.value.trim() !== '' && !historySearches.value.includes(searchText.value)) {
    historySearches.value.unshift(searchText.value);
    // 限制历史记录最多保存10个
    if (historySearches.value.length > 10) {
      historySearches.value.pop();
    }
  }
};

const onSelectSuggestion = (item) => {
  console.log('用户选择了推荐结果:', item);
  searchText.value = item;
};

const onSelectHistory = (item) => {
  console.log('用户选择了历史搜索:', item);
  searchText.value = item;
  onSearch();
};

const onSelectPopular = (item) => {
  console.log('用户选择了热门搜索:', item);
  searchText.value = item;
  onSearch();
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
.history-tags, .popular-tags {
  margin-top: 16px;
}
.section-title {
  font-weight: bold;
  margin-bottom: 8px;
  display: block;
}
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}
.tag {
  padding: 6px 12px;
  background-color: #e0e0e0;
  border-radius: 16px;
  cursor: pointer;
}
.tag:hover {
  background-color: #d0d0d0;
}
</style>
