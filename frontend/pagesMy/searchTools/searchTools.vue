<template>
  <view class="search-container">
    <view class="header">
      <input
        class="search-box"
        v-model="searchText"
        @input="onSearchInput"
        :placeholder="placeholderText"
      />
      <button class="search-button" @click="onSearch">
        {{ t('text_search') }}
      </button>
    </view>

    <!-- 联想搜索列表 -->
    <view v-if="suggestions.length > 0" class="suggestions">
      <view
        v-for="(item, index) in suggestions"
        :key="index"
        class="suggestion-item"
        @click="onSelectSuggestion(item)"
      >
        <span v-html="highlightMatch(item.title)"></span>
      </view>
    </view>

    <!-- 历史搜索 -->
    <view v-if="historySearches.length > 0" class="history-tags">
      <text class="section-title">历史搜索</text>
      <view class="tags">
        <view
          v-for="(item, index) in historySearches"
          :key="index"
          class="tag"
          @click="onSelectHistory(item)"
        >
          {{ item }}
        </view>
      </view>
    </view>

    <!-- 热门搜索 -->
    <view v-if="popularSearches.length > 0" class="popular-tags">
      <text class="section-title">热门搜索</text>
      <view class="tags">
        <view
          v-for="(item, index) in popularSearches"
          :key="index"
          class="tag"
          @click="onSelectPopular(item)"
        >
          {{ item }}
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

/* ----------------- Setup ----------------- */
const { t } = useI18n()

/* ----------------- Reactive & State ----------------- */
const searchText = ref('')
const suggestions = ref([])
const placeholderText = ref(t('text_search_label'))

// 历史搜索数据
const historySearches = ref([])
// 热门搜索数据
const popularSearches = ref(['家庭管理', '碳计算器', '营养日历'])

/* ----------------- Methods ----------------- */
// 模拟后台搜索
async function fetchSuggestions(query) {
  const pagesData = [
    { path: "/pagesTool/home_servant/home_servant", title: "家庭管理" },
    { path: "/pagesTool/carbon_calculator/carbon_calculator", title: "碳计算器" },
    { path: "/pagesTool/food_recommend/food_recommend", title: "食谱推荐" },
    { path: "/pagesTool/nutrition_calendar/nutrition_calendar", title: "营养日历"},
    { path: "/pagesNews/create_news/create_news", title: "图文编辑"},
    { path: "/pagesNews/create_news/create_news", title: "发表图文"},
    { path: "/pagesSetting/ConnectUs/ConnectUs", title: "联系我们"},
    { path: "/pagesSetting/SoftwareInfo/SoftwareInfo", title: "软件信息"},
    { path: "/pagesSetting/language/language", title: "多语言"},
  ]
  return new Promise((resolve) => {
    setTimeout(() => {
      const matchedPages = pagesData.filter(page => page.title.includes(query))
      resolve(matchedPages)
    }, 500)
  })
}

async function onSearchInput() {
  if (searchText.value.trim() !== '') {
    suggestions.value = await fetchSuggestions(searchText.value)
  } else {
    suggestions.value = []
  }
}

async function onSearch() {
  const query = searchText.value.trim()
  if (query === '') {
    uni.showToast({
      title: '请输入搜索内容',
      icon: 'none',
      duration: 2000
    })
    return
  }

  // 执行搜索
  const matchedPages = await fetchSuggestions(query)

  if (matchedPages.length > 0) {
    // 更新历史搜索
    if (!historySearches.value.includes(query)) {
      historySearches.value.unshift(query)
      if (historySearches.value.length > 10) {
        historySearches.value.pop()
      }
    }

    if (matchedPages.length === 1) {
      // 只有一个匹配，直接跳转
      const selectedPage = matchedPages[0]
      uni.navigateTo({
        url: selectedPage.path
      })
    } else {
      // 多个匹配，展示建议列表
      suggestions.value = matchedPages
      uni.showToast({
        title: '请选择一个',
        icon: 'none',
        duration: 2000
      })
    }
  } else {
    // 没有匹配，提示用户
    uni.showToast({
      title: '未找到匹配页面',
      icon: 'none',
      duration: 2000
    })
    suggestions.value = []
  }
}

function onSelectSuggestion(item) {
  searchText.value = item.title
  if (searchText.value.trim() !== '' && !historySearches.value.includes(searchText.value)) {
    historySearches.value.unshift(searchText.value)
    if (historySearches.value.length > 10) {
      historySearches.value.pop()
    }
  }
  uni.navigateTo({
    url: item.path
  })
}

function onSelectHistory(item) {
  searchText.value = item
  onSearch()
}

function onSelectPopular(item) {
  searchText.value = item
  onSearch()
}

function highlightMatch(item) {
  const query = searchText.value.trim()
  if (!query) return item
  const chars = query.split('').map(char => char.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'))
  const regex = new RegExp(`(${chars.join('|')})`, 'gi')
  return item.replace(regex, '<span style="color: green;">$1</span>')
}
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

.history-tags,
.popular-tags {
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