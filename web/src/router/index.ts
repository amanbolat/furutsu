import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'
import ProductList from '../views/ProductList.vue'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'ProductList',
    component: ProductList
  },
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
