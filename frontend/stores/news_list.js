// stores/news_list.js
import { defineStore } from 'pinia'
import { computed } from 'vue'
import { useUserStore } from './user.js'

const BASE_URL = 'http://xcxcs.uwdjl.cn:8080'

// 统一使用与 carbon_and_nutrition_data.js 类似的请求封装
function createRequestConfig(config) {
  const userStore = useUserStore()
  return {
    ...config,
    header: {
      Authorization: `Bearer ${userStore.user.token || ''}`,
      ...(config.header || {})
    }
  }
}

function request(config) {
  const userStore = useUserStore()
  return new Promise((resolve, reject) => {
    uni.request({
      ...config,
      success: async res => {
        if (res.statusCode === 401) {
          try {
            await userStore.refreshToken()
            uni.request({
              ...config,
              header: {
                ...config.header,
                Authorization: `Bearer ${userStore.user.token}`
              },
              success: res2 => {
                if (res2.statusCode === 401) reject(new Error('Unauthorized'))
                else resolve(res2)
              },
              fail: err2 => reject(err2)
            })
          } catch {
            uni.navigateTo({ url: '/pagesMy/login/login' })
            reject(new Error('Unauthorized'))
          }
        } else {
          resolve(res)
        }
      },
      fail: err => reject(err)
    })
  })
}

export const useNewsStore = defineStore('news', {
  state: () => ({
    allNewsItems: [],
    rawNewsData: [],
    isRefreshing: false,
    isLoading: true
  }),
  getters: {
    filteredNewsItems(state) {
      if (state.isLoading) return []
      return state.allNewsItems
    }
  },
  actions: {
    refreshNews() {
      this.isRefreshing = true
      setTimeout(() => {
        this.allNewsItems = this.allNewsItems.sort(() => Math.random() - 0.5)
        this.isRefreshing = false
      }, 1000)
    },
    async fetchNews(page = 1, type = 'top-views', text = '') {
      this.allNewsItems = []
      this.isLoading = true
      const loadingTimeout = setTimeout(() => {
        this.isLoading = false
      }, 3000)

      try {
        let url = ''
        let data = {}
        let ifSearch = false
        switch (type) {
          case 'top-views':
            url = `${BASE_URL}/news/paginated/view_count?page=${page}`
            break
          case 'top-likes':
            url = `${BASE_URL}/news/paginated/like_count?page=${page}`
            break
          case 'latest':
            url = `${BASE_URL}/news/paginated/upload_time?page=${page}`
            break
          case 'favorite':
            url = `${BASE_URL}/users/favorited`
            break
          case 'viewed':
            url = `${BASE_URL}/users/viewed`
            break
          case 'search':
            url = `${BASE_URL}/news/search`
            data = { query: text }
            ifSearch = true
            break
          default:
            throw new Error('Invalid news type')
        }

        if (ifSearch) {
          const res = await request(
            createRequestConfig({
              url,
              method: 'POST',
              data
            })
          )
          const newsResults = res.data?.results
          if (Array.isArray(newsResults) && newsResults.length) {
            const newsIds = newsResults.map(item => item.id)
            const newsDetails = await Promise.all(newsIds.map(id => this.getNewsDetails(id)))
            this.rawNewsData = newsDetails.filter(detail => detail)
            this.rawNewsData.forEach(this.convertNewsToItems)
          }
        } else {
          const res = await request(createRequestConfig({ url, method: 'GET' }))
          const newsIds = res.data?.news_ids || []
          if (newsIds.length) {
            const newsDetails = await Promise.all(newsIds.map(id => this.getNewsDetails(id)))
            this.rawNewsData = newsDetails
            this.rawNewsData.forEach(this.convertNewsToItems)
          }
        }
      } catch (error) {
        console.error('Error fetching data:', error)
      } finally {
        this.isLoading = false
        clearTimeout(loadingTimeout)
      }
    },
    async getNewsDetails(id) {
      try {
        const res = await request(
          createRequestConfig({
            url: `${BASE_URL}/news/details/news/${id}`,
            method: 'GET'
          })
        )
        return res.data
      } catch (error) {
        console.error('Error fetching article details', error)
        return null
      }
    },
    convertNewsToItems(news) {
      const formatPublishTime = publishTime => {
        const date = new Date(publishTime)
        const now = new Date()
        const isSameDay =
          date.getFullYear() === now.getFullYear() &&
          date.getMonth() === now.getMonth() &&
          date.getDate() === now.getDate()
        if (isSameDay) {
          const hours = String(date.getHours()).padStart(2, '0')
          const minutes = String(date.getMinutes()).padStart(2, '0')
          return `今天 ${hours}:${minutes}`
        } else {
          const year = date.getFullYear()
          const month = String(date.getMonth() + 1).padStart(2, '0')
          const day = String(date.getDate()).padStart(2, '0')
          return `${year}-${month}-${day}`
        }
      }

      const formattedNews = {
        id: news.id,
        link: 'news_detail',
        title: news.title,
        description: `${formatPublishTime(news.upload_time)}`,
        info: `阅读量: ${news.view_count} | 点赞量: ${news.like_count}`,
        form: news.form
      }
      this.allNewsItems.push(formattedNews)
    }
  }
})