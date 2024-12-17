// stores/news_list.js
import { defineStore } from 'pinia';
import { computed } from "vue";
import { useUserStore } from "./user.js";

const userStore = useUserStore(); // 使用用户存储
const jwtToken = computed(() => userStore.user.token); // Replace with actual token
const systemDate = new Date();
const systemDateStr = systemDate.toISOString().slice(0, 10); // 获取当前系统日期，格式：YYYY-MM-DD

const BASE_URL = 'http://122.51.231.155:8080'; //url基本路径

// Helper function to format publish time
const formatPublishTime = (publishTime) => {
  const date = new Date(publishTime);
  const dateStr = date.toISOString().slice(0, 10);
  if (dateStr === systemDateStr) {
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `今天 ${hours}:${minutes}`;
  } else {
    return dateStr;
  }
};

export const useNewsStore = defineStore('news', {
  state: () => ({
    allNewsItems: [],
    rawNewsData: [], // Raw backend data
    isRefreshing: false, // Whether the page is refreshing
    isLoading: true, // Flag to track loading state
  }),
  getters: {
    filteredNewsItems(state) {
      if (state.isLoading) {
        // You can return a loading state or an empty array while data is being fetched
        console.log('新闻项正在加载...');
        return []; // or you can return some loading indicator here
      }
      console.log('新闻项：', state.allNewsItems);
      return state.allNewsItems;
    },
  },
  actions: {
    // Refresh the news items with random order
    refreshNews() {
      this.isRefreshing = true;
      console.log("in refreshing");
      setTimeout(() => {
        this.allNewsItems = this.allNewsItems.sort(() => Math.random() - 0.5);
        this.isRefreshing = false;
      }, 1000);
    },

    // Fetch the news based on type (by views, likes, or latest)
    async fetchNews(page = 1, type = 'top-views') {
      console.log(`正在获取新闻, type: ${type}, page: ${page}`);
      this.allNewsItems = [];
      this.isLoading = true; // Start loading
      const loadingTimeout = setTimeout(() => {
        // If the data hasn't been loaded after 3 seconds, set isLoading to false
        this.isLoading = false;
        console.log('加载超时');
      }, 3000); // Set the timeout to 3 seconds

      try {
        let url = '';
        switch (type) {
          case 'top-views':
            url = `${BASE_URL}/news/paginated/view_count?page=${page}`;
            break;
          case 'top-likes':
            url = `${BASE_URL}/news/paginated/like_count?page=${page}`;
            break;
          case 'latest':
            url = `${BASE_URL}/news/paginated/upload_time?page=${page}`;
            break;
          default:
            throw new Error('Invalid news type');
        }

        const res = await uni.request({
          url: url,
          method: 'GET',
          header: {
            "Content-type": "application/json",
            "Authorization": `Bearer ${jwtToken.value}`,
          },
        });
        console.log('后端返回：', res);
        const newsIds = res.data.news_ids;
        console.log('获取的新闻id:', newsIds);
        if (newsIds.length) {
          // Get details for each news ID
          const newsDetails = await Promise.all(newsIds.map(id => this.getNewsDetails(id)));
          this.rawNewsData = newsDetails;
          this.rawNewsData.forEach(this.convertNewsToItems);
        }
        this.isLoading = false; // Finish loading
        clearTimeout(loadingTimeout); // Clear the timeout if data is loaded
      } catch (error) {
        console.error('Error fetching data:', error);
        this.isLoading = false; // Finish loading even on error
        clearTimeout(loadingTimeout); // Clear the timeout on error
      }
    },

    // Fetch detailed information for a specific news item by ID
    async getNewsDetails(id) {
      const url = `${BASE_URL}/news/details/news/${id}`;
      try {
        const res = await uni.request({
          url: url,
          method: 'GET',
          header: {
            'Authorization': `Bearer ${jwtToken.value}`,
          },
        });
        console.log('获取详细信息:', res);
        return res.data;
      } catch (error) {
        console.error('Error fetching article details', error);
        return null;
      }
    },

    // Convert raw news data into a usable format for the UI
    convertNewsToItems(news) {
        const formattedNews = {
          id: news.id,
          link: 'news_detail',
          title: news.title,
          description: `${formatPublishTime(news.upload_time)}`,
          info: `阅读量: ${news.followCount} | 点赞量: ${news.likeCount}`,
          form: news.form,
        };
		 console.log('正在将新闻放入数组...',formattedNews);
        this.allNewsItems.push(formattedNews);
    },
  },
});
