<template>
  <view>
    <!-- Header Section -->
	<image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view class="header">
      <button @click="showSection('全部')">全部</button>
      <button @click="showSection('环保科普')">环保科普</button>
      <button @click="showSection('环保要闻')">环保要闻</button>
    </view>

    <!-- News Section -->
    <view class="news-section">
      <view v-for="(item, index) in newsItems" :key="index" class="news-item" @click="navigateTo(item.link)">
        <view class="news-title">{{ item.title }}</view>
        <view v-if="item.image" class="news-image">
          <image :src="item.image" :alt="item.title" mode="widthFix" />
        </view>
        <view class="news-description">{{ item.description }}</view>
      </view>
    </view>
  </view>
</template>

<script>
export default {
  data() {
    return {
      newsItems: [],
    };
  },
  methods: {
    // 控制新闻分类的显示
    showSection(section) {
      if (section === "全部") {
        this.newsItems = [
          {
            title: "国际氢能联盟和麦肯锡联合发布《氢能洞察2024》",
            description: "环保要闻  |  双碳指挥  |  刚刚",
            link: "",
          },
          {
            title: "把自然讲给你听 | 什么是森林？",
            description: "环保科普  |  环保科普365  |  1小时前",
            image: "",
            link: "https://mp.weixin.qq.com/s/mzFR2d-17men_Lm297fweQ",
          },
          {
            title: "视频 | 垃圾分类",
            description: "环保科普  |  环保科普365  |  2024-10-14",
            video: true,
            link: "/pages/video_details/video_details",
          },
        ];
      } else if (section === "环保科普") {
        this.newsItems = [
          {
            title: "把自然讲给你听 | 什么是森林？",
            description: "环保科普  |  环保科普365  |  1小时前",
            image: "@/static/images/nature.jpg",
            link: "https://mp.weixin.qq.com/s/mzFR2d-17men_Lm297fweQ",
          },
          {
            title: "视频 | 垃圾分类",
            description: "环保科普  |  环保科普365  |  2024-10-14",
            video: true,
            link: "/pages/video_details/video_details",
          },
        ];
      } else if (section === "环保要闻") {
        this.newsItems = [
          {
            title: "国际氢能联盟和麦肯锡联合发布《氢能洞察2024》",
            description: "环保要闻  |  双碳指挥  |  刚刚",
            link: "/pages/details/details",
          },
        ];
      }
    },

    // 页面跳转方法
    navigateTo(link) {
      if (link.startsWith("http")) {
        // 外部链接跳转
        uni.navigateTo({
          url: `/pages/webview/webview?url=${encodeURIComponent(link)}`,
        });
      } else {
        // 内部页面跳转
        uni.navigateTo({
          url: link,
        });
      }
    },
  },

  mounted() {
    // 默认加载“全部”新闻
    this.showSection("全部");
  },
};
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
		/* 将背景图片置于最底层 */
		opacity: 0.1;
		/* 调整透明度以不干扰内容 */
	}

/* Header Section */
.header {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #ffffff;
  border-bottom: 1px solid #e0e0e0;
  justify-content: space-around;
}

.header button {
  border: none;
  background-color: transparent;
  font-size: 16px;
  cursor: pointer;
}

/* News Section */
.news-section {
  padding: 20px;
}

.news-item {
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  padding: 15px;
  margin-bottom: 20px;
  cursor: pointer;
}

.news-item image {
  width: 100%;
  height: auto;
  border-radius: 5px;
}

.news-title {
  font-size: 18px;
  font-weight: bold;
  margin-bottom: 10px;
}

.news-description {
  font-size: 14px;
  margin-bottom: 10px;
}

/* Footer Section */
.footer {
  background-color: #ffffff;
  padding: 10px 0;
  border-top: 1px solid #e0e0e0;
  position: fixed;
  bottom: 0;
  width: 100%;
}

.footer-nav {
  display: flex;
  justify-content: space-around;
}

.nav-item {
  text-decoration: none;
  color: #333;
  font-weight: bold;
  cursor: pointer;
}

.nav-item:hover {
  color: #4caf50;
}
</style>
