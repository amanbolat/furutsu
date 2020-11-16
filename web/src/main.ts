import Vue from 'vue'
import App from './App.vue'
import createRouter from './router'
import store from './store'
import vuetify from './plugins/vuetify'
import './style/index.scss'
import Notifications  from 'vue-notification'
import './filters'
import NumberInput from '@chenfengyuan/vue-number-input'

Vue.config.productionTip = false
Vue.use(Notifications)
Vue.use(NumberInput)

new Vue({
  router: createRouter(),
  store,
  vuetify,
  render: h => h(App),
  async created() {
    (this.$router as any).start()
  }
}).$mount('#app')
