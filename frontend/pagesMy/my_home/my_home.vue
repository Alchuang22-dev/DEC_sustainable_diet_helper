<template>
	<image src="/static/images/index/background_img.jpg" class="background-image"></image>
	<view class="header">
	  <text class="title">{{$t('menu_creations')}}</text>
	  <view class="header-actions">
	    <button class="menu-icon"></button>
	    <button class="camera-icon"></button>
	  </view>
	</view>
  <view class="container">
    <!-- 图文卡片列表 -->
    <view class="card-list">
      <view
        v-for="(item, index) in articles"
        :key="index"
        class="card"
        :style="{ backgroundColor: item.bgColor }"
      >
        <view class="card-header">
          <view class="title">{{ item.title }}</view>
          <view class="status">{{ item.status }}</view>
        </view>
        <view class="card-body">
          <view class="description">{{ item.description }}</view>
          <view class="info">
            <text class="publish-time">{{ item.publishTime }}</text>
            <view class="stats">
              <text class="like-count">👍 {{ item.likes }}</text>
              <text class="favorite-count">⭐ {{ item.favorites }}</text>
              <text class="share-count">🔗 {{ item.shares }}</text>
            </view>
          </view>
        </view>
        <view class="card-footer">
          <button @click="viewArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/view.svg" class="icon" alt="View" ></image>
          </button>
          <button @click="editArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/edit.svg" class="icon" alt="Edit" ></image>
          </button>
          <button @click="deleteArticle(index)" class="action-btn">
            <image src="@/pagesMy/static/delete.svg" class="icon" alt="Delete" ></image>
          </button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n';
//import uni from '@dcloudio/uni-app';
import { useStore } from 'vuex'; // 引入 Vuex 的 useStore

const { t, locale, messages } = useI18n();

// 模拟数据
const articles = ref([
  {
    title: '如何保持健康饮食',
    description: '这是一篇关于健康饮食的文章。',
    publishTime: '2024-12-06 12:30',
    likes: 120,
    favorites: 50,
    shares: 30,
    status: '已发布', // 文章状态
    bgColor: 'rgba(0, 123, 255, 0.1)' // 背景颜色
  },
  {
    title: '学习Vue3的基本概念',
    description: '本文介绍了Vue3的一些基础概念。',
    publishTime: '2024-12-05 10:00',
    likes: 80,
    favorites: 30,
    shares: 20,
    status: '草稿', // 文章状态
    bgColor: 'rgba(255, 193, 7, 0.1)' // 背景颜色
  },
  // 更多文章数据...
])

// 查看文章
const viewArticle = (index) => {
  console.log('查看文章:', articles.value[index])
  // 跳转到文章详情页
}

// 编辑文章
const editArticle = (index) => {
  console.log('编辑文章:', articles.value[index])
  uni.navigateTo({
  	url: "/pagesNews/edit_draft/edit_draft",
  })
  // 跳转到编辑页面
}

// 删除文章
const deleteArticle = (index) => {
  console.log('删除文章:', articles.value[index])
  // 执行删除操作
  articles.value.splice(index, 1) // 从数据中删除
}
</script>

<style scoped>
/* Body */
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

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 16px;
  height: 60px;
  background-color: #fff;
  border-bottom: 1px solid #ebebeb;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.header-actions button {
  background: none;
  border: none;
}

.container {
  padding: 20px;
}

.card-list {
  display: flex;
  flex-direction: column;
}

.card {
  margin-bottom: 20px;
  border-radius: 8px;
  padding: 15px;
  background-color: #fff;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1); /* 添加阴影效果 */
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.title {
  font-size: 18px;
  font-weight: bold;
}

.status {
  font-size: 14px;
  color: #007bff;
}

.card-body {
  margin-bottom: 15px;
}

.description {
  font-size: 14px;
  color: #555;
  margin-bottom: 10px;
}

.info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #777;
}

.stats {
  display: flex;
  gap: 10px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.action-btn {
  background-color: transparent;
  border: none;
  cursor: pointer;
  padding: 5px;
}

.icon {
  width: 24px;
  height: 24px;
  transition: transform 0.2s ease;
}

.icon:hover {
  transform: scale(1.2); /* 鼠标悬浮时放大图标 */
}

.publish-time {
  font-size: 12px;
  color: #777;
}

.like-count, .favorite-count, .share-count {
  font-size: 12px;
  color: #777;
}
</style>
