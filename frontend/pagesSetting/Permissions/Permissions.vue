<template>
  <view class="container">
    <view v-for="(section, sectionIndex) in sections" :key="sectionIndex">
      <view class="divider">
        <text class="divider-text">{{ section.title }}</text>
      </view>
      <view
        v-for="(item, index) in section.items"
        :key="index"
        class="item-container"
      >
        <view class="item" @click="toggleItem(sectionIndex, index)">
          <text class="item-title">{{ item.name }}</text>
          <text class="arrow">{{ expanded[sectionIndex][index] ? '▲' : '▼' }}</text>
        </view>
        <view v-if="expanded[sectionIndex][index]" class="item-details">
          <view class="detail">
            <text class="detail-label">内容:</text>
            <text class="detail-value">[内容]</text>
          </view>
          <view class="detail">
            <text class="detail-label">获取目的:</text>
            <text class="detail-value">[获取目的]</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
/* ----------------- Imports ----------------- */
import { reactive } from 'vue'

/* ----------------- Setup ----------------- */

const sections = [
  {
    title: '用户基本信息',
    items: [
      { name: '头像' },
      { name: '昵称' },
      { name: '手机号' },
      { name: '地区' },
      { name: '地址' }
    ]
  },
  {
    title: '用户使用过程信息',
    items: [{ name: '位置' }, { name: '图片和视频' }]
  },
  {
    title: '用户输入信息',
    items: [{ name: '偏好' }, { name: '上传图文' }]
  }
]

const expanded = reactive(sections.map((section) => section.items.map(() => false)))

/* ----------------- Methods ----------------- */
function toggleItem(sectionIndex, itemIndex) {
  expanded[sectionIndex][itemIndex] = !expanded[sectionIndex][itemIndex]
}
</script>

<style scoped>
.container {
  padding: 20px;
  background-color: #f5f5f5;
}

.divider {
  position: relative;
  height: 20px;
  margin: 20px 0;
}
.divider::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 0;
  right: 0;
  height: 1px;
  background-color: #ebebeb;
  transform: translateY(-50%);
}
.divider-text {
  position: relative;
  z-index: 1;
  background-color: #f8f8f8;
  padding: 0 10px;
  text-align: center;
  font-size: 14px;
  color: #888;
}

.item-container {
  margin: 10px 0;
}

.item {
  display: flex;
  justify-content: space-between;
  padding: 12px;
  background-color: #fff;
  border-radius: 8px;
  align-items: center;
}

.item:hover {
  background-color: #f0f0f0;
}

.item-title {
  font-size: 16px;
  color: #333;
}

.arrow {
  font-size: 16px;
  color: #999;
}

.item-details {
  padding: 10px 12px;
  background-color: #fafafa;
  border-left: 2px solid #007aff;
  border-radius: 0 8px 8px 0;
  margin-top: 5px;
}

.detail {
  display: flex;
  margin-bottom: 8px;
}

.detail:last-child {
  margin-bottom: 0;
}

.detail-label {
  width: 100px;
  font-weight: bold;
  color: #555;
}

.detail-value {
  color: #777;
}
</style>