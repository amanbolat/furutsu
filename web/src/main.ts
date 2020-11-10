import Vue from 'vue'
import App from './App.vue'
import createRouter from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import './style/index.scss'
import Notifications  from 'vue-notification'
import {InputNumber} from 'ant-design-vue'
import 'ant-design-vue/lib/input-number/style/index.css'
import './filters'

Vue.config.productionTip = false
Vue.use(Notifications)
Vue.use(InputNumber)

new Vue({
  router: createRouter(),
  store,
  vuetify,
  render: h => h(App),
  async created() {
    (this.$router as any).start()
  }
}).$mount('#app')
