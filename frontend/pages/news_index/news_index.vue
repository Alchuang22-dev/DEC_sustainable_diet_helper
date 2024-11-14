<template>
  <view @touchstart="refreshPage">
    <!-- Header Section -->
    <image src="/static/images/index/background_img.jpg" class="background-image"></image>
    <view class="header">
      <button
        @click="showSection('全部')"
        :class="{ active: selectedSection === '全部' }"
      >
        全部
      </button>
      <button
        @click="showSection('环保科普')"
        :class="{ active: selectedSection === '环保科普' }"
      >
        环保科普
      </button>
      <button
        @click="showSection('环保要闻')"
        :class="{ active: selectedSection === '环保要闻' }"
      >
        环保要闻
      </button>
	  <button
	    @click="showSection('环保专栏')"
	    :class="{ active: selectedSection === '环保专栏' }"
	  >
	    环保专栏
	  </button>
    </view>

    <!-- News Section -->
    <view class="news-section">
      <view
        v-for="(item, index) in newsItems"
        :key="index"
        :class="['news-item', { active: activeIndex === index }]"
        @click="navigateTo(item.link,item.title)"
        @touchstart="pressFeedback(index)"
        @touchend="releaseFeedback()"
      >
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
	      allNewsItems: [
	        {
	          title: "国际氢能联盟和麦肯锡联合发布《氢能洞察2024》",
	          description: "环保要闻  |  双碳指挥  |  刚刚",
	          link: "news_detail",
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
	          link: "video_detail",
	        },
	        {
	          title: "联合国发布2024气候计划",
	          description: "环保要闻  |  环保科普365  |  2024-10-14",
	          video: true,
	          link: "news_detail",
	        },
			{
			  title: "专栏 | 寂静的春天",
			  description: "环保专栏  |  爱读夜  |  2024-10-14",
			  video: true,
			  link: "news_detail",
			},
			{
			  title: "社团招新 | 根与芽2025",
			  description: "环保专栏  |  公益事业  |  2024-6-16",
			  video: true,
			  link: "news_detail",
			},
	      ],
	      newsItems: [],
	      activeIndex: null,
	      selectedSection: '全部', // 默认选择“全部”
	      isRefreshing: false, // 用于显示正在更新的状态
	    };
	  },
	  async created() {
	    // 在组件创建时调用后端获取数据
	    await this.fetchData();
	  },
	  methods: {
		async fetchData() {
		    try {
		        uni.request({ // 模拟从后端获取数据
		            url: 'https://122.51.231.155/news/{id}', // 模拟的后端接口URL
		            method: 'GET',
		            data: {
		                id: 1,
		            },
		            success: (res) => {
		                const mockResponse = {
		                    data: [{
		                            id: 1,
		                            form: 'news',
		                            newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
		                            imgsSrc: '',
		                            tabs: ['环境保护', '环保要闻'],
		                            time: '2024-4-17',
		                            newsName: '测试视频',
		                            authorName: 'user_test',
		                            authorAvatar: '',
		                            newsinfo: '测试测试测试测试测试',
		                            newsbody: '9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。',
		                            likeCount: 1001,
		                            shareCount: 37,
		                            favoriteCount: 897,
		                            followCount: 189,
		                            dislikeCount: 100,
		                            type: 'main'
		                        },
		                        {
		                            id: 2,
		                            form: 'news',
		                            newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
		                            imgsSrc: '',
		                            tabs: ['环境保护', '环保要闻'],
		                            time: '2024-4-17',
		                            newsName: '测试新闻',
		                            authorName: '中野梓',
		                            authorAvatar: '',
		                            newsinfo: '测试测试测试测试测试',
		                            newsbody: '',
		                            likeCount: 1001,
		                            shareCount: 37,
		                            favoriteCount: 897,
		                            followCount: 189,
		                            dislikeCount: 100,
		                            type: 'reco'
		                        },
		                    ]
		                };
		                this.videoData = mockResponse.data;
		                this.recommendations = [];
		                this.videoData.forEach(video => this.convertVideoToItems(video));
		            },
		            fail: (err) => {
		                const mockResponse = {
		                    data: [{
		                            id: 1,
		                            form: 'news',
		                            newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
		                            imgsSrc: '',
		                            tabs: ['环境保护', '环保要闻'],
		                            time: '2024-4-17',
		                            newsName: '测试视频',
		                            authorName: 'user_test',
		                            authorAvatar: '',
		                            newsinfo: '测试测试测试测试测试',
		                            newsbody: '9月17日，国际氢能联盟与麦肯锡联合发布《氢能洞察2024》，分析了全球氢能行业在过去一年的重要进展。该报告显示，全球氢能项目投资显著增长，氢能在清洁能源转型中扮演了重要角色。',
		                            likeCount: 1001,
		                            shareCount: 37,
		                            favoriteCount: 897,
		                            followCount: 189,
		                            dislikeCount: 100,
		                            type: 'main'
		                        },
		                        {
		                            id: 2,
		                            form: 'news',
		                            newsSrc: 'http://vjs.zencdn.net/v/oceans.mp4',
		                            imgsSrc: '',
		                            tabs: ['环境保护', '环保要闻'],
		                            time: '2024-4-17',
		                            newsName: '测试新闻',
		                            authorName: '平泽唯',
		                            authorAvatar: '',
		                            newsinfo: '测试测试测试测试测试',
		                            newsbody: '',
		                            likeCount: 1001,
		                            shareCount: 37,
		                            favoriteCount: 897,
		                            followCount: 189,
		                            dislikeCount: 100,
		                            type: 'reco'
		                        },
		                    ]
		                };
		                this.videoData = mockResponse.data;
		                this.recommendations = [];
		                this.videoData.forEach(video => this.convertVideoToItems(video));
		            },
		        });
		    } catch (error) {
		        console.error('Error fetching data:', error);
		    }
		},
		convertVideoToItems(video) {
		  if (video.type === 'main') {
			if(video.form === 'web') {
				this.allNewsItems.push({
				  id: video.id,
				  link: video.newsSrc,
				  image: '',
				  title: video.newsName,
				  description: video.tabs.join(' | ') + ' | '+ video.time,
				  info: '阅读量: ' + video.followCount + ' | 点赞量: ' + video.likeCount,
				  form: video.form,
				});
			}
			else if(video.form === 'news') {
				this.allNewsItems.push({
				  id: video.id,
				  link: 'news_detail',
				  image: '',
				  title: video.newsName,
				  description: video.tabs.join(' | ') + ' | '+ video.time,
				  info: '阅读量: ' + video.followCount + ' | 点赞量: ' + video.likeCount,
				  form: video.form,
				});
			}
			else if(video.form === 'video') {
				this.allNewsItems.push({
				  id: video.id,
				  link: 'video_detail',
				  image: '',
				  title: video.newsName,
				  description: video.tabs.join(' | ') + ' | '+ video.time,
				  info: '阅读量: ' + video.followCount + ' | 点赞量: ' + video.likeCount,
				  form: video.form,
				});
			}
		  }
		},
	    // 控制新闻分类的显示
	    showSection(section) {
	      this.selectedSection = section; // 更新选中的分类
	      if (section === "全部") {
	        this.newsItems = this.allNewsItems;
	      } else {
	        this.newsItems = this.allNewsItems.filter(item => item.description.includes(section));
	      }
	    },

	    // 页面跳转方法
	    navigateTo(link,name) {
	      setTimeout(() => {
	        if (link.startsWith("http")) {
	          // 外部链接跳转
	          uni.navigateTo({
	            url: `/pagesNews/web_detail/web_detail?url=${encodeURIComponent(link)}`,
	          });
	        } else {
	          // 内部页面跳转
	          uni.navigateTo({
	            url: `/pagesNews/${link}/${link}?title=${name}`,
	          });
	        }
	      }, 100); // 延迟 100 毫秒
	    },

	    // 触摸反馈
	    pressFeedback(index) {
	      console.log("Press feedback:", index);
	      this.activeIndex = index;
	    },

	    releaseFeedback() {
	      console.log("Release feedback");
	      this.activeIndex = null;
	    },

	    // 页面更新方法
	    refreshPage() {
	      this.isRefreshing = true;
	      setTimeout(() => {
	        this.newsItems = this.newsItems.sort(() => Math.random() - 0.5);
	        this.isRefreshing = false;
	      }, 1000); // 模拟1秒的加载时间
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
  opacity: 0.1;
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
  transition: color 0.3s;
}

.header button.active {
  color: #4caf50; /* 选中状态颜色 */
  font-weight: bold; /* 选中状态加粗 */
}

/* News Section */
.news-section {
  padding: 20px;
  position: relative;
}

.news-section::before {
  content: "正在更新...";
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  font-size: 16px;
  color: #333;
  display: none;
}

.news-section[data-refreshing="true"]::before {
  display: block;
}

.news-item {
  background-color: #ffffff;
  border-radius: 10px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  padding: 15px;
  margin-bottom: 20px;
  cursor: pointer;
  transition: transform 0.1s, box-shadow 0.1s;
  position: relative;  /* 确保其层级设置有效 */
  z-index: 1;          /* 确保点击事件可以被接收 */
}

.news-item.active {
  transform: scale(0.98);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  background-color: #e6f7ff;
}

.news-image {
  pointer-events: none; /* 确保图片不会阻止父元素的点击事件 */
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
