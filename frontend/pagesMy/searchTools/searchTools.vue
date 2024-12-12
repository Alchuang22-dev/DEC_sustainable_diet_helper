<template>
  <view class="search-container">
    <view class="header">
      <input class="search-box" v-model="searchText" @input="onSearchInput" :placeholder="placeholderText" />
      <button class="search-button" @click="onSearch"> {{$t('text_search')}} </button>
    </view>
    <view v-if="suggestions.length > 0" class="suggestions">
      <view v-for="(item, index) in suggestions" :key="index" class="suggestion-item" @click="onSelectSuggestion(item)">
        <span v-html="highlightMatch(item.title)"></span>
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
const historySearches = ref([]);
// 热门搜索数据（后端更新）
const popularSearches = ref(['热门搜索1', '热门搜索2', '热门搜索3']);

// 模拟与后端通信的函数
const fetchSuggestions = async (query) => {
  // 模拟从 pages.json 获取页面数据
  const pagesData = [
    { path: "/pagesTool/home_servant/home_servant", title: "家庭管理" },
    { path: "/pagesTool/carbon_calculator/carbon_calculator", title: "碳计算器" },
	{ path: "/pagesTool/food_recommend/food_recommend", title: "食谱推荐" },
	{ path: "/pagesTool/nutrition_calculator/nutrition_calculator", title: "营养日历"},
	{ path: "/pagesNews/create_news/create_news", title: "图文编辑"},
	{ path: "/pagesNews/create_news/create_news", title: "发表图文"},
	{ path: "/pagesSetting/Bend/Bend", title: "绑定账号"},
	{ path: "/pagesSetting/ConnectUs/ConnectUs", title: "联系我们"},
	{ path: "/pagesSetting/SoftwareInfo/SoftwareInfo", title: "软件信息"},
	{ path: "/pagesSetting/language/language", title: "多语言"},
	{ path: "/pagesSetting/DeleteData/DeleteData", title: "删除数据"},
	{ path: "/pagesSetting/DeleteId/DeleteId", title: "删除账号"},
    // 添加更多页面数据
	// 此处未来将添加新闻数据
  ];

  // 模拟后端通信，您可以根据实际需求替换为后端API调用
  return new Promise((resolve) => {
    setTimeout(() => {
      const matchedPages = pagesData.filter(page => page.title.includes(query)); // 根据标题进行模糊匹配
      resolve(matchedPages);
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
  console.log('搜索：', searchText.value);
  if (searchText.value.trim() !== '' && !historySearches.value.includes(searchText.value)) {
    historySearches.value.unshift(searchText.value);
    // 限制历史记录最多保存10个
    if (historySearches.value.length > 10) {
      historySearches.value.pop();
    }
  }

  // 如果用户点击了建议，执行跳转逻辑
  if (suggestions.value.length === 1) {
    const selectedPage = suggestions.value[0]; // 获取选中的页面
    uni.navigateTo({
      url: `${selectedPage.path}`  // 根据路径跳转
    });
  }
};

const onSelectSuggestion = (item) => {
  console.log('用户选择了推荐结果:', item);
  searchText.value = item.title;  // 填充输入框内容为页面标题
  // 触发搜索并跳转到相应页面
  if (searchText.value.trim() !== '' && !historySearches.value.includes(searchText.value)) {
    historySearches.value.unshift(searchText.value);
    // 限制历史记录最多保存10个
    if (historySearches.value.length > 10) {
      historySearches.value.pop();
    }
  }
  
  uni.navigateTo({
    url: `${item.path}`
  });
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
