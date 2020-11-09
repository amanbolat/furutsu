import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'
import ProductList from '../views/ProductList.vue'
import Cart from '../views/Cart.vue'
import OrderForm from '@/views/OrderForm.vue'
import OrderList from '@/views/OrderList.vue'
import PaymentForm from '@/views/PaymentForm.vue'
import store from '@/store'
import eventBus from '@/utils/event_bus'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'ProductList',
    component: ProductList
  },
  {
    path: '/cart',
    name: 'Cart',
    component: Cart,
    meta: {needAuth: true}
  },
  {
    path: '/order/:id',
    name: 'OrderForm',
    component: OrderForm,
    meta: {needAuth: true}
  },
  {
    path: '/order',
    name: 'OrderList',
    component: OrderList,
    meta: {needAuth: true}
  },
  {
    path: '/payment',
    name: 'PaymentForm',
    component: PaymentForm,
    meta: {needAuth: true}
  }
]

function createRouter(): VueRouter {
  const r = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
  })

  const promise = new Promise(resolve => {
    Object.assign(r, {
      start: resolve
    })
  })

  r.beforeEach(async (to, from, next) => {
    await promise
    const isAuthenticated = store.getters.isAuthenticated
    const needAuth = !!to.meta.needAuth
    if (needAuth) {
      if (!isAuthenticated) {
        console.log('Router: not authed')
        eventBus.$emit('prompt_auth', () => {
          next()
        })
      } else {
        next()
      }
    } else {
      next()
    }
  })

  return r
}


export default createRouter
