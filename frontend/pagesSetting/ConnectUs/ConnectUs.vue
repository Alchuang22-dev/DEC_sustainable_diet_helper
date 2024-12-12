<template>
  <view class="author-page">
    <view v-for="(author, index) in authors" :key="index" class="author-card">
      <image :src="author.avatar" class="avatar"></image>
      <view class="author-info">
        <text class="name">{{ author.name }}</text>
        <text class="organization">{{ author.organization }}</text>
        <text class="position">{{ author.position }}</text>
        <text class="email">{{ author.email }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'

// 作者信息列表，可以根据需要添加更多作者
const authors = ref([
  {
    avatar: 'https://avatars.githubusercontent.com/u/111627783?s=400&u=3d7dc2dc83fa9a03e03c71cac028ecec00ce01d0&v=4', // 头像占位符
    name: 'Zeyu Zhang', // 姓名占位符
    organization: 'Tsinghua University', // 所在机构占位符
    position: '    Frontend Developer', // 工作岗位占位符
    email: 'alchuang671@gmail.com' // 联系邮箱占位符
  },
  {
    avatar: 'https://avatars.githubusercontent.com/u/108730633?v=4',
    name: 'Yuanxu Sun',
    organization: 'Tsinghua University',
    position: '    Data Manager',
    email: 'lisi@example.com'
  },{
    avatar: 'https://avatars.githubusercontent.com/u/125715741?v=4',
    name: 'Shenyu Qiao',
    organization: 'Tsinghua University',
    position: '    Backend Main Developer',
    email: 'lisi@example.com'
  },
  {
    avatar: 'https://avatars.githubusercontent.com/u/116412601?v=4',
    name: 'Xingjian Su',
    organization: 'Tsinghua University',
    position: '    Frontend Main Developer',
    email: 'lisi@example.com'
  },
  // 可以继续添加更多作者信息
])
</script>

<style>
.author-page {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px;
  background-color: #f5f5f5;
}

.author-card {
  width: 90%;
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
  padding: 15px;
  display: flex;
  align-items: center;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 40px;
  margin-right: 15px;
}

.author-info {
  flex: 1;
  display: flex; /* 添加 flex 布局 */
  flex-direction: column; /* 设置为列方向 */
}

.name {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 5px;
  color: #333333;
}

.organization {
  font-size: 14px;
  color: #666666;
  margin-bottom: 3px;
}

.position {
  font-size: 14px;
  color: #666666;
  margin-bottom: 3px;
}

.email {
  font-size: 14px;
  color: #1E90FF;
}
</style>
